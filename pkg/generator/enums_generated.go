// generated by enum2go. DO NOT EDIT.

package generator

import (
	"errors"
	"fmt"
)

var (
	_ = errors.New
	_ = fmt.Errorf
)
var enumFormatEnum _enumFormatEnum

type _enumFormatEnum struct{}

func (_enumFormatEnum) AllValues() []enumFormat { return []enumFormat{1, 2, 3} }

func (_enumFormatEnum) AllNames() []string {
	return []string{"Strict", "Snake", "Kebab"}
}
func (_enumFormatEnum) Parse(str string) (enumFormat, error) {
	var empty enumFormat
	switch str {
	case "Strict":
		return 1, nil
	case "Snake":
		return 2, nil
	case "Kebab":
		return 3, nil
	default:
		return empty, fmt.Errorf("unexpected value %q. Valid inputs: %v", str, []string{"Strict", "Snake", "Kebab"})
	}
}
func (_enumFormatEnum) Empty() enumFormat {
	var empty enumFormat
	return empty
}
func (v enumFormat) String() (str string) {
	switch v {
	case 1:
		return "Strict"
	case 2:
		return "Snake"
	case 3:
		return "Kebab"
	default:
		return fmt.Sprintf("unexpected value %v. Valid values: %v", int(v), []string{"Strict", "Snake", "Kebab"})
	}
}
func (v enumFormat) IsZero() bool {
	var empty enumFormat
	return v == empty
}
func (v enumFormat) IsValid() bool {
	return 1 <= v && v <= 2
}
func (v enumFormat) MarshalText() (
	[]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}
func (v *enumFormat) UnmarshalText(data []byte) error {
	var parsed,
		err = enumFormatEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
func (_enumFormatEnum) Strict() enumFormat {
	return 1
}
func (_enumFormatEnum) Snake() enumFormat {
	return 2
}
func (_enumFormatEnum) Kebab() enumFormat {
	return 3
}

var (
	_ = errors.New
	_ = fmt.Errorf
)
var enumKindEnum _enumKindEnum

type _enumKindEnum struct{}

func (_enumKindEnum) AllValues() []enumKind { return []enumKind{1, 2} }

func (_enumKindEnum) AllNames() []string {
	return []string{"String", "Int"}
}
func (_enumKindEnum) Parse(str string) (enumKind, error) {
	var empty enumKind
	switch str {
	case "String":
		return 1, nil
	case "Int":
		return 2, nil
	default:
		return empty, fmt.Errorf("unexpected value %q. Valid inputs: %v", str, []string{"String", "Int"})
	}
}
func (_enumKindEnum) Empty() enumKind {
	var empty enumKind
	return empty
}
func (v enumKind) String() (str string) {
	switch v {
	case 1:
		return "String"
	case 2:
		return "Int"
	default:
		return fmt.Sprintf("unexpected value %v. Valid values: %v", int(v), []string{"String", "Int"})
	}
}
func (v enumKind) IsZero() bool {
	var empty enumKind
	return v == empty
}
func (v enumKind) IsValid() bool {
	return 1 <= v && v <= 1
}
func (v enumKind) MarshalText() (
	[]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}
func (v *enumKind) UnmarshalText(data []byte) error {
	var parsed,
		err = enumKindEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
func (_enumKindEnum) String() enumKind {
	return 1
}
func (_enumKindEnum) Int() enumKind {
	return 2
}