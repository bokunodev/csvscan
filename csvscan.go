package csvscan

import (
	"encoding"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type Scanner struct {
	reader  *csv.Reader
	indexes []int
}

func New() *Scanner {
	return &Scanner{}
}

func (scn *Scanner) Init(file io.Reader, ptr any) error {
	scn.reader = nil
	scn.indexes = scn.indexes[:0]

	scn.reader = csv.NewReader(file)
	dst_typ := reflect.TypeOf(ptr)
	if dst_typ.Kind() != reflect.Pointer || dst_typ.Elem().Kind() != reflect.Struct {
		return errors.New("`ptr` must be a pointer to struct")
	}
	dst_typ = dst_typ.Elem()

	headers, err := scn.reader.Read()
	if err != nil {
		return err
	}

	// read csv header
	for _, header := range headers {
		idx, ok := find_field_index(dst_typ, strings.TrimSpace(header))
		if ok {
			scn.indexes = append(scn.indexes, idx)
			continue
		}

		scn.indexes = append(scn.indexes, -1)
	}

	return nil
}

func (scn *Scanner) Scan(dst any) error {
	dst_val := reflect.ValueOf(dst).Elem()

	cols, err := scn.reader.Read()
	if err != nil {
		return err
	}

	for i, col := range cols {
		idx := scn.indexes[i]
		if idx == -1 {
			continue
		}

		field := dst_val.Field(idx)

		if field.Kind() == reflect.Pointer {
			value := reflect.New(field.Type().Elem())
			if err := set_value(value.Elem(), col); err != nil {
				return err
			}

			field.Set(value)
			continue
		}

		if err := set_value(field, col); err != nil {
			return err
		}
	}

	return nil
}

func set_value(dst reflect.Value, str string) error {
	if dst.CanAddr() {
		if tu, ok := dst.Addr().Interface().(encoding.TextUnmarshaler); ok {
			if err := tu.UnmarshalText([]byte(str)); err != nil {
				return err
			}
			return nil
		}
	}

	switch dst.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		dst.SetBool(val)
		return nil
	case reflect.Int8:
		val, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return err
		}
		dst.SetInt(val)
		return nil
	case reflect.Int16:
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return err
		}
		dst.SetInt(val)
		return nil
	case reflect.Int32:
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err
		}
		dst.SetInt(val)
		return nil
	case reflect.Int, reflect.Int64:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		dst.SetInt(val)
		return nil
	case reflect.Uint8:
		val, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return err
		}
		dst.SetUint(val)
		return nil
	case reflect.Uint16:
		val, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return err
		}
		dst.SetUint(val)
		return nil
	case reflect.Uint32:
		val, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return err
		}
		dst.SetUint(val)
		return nil
	case reflect.Uint, reflect.Uint64:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		dst.SetUint(val)
		return nil
	case reflect.Float32:
		val, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}
		dst.SetFloat(val)
		return nil
	case reflect.Float64:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		dst.SetFloat(val)
		return nil
	case reflect.String:
		dst.SetString(str)
		return nil
	}

	return fmt.Errorf("cound not scan into unsuported type: %w", errors.ErrUnsupported)
}

func find_field_index(typ reflect.Type, name string) (int, bool) {
	for i := range typ.NumField() {
		f := typ.Field(i)
		if tag := f.Tag.Get("csv"); tag == name || f.Name == name {
			return f.Index[0], true
		}
	}

	return -1, false
}
