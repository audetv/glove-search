package search

import (
	"glove-search/vectorizer"
	"sort"
)

// SearchResult представляет результат поиска.
type SearchResult struct {
	Line       string
	Similarity float64
}

// Search выполняет поиск по корпусу.
func Search(queryVector []float64, corpus [][]float64, lines []string, topN int) ([]SearchResult, error) {
	var results []SearchResult
	for i, vec := range corpus {
		if isZeroVector(vec) {
			continue
		}
		sim, err := vectorizer.CosineSimilarity(queryVector, vec)
		if err != nil {
			return nil, err
		}
		results = append(results, SearchResult{Line: lines[i], Similarity: sim})
	}

	// Сортировка по убыванию сходства
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Выбор topN результатов
	if len(results) > topN {
		results = results[:topN]
	}

	return results, nil
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
