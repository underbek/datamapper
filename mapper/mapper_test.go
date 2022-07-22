package mapper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/options"
)

const (
	modelsPath  = "test_data/datamapper/models.go"
	modelName   = "TestModel"
	modelTag    = "map"
	toModelName = "TestModelTo"
	toModelTag  = "map"

	mapperDomainSource    = "../_test_data/mapper/domain"
	mapperTransportSource = "../_test_data/mapper/transport"
	customCFPath          = "../_test_data/mapper/convertors"

	destination = "../_test_data/generated/mapper/user_convertor.go"
)

func Test_IncorrectOptions(t *testing.T) {
	tests := []struct {
		name string
		opts options.Options
	}{
		{
			name: "Incorrect source",
			opts: options.Options{
				FromSource: "incorrect_source",
				ToSource:   modelsPath,
				FromName:   modelName,
				ToName:     toModelName,
				FromTag:    modelTag,
				ToTag:      toModelTag,
			},
		},
		{
			name: "Incorrect source",
			opts: options.Options{
				FromSource: modelsPath,
				ToSource:   "incorrect_source",
				FromName:   modelName,
				ToName:     toModelName,
				FromTag:    modelTag,
				ToTag:      toModelTag,
			},
		},
		{
			name: "Incorrect model name",
			opts: options.Options{
				FromSource: modelsPath,
				ToSource:   modelsPath,
				FromName:   "IncorrectName",
				ToName:     toModelName,
				FromTag:    modelTag,
				ToTag:      toModelTag,
			},
		},
		{
			name: "Incorrect to model name",
			opts: options.Options{
				FromSource: modelsPath,
				ToSource:   modelsPath,
				FromName:   modelName,
				ToName:     "IncorrectName",
				FromTag:    modelTag,
				ToTag:      toModelTag,
			},
		},
		{
			name: "Incorrect model tag",
			opts: options.Options{
				FromSource: modelsPath,
				ToSource:   modelsPath,
				FromName:   modelName,
				ToName:     toModelName,
				FromTag:    "incorrect tag",
				ToTag:      toModelTag,
			},
		},
		{
			name: "Incorrect to model tag",
			opts: options.Options{
				FromSource: modelsPath,
				ToSource:   modelsPath,
				FromName:   modelName,
				ToName:     toModelName,
				FromTag:    modelTag,
				ToTag:      "incorrect tag",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapModels(tt.opts)
			require.Error(t, err)
		})
	}
}

func Test_MapModels(t *testing.T) {
	opts := options.Options{
		Destination:  destination,
		UserCFSource: customCFPath,
		FromSource:   mapperTransportSource,
		ToSource:     mapperDomainSource,
		FromName:     "User",
		ToName:       "User",
		FromTag:      modelTag,
		ToTag:        toModelTag,
	}
	err := MapModels(opts)
	require.NoError(t, err)
}
