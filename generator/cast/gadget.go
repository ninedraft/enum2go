package cast

var ΘEnum _ΘEnum

type _ΘEnum struct{}

func (_ΘEnum) AllValues() []Θ { return nil }

func (_ΘEnum) AllNames() []string { return nil }

func (_ΘEnum) ToStrMap() map[string]Θ { return nil }

func (_ΘEnum) FromStrMap() map[Θ]string { return nil }

func (_ΘEnum) Parse(str string) (Θ, error) {
	return 0, nil
}
