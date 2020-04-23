package cast

import (
	"errors"
)

func (v Θ) String() (str string) {
	return
}

// IsValid returns true if the value of Θ is valid.
func (v Θ) IsValid() bool {
	return true
}

// MarshalText encodes Θ to a text representation.
func (v Θ) MarshalText() ([]byte, error) {
	var str = v.String()
	if v.IsValid() {
		return []byte(str), nil
	}
	return nil, errors.New(str)
}

// UnmarshalText decodes Θ from a text representation.
func (v *Θ) UnmarshalText(data []byte) error {
	var parsed, err = ΘEnum.Parse(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}
