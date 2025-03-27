.PHONY: build run docker-build docker-run docker-push

build:
	go build -o bin/vectorsynth ./cmd/server

run:
	./bin/vectorsynth -vectors data/vectors.txt

docker-build:
	docker build -t ghcr.io/terratensor/vectorsynth:latest .

docker-run:
	docker run -p 8080:8080 -v ./data/vectors:/data/vectors ghcr.io/terratensor/vectorsynth:latest

docker-push:
	docker push ghcr.io/terratensor/vectorsynth:latest

deploy:
	scp docker-compose.yml ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }}:/path/to/deploy
	ssh ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }} "cd /path/to/deploy && docker-compose pull && docker-compose up -d"