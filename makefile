BUILD_DIR=dist
BUILD_NAME=server

build: clean
	go build -o ${BUILD_DIR}/${BUILD_NAME} main.go handlers.go

clean:
	rm -rf ${BUILD_DIR}

run: build
	${BUILD_DIR}/${BUILD_NAME}
