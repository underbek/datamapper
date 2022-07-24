package parser

import "github.com/google/uuid"

type String string

type StringAlias = string

type ModelRedefinition Model

type Array [16]uint

type Slice []uint

type Map map[int]string

type WithAlias struct {
	String            String
	StringAlias       StringAlias
	Array             Array
	Slice             Slice
	Map               Map
	RawArray          [16]uint
	RawSlice          []uint
	RawMap            map[int]string
	ModelRedefinition ModelRedefinition
	UUID              uuid.UUID
}
