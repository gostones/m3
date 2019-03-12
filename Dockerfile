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

RUN git clone https://github.com/dhnt/home.git /app/home

# RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' ./...
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...

###
FROM alpine

ARG DHNT_USER=dhnt
ARG DHNT_PWD=password

RUN apk add --no-cache curl git sudo wget

# sudo
RUN adduser ${DHNT_USER} -D \
    && echo "${DHNT_USER}:${DHNT_PWD}" | chpasswd \
    && sed -e 's;^# \(%sudo.*ALL\);\1;g' -i /etc/sudoers \
    && addgroup sudo \
    && adduser ${DHNT_USER} sudo

# RUN curl -L -o /usr/bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
#     && chmod +x /usr/bin/kubectl \
#     && kubectl version --client

#https://github.com/docker/compose/issues/3465
# RUN curl -L -o /usr/bin/docker-compose "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" \
#     && chmod +x /usr/bin/docker-compose \
#     && docker-compose --version

##
VOLUME /dhnt/etc /home/dhnt
EXPOSE 18080 8080 1080
WORKDIR /

COPY --from=builder /dist /dhnt
COPY --from=builder /go/bin/* /dhnt/bin/
COPY --from=builder /app/home/public /var/caddy/hugo/www

ENV PATH="/dhnt/bin:${PATH}"
ENV DHNT_BASE=/dhnt

CMD ["/dhnt/bin/m3", "run", "--base", "/dhnt"]
##