package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"mrshanahan.com/notes-indexer/internal/util"
	"mrshanahan.com/notes-indexer/pkg/stemmer"
	"mrshanahan.com/notes-indexer/pkg/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: expected command")
		os.Exit(1)
	}

	command := os.Args[1]
	if strings.ToLower(command) == "stemmer" {
		stem()
	} else if strings.ToLower(command) == "tokenizer" {
		tokenize()
	} else if strings.ToLower(command) == "markdown" {
		parseMarkdown()
	} else {
		fmt.Fprintf(os.Stderr, "error: invalid command: %s", command)
		os.Exit(1)
	}
}

func parseMarkdown() {
	// if len(os.Args) > 2 {
	// 	for _, f := range os.Args[2:] {
	// 		bs, err := os.ReadFile(f)
	// 		if err != nil {
	// 			log.Fatalf("error: failed to read file %s: %v", f, err)
	// 		}
	// 		text := string(bs)
	// 		doc, err := markdown.Parse(text)
	// 		if err != nil {
	// 			log.Fatalf("error: failed to parse document in %s: %v", f, err)
	// 		}
	// 		log.Printf("%v", doc)
	// 	}
	// } else {
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	lines := []string{}

	// 	for scanner.Scan() {
	// 		lines = append(lines, scanner.Text())
	// 	}
	// 	if err := scanner.Err(); err != nil {
	// 		log.Fatalf("error: failed to read from stdin: %v", err)
	// 	}

	// 	text := strings.Join(lines, "\n")
	// 	doc, err := markdown.Parse(text)
	// 	if err != nil {
	// 		log.Fatalf("error: failed to parse document: %v", err)
	// 	}
	// 	log.Printf("%v", doc)
	// }
}

func tokenize() {
	var text string
	if len(os.Args) > 2 {
		text = strings.Join(os.Args[2:], " ")
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		lines := []string{}

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		text = strings.Join(lines, "\n")
	}

	t := tokenizer.NewDefault()
	tokens, err := t.Tokenize(text)
	if err != nil {
		log.Fatalf("error: failed to tokenize text: %s", err)
	}
	output := strings.Join(util.Map(tokens, func(t tokenizer.Token) string { return t.Value }), "\n")
	fmt.Println(output)
}

func stem() {
	if len(os.Args) > 2 {
		for _, t := range os.Args[2:] {
			stemmed := stemmer.Stem(t)
			fmt.Println(stemmed)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			stemmed := stemmer.Stem(line)
			fmt.Println(stemmed)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}
