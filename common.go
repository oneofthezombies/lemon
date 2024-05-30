package lemon

type (
	TaskID = int
	Result = int
)

type Task struct {
	ID                  TaskID
	Result              Result
	Throw               bool
	PrerequisiteTaskIDs []TaskID
}

type Tasks = []*Task

type TaskResult struct {
	ID        TaskID
	IsSuccess bool
	Result    Result // Valid if IsSuccess == true
	Error     error  // Valid if IsSuccess == false
}

type BatchTaskResult struct {
	CompletedTasks             []*TaskResult
	FailedTasks                []*TaskResult
	CyclicDependenciesDetected bool
}
