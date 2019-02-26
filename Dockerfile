###
FROM golang:1.11.5-alpine3.9 as builder

RUN apk add --no-cache git

RUN mkdir /app
ADD . /app/
WORKDIR /app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' ./...
RUN CGO_ENABLED=0 GOOS=linux GOPATH=/go go install ./cmd/...

###
FROM alpine

COPY --from=builder /go/bin /app/bin
COPY --from=builder /app/build/go/bin/linux_amd64/* /app/bin/

ENV PATH="/app/bin:${PATH}"

EXPOSE 80
EXPOSE 443
EXPOSE 18080

CMD ["/app/bin/m3d", "--base", "/app/dhnt"]

##