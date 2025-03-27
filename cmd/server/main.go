package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/terratensor/vectorsynth/internal/glove"
)

func main() {
	vectorFile := flag.String("vectors", "data/vectors.txt", "Путь к файлу векторов")
	port := flag.String("port", "8080", "Порт для сервера")
	webDir := flag.String("web", "internal/web", "Директория с фронтендом")
	flag.Parse()

	// Получаем абсолютный путь к директории с веб-файлами
	absWebDir, err := filepath.Abs(*webDir)
	if err != nil {
		log.Fatalf("Ошибка получения абсолютного пути: %v", err)
	}

	// Проверяем существование директории
	if _, err := os.Stat(absWebDir); os.IsNotExist(err) {
		log.Fatalf("Директория с фронтендом не существует: %s", absWebDir)
	}

	engine, err := glove.NewEngine(*vectorFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации движка: %v", err)
	}

	// API endpoint
	http.HandleFunc("/api/similar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Expression string `json:"expression"`
			TopN       int    `json:"topN"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		if request.TopN == 0 {
			request.TopN = 20
		}

		results, err := engine.FindSynonyms(request.Expression, request.TopN)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка обработки: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	// Статические файлы фронтенда
	fs := http.FileServer(http.Dir(absWebDir))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("Сервер запущен на порту %s", *port)
	log.Printf("Фронтенд доступен из директории: %s", absWebDir)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
