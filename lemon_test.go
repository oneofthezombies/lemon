package lemon

import (
	"fmt"
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

	_ = tasks

	expected := `
completed tasks: [0:10, 1:20]
failed tasks: []
`
	fmt.Println(expected)

	// result := Print()
	// expected := 5

	// if result != expected {
	// 	t.Errorf("Print() = %d; want %d", result, expected)
	// }
}

func TestExample2(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 5      | false | []
1       | 15     | true  | []
2       | 25     | false | [0, 1]
3       | 35     | true  | [0]
`
	fmt.Println(input)

	expected := `
completed tasks: [0:5]
failed tasks: [1, 2, 3]
`
	fmt.Println(expected)
}

func TestExample3(t *testing.T) {
	input := `
task id | result | throw | prerequisite task ids
0       | 10     | false | [2]
1       | 20     | false | [0]
2       | 30     | false | [1]
`
	fmt.Println(input)

	expected := `
cyclic dependencies detected: true
`
	fmt.Println(expected)
}
