export GOPATH:=$(shell pwd)/gopath


test: deps
	go test github.com/bearstech/webpubsub/pathmatch
	go test github.com/bearstech/webpubsub/mailbox

gopath/src/github.com/bearstech/webpubsub:
	mkdir -p gopath/src/github.com/bearstech/
	ln -s ../../../..  gopath/src/github.com/bearstech/webpubsub

deps: gopath/src/github.com/bearstech/webpubsub
