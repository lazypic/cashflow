#!/bin/sh
APP="cashflow"
# OS별로 빌드함.
GOOS=linux GOARCH=amd64 go build -o ./bin/linux/${APP} cashflow.go dbapi.go struct.go timefunc.go
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/${APP}.exe cashflow.go dbapi.go struct.go timefunc.go
GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/${APP} cashflow.go dbapi.go struct.go timefunc.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/windows/ && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -
cd ./bin/darwin/ && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/windows
rm -rf ./bin/darwin
