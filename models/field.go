package models

type Field struct {
	Name          string
	Type          Type
	SkippedStruct *Struct
	CurrentStruct *Struct
	Head          *Field
	Tags          []Tag
}

type Fields struct {
	fields []Field
}

func NewFields(fields []Field) Fields {
	return Fields{fields: fields}
}

func (f Fields) Len() int {
	count := 0
	f.Range(func(field Field) {
		count++
	})

	return count
}

type EachFunc func(field *Field) error
type RangeFunc func(field Field)
type FilterFunc func(field *Field) bool

func (f *Fields) Range(fn RangeFunc) {
	for _, field := range f.fields {
		if field.SkippedStruct != nil {
			field.SkippedStruct.Fields.Range(fn)
			continue
		}
		fn(field)
	}
}

func (f *Fields) Each(fn EachFunc) error {
	var err error
	for i := range f.fields {
		if f.fields[i].SkippedStruct != nil {
			err = f.fields[i].SkippedStruct.Fields.Each(fn)
			if err != nil {
				return err
			}
			continue
		}
		err = fn(&f.fields[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Fields) SetFields(fields []Field) {
	f.fields = fields
}

func (f Fields) Filter(fn FilterFunc) Fields {
	var res []Field
	for _, field := range f.fields {
		if field.SkippedStruct != nil {
			field.SkippedStruct.Fields = field.SkippedStruct.Fields.Filter(fn)
			if field.SkippedStruct.Fields.Len() > 0 {
				res = append(res, field)
			}

			continue
		}

		if fn(&field) {
			res = append(res, field)
		}
	}

	return NewFields(res)
}

func (f *Fields) Add(field ...Field) {
	f.fields = append(f.fields, field...)
}
