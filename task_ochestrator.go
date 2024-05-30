package lemon

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type CyclicDependenciesError struct {
	Cycles []TaskID
}

func (e *CyclicDependenciesError) Error() string {
	return fmt.Sprintf("cyclic dependencies: %v", e.Cycles)
}

type TaskOchestrator struct {
	EventReporter      EventReporter
	StatusReportPeriod time.Duration
}

func NewTaskOchestrator() *TaskOchestrator {
	return &TaskOchestrator{
		EventReporter:      &ConsoleEventReporter{},
		StatusReportPeriod: 500 * time.Millisecond,
	}
}

func (o *TaskOchestrator) Run(
	ctx context.Context,
	tasks Tasks) (*BatchTaskResult, error) {
	taskMap := make(map[TaskID]*Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}

	independentTaskIDs := make([]TaskID, 0)
	err := o.parseCyclicDependencies(taskMap, &independentTaskIDs)
	if err != nil {
		if _, ok := err.(*CyclicDependenciesError); ok {
			return &BatchTaskResult{
				CyclicDependenciesDetected: true,
			}, nil
		} else {
			return nil, err
		}
	}

	readySet := make(map[TaskID]any)
	for taskID := range taskMap {
		readySet[taskID] = struct{}{}

	}

	runningSet := make(map[TaskID]any)
	taskResultChan := make(chan *TaskResult, len(tasks))
	for independentTaskID := range independentTaskIDs {
		task, exist := taskMap[independentTaskID]
		if !exist {
			return nil, fmt.Errorf("task %d must exist", independentTaskID)
		}

		delete(readySet, task.ID)
		runningSet[task.ID] = struct{}{}
		go o.runTask(ctx, task, taskResultChan)
	}

	statusReportTicker := time.NewTicker(o.StatusReportPeriod)
	defer statusReportTicker.Stop()

	taskResultMap := make(map[TaskID]*TaskResult)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-statusReportTicker.C:
			// 설정값 (기본 500 밀리초)에 따라 주기적으로 현재 진행 상태 이벤트를 리포팅합니다.
			data := map[string]any{
				"message": "status",
			}

			runningTasks := make([]TaskID, 0)
			for runningTaskID := range runningSet {
				runningTasks = append(runningTasks, runningTaskID)
			}
			data["runningTasks"] = runningTasks

			completedTasks := make([]TaskID, 0)
			failedTasks := make([]TaskID, 0)
			for _, taskResult := range taskResultMap {
				if taskResult.IsSuccess {
					completedTasks = append(completedTasks, taskResult.ID)
				} else {
					failedTasks = append(failedTasks, taskResult.ID)
				}
			}
			data["completedTasks"] = completedTasks
			data["failedTasks"] = failedTasks

			readyTasks := make([]TaskID, 0)
			for readyTaskID := range readySet {
				readyTasks = append(readyTasks, readyTaskID)
			}
			data["readyTasks"] = readyTasks

			o.EventReporter.Report(data)

		case taskResult := <-taskResultChan:
			taskResultMap[taskResult.ID] = taskResult
			delete(runningSet, taskResult.ID)

			// 작업 완료 결과를 받을 때마다 이벤트를 리포팅합니다.
			data := map[string]any{
				"message":   "task done",
				"taskId":    taskResult.ID,
				"isSuccess": taskResult.IsSuccess,
			}
			if taskResult.IsSuccess {
				data["result"] = taskResult.Result
			} else {
				data["error"] = taskResult.Error
			}

			o.EventReporter.Report(data)

			// 레디 셋을 순회해, 의존 작업이 실패했다면, 즉시 실패 결과에 추가하거나, 실행 가능하다면 실행합니다.
			for readyTaskID := range readySet {
				task, exist := taskMap[readyTaskID]
				if !exist {
					return nil, fmt.Errorf("task %d must exist", readyTaskID)
				}

				runnable := true
				failedPrerequisiteTaskIDs := make([]TaskID, 0)
				for prerequisiteTaskID := range task.PrerequisiteTaskIDs {
					taskResult, exist := taskResultMap[prerequisiteTaskID]

					// 의존 작업이 완료되지 않았습니다.
					if !exist {
						runnable = false
						continue
					}

					// 실패한 의존 작업 목록에 추가합니다.
					if !taskResult.IsSuccess {
						failedPrerequisiteTaskIDs = append(failedPrerequisiteTaskIDs, prerequisiteTaskID)
						runnable = false
					}
				}

				// 의존 작업이 하나라도 실패했다면, 실패 결과를 레디 셋에서 즉시 작업 결과 맵으로 이동시킵니다.
				if len(failedPrerequisiteTaskIDs) > 0 {
					delete(readySet, readyTaskID)
					taskResult := &TaskResult{
						ID:        readyTaskID,
						IsSuccess: false,
						Error:     fmt.Errorf("prerequisite tasks %v failed", failedPrerequisiteTaskIDs),
					}
					taskResultMap[readyTaskID] = taskResult

					o.EventReporter.Report(map[string]any{
						"message":   "task skipped",
						"taskId":    taskResult.ID,
						"isSuccess": taskResult.IsSuccess,
						"error":     taskResult.Error,
					})
				}

				if runnable {
					delete(readySet, task.ID)
					runningSet[task.ID] = struct{}{}
					go o.runTask(ctx, task, taskResultChan)
				}
			}

			// 작업이 모두 완료 됐습니다.
			if len(taskResultMap) == len(taskMap) {
				completedTasks := make([]*TaskResult, 0)
				failedTasks := make([]*TaskResult, 0)
				for _, taskResult := range taskResultMap {
					if taskResult.IsSuccess {
						completedTasks = append(completedTasks, taskResult)
					} else {
						failedTasks = append(failedTasks, taskResult)
					}
				}

				return &BatchTaskResult{
					CompletedTasks: completedTasks,
					FailedTasks:    failedTasks,
				}, nil
			}
		}
	}
}

