BUILD_DIR = ./build

all:

build:
	go build -o ${BUILD_DIR}/replacer

install:
	go install

clean:
	rm -rf ${BUILD_DIR}
