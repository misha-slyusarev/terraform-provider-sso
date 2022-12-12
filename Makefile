NAME=mysso
VERSION=v0.0.1
#TEST=$$(go list ./... | grep -v 'vendor')
BINARY=terraform-provider-${NAME}_${VERSION}

.PHONY: build #test testacc

default: build

build:
	go build -o ${BINARY}

#test: 
#	go test -i $(TEST) || exit 1                                                   
#	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    
#
#testacc: 
#	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
#
