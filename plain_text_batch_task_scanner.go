package lemon

import (
	"fmt"
	"strconv"
	"strings"
)

type PlainTextBatchTaskScanner struct{}

var _ BatchTaskScanner = (*PlainTextBatchTaskScanner)(nil)

func (s *PlainTextBatchTaskScanner) Scan(input string) (Tasks, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("input table header line does not exist")
	}

	if strings.TrimSpace(lines[0]) != "task id | result | throw | prerequisite task ids" {
		return nil, fmt.Errorf("input table header line mismatched")
	}

	taskLines := lines[1:]
	tasks := make(Tasks, len(taskLines))
	for rowIdx, line := range taskLines {
		rows := strings.Split(strings.TrimSpace(line), "|")
		if len(rows) != 4 {
			return nil, fmt.Errorf("the number of columns in row %d is not 4", rowIdx)
		}

		taskId, err := strconv.Atoi(strings.TrimSpace(rows[0]))
		if err != nil {
			return nil, fmt.Errorf("task id parsing failed in row %d", rowIdx)
		}

		result, err := strconv.Atoi(strings.TrimSpace(rows[1]))
		if err != nil {
			return nil, fmt.Errorf("result parsing failed in row %d", rowIdx)
		}

		throw, err := strconv.ParseBool(strings.TrimSpace(rows[2]))
		if err != nil {
			return nil, fmt.Errorf("throw parsing failed in row %d", rowIdx)
		}

		prerequisitesFormat := strings.TrimSpace(rows[3])
		if prerequisitesFormat[0] != '[' || prerequisitesFormat[len(prerequisitesFormat)-1] != ']' {
			return nil, fmt.Errorf("invalid prerequisites format in row %d", rowIdx)
		}

		unwrappedPrerequisitesFormat := prerequisitesFormat[1 : len(prerequisitesFormat)-1]
		var prerequisites []string
		if len(unwrappedPrerequisitesFormat) > 0 {
			prerequisites = strings.Split(unwrappedPrerequisitesFormat, ",")
		} else {
			prerequisites = make([]string, 0)
		}

		prerequisiteTaskIDs := make([]TaskID, len(prerequisites))
		for prerequisiteIdx, preprerequisite := range prerequisites {
			prerequisiteTaskID, err := strconv.Atoi(strings.TrimSpace(preprerequisite))
			if err != nil {
				return nil, fmt.Errorf("prerequisite task id parsing failed in row %d", rowIdx)
			}

			prerequisiteTaskIDs[prerequisiteIdx] = prerequisiteTaskID
		}

		tasks[rowIdx] = &Task{
			ID:                  taskId,
			Result:              result,
			Throw:               throw,
			PrerequisiteTaskIDs: prerequisiteTaskIDs,
		}
	}

	return tasks, nil
}
