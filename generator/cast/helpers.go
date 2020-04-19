package cast

import (
	"errors"
)

func (v Θ) String() (str string) {
	return
}

func (v Θ) IsValid() bool {
	return true
}

func (v Θ) MarshalText() ([]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}

func (v *Θ) UnmarshalText(data []byte) error {
	var parsed, err = ΘEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
