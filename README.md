[![Go Report Card](https://goreportcard.com/badge/github.com/ninedraft/enum2go)](https://goreportcard.com/report/github.com/ninedraft/enum2go) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![go.dev reference](https://img.shields.io/badge/doc-go.dev-blue)](https://pkg.go.dev/mod/github.com/ninedraft/enum2go)

- [Installation](#installation)
- [Main idea](#main-idea)
- [Usage](#usage)
  - [CLI](#cli)
  - [Flags](#flags)
- [Licensing](#licensing)
  - [Code generation output](#code-generation-output)
  - [Tool source code](#tool-source-code)

## Installation

```sh
go get github.com/ninedraft/enum2go/cmd/enum2go
```

## Main idea

Tool generates enum definitions for Golang. It supports `string`, `byte`, `int(*)` and `uint(*)` base types. To generate enum specs for a specific type, the tool searches for enum specs in the source code.
  
Example spec:

```go  
      type (
          // enum type definition
          Fruit int
  
          // enum values definition. Must be in the same package.
          _ struct {
              Enum struct { Apple, Banana, Peach, Orange Fruit }
          }
      )
```

For each spec and type tool will generate a singleton object `FruitEnum` of type `_FruitEnum`. with static methods `Apple`, `Banana`, `Peach` and `Orange`. Each of methods returns an unique value of type `Fruit`.

```go  
      var EnumFruit _EnumFruit
  
      type _EnumFruit struct{}
  
      func(_EnumFruit) Apple() Fruit { return 1 }
      func(_EnumFruit) Banana() Fruit { return 2 }
      func(_EnumFruit) Peach() Fruit { return 3 }
      func(_EnumFruit) Orange() Fruit { return 4 }
```

For the Fruit type tool will generate util methods such as `.String`, `.MarshalText`, `.IsValid`, etc.

## Usage

### CLI

```sh
enum2go -dir ./
```

### Flags

| Flag      | Type   | Description                                             |
|-----------|--------|---------------------------------------------------------|
| -d, --dir | string | package dir to parse                                    |
| -o, --out | string | file to generated result (default "enums_generated.go") |

## Licensing

### Code generation output

[![](https://upload.wikimedia.org/wikipedia/commons/6/69/CC0_button.svg)](https://creativecommons.org/publicdomain/zero/1.0/deed.en)

You can copy, modify, distribute and use in the other ways the generated code, even for commercial purposes, all without asking permission.

### Tool source code

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

```
MIT License

Copyright (c) 2020 Petrukhin Pavel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
