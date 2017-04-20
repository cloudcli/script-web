start:
	gin -e .go -e .html web start

seed: create-table
	go run main.go --config=./app.ini bootstrap seed

create-table:
	go run main.go --config=./app.ini bootstrap create-tables

build:
	go build -i -v

run:
	# govendor build
	go build -i -v
	./script-web --debug web start --port=3000

pack:build
	rm -rf /tmp/scriptweb
	mkdir /tmp/scriptweb
	cp internals/app.example.ini /tmp/scriptweb/app.ini
	cp script-web /tmp/scriptweb/
	cp -r frontend /tmp/scriptweb/frontend
	(cd /tmp/ && tar -czf scriptweb.tar.gz scriptweb)
	cp /tmp/scriptweb.tar.gz .
