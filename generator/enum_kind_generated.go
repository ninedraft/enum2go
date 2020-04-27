// code generated. DO NOT EDIT.

package generator

import (
	"errors"
	"fmt"
)

var (
	_ = errors.New
	_ = fmt.Errorf
)

var enumKindEnum _enumKindEnum

type _enumKindEnum struct{}

func (_enumKindEnum) Empty() enumKind {
	var empty enumKind
	return empty
}

func (_enumKindEnum) AllValues() []enumKind { return []enumKind{1, 2} }

func (_enumKindEnum) AllNames() []string { return []string{"int", "string"} }

func (_enumKindEnum) Parse(str string) (enumKind, error) {
	var empty enumKind
	switch str {
	case "int":
		return 1, nil
	case "string":
		return 2, nil
	default:
		return empty, fmt.Errorf("unexpected value %q. Valid inputs: %v", str, []string{"int", "string"})
	}
}
func (v enumKind) String() (str string) {
	switch v {
	case 1:
		return "int"
	case 2:
		return "string"
	default:
		return fmt.Sprintf("unexpected value %v. Valid values: %v", int(v), []string{"int", "string"})
	}
}

func (v enumKind) IsValid() bool { return 1 <= v && v <= 1 }

func (v enumKind) MarshalText() ([]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}

func (v *enumKind) UnmarshalText(data []byte) error {
	var parsed, err = enumKindEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
func (_enumKindEnum) Int() enumKind {
	return 1
}
func (_enumKindEnum) String() enumKind {
	return 2
}