// 순환 참조를 검사합니다. 그리고 의존 작업이 없는 "독립" 작업 목록을 채웁니다.
func (o *TaskOchestrator) parseCyclicDependencies(
	taskMap map[TaskID]*Task,
	independentTaskIDs *[]TaskID) error {
	if independentTaskIDs == nil {
		return fmt.Errorf("independentTaskIDs must not nil")
	}

	visitSet := make(map[TaskID]any)
	traverseStack := make([]TaskID, 0)

	for taskID := range taskMap {
		if err := o.traverseTask(
			taskMap,
			independentTaskIDs,
			visitSet,
			&traverseStack,
			taskID); err != nil {
			return err
		}
	}

	return nil
}

// 순환 참조 검사와 독립 작업 목록을 채우는 재귀 함수입니다.
func (o *TaskOchestrator) traverseTask(
	taskMap map[TaskID]*Task,
	independentTaskIDs *[]TaskID,
	visitSet map[TaskID]any,
	traverseStack *[]TaskID,
	taskID TaskID) error {
	if independentTaskIDs == nil {
		return fmt.Errorf("independentTaskIDs must not nil")
	}

	if traverseStack == nil {
		return fmt.Errorf("traverseStack must not nil")
	}

	if _, exist := visitSet[taskID]; exist {
		return nil
	}

	// 방문 플래그를 업데이트합니다.
	visitSet[taskID] = struct{}{}

	// 방문 스택에 추가합니다.
	*traverseStack = append(*traverseStack, taskID)
	task, exist := taskMap[taskID]
	if !exist {
		return fmt.Errorf("task %d not found", taskID)
	}

	// 의존 작업이 없다면, 최초로 실행될 독립 작업 목록에 추가합니다.
	if len(task.PrerequisiteTaskIDs) == 0 {
		*independentTaskIDs = append(*independentTaskIDs, taskID)
	}

	for _, prerequisiteTaskID := range task.PrerequisiteTaskIDs {
		if _, visited := visitSet[prerequisiteTaskID]; visited {
			// 방문 스택에 이미 존재하는 경우, 순환 참조를 의미합니다.
			foundIdx := slices.Index(*traverseStack, prerequisiteTaskID)
			if foundIdx != -1 {
				cycles := make([]TaskID, 0)
				for i := foundIdx; i < len(*traverseStack); i++ {
					cycles = append(cycles, (*traverseStack)[i])
				}
				cycles = append(cycles, prerequisiteTaskID)
				o.EventReporter.Report(map[string]any{
					"message": "task cycle detected",
					"cycles":  cycles,
				})
				return &CyclicDependenciesError{
					Cycles: cycles,
				}
			}
		} else {
			if err := o.traverseTask(
				taskMap,
				independentTaskIDs,
				visitSet,
				traverseStack,
				prerequisiteTaskID); err != nil {
				return err
			}
		}
	}

	// 방문 스택에서 pop합니다.
	*traverseStack = (*traverseStack)[:len(*traverseStack)-1]
	return nil
}

// 작업을 병렬로 실행시 호출되는 함수입니다.
// 반드시 하나의 결과만 전달해야 합니다.
func (o *TaskOchestrator) runTask(ctx context.Context, task *Task, taskResultChan chan<- *TaskResult) {
	// 전달받은 컨텍스트가 종료됐기 때문에, 컨텍스트 에러 결과를 채널에 전달하고 함수를 반환합니다.
	select {
	case <-ctx.Done():
		taskResultChan <- &TaskResult{
			ID:        task.ID,
			IsSuccess: false,
			Error:     ctx.Err(),
		}
		return
	default:
		break
	}

	// 1~3 초 랜덤하게 대기합니다.
	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)

	// Throw일 경우, 예외 발생 에러 결과를 채널에 전달하고 함수를 반환합니다.
	if task.Throw {
		taskResultChan <- &TaskResult{
			ID:        task.ID,
			IsSuccess: false,
			Error:     fmt.Errorf("exception has been thrown"),
		}
		return
	}

	// 성공 결과를 채널에 전달하고 함수를 반환합니다.
	taskResultChan <- &TaskResult{
		ID:        task.ID,
		IsSuccess: true,
		Result:    task.Result,
	}
}
