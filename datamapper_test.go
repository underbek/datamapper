package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	modelsPath  = "test_data/datamapper/models.go"
	modelName   = "TestModel"
	modelTag    = "map"
	toModelName = "TestModelTo"
	toModelTag  = "map"
)

func Test_IncorrectOptions(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "Incorrect source",
			opts: Options{
				Source:   "incorrect_source",
				FromName: modelName,
				ToName:   toModelName,
				FromTag:  modelTag,
				ToTag:    toModelTag,
			},
		},
		{
			name: "Incorrect model name",
			opts: Options{
				Source:   modelsPath,
				FromName: "IncorrectName",
				ToName:   toModelName,
				FromTag:  modelTag,
				ToTag:    toModelTag,
			},
		},
		{
			name: "Incorrect to model name",
			opts: Options{
				Source:   modelsPath,
				FromName: modelName,
				ToName:   "IncorrectName",
				FromTag:  modelTag,
				ToTag:    toModelTag,
			},
		},
		{
			name: "Incorrect model tag",
			opts: Options{
				Source:   modelsPath,
				FromName: modelName,
				ToName:   toModelName,
				FromTag:  "incorrect tag",
				ToTag:    toModelTag,
			},
		},
		{
			name: "Incorrect to model tag",
			opts: Options{
				Source:   modelsPath,
				FromName: modelName,
				ToName:   toModelName,
				FromTag:  modelTag,
				ToTag:    "incorrect tag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapModels(tt.opts)
			require.Error(t, err)
		})
	}
}
