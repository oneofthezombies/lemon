package lemon

import (
	"testing"
)

func TestPlainTextBatchTaskResultPrinterExample1(t *testing.T) {
	expected := `completed tasks: [0:10, 1:20]
failed tasks: []`

	printer := &PlainTextBatchTaskResultPrinter{}
	output := printer.Print(&BatchTaskResult{
		CompletedTasks: []*TaskResult{
			{
				ID:     0,
				Result: 10,
			},
			{
				ID:     1,
				Result: 20,
			},
		},
	})
	if output != expected {
		t.Error("output mismatched")
		return
	}
}

func TestPlainTextBatchTaskResultPrinterExample2(t *testing.T) {
	expected := `completed tasks: [0:5]
failed tasks: [1, 2, 3]`

	printer := &PlainTextBatchTaskResultPrinter{}
	output := printer.Print(&BatchTaskResult{
		CompletedTasks: []*TaskResult{
			{
				ID:     0,
				Result: 5,
			},
		},
		FailedTasks: []*TaskResult{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
			{
				ID: 3,
			},
		},
	})
	if output != expected {
		t.Error("output mismatched")
		return
	}
}

func TestPlainTextBatchTaskResultPrinterExample3(t *testing.T) {
	expected := `cyclic dependencies detected: true`

	printer := &PlainTextBatchTaskResultPrinter{}
	output := printer.Print(&BatchTaskResult{
		CyclicDependenciesDetected: true,
	})
	if output != expected {
		t.Error("output mismatched")
		return
	}
}
