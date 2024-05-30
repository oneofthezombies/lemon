package lemon

import "fmt"

type ConsoleEventReporter struct{}

var _ EventReporter = (*ConsoleEventReporter)(nil)

func (r *ConsoleEventReporter) Report(data map[string]any) {
	message, messageExist := data["message"]
	delete(data, "message")

	pairs := make([]string, 0)
	for k, v := range data {
		pairs = append(pairs, fmt.Sprintf("%v: %v", k, v))
	}

	str := ""
	if messageExist {
		str += fmt.Sprintf("message: %v", message)
		if len(pairs) > 0 {
			str += ", "
		}
	}

	for i, pair := range pairs {
		str += pair
		if i != len(pairs)-1 {
			str += ", "
		}
	}

	fmt.Println(str)
}
