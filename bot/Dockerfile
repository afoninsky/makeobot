FROM golang:1.11 AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV GOBIN /bin
WORKDIR /src
COPY . /src
RUN go get -v
RUN go build -a -installsuffix nocgo -o /service .

FROM alpine AS certificates
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /service /service
ENTRYPOINT ["/service"]
