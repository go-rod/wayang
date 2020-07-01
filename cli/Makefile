GO_FILES? = $$(find . -name '*.go' |grep -v vendor)
ifeq ($(OS),Windows_NT)
    output := wayang.exe
else
    output := wayang
endif

all: build help

build:
	go build -o ${output} ${GO_FILES?}

help:
	@echo "The wayang command line tool has now been built."
	@cat usage.txt
	@echo "Here is a sample json script that will navigate to google and search wayang:"
	@cat example.json
	@cat ${output} -file

