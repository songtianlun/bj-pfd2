SHELL := /bin/bash
BASEDIR = $(shell pwd)
GOPROXY=https://mirrors.cloud.tencent.com/go/

# build with version infos
versionDir = "bj-pfd2/com/v"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q -E 'clean|干净';then echo clean; else echo dirty; fi)

ldflags=" -extldflags '-static' -s -w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"
# -tags netgo 解决 exec user process caused: no such file or directory
all: build
	@go build -a -v -tags netgo -ldflags ${ldflags} -o app .
.PHONY: clean
clean:
	rm -f app
	find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}
.PHONY: build
build: test
	gofmt -w .
	go vet . | grep -v vendor;true
.PHONY: ca
ca:
	openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"
.PHONY: run
run:
	if [ -e ./config.yaml ] ; then cp config.yaml ./tmp/; fi
	air
.PHONY: test
test:
	go test -v ./...
.PHONY: help
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"

