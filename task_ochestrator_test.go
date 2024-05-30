package lemon

import (
	"context"
	"testing"
)

func TestTaskOchestratorExample1(t *testing.T) {
	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), Tasks{
		{
			ID:     0,
			Result: 10,
			Throw:  false,
		},
		{
			ID:     1,
			Result: 20,
			Throw:  false,
		},
	})
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	if len(res.CompletedTasks) != 2 {
		t.Error("unexpected completed tasks length")
		return
	}
}

func TestTaskOchestratorExample2(t *testing.T) {
	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), Tasks{
		{
			ID:     0,
			Result: 5,
			Throw:  false,
		},
		{
			ID:     1,
			Result: 15,
			Throw:  true,
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
	})
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	if len(res.CompletedTasks) != 1 {
		t.Error("unexpected completed tasks length")
		return
	}

	if len(res.FailedTasks) != 3 {
		t.Error("unexpected failed tasks length")
		return
	}
}

func TestTaskOchestratorExample3(t *testing.T) {
	ochestrator := NewTaskOchestrator()
	res, err := ochestrator.Run(context.Background(), Tasks{
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
	})
	if err != nil {
		t.Error("ochestrator running failed")
		return
	}

	if !res.CyclicDependenciesDetected {
		t.Error("it must be cyclic dependencies detected")
		return
	}
}
