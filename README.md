## Package task
## About
task is a package for periodic tasks
it contains a task interface and a scheduler


## Usage
```
import "10.204.28.137/samzho/task"

t := YourTaskImpl{}
task.NewScheduler(3 * time.Second, t)
task.Start()


// Trigger the task once manually
task.Trigger()
```
