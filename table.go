package qgb

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/GoWebProd/gip/types/iface"
	"github.com/pkg/errors"
)

type field struct {
	name         string
	offset       uintptr
	fType        unsafe.Pointer
	isPrimaryKey bool
}

type table struct {
	name string

	fields    []*field
	fieldsMap map[string]*field

	createdAt  *field
	updatedAt  *field
	primaryKey *field
}

func buildTable[T any](name string) (table, error) {
	var (
		table table
		t     T
	)

	table.name = name
	table.fieldsMap = make(map[string]*field)

	rType := reflect.TypeOf(t)

	for i := range rType.NumField() {
		f := rType.Field(i)
		if !f.IsExported() {
			continue
		}

		tag := f.Tag.Get("db")
		if tag == "" {
			continue
		}

		options := strings.Split(tag, ",")
		name := options[0]
		options = options[1:]

		if name == "-" {
			continue
		}

		t, _ := iface.Unpack(reflect.New(f.Type).Interface())
		field := &field{name: name, offset: f.Offset, fType: t, isPrimaryKey: hasPrimaryKeyOption(options)}

		if field.isPrimaryKey {
			table.primaryKey = field
		}

		switch name {
		case "created_at":
			table.createdAt = field
		case "updated_at":
			table.updatedAt = field
		default:
			table.fields = append(table.fields, field)
			table.fieldsMap[name] = field
		}
	}

	if table.primaryKey == nil {
		return table, errors.New("no has primary key")
	}

	return table, nil
}

func hasPrimaryKeyOption(options []string) bool {
	for _, o := range options {
		if o == "primaryKey" {
			return true
		}
	}

	return false
}

func transformArgs(table *table, args []placeholderValue) ([]placeholderValue, error) {
	for idx := range args {
		ph, ok := args[idx].value.(placeholder)
		if !ok || ph.name == "" {
			continue
		}

		fieldName := args[idx].field
		if ok {
			fieldName = ph.name
		}

		f, ok := table.fieldsMap[fieldName]
		if !ok {
			return nil, errors.Errorf("field %s not found", fieldName)
		}

		args[idx].value = f
	}

	return args, nil
}
