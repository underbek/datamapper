package _test_data

import (
	"embed"
	"io"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed generator
var generatorExpected embed.FS

func Generator(t *testing.T, fileName string) string {
	return readFile(t, generatorExpected, "generator/"+fileName)
}

func readFile(t *testing.T, fs fs.FS, path string) string {
	f, err := fs.Open(path)
	require.NoError(t, err)

	defer func() { require.NoError(t, f.Close()) }()

	body, err := io.ReadAll(f)
	require.NoError(t, err)

	return string(body)
}
