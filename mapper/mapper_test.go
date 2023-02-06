package mapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/_test_data"
	"github.com/underbek/datamapper/logger"
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
				Options: []options.Option{
					{
						From: options.Model{
							Source: "incorrect_source",
							Name:   modelName,
							Tag:    modelTag,
						},
						To: options.Model{
							Source: modelsPath,
							Name:   toModelName,
							Tag:    toModelTag,
						},
					},
				},
			},
		},
		{
			name: "Incorrect source",
			opts: options.Options{
				Options: []options.Option{
					{
						From: options.Model{
							Source: modelsPath,
							Name:   modelName,
							Tag:    modelTag,
						},
						To: options.Model{
							Source: "incorrect_source",
							Name:   toModelName,
							Tag:    toModelTag,
						},
					},
				},
			},
		},
		{
			name: "Incorrect model name",
			opts: options.Options{
				Options: []options.Option{
					{
						From: options.Model{
							Source: modelsPath,
							Name:   "IncorrectName",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: modelsPath,
							Name:   toModelName,
							Tag:    toModelTag,
						},
					},
				},
			},
		},
		{
			name: "Incorrect to model name",
			opts: options.Options{
				Options: []options.Option{
					{
						From: options.Model{
							Source: modelsPath,
							Name:   modelName,
							Tag:    modelTag,
						},
						To: options.Model{
							Source: modelsPath,
							Name:   "IncorrectName",
							Tag:    toModelTag,
						},
					},
				},
			},
		},
		{
			name: "Incorrect model tag",
			opts: options.Options{
				Options: []options.Option{
					{
						From: options.Model{
							Source: modelsPath,
							Name:   modelName,
							Tag:    "incorrect tag",
						},
						To: options.Model{
							Source: modelsPath,
							Name:   toModelName,
							Tag:    toModelTag,
						},
					},
				},
			},
		},
		{
			name: "Incorrect to model tag",
			opts: options.Options{
				Options: []options.Option{
					{
						From: options.Model{
							Source: modelsPath,
							Name:   modelName,
							Tag:    modelTag,
						},
						To: options.Model{
							Source: modelsPath,
							Name:   toModelName,
							Tag:    "incorrect tag",
						},
					},
				},
			},
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapModels(lg, tt.opts)
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
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "User",
							Tag:    toModelTag,
						},
					},
				},
			},
			expectedPath: "map_models",
		},
		{
			name: "With some cf sources",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
					{Source: otherCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "User",
							Tag:    toModelTag,
						},
					},
				},
			},
			expectedPath: "with_some_cf_sources",
		},
		{
			name: "With aliases",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath, Alias: "customCf"},
					{Source: otherCFPath, Alias: "otherCf"},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Alias:  "from",
							Name:   "User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Alias:  "to",
							Name:   "User",
							Tag:    toModelTag,
						},
					},
				},
			},
			expectedPath: "with_aliases",
		},
		{
			name: "With invert",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "User",
							Tag:    toModelTag,
						},
						Inverse: true,
					},
				},
			},
			expectedPath: "with_invert",
		},
		{
			name: "With pointers",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "*User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "*User",
							Tag:    toModelTag,
						},
					},
				},
			},
			expectedPath: "with_pointers",
		},
		{
			name: "With generated functions",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "*User",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "*User",
							Tag:    toModelTag,
						},
					},
					{
						Destination: destination,
						From: options.Model{
							Source: mapperTransportSource,
							Name:   "Order",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "Order",
							Tag:    toModelTag,
						},
					},
				},
			},
			expectedPath: "with_generated",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapModels(lg, tt.opts)
			require.NoError(t, err)

			actual := readActual(t)
			expected := _test_data.MapperExpected(t, tt.expectedPath)
			assert.Equal(t, expected, actual)
		})
	}
}
