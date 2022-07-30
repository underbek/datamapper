package parser

type ModelWithCollections struct {
	Array [12]uint64
	Slice []string
	Map   map[int]string

	PointerArray *[12]uint64
	PointerSlice *[]string
	PointerMap   *map[int]string

	ArrayPointers [12]*uint64
	SlicePointers []*string
	MapPointers   map[*int]*string

	ArrayModel [12]Model
	SliceModel []Model
	MapModel   map[Model]Model
}
