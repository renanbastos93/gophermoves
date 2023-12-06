# Gopher Moves

GopherMoves is a Go package that provides an interface and an example implementation for simulating the movements of a character, particularly designed for a Gopher character.

https://github.com/renanbastos93/gophermoves/assets/8202898/4ecaa097-08e4-4a67-a80f-a311b6e1af67


## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Interface Methods](#interface-methods)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use GopherMoves in your Go project, you can simply run:

```bash
go get -u github.com/renanbastos93/gophermoves
```

## Usage

```go
package main

import (
	"github.com/renanbastos93/gophermoves"
)

func main() {
    orderMatrix := 5
    g := gophermoves.New(orderMatrix)
    g.Start()
}
```

## Interface Methods

The GopherMoves interface provides the following methods:

- `Start()`: Initiates movements for the character.
- `Reset()`: Redefines the default states of the X and Y positions.
- `Up()`: Moves the character upward.
- `Down()`: Moves the character downward.
- `Left()`: Turns the character to the left.
- `Right()`: Turns the character to the right.


## Contributing

Feel free to contribute to this project. If you find any issues or have suggestions, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
