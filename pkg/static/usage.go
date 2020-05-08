// code generate by mage script. DO NOT EDIT.

package static

const Usage = "Generates enum definitions for Golang.\nIt supports string, byte, int(*) and uint(*) base types.\nTo generate enum specs for a specific type, the tool searches\nfor enum specs in the source code.\n\nExample spec:\n\n    type (\n        // enum type definition\n        Fruit int\n\n        // enum values definition. Must be in the same package.\n        _ struct {\n            Enum struct { Apple, Banana, Peach, Orange Fruit }\n        }\n    )\n\nFor each spec and type tool will generate a singleton object FruitEnum of type _FruitEnum.\nwith static methods Apple, Banana, Peach and Orange.\nEach of methods returns an unique value of type Fruit.\n\n    var EnumFruit _EnumFruit\n\n    type _EnumFruit struct{}\n\n    func(_EnumFruit) Apple() Fruit { return 1 }\n    func(_EnumFruit) Banana() Fruit { return 2 }\n    func(_EnumFruit) Peach() Fruit { return 3 }\n    func(_EnumFruit) Orange() Fruit { return 4 }\n\nFor Fruit type tool will generate util methods such as .String, \n.MarshalText, .IsValid, etc."

