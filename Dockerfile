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
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  " -w -s \
    -X main.gbuild=$GBUILD \
    -X main.hash=$LONGHASH \
    -X main.short=$SHORTHASH \
    -X main.date=$COMMITDATE \
    -X main.count=$COMMITCOUNT \
    -X main.build=$BUILDDATE" -a -o proxy .

# Final stage
FROM gcr.io/distroless/base
COPY --from=0 /src/proxy /
ENV PORT="8000"
EXPOSE 8000
ENTRYPOINT ["/proxy"]