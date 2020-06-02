all: compile run

compile:
	@echo "building into ${pwd}/build directory...\n"
	@go build -o ./build/parasut-client

run:
	@./build/parasut-client
