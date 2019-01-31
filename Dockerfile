# Build stage
FROM golang
ENV GO111MODULE=on
WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download
COPY . /src
RUN GBUILD=$(go version | sed -e "s/ /-/g") && \
    LONGHASH=$(git log -n1 --pretty="format:%H") && \ 
    SHORTHASH=$(git log -n1 --pretty="format:%h") && \
    COMMITDATE=$(git log -n1 --pretty="format:%cd"| sed -e "s/ /-/g") && \
    COMMITCOUNT=$(git rev-list HEAD --count| cat) && \
    BUILDDATE=$(date| sed -e "s/ /-/g") && \
    PORT=8000 && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  " -w -s \
    -X github.com/g0lang/proxy/src/config.gbuild=$GBUILD \
    -X github.com/g0lang/proxy/src/config.hash=$LONGHASH \
    -X github.com/g0lang/proxy/src/config.short=$SHORTHASH \
    -X github.com/g0lang/proxy/src/config.date=$COMMITDATE \
    -X github.com/g0lang/proxy/src/config.count=$COMMITCOUNT \
    -X github.com/g0lang/proxy/src/config.port=$PORT \
    -X github.com/g0lang/proxy/src/config.build=$BUILDDATE" -a -o proxy .

# Final stage
FROM gcr.io/distroless/base
COPY --from=0 /src/proxy /
ENV PORT="8000"
EXPOSE 8000
ENTRYPOINT ["/proxy"]