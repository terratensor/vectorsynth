package glove

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Engine struct {
	Vectors map[string][]float64
}

func NewEngine(vectorFile string) (*Engine, error) {
	vectors, err := loadVectors(vectorFile)
	if err != nil {
		return nil, err
	}
	return &Engine{Vectors: vectors}, nil
}

func loadVectors(filename string) (map[string][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	vectors := make(map[string][]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			continue
		}
		word := parts[0]
		vector := make([]float64, len(parts)-1)
		for i := 1; i < len(parts); i++ {
			val, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				return nil, err
			}
			vector[i-1] = val
		}
		vectors[word] = vector
	}

	return vectors, scanner.Err()
}

func (e *Engine) FindSynonyms(expr string, topN int) ([]Similarity, error) {
	targetVector, err := e.ParseVectorExpression(expr)
	if err != nil {
		return nil, err
	}

	var similarities []Similarity
	for word, vec := range e.Vectors {
		sim, err := cosineSimilarity(targetVector, vec)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, Similarity{Word: word, Similarity: sim})
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	if topN > len(similarities) {
		topN = len(similarities)
	}
	return similarities[:topN], nil
}

func (e *Engine) ParseVectorExpression(expr string) ([]float64, error) {
	parts := strings.Fields(expr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	var vectorLength int
	for _, vec := range e.Vectors {
		vectorLength = len(vec)
		break
	}

	result := make([]float64, vectorLength)
	operation := "+"

	for _, part := range parts {
		switch part {
		case "+", "-":
			operation = part
		default:
			vec, ok := e.Vectors[part]
			if !ok {
				return nil, fmt.Errorf("word '%s' not found", part)
			}

			switch operation {
			case "+":
				for i := range result {
					result[i] += vec[i]
				}
			case "-":
				for i := range result {
					result[i] -= vec[i]
				}
			}
		}
	}

	return result, nil
}

func cosineSimilarity(vec1, vec2 []float64) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("vectors must have same length")
	}

	var dotProduct, magnitude1, magnitude2 float64
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		magnitude1 += vec1[i] * vec1[i]
		magnitude2 += vec2[i] * vec2[i]
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, fmt.Errorf("one of the vectors has zero length")
	}

	return dotProduct / (magnitude1 * magnitude2), nil
}

type Similarity struct {
	Word       string  `json:"word"`
	Similarity float64 `json:"similarity"`
}
