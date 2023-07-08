build:
	go build -o output/backend .

run:
	go run .

docker-build:
	docker build -t backend .