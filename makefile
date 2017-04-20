start:
	gin -e .go -e .html web start

seed: create-table
	go run main.go --config=./app.ini bootstrap seed

create-table:
	go run main.go --config=./app.ini bootstrap create-tables

build:
	go build

run:
	# govendor build
	go build -i -v
	./script-web --debug web start --port=3002
