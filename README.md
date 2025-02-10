# lispish

## Lexer for Scheme

This repository contains a lexer implementation for the Scheme programming language, written in Go.

### Features

- Tokenizes Scheme code into various token types: Identifiers, Keywords, Literals, Operators, Delimiters, and Comments.
- Handles string literals, including escape sequences.
- Handles numeric literals in different formats.
- Provides detailed error messages with line number, column number, and input snippet.
- Categorizes errors into different types: syntax errors, invalid character errors, and unexpected end of input errors.
- Includes unit tests to ensure correctness.

### Usage

To use the lexer, follow these steps:

1. Clone the repository:

   ```sh
   git clone https://github.com/Warashi/lispish.git
   cd lispish
   ```

2. Run the lexer on a Scheme file:

   ```sh
   go run lexer/lexer.go <path-to-scheme-file>
   ```

3. Run the unit tests:

   ```sh
   go test ./lexer -v
   ```

### Example

Here is an example of using the lexer to tokenize a Scheme code snippet:

```go
package main

import (
    "fmt"
    "github.com/Warashi/lispish/lexer"
)

func main() {
    input := "(define (square x) (* x x))"
    l := lexer.NewLexer(input)

    for {
        tok, err := l.NextToken()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }
        if tok.Type == lexer.EOF {
            break
        }
        fmt.Printf("Token: %+v\n", tok)
    }
}
```

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
