package vectorizer

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// LoadVectors загружает векторы из файла.
func LoadVectors(filename string) (map[string][]float64, error) {
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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return vectors, nil
}

// TextToVector преобразует текст в средний вектор слов.
func TextToVector(text string, vectors map[string][]float64) []float64 {
	words := strings.Fields(text)
	if len(words) == 0 || len(vectors) == 0 {
		return nil
	}

	var vectorLength int
	for _, vec := range vectors {
		vectorLength = len(vec)
		break
	}

	sumVector := make([]float64, vectorLength)
	count := 0

	for _, word := range words {
		if vec, ok := vectors[word]; ok {
			for i := range vec {
				sumVector[i] += vec[i]
			}
			count++
		}
	}

	if count > 0 {
		for i := range sumVector {
			sumVector[i] /= float64(count)
		}
		return sumVector
	}

	return nil
}

// CosineSimilarity вычисляет косинусное сходство между двумя векторами.
func CosineSimilarity(vec1, vec2 []float64) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("векторы должны быть одинаковой длины")
	}

	var dotProduct, magnitude1, magnitude2 float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		magnitude1 += vec1[i] * vec1[i]
		magnitude2 += vec2[i] * vec2[i]
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, fmt.Errorf("один из векторов имеет нулевую длину")
	}

	return dotProduct / (magnitude1 * magnitude2), nil
}
