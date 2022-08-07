package mapper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/_test_data"
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
	otherCFPath           = "../_test_data/mapper/other_convertors"

	destination = "../_test_data/generated/mapper/user_convertor.go"
)

func readActual(t *testing.T) string {
	data, err := os.ReadFile(destination)
	require.NoError(t, err)

	return string(data)
}

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
	tests := []struct {
		name         string
		opts         options.Options
		expectedPath string
	}{
		{
			name: "Map models",
			opts: options.Options{
				Destination:   destination,
				UserCFSources: []string{customCFPath},
				FromSource:    mapperTransportSource,
				ToSource:      mapperDomainSource,
				FromName:      "User",
				ToName:        "User",
				FromTag:       modelTag,
				ToTag:         toModelTag,
			},
			expectedPath: "map_models",
		},
		{
			name: "With some cf sources",
			opts: options.Options{
				Destination:   destination,
				UserCFSources: []string{customCFPath, otherCFPath},
				FromSource:    mapperTransportSource,
				ToSource:      mapperDomainSource,
				FromName:      "User",
				ToName:        "User",
				FromTag:       modelTag,
				ToTag:         toModelTag,
			},
			expectedPath: "with_some_cf_sources",
		},
		{
			name: "With aliases",
			opts: options.Options{
				Destination: destination,
				UserCFSources: []string{
					fmt.Sprintf("%s:%s", customCFPath, "customCf"),
					fmt.Sprintf("%s:%s", otherCFPath, "otherCf"),
				},
				FromSource: fmt.Sprintf("%s:%s", mapperTransportSource, "from"),
				ToSource:   fmt.Sprintf("%s:%s", mapperDomainSource, "to"),
				FromName:   "User",
				ToName:     "User",
				FromTag:    modelTag,
				ToTag:      toModelTag,
			},
			expectedPath: "with_aliases",
		},
		{
			name: "With invert",
			opts: options.Options{
				Destination:   destination,
				UserCFSources: []string{customCFPath},
				FromSource:    mapperTransportSource,
				ToSource:      mapperDomainSource,
				FromName:      "User",
				ToName:        "User",
				FromTag:       modelTag,
				ToTag:         toModelTag,
				Invert:        true,
			},
			expectedPath: "with_invert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapModels(tt.opts)
			require.NoError(t, err)

			actual := readActual(t)
			expected := _test_data.MapperExpected(t, tt.expectedPath)
			assert.Equal(t, expected, actual)
		})
	}
}
