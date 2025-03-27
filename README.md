# VectorSynth

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/terratensor/vectorsynth.svg)](https://github.com/terratensor/vectorsynth/releases)

Инструмент для поиска семантически близких слов с использованием векторных представлений GloVe.

## Особенности

- Поиск семантически близких слов по векторному сходству
- Поддержка векторной арифметики (например: "царь - мужчина + женщина")
- Веб-интерфейс и REST API
- Поддержка Docker и Docker Swarm
- Интеграция с Traefik для автоматического HTTPS

## Быстрый старт

### Требования

- Docker и Docker Compose
- Файл с векторами (по умолчанию `/data/vectors/vectors.txt`)

### Запуск через Docker

```bash
docker run -d \
-p 8080:8080 \
-v /path/to/vectors:/data/vectors \
ghcr.io/terratensor/vectorsynth:latest
```

### Docker Swarm

```bash
docker stack deploy -c docker-compose.yml vectorsynth
```

## Использование

### Веб-интерфейс

Откройте в браузере:
`http://localhost:8080` (для локального запуска)
или
`https://vectorsynth.gmtx.ru` (для production)

### API

```bash
curl -X POST https://vectorsynth.gmtx.ru/api/similar \
-H "Content-Type: application/json" \
-d '{"expression":"компьютер", "topN":5}'
```

## Источник векторов

Векторные представления слов получены с использованием [glove-pipeline](https://github.com/terratensor/glove-pipeline) на основе текстов с сайта [SVODD](https://svodd.ru).

## Лицензия

Этот проект распространяется под лицензией MIT. См. файл [LICENSE](LICENSE).

## Разработка

### Сборка

```bash
go build -o bin/vectorsynth ./cmd/server
```

### Запуск

```bash
./bin/vectorsynth -vectors data/vectors.txt
```

### Тестирование

```bash
go test ./...
```

## Вклад в проект

PR и issues приветствуются!
