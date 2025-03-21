package main

import (
	"bufio"
	"fmt"
	"glove-search/corpus"
	"glove-search/search"
	"glove-search/vectorizer"
	"os"
)

func main() {
	// Загрузка векторов
	vectors, err := vectorizer.LoadVectors("data/vectors.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки векторов:", err)
		return
	}

	// Загрузка корпуса текста
	corpus, err := corpus.LoadCorpus("data/cleaned_corpus.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки корпуса:", err)
		return
	}

	// Попытка загрузить векторизированный корпус из файла
	vectorizedCorpusFile := "data/vectorized_corpus.gob"
	if err := corpus.LoadVectorizedCorpus(vectorizedCorpusFile); err != nil {
		fmt.Println("Векторизированный корпус не найден, начинаем векторизацию...")
		// Векторизация корпуса
		if err := corpus.VectorizeCorpus(vectors); err != nil {
			fmt.Println("Ошибка векторизации корпуса:", err)
			return
		}
		// Сохранение векторизированного корпуса
		if err := corpus.SaveVectorizedCorpus(vectorizedCorpusFile); err != nil {
			fmt.Println("Ошибка сохранения векторизированного корпуса:", err)
			return
		}
		fmt.Println("Векторизация корпуса завершена и сохранена.")
	} else {
		fmt.Println("Векторизированный корпус успешно загружен.")
	}

	// Выбор метода поиска
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Выберите метод поиска (cosine/knn):")
	scanner.Scan()
	method := scanner.Text()

	if method != "cosine" && method != "knn" {
		fmt.Println("Неправильный метод. Используйте 'cosine' или 'knn'.")
		return
	}

	// Интерактивный поиск
	for {
		fmt.Println("Введите фразу для поиска (или 'exit' для выхода):")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Преобразуем фразу в вектор
		queryVector := vectorizer.TextToVector(input, vectors)
		if queryVector == nil || isZeroVector(queryVector) {
			fmt.Println("Фраза не содержит слов из векторов или вектор имеет нулевую длину.")
			continue
		}

		// Поиск по корпусу
		topN := 10
		var results []search.SearchResult
		var err error

		switch method {
		case "cosine":
			results, err = search.Search(queryVector, corpus.VectorizedLines, corpus.Lines, topN)
		case "knn":
			results, err = search.KNNSearch(queryVector, corpus.VectorizedLines, corpus.Lines, topN)
		}

		if err != nil {
			fmt.Println("Ошибка поиска:", err)
			continue
		}

		// Вывод результатов
		fmt.Printf("Топ-%d результатов:\n", topN)
		for i, result := range results {
			fmt.Printf("%d. Сходство: %.4f, Текст: %s\n", i+1, result.Similarity, result.Line)
		}
		fmt.Println()
	}
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
