package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/terratensor/vectorsynth/internal/glove"
)

func main() {
	vectorFile := flag.String("vectors", "data/vectors.txt", "Path to vectors file")
	flag.Parse()

	engine, err := glove.NewEngine(*vectorFile)
	if err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	fmt.Println("VectorSynth - Word Vector Search")
	fmt.Println("===============================")
	fmt.Println("Enter a vector expression (e.g.: 'king - man + woman')")
	fmt.Println("Or just a word to find similar words")
	fmt.Println("Type 'exit' to quit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		results, err := engine.FindSynonyms(input, 20)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("Results:")
		for i, item := range results {
			fmt.Printf("%d. %s (%.4f)\n", i+1, item.Word, item.Similarity)
		}
		fmt.Println()
	}
}
