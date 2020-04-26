package cast

import (
	"errors"
	"fmt"
)

var (
	_ = errors.New
	_ = fmt.Errorf
)

// ΘEnum is a template singletone for code generation.
var ΘEnum _ΘEnum

// _ΘEnum is a template type for code generation.
type _ΘEnum struct{}

func (_ΘEnum) AllValues() []Θ { return nil }

func (_ΘEnum) AllNames() []string { return nil }

func (_ΘEnum) Parse(str string) (Θ, error) {
	return 0, nil
}

func (_ΘEnum) Empty() Θ {
	var empty Θ
	return empty
}
