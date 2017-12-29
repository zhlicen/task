package task

import (
	"reflect"
	"sutra/utils/debug/logs"
	"sutra/utils/msg/errors"
	"sync"
	"time"
)

type taskRunner struct {
	runChan    chan []interface{}
	quitChan   chan interface{}
	quitWg     sync.WaitGroup
	onceRunner sync.Once
	task       Task
	mutex      sync.Mutex
}

func newTaskRunner(task Task) *taskRunner {
	runner := &taskRunner{task: task}
	return runner
}

// start start of a runner can be called only once
// and calls over 1 time will be ignores
func (r *taskRunner) start() {
	taskType := reflect.TypeOf(r.task)
	runnerFun := func() {
		initChan := make(chan int)
		go func() {
			r.quitWg.Add(1)
			{
				r.mutex.Lock()
				r.runChan = make(chan []interface{})
				r.quitChan = make(chan interface{})
				r.mutex.Unlock()
			}
			initChan <- 0
			logs.Debug("Task %s[%s] started running.", taskType, r.task.Name())
		loop:
			for {
				select {
				case args := <-r.runChan:
					r.task.Run(args...)
				case <-r.quitChan:
					logs.Debug("Task %s[%s] quit msg received.", taskType, r.task.Name())
					break loop
				}
			}
			r.quitWg.Done()
			logs.Debug("Task %s[%s] quited successfully.", taskType, r.task.Name())
		}()
		<-initChan
	}
	r.onceRunner.Do(runnerFun)
}

// run task once
// may be blocked if the task is channel is occupied
// run call may cause panic if the task is stopped
func (r *taskRunner) runOnce(args ...interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	logs.Debug("Task %s[%s] triggered.", reflect.TypeOf(r.task), r.task.Name())
	select {
	case _, ok := <-r.quitChan:
		if !ok {
			// return err if the chan is closed
			return errors.New("runner is not started")
		}
	default:
	}
	r.runChan <- args
	logs.Debug("Task %s[%s] scheduled.", reflect.TypeOf(r.task), r.task.Name())
	return nil

}

// stop stop function can only be called once
// stop call on an stopped runner may cause panic
func (r *taskRunner) stop() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.quitChan <- ""
	r.quitWg.Wait()
	close(r.quitChan)
	close(r.runChan)
	logs.Debug("Task %s[%s] stopped.", reflect.TypeOf(r.task), r.task.Name())
}

type scheduler struct {
	cycle    time.Duration
	runner   *taskRunner
	quitChan chan interface{}
	loopOnce sync.Once
	stopWg   sync.WaitGroup
}

// NewScheduler create scheduler
func NewScheduler(cycle time.Duration, task Task) *scheduler {
	s := &scheduler{cycle: cycle}
	s.runner = newTaskRunner(task)
	s.quitChan = make(chan interface{})
	return s
}

// Start start loop
func (s *scheduler) Start() {
	s.runner.start()
	s.startTimeLoop()
}

// Stop stop loop
func (s *scheduler) Stop() {
	logs.Debug("Stop scheduler")
	s.quitChan <- ""
	s.runner.stop()
	s.stopWg.Wait()
}

// Trigger trigger the task manually
func (s *scheduler) Trigger(args ...interface{}) error {
	return s.runner.runOnce(args...)
}

func (s *scheduler) startTimeLoop() {
	timerFunc := func() {
		initChan := make(chan int)
		go func() {
			s.stopWg.Add(1)
			ticker := time.NewTicker(s.cycle)
			looperChan := ticker.C
			initChan <- 0
		loop:
			for {
				select {
				case <-looperChan:
					err := s.runner.runOnce()
					if err != nil {
						logs.Info(err.Error())
					}
				case <-s.quitChan:
					break loop
				}
			}
			ticker.Stop()
			s.stopWg.Done()
		}()
		<-initChan
	}
	s.loopOnce.Do(timerFunc)
}