package task

import (
	"log"
	"runtime"
	"testing"
	"time"
)

type TestTask struct {
	TaskName string
}

func (t *TestTask) Run(a ...interface{}) {
	log.Println("Test task is running")
	time.Sleep(500 * time.Millisecond)
}
func (t *TestTask) Name() string {
	return t.TaskName
}
func Test_taskRunner(t *testing.T) {
	r := newTaskRunner(&TestTask{TaskName: "TestTask"})
	go func() {
		r.start()
		for {
			err := r.runOnce()
			if err != nil {
				log.Println(err.Error())
				break
			}
			runtime.Gosched()
		}
	}()

	time.Sleep(3 * time.Second)
	r.stop()
}

func Test_taskScheduler(t *testing.T) {
	s := NewScheduler(1*time.Second, &TestTask{TaskName: "TestTask"})
	s.Start()
	time.Sleep(10 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
}
