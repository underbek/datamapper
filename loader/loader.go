package loader

import (
	"embed"
	"os"

	"github.com/underbek/datamapper/models"
	"gopkg.in/yaml.v3"
)

const (
	embedFileName  = "data/converts.yaml"
	fileNameByRoot = "loader/data/converts.yaml"
)

//go:embed data
var data embed.FS

func Save(funcs models.Functions) error {
	funcSlice := values(funcs)
	data, err := yaml.Marshal(funcSlice)
	if err != nil {
		return err
	}

	return os.WriteFile(fileNameByRoot, data, 0644)
}

func Read() (models.Functions, error) {
	body, err := data.ReadFile(embedFileName)
	if err != nil {
		return nil, err
	}

	var cf []models.ConversionFunction
	err = yaml.Unmarshal(body, &cf)
	if err != nil {
		return nil, err
	}

	return makeMap(cf), nil
}

func values(m models.Functions) []models.ConversionFunction {
	r := make([]models.ConversionFunction, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func makeMap(cf []models.ConversionFunction) models.Functions {
	r := make(models.Functions, len(cf))
	for _, v := range cf {
		r[models.ConversionFunctionKey{
			FromType: v.FromType,
			ToType:   v.ToType,
		}] = v
	}
	return r
}
