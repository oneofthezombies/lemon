# lemon

## Prerequisites

- Go 1.22 or higher

## Highlights

- `mutex` 사용 없이, 고루틴과 채널을 사용해 구현
  - 경량 스레드간 큐를 사용해 데이터를 주고 받음
- 순환 참조 감지
- 결과 이벤트 리포팅 (`EventReporter` 인터페이스를 구현한 `ConsoleEventReporter`가 stdout에 작성)
- 주기적 상태 리포팅 (기본 0.5 밀리초)
- 하위 고루틴들의 라이프타임 관리를 위한 `context` 인자
- `BatchTaskScanner` - `Task` 타입 -> `TaskOchestrator` - `TaskResult` 타입 -> `BatchTaskResultPrinter` 역할 및 코드 인터페이스 분리
  - `JsonBatchTaskScanner`, `JsonBatchTaskResultPrinter`과 같은 추후 확장 고려한 설계
- 예제용 `BatchTaskScanner` 인터페이스를 구현한 `PlainTextBatchTaskScanner` 구현
- 예제용 `BatchTaskResultPrinter` 인터페이스를 구현한 `PlainTextBatchTaskResultPrinter` 구현
- 유닛 테스트 및 통합 테스트 작성
- 실제 시나리오 모사를 위해 각 병렬 작업 실행시 랜덤하게 1-3초 대기 후 결과 전송

## How to Test

```sh
go test
```

`lemon_test.go`  
예제 1, 2, 3의 입력 문자열을 파싱, 작업 실행, 출력 문자열 구성 테스트
- TestExample1
- TestExample2
- TestExample3

`task_ochestrator_test.go`  
예제 1, 2, 3의 `Task` 구조체를 사용해 작업 실행, `TaskResult` 구조체 반환 테스트
- TestTaskOchestratorExample1
- TestTaskOchestratorExample2
- TestTaskOchestratorExample3


`plain_text_batch_task_scanner_test.go`  
예제 1, 2, 3의 입력 문자열을 파싱 후 구조체 비교 테스트
- TestPlainTextBatchTaskScannerExample1
- TestPlainTextBatchTaskScannerExample2
- TestPlainTextBatchTaskScannerExample3

`plain_text_batch_task_result_printer_test.go`  
예제 1, 2, 3의 구조체로부터 출력 문자열 구성 비교 테스트
- TestPlainTextBatchTaskResultPrinterExample1
- TestPlainTextBatchTaskResultPrinterExample2
- TestPlainTextBatchTaskResultPrinterExample3
