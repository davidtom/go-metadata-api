BUILD_DIR=dist
BUILD_NAME=server

PORT=""

build: clean
	go build -o ${BUILD_DIR}/${BUILD_NAME} main.go handlers.go storage.go

clean:
	rm -rf ${BUILD_DIR}

run: build
	PORT=${PORT} ${BUILD_DIR}/${BUILD_NAME}
