package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "mrshanahan.com/notes-indexer/pkg/tokenizer"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    lines := []string{}

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    text := strings.Join(lines, "\n")
    t := tokenizer.New()
    tokens := t.Tokenize(text)
    fmt.Println(tokens)
}
