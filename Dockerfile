##
FROM golang:1.11.5-alpine3.9 as builder

RUN apk add --no-cache git

#
WORKDIR /
ARG M3EXT_VERSION=v0.0.4
RUN wget --no-check-certificate -qO- "https://github.com/dhnt/m3-ext/releases/download/${M3EXT_VERSION}/m3-ext.tar.gz" \
    | tar -xzv
#
COPY . /app
WORKDIR /app

#https://github.com/moby/moby/issues/15858
COPY dhnt/etc /dist/etc

#
# RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' ./...
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...

###
FROM alpine

ARG DHNT_USER=dhnt
ARG DHNT_PWD=password

RUN apk add --no-cache git curl openssl sudo

# sudo
RUN adduser ${DHNT_USER} -D \
    && echo "${DHNT_USER}:${DHNT_PWD}" | chpasswd \
    && sed -e 's;^# \(%sudo.*ALL\);\1;g' -i /etc/sudoers \
    && addgroup sudo \
    && adduser ${DHNT_USER} sudo

VOLUME /dhnt/etc
EXPOSE 18080
WORKDIR /

COPY --from=builder /dist /dhnt
COPY --from=builder /go/bin/* /dhnt/bin/

ENV PATH="/dhnt/bin:${PATH}"
ENV DHNT_BASE=/dhnt

CMD ["/dhnt/bin/m3", "run", "--base", "/dhnt"]
##