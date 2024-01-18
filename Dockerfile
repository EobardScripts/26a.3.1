FROM golang:latest AS compiling_stage
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main main.go

FROM alpine:latest
LABEL version="1.0.0"
WORKDIR /app
COPY --from=compiling_stage main /bin/main
ENTRYPOINT ["/bin/main"]