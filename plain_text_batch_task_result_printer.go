package lemon

type PlainTextBatchTaskResultPrinter struct{}

var _ BatchTaskResultPrinter = (*PlainTextBatchTaskResultPrinter)(nil)

func (p *PlainTextBatchTaskResultPrinter) Print() string {
	return ""
}
