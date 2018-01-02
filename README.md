## Package task
## About
The task is a package for periodic tasks
it contains a task interface and a scheduler


## Usage
```
import "github.com/zhlicen/task"

t := YourTaskImpl{}
s := task.NewScheduler(3 * time.Second, t)
s.Start()


// Trigger the task once manually
s.Trigger()
```
