GVER		:= $$(go version | sed -e "s/ /-/g")
LONGHASH	:= $$(git log -n1 --pretty="format:%H")
SHORTHASH	:= $$(git log -n1 --pretty="format:%h")
COMMITDATE	:= $$(git log -n1 --pretty="format:%cd"| sed -e "s/ /-/g")
COMMITCOUNT	:= $$(git rev-list HEAD --count| cat)
BUILDDATE	:= $$(date| sed -e "s/ /-/g")
CGO_ENABLED	:= 0
GOOS		:= linux
GOARCH		:= amd64
GO111MODULE	:= on
IMGNAME		:= proxy
IMGTAG		:= ${IMGNAME}:${SHORTHASH}
LATEST		:= ${IMGNAME}:latest

# Run commands with the debugger. (default: false)
DEBUG ?= false

# Show this help prompt.
help:
	@ echo
	@ echo '  Usage:'
	@ echo ''
	@ echo '    make <target> [flags...]'
	@ echo ''
	@ echo '  Targets:'
	@ echo ''
	@ awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':'
	@ echo ''
	@ echo '  Flags:'
	@ echo ''
	@ awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?\?=/{ print "   ", $$1, $$2, comment }' $(MAKEFILE_LIST) | column -t -s '?=' | sort
	@ echo ''

# Show variable.
vars:
	@ echo '  Variable:'
	@ echo ''
	@ awk '/^[A-Z]+\t*[^\t]+?\:=/{ print "   ", $$1 }' $(MAKEFILE_LIST) | sort
	@ echo ''

# Build app 
build:
	@echo "Building Project Binary To ./bin"
	@GOARC=${GOARCH} GOOS=${GOOS} CGO_ENABLED=${CGO_ENABLED} go build -ldflags  " -w -s \
    -X main.gver=${GVER} \
    -X main.hash=${LONGHASH} \
    -X main.short=${SHORTHASH} \
    -X main.date=${COMMITDATE} \
    -X main.count=${COMMITCOUNT} \
    -X main.build=${BUILDDATE}" -a -o bin/proxy .

# Run app
run:
	@go run -ldflags  " \
    -X main.gver=${GVER} \
    -X main.hash=${LONGHASH} \
    -X main.short=${SHORTHASH} \
    -X main.date=${COMMITDATE} \
    -X main.count=${COMMITCOUNT} \
    -X main.build=${BUILDDATE}" . 

# Run Test
coverage:
	@go test ./... -cover

# Run Test
test:
	@go test ./...

# Build docker image.
ibuild:
	@ docker build -t ${IMGTAG} .
	@ docker tag ${IMGTAG} ${LATEST}

# Build docker no cache
ibuild-nc:
	@ docker build --no-cache -t ${IMGTAG} .
	@ docker tag ${IMGTAG} ${LATEST}

# Run docker image.
irun:
	@ docker run -p 8000:8000 ${LATEST}

# Clean Docker.
iclean:
	@ docker container prune -f 
	@ docker image prune -f 
	@ docker image rm ${IMGTAG} ${LATEST}