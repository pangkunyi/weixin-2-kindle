export GOPATH:=$(shell pwd):$(GOPATH)
export GOBIN:=$(shell pwd)/bin

install:
	@go install weixin2kindle
run:
	@./bin/weixin2kindle
stop:
	@pkill weixin2kindle || echo "no weixin2kindle process"
