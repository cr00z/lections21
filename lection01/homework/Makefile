CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=github.com/cr00z/goTinkoffFintech

all: format build test lint

precommit: format build test lint
	@echo "OK"

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

build: bindir
	go build -o ${BINDIR}/homework1

test:
	go test ./...

lint: install-lint
	${LINTBIN} run

run:
	go run ${PACKAGE}

demo: build
	${BINDIR}/homework1

# service

bindir:
	@mkdir -p ${BINDIR}

install-lint: bindir
	@test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	@test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})
