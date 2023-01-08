CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=github.com/cr00z/goSimpleChat/cmd/server

all: format build test lint

precommit: format build test lint
	echo "OK"

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

build: bindir
	go build -o ${BINDIR}/server ${PACKAGE}

test:
	go test ./...

lint: install-lint
	${LINTBIN} run	

run: build
	${BINDIR}/server

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/mocks/messages/messages_mocks.go

# service

bindir:
	mkdir -p ${BINDIR}

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})		

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})
