// code generate by mage script. DO NOT EDIT.

package static

const Usage = "Generates enum definitions for Golang. It supports string, byte, int(*) and uint(*) base types. To generate enum specs for a specific type, the tool searches for enum specs in the source code.\n\nExample spec:\n\n    type (\n        // enum type definition\n        Fruit int\n\n        // enum values definition. Must be in the same package.\n        _ struct {\n            Enum struct { Apple, Banana, Peach, Orange Fruit }\n        }\n    )\n\nFor each spec and type tool will generate a singleton object FruitEnum of type _FruitEnum. with static methods Apple, Banana, Peach and Orange. Each of methods returns an unique value of type Fruit.\n\n    var EnumFruit _EnumFruit\n\n    type _EnumFruit struct{}\n\n    func(_EnumFruit) Apple() Fruit { return 1 }\n    func(_EnumFruit) Banana() Fruit { return 2 }\n    func(_EnumFruit) Peach() Fruit { return 3 }\n    func(_EnumFruit) Orange() Fruit { return 4 }\n\nFor the Fruit type tool will generate util methods such as .String, .MarshalText, .IsValid, etc."

