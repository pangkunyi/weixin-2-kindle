export GOPATH:=$(shell pwd):$(GOPATH)
install:
	@go install weixin-2-kindle 
run:
	@./bin/kindle-delivery
stop:
	@pkill weixin-2-kindle || echo "no kindle-delivery process"
