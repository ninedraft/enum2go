Generates enum definitions for Golang.
It supports string, byte, int(*) and uint(*) base types.
To generate enum specs for a specific type, the tool searches
for enum specs in the source code.

Example spec:

    type (
        // enum type definition
        Fruit int

        // enum values definition. Must be in the same package.
        _ struct {
            Enum struct { Apple, Banana, Peach, Orange Fruit }
        }
    )

For each spec and type tool will generate a singleton object FruitEnum of type _FruitEnum.
with static methods Apple, Banana, Peach and Orange.
Each of methods returns an unique value of type Fruit.

    var EnumFruit _EnumFruit

    type _EnumFruit struct{}

    func(_EnumFruit) Apple() Fruit { return 1 }
    func(_EnumFruit) Banana() Fruit { return 2 }
    func(_EnumFruit) Peach() Fruit { return 3 }
    func(_EnumFruit) Orange() Fruit { return 4 }

For Fruit type tool will generate util methods such as .String, 
.MarshalText, .IsValid, etc. 