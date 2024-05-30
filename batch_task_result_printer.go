package lemon

type BatchTaskResultPrinter interface {
	Print(result *BatchTaskResult) string
}
