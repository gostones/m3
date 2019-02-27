#!/usr/bin/env bash

set -x
#linux
export GOOS=linux

export GOARCH=amd64
export CGO_ENABLED=0

#
export DHNT_BASE=$PWD/build

##
function set_env() {
    #
    mkdir -p $DHNT_BASE

    #
    export GO111MODULE=auto
    export GOPATH=$DHNT_BASE/go
    export PATH=$GOPATH/bin:$DHNT_BASE/bin:$PATH

    #
    mkdir -p $DHNT_BASE/go/bin
    mkdir -p $DHNT_BASE/home
    mkdir -p $DHNT_BASE/etc
}

# ipfs
function install_ipfs() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/ipfs
    cd $GOPATH/src/github.com/ipfs
    # rm -rf go-ipfs

    git clone https://github.com/ipfs/go-ipfs.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd go-ipfs
    
    #
    # some dependant tools need to be executable on the build system
    (GOOS=  GOARCH= make clean build)

    #
    make install

    # initialization example:
    # ipfs init 
    # #optional - change default ports
    # #ipfs config Addresses
    # ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/9001
    # #ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
    # ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://ipfs.home", "http://127.0.0.1:5001", "https://webui.ipfs.io"]'
    # ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
    # #$GOPATH/bin/ipfs
}

# git server
function install_gogs() {
    export CC=x86_64-linux-musl-gcc
    export CXX=x86_64-linux-musl-g++ 

    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/gogs
    cd $GOPATH/src/github.com/gogs
    git clone https://github.com/gogs/gogs.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd gogs
    git fetch && git fetch --tags
    git checkout v0.11.79

    #
    case "$GOOS" in
        darwin)
            go install -tags "sqlite cert netgo" -a
            ;;
        linux)
            CGO_ENABLED=1 go install -tags "sqlite cert netgo" -a -ldflags '-w -linkmode external -extldflags "-static"'
            ;;
        *)
            echo "not supported"
            exit 1
    esac

    #
    mkdir -p $DHNT_BASE/home/gogs
    rm -rf $DHNT_BASE/home/gogs/templates
    rm -rf $DHNT_BASE/home/gogs/public
    cp -R $GOPATH/src/github.com/gogs/gogs/templates $DHNT_BASE/home/gogs
    cp -R $GOPATH/src/github.com/gogs/gogs/public $DHNT_BASE/home/gogs
    #
}

# web terminal
function install_gotty() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/yudai
    cd $GOPATH/src/github.com/yudai
    git clone https://github.com/yudai/gotty.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd gotty
    git pull

    go install -a -ldflags '-w -extldflags "-static"'
}

# traefik
function install_traefik() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/containous
    cd $GOPATH/src/github.com/containous
    git clone https://github.com/containous/traefik.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd traefik
    git checkout v1.7

    # go-bindata needs to be executable on the build system
    (GOOS=  GOARCH= go get github.com/containous/go-bindata/...)

    go generate
    go install -a -ldflags '-w -extldflags "-static"' ./cmd/traefik

    #web ui
    cd $GOPATH/src/github.com/containous/traefik/webui
    yarn install
    yarn run build

    mkdir -p $DHNT_BASE/home/traefik
    rm -rf $DHNT_BASE/home/traefik/*
    cp -R $GOPATH/src/github.com/containous/traefik/static $DHNT_BASE/home/traefik
}

# reverse proxy
function install_frp() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=on

    mkdir -p $GOPATH/src/github.com/fatedier
    cd $GOPATH/src/github.com/fatedier
    git clone -b m3 https://github.com/gostones/frp.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd frp

    make
    cp bin/* $GOPATH/bin/linux_amd64
    mkdir -p $DHNT_BASE/etc/frp
    cp $GOPATH/src/github.com/fatedier/frp/conf/* $DHNT_BASE/etc/frp/
}

# gost
function install_gost() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/ginuerzh
    cd $GOPATH/src/github.com/ginuerzh
    git clone -b v2.7.2 https://github.com/gostones/gost.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd gost

    go install ./cmd/...
}

# etcd
function install_etcd() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/etcd-io
    cd $GOPATH/src/github.com/etcd-io
    git clone -b v3.3.12 https://github.com/gostones/etcd.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd etcd

    ./build
    cp bin/* $GOPATH/bin/linux_amd64
}

# caddy
function install_caddy() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/mholt
    cd $GOPATH/src/github.com/mholt
    git clone -b v0.11.4 https://github.com/gostones/caddy.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    go get github.com/caddyserver/builds
    cd caddy

    # go run build.go --goos=$GOOS --goarch=$GOARCH
    go install ./caddy/...
}

# chisel
function install_chisel() {
    export GOPATH=$DHNT_BASE/go
    export GO111MODULE=on

    mkdir -p $GOPATH/src/github.com/jpillora
    cd $GOPATH/src/github.com/jpillora
    git clone -b 1.3.1 https://github.com/gostones/chisel.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd chisel

    go install
}

function install_all() {
    install_ipfs
    install_gogs
    install_gotty
    install_traefik
    install_frp
    install_gost
    install_etcd
    install_caddy
    install_chisel
}

## setup

set_env

#
case "$1" in
        ipfs)
            install_ipfs
            ;;
        gogs)
            install_gogs
            ;;  
        gotty)
            install_gotty            
            ;;
        traefik)
            install_traefik
            ;;
        frp)
            install_frp
            ;;
        gost)
            install_gost
            ;;
        etcd)
            install_etcd
            ;;
        caddy)
            install_caddy
            ;;
        chisel)
            install_chisel
            ;;
        help)
            echo $"Usage: $0 {ipfs|gogs|gotty|traefik|frp|gost|etcd|caddy|chisel|help|_all_}"
            exit 1
            ;;
        *)
            install_all
esac

chmod -R 755 $GOPATH/bin/

echo "Done!"

##
