package lemon

import (
	"slices"
	"strconv"
)

type PlainTextBatchTaskResultPrinter struct{}

var _ BatchTaskResultPrinter = (*PlainTextBatchTaskResultPrinter)(nil)

func (p *PlainTextBatchTaskResultPrinter) Print(result *BatchTaskResult) string {
	if result.CyclicDependenciesDetected {
		return "cyclic dependencies detected: true"
	}

	str := "completed tasks: ["
	slices.SortFunc(result.CompletedTasks, func(a *TaskResult, b *TaskResult) int {
		return a.ID - b.ID
	})
	for i, task := range result.CompletedTasks {
		str += strconv.Itoa(task.ID)
		str += ":"
		str += strconv.Itoa(task.Result)
		if i != len(result.CompletedTasks)-1 {
			str += ", "
		}
	}
	str += "]\n"

	str += "failed tasks: ["
	slices.SortFunc(result.FailedTasks, func(a *TaskResult, b *TaskResult) int {
		return a.ID - b.ID
	})
	for i, task := range result.FailedTasks {
		str += strconv.Itoa(task.ID)
		if i != len(result.FailedTasks)-1 {
			str += ", "
		}
	}
	str += "]"

	return str
}
