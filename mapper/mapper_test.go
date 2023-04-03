package mapper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/underbek/datamapper/_test_data"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/options"
)

const (
	modelsPath   = "test_data/datamapper/models.go"
	modelName    = "TestModel"
	modelTag     = "map"
	toModelName  = "TestModelTo"
	toModelTag   = "map"
	recursiveTag = "recursive"

	mapperDomainSource    = "../_test_data/mapper/domain"
	mapperTransportSource = "../_test_data/mapper/transport"
	customCFPath          = "../_test_data/mapper/convertors"
	otherCFPath           = "../_test_data/mapper/other_convertors"
	recursiveFrom         = "../_test_data/mapper/recursive/from"
	recursiveTo           = "../_test_data/mapper/recursive/to"

	destination     = "../_test_data/generated/mapper/user_convertor.go"
	destinationPath = "../_test_data/generated/mapper"
)

func readActual(t *testing.T) string {
	data, err := os.ReadFile(destination)
	require.NoError(t, err)

	return string(data)
}

func readFile(t *testing.T, fileName string) string {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", destinationPath, fileName))
	require.NoError(t, err)

	return string(data)
}

func clearDestination(t *testing.T, dest string) {
	assert.NoError(t, os.RemoveAll(dest))
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
		{
			name: "With slice",
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
						WithSlice: true,
					},
				},
			},
			expectedPath: "with_slice",
		},
		{
			name: "With invert and slice",
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
							Alias:  "api",
						},
						To: options.Model{
							Source: mapperDomainSource,
							Name:   "User",
							Tag:    toModelTag,
						},
						WithSlice: true,
						Inverse:   true,
					},
				},
			},
			expectedPath: "with_invert_and_slice",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clearDestination(t, destinationPath)

			err := MapModels(lg, tt.opts)
			require.NoError(t, err)

			actual := readActual(t)
			expected := _test_data.MapperExpected(t, tt.expectedPath)
			assert.Equal(t, expected, actual)
		})
	}
}

func Test_MapRecursiveModels(t *testing.T) {
	from := options.Model{
		Source: recursiveFrom,
		Name:   "Order",
		Tag:    recursiveTag,
		Alias:  "f",
	}

	to := options.Model{
		Source: recursiveTo,
		Name:   "Order",
		Tag:    recursiveTag,
		Alias:  "t",
	}

	destination := "../_test_data/generated/mapper/order.go"

	tests := []struct {
		name         string
		opts         options.Options
		isError      bool
		expectedPath string
	}{
		{
			name:    "without recursive",
			isError: true,
			opts: options.Options{
				Options: []options.Option{
					{
						Destination: destination,
						From:        from,
						To:          to,
					},
				},
			},
		},
		{
			name:         "recursive",
			expectedPath: "recursive",
			opts: options.Options{
				Options: []options.Option{
					{
						Destination: destination,
						Recursive:   true,
						From:        from,
						To:          to,
					},
				},
			},
		},
		{
			name:         "recursive with inverse",
			expectedPath: "recursive_with_inverse",
			opts: options.Options{
				Options: []options.Option{
					{
						Destination: destination,
						Recursive:   true,
						Inverse:     true,
						From:        from,
						To:          to,
					},
				},
			},
		},
		{
			name:         "recursive with pointers",
			expectedPath: "recursive_with_pointers",
			opts: options.Options{
				Options: []options.Option{
					{
						Destination:  destination,
						Recursive:    true,
						WithPointers: true,
						From:         from,
						To:           to,
					},
				},
			},
		},
		{
			name:         "recursive with inverse and pointers",
			expectedPath: "recursive_with_inverse_and_pointers",
			opts: options.Options{
				Options: []options.Option{
					{
						Destination:  destination,
						Recursive:    true,
						Inverse:      true,
						WithPointers: true,
						From:         from,
						To:           to,
					},
				},
			},
		},
	}

	lg := logger.New()

	converters := []string{
		"account_converter.go",
		"operation_converter.go",
		"order.go",
		"user_converter.go",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clearDestination(t, destinationPath)

			err := MapModels(lg, tt.opts)
			if tt.isError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			for _, converterName := range converters {
				actual := readFile(t, converterName)
				expected := _test_data.MapperExpectedFile(t, tt.expectedPath, converterName)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func Test_MapWithDash(t *testing.T) {
	tests := []struct {
		name         string
		opts         options.Options
		expectedPath string
	}{
		{
			name: "With dash",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
				},
				Options: []options.Option{
					{
						Destination: destination,
						Inverse:     true,
						From: options.Model{
							Source: "../_test_data/mapper/with_dash/domain",
							Name:   "Order",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: "../_test_data/mapper/with_dash/dao",
							Name:   "OrderData",
							Tag:    "db",
							Alias:  "db",
						},
					},
				},
			},
			expectedPath: "with_dash",
		},
		{
			name: "With dash and pointers",
			opts: options.Options{
				ConversionFunctions: []options.ConversionFunction{
					{Source: customCFPath},
					{Source: "../_test_data/mapper/with_dash_and_pointers/convertors/additional_convertor.go"},
				},
				Options: []options.Option{
					{
						Destination: destination,
						Inverse:     true,
						From: options.Model{
							Source: "../_test_data/mapper/with_dash_and_pointers/domain",
							Name:   "*Order",
							Tag:    modelTag,
						},
						To: options.Model{
							Source: "../_test_data/mapper/with_dash_and_pointers/dao",
							Name:   "OrderData",
							Tag:    "db",
							Alias:  "db",
						},
					},
				},
			},
			expectedPath: "with_dash_and_pointers",
		},
	}

	lg := logger.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clearDestination(t, destinationPath)

			err := MapModels(lg, tt.opts)
			require.NoError(t, err)

			actual := readActual(t)
			expected := _test_data.MapperExpected(t, tt.expectedPath)
			assert.Equal(t, expected, actual)
		})
	}
}
