package lemon

import (
	"context"
	"testing"
)

func TestExample1(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 10     | false | []
1       | 20     | false | []
`
	scanner := &PlainTextBatchTaskScanner{}
	tasks, err := scanner.Scan(input)
	if err != nil {
		t.Error("input parsing failed")
	}

	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), tasks)
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	printer := PlainTextBatchTaskResultPrinter{}
	output := printer.Print(res)
	expected := `completed tasks: [0:10, 1:20]
failed tasks: []`
	if output != expected {
		t.Error("unexpected output")
		return
	}
}

func TestExample2(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 5      | false | []
1       | 15     | true  | []
2       | 25     | false | [0, 1]
3       | 35     | true  | [0]
`
	scanner := &PlainTextBatchTaskScanner{}
	tasks, err := scanner.Scan(input)
	if err != nil {
		t.Error("input parsing failed")
	}

	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), tasks)
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	printer := PlainTextBatchTaskResultPrinter{}
	output := printer.Print(res)
	expected := `completed tasks: [0:5]
failed tasks: [1, 2, 3]`
	if output != expected {
		t.Error("unexpected output")
		return
	}
}

func TestExample3(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 10     | false | [2]
1       | 20     | false | [0]
2       | 30     | false | [1]
`
	scanner := &PlainTextBatchTaskScanner{}
	tasks, err := scanner.Scan(input)
	if err != nil {
		t.Error("input parsing failed")
	}

	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), tasks)
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	printer := PlainTextBatchTaskResultPrinter{}
	output := printer.Print(res)
	expected := `cyclic dependencies detected: true`
	if output != expected {
		t.Error("unexpected output")
		return
	}
}
