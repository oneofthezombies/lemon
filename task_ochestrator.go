package lemon

import "errors"

type TaskOchestrator struct {
	taskEventReporter       *TaskEventReporter
	batchTaskStatusReporter *BatchTaskStatusReporter
}

func (o *TaskOchestrator) Run() (*BatchTaskResult, error) {
	return nil, errors.New("not implemented")
}
