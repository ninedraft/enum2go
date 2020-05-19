package sample

type (
	Foo int
	_   struct {
		Enum struct {
			A, B, C Foo // will emit enum values
			D, E    Foo // will be ignored
		}
	}
)
