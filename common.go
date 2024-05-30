package lemon

type TaskID = int
type Result = int

type Task struct {
	ID                  TaskID
	Result              Result
	Throw               bool
	PrerequisiteTaskIDs []TaskID
}

type Tasks = []*Task

type BatchTaskResult struct {
	CompletedTasks             []*Task
	FailedTasks                []*Task
	CyclicDependenciesDetected bool
}
