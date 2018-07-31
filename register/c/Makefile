# Author: Zhongwen Lan(lanzhongwen@jd.com)
# Created: 2018/05/28
TOPDIR=$(shell while [ ! -d src ] && [ `pwd` != / ] ; do cd ..; done ; echo `pwd`)
GOPATH=$(TOPDIR)

all:
	export GOPATH=$(GOPATH)
	go build -buildmode=c-archive register.go

test:
	export GOPATH=$(GOPATH)
	go test
