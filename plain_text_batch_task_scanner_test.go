package lemon

import (
	"reflect"
	"testing"
)

func TestPlainTextBatchTaskScanner(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 10     | false | []
1       | 20     | false | []
`

	scanner := &PlainTextBatchTaskScanner{}
	tasks, err := scanner.Scan(input)
	if err != nil {
		t.Error("scan failed")
		return
	}

	expected := []*Task{
		{
			ID:                  0,
			Result:              10,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{},
		},
		{
			ID:                  1,
			Result:              20,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{},
		},
	}

	if len(tasks) != len(expected) {
		t.Error("tasks length mismatched")
		return
	}

	for i := range tasks {
		if !reflect.DeepEqual(tasks[i], expected[i]) {
			t.Errorf("%d task mismatched", i)
			return
		}
	}
}
