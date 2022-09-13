#!/bin/bash


build(){
	go build -ldflags="-s -w" main.go
	mkdir -p $name/logs
	cp main $name/$name
	# cp -rf conf logs log-collect tcs-tools $name/
	cp -rf conf logs tcs-tools data $name/

	tar zcf $name-$GOOS-$GOARCH.tar.gz $name 
	rm -rf $name/ main
	# rm -rf $name/
}

name="tianxun-lite"
export CGO_ENABLED=0
export GOOS=linux

# amd64
export GOARCH=amd64
echo $GOOS-$GOARCH
build
