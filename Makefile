all: clean
	go build -o build/bin/main_app cmd/task.go
	./build/bin/main_app

run:
	docker-compose up

clean:
	rm build/bin/*