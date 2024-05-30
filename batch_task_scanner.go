package lemon

type BatchTaskScanner interface {
	Scan(input string) (Tasks, error)
}
