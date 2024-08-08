# 程序名称
NAME = nav-green-download

# 主版本
VERSION ?= v0.0.2

# 目标输出目录
DIST_FOLDER := dist

# 版本构建目录
RELEASE_FOLDER := resources


.PHONY: build container clean
build: 
	go mod tidy
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./${DIST_FOLDER}/${NAME} ./cmd/${NAME}
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./${DIST_FOLDER}/scripts ./cmd/scripts


container: build
	docker build -t ${NAME}:${VERSION} -f ${RELEASE_FOLDER}/Dockerfile .;


clean:
	-rm -rf ${DIST_FOLDER}
	-go clean
	-go clean -cache
