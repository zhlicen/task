package task

// Task engine task
type Task interface {
	Name() string
	Run(args ...interface{})
}
