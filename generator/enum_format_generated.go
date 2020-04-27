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

var enumFormatEnum _enumFormatEnum

type _enumFormatEnum struct{}

func (_enumFormatEnum) AllValues() []enumFormat { return []enumFormat{1, 2, 3} }

func (_enumFormatEnum) AllNames() []string { return []string{"strict", "kebab", "snake"} }

func (_enumFormatEnum) Parse(str string) (enumFormat, error) {
	var empty enumFormat
	switch str {
	case "strict":
		return 1, nil
	case "kebab":
		return 2, nil
	case "snake":
		return 3, nil
	default:
		return empty, fmt.Errorf("unexpected value %q. Valid inputs: %v", str, []string{"strict", "kebab", "snake"})
	}
}
func (v enumFormat) String() (str string) {
	switch v {
	case 1:
		return "strict"
	case 2:
		return "kebab"
	case 3:
		return "snake"
	default:
		return fmt.Sprintf("unexpected value %v. Valid values: %v", int(v), []string{"strict", "kebab", "snake"})
	}
}

func (v enumFormat) IsValid() bool { return 1 <= v && v <= 2 }

func (v enumFormat) MarshalText() ([]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}

func (v *enumFormat) UnmarshalText(data []byte) error {
	var parsed, err = enumFormatEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
func (_enumFormatEnum) Strict() enumFormat {
	return 1
}
func (_enumFormatEnum) Kebab() enumFormat {
	return 2
}
func (_enumFormatEnum) Snake() enumFormat {
	return 3
}
