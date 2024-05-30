package lemon

type ConsoleBatchTaskStatusReporter struct{}

var _ BatchTaskStatusReporter = (*ConsoleBatchTaskStatusReporter)(nil)
