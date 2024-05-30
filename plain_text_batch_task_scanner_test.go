package lemon

import (
	"reflect"
	"testing"
)

func TestPlainTextBatchTaskScannerExample1(t *testing.T) {
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

func TestPlainTextBatchTaskScannerExample2(t *testing.T) {
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
		t.Error("scan failed")
		return
	}

	expected := []*Task{
		{
			ID:                  0,
			Result:              5,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{},
		},
		{
			ID:                  1,
			Result:              15,
			Throw:               true,
			PrerequisiteTaskIDs: []TaskID{},
		},
		{
			ID:                  2,
			Result:              25,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{0, 1},
		},
		{
			ID:                  3,
			Result:              35,
			Throw:               true,
			PrerequisiteTaskIDs: []TaskID{0},
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

func TestPlainTextBatchTaskScannerExample3(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 10     | false | [2]
1       | 20     | false | [0]
2       | 30     | false | [1]
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
			PrerequisiteTaskIDs: []TaskID{2},
		},
		{
			ID:                  1,
			Result:              20,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{0},
		},
		{
			ID:                  2,
			Result:              30,
			Throw:               false,
			PrerequisiteTaskIDs: []TaskID{1},
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
