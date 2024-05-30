package lemon

type EventReporter interface {
	Report(data map[string]any)
}
