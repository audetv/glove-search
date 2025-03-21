package corpus

import (
	"bufio"
	"encoding/gob"
	"glove-search/vectorizer"
	"os"
)

// Corpus представляет корпус текста.
type Corpus struct {
	Lines           []string
	VectorizedLines [][]float64
}

// LoadCorpus загружает корпус текста из файла.
func LoadCorpus(filename string) (*Corpus, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Corpus{Lines: lines}, nil
}

// VectorizeCorpus векторизирует корпус текста.
func (c *Corpus) VectorizeCorpus(vectors map[string][]float64) error {
	for _, line := range c.Lines {
		vector := vectorizer.TextToVector(line, vectors)
		if vector != nil && !isZeroVector(vector) {
			c.VectorizedLines = append(c.VectorizedLines, vector)
		}
	}
	return nil
}

// SaveVectorizedCorpus сохраняет векторизированный корпус в файл.
func (c *Corpus) SaveVectorizedCorpus(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(c.VectorizedLines)
}

// LoadVectorizedCorpus загружает векторизированный корпус из файла.
func (c *Corpus) LoadVectorizedCorpus(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&c.VectorizedLines)
}

// isZeroVector проверяет, является ли вектор нулевым.
func isZeroVector(vector []float64) bool {
	for _, val := range vector {
		if val != 0 {
			return false
		}
	}
	return true
}
