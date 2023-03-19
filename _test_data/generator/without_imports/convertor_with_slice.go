// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package without_imports is a generated datamapper package.
package without_imports

// ConvertFromSliceToToSlice convert []From to []To
func ConvertFromSliceToToSlice(fromSlice []From) []To {
	if fromSlice == nil {
		return nil
	}

	toSlice := make([]To, 0, len(fromSlice))
	for _, from := range fromSlice {
		toSlice = append(toSlice, ConvertFromToTo(from))
	}

	return toSlice
}