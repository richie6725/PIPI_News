package gormDao

import "fmt"

type FieldName string

func (f FieldName) AddField() string {
	return fmt.Sprintf("%s = ?", f)
}

func (f FieldName) ContainAllField() string {
	return fmt.Sprintf("%s @> ?", f)
}

func (f FieldName) ContainInField() string {
	return fmt.Sprintf("%s <@ ?", f)
}

func (f FieldName) OverlapsField() string {
	return fmt.Sprintf("%s && ?", f)
}

func (f FieldName) MoreThanField() string {
	return fmt.Sprintf("%s > ?", f)
}

func (f FieldName) MoreThanEqualField() string {
	return fmt.Sprintf("%s >= ?", f)
}

func (f FieldName) LessThanField() string {
	return fmt.Sprintf("%s < ?", f)
}

func (f FieldName) LessThanEqualField() string {
	return fmt.Sprintf("%s <= ?", f)
}

func (f FieldName) String() string { return string(f) }
