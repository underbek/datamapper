// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package with_from_pointer is a generated datamapper package.
package with_from_pointer

// ConvertFromSliceToToSlice convert []*From to []To
func ConvertFromSliceToToSlice(fromSlice []*From) ([]To, error) {
	if fromSlice == nil {
		return nil, nil
	}

	toSlice := make([]To, 0, len(fromSlice))
	for _, from := range fromSlice {
		to, err := ConvertFromToTo(from)
		if err != nil {
			return nil, err
		}
		toSlice = append(toSlice, to)
	}

	return toSlice, nil
}
