package parser

type CFTestType = string

type CFEmptyInterface interface{}

type CFTestModel struct {
	Name string
}

type CFTestFuncType = func(from int) string

var CFTestFuncVar = func(from int) string {
	return "test"
}
