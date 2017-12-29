## Package task
## About
The task is a package for periodic tasks
it contains a task interface and a scheduler


## Usage
```
import "github.com/zhlicen/task"

t := YourTaskImpl{}
task.NewScheduler(3 * time.Second, t)
task.Start()


// Trigger the task once manually
task.Trigger()
```
