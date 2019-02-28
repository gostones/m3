##
FROM golang:1.11.5-alpine3.9 as builder

RUN apk add --no-cache git

#
WORKDIR /
ARG M3EXT_VERSION=v0.0.1
RUN wget -qO- "https://github.com/dhnt/m3-ext/releases/download/${M3EXT_VERSION}/m3-ext.tar.gz" \
    | tar -xzv

#
COPY . /app
WORKDIR /app

#
# RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' ./...
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...

###
FROM alpine

COPY --from=builder /dist /dhnt
COPY --from=builder /go/bin/* /dhnt/bin/

ENV PATH="/dhnt/bin:${PATH}"

EXPOSE 80 443 18080

CMD ["/dhnt/bin/m3", "run", "--base", "/dhnt"]
##