package lemon

type ConsoleTaskEventReporter struct{}

var _ TaskEventReporter = (*ConsoleTaskEventReporter)(nil)
