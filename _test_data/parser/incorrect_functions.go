package parser

func Void() {}

func WithoutResult(from string) {}

func WithoutIn() string {
	return "test"
}

func ManyArguments(fromA, fromB int) string {
	return "test"
}

func ManyArguments2(fromA float32, fromB int) string {
	return "test"
}

func ManyResults(from int) (string, string) {
	return "test", "test"
}

func ManyResults2(from int) (a, b string) {
	return "test", "test"
}
