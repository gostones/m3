#!/usr/bin/env bash

#linux
# export CC=x86_64-linux-musl-gcc
# export CXX=x86_64-linux-musl-g++ 

##
function set_env() {
    export DHNT_BASE=~/dhnt

    #
    export GOPATH=$DHNT_BASE/go
    export PATH=$GOPATH/bin:$DHNT_BASE/bin:$PATH

    #
    # export IPFS_PATH=$DHNT_BASE/home/ipfs
    # export GOGS_WORK_DIR=$DHNT_BASE/home/gogs
}
##
function install_m3() {
    export GO111MODULE=on

    mkdir -p $DHNT_BASE
    cd $DHNT_BASE
    git clone https://github.com/dhnt/m3.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd m3
    #no auto pull here
    ./build.sh

    #traefik config
    mkdir -p $DHNT_BASE/etc
    cp -R $DHNT_BASE/m3/internal/rp/traefik $DHNT_BASE/etc
}

# p2p
function install_ipfs() {
    export GOPATH=~/dhnt/go
    export GO111MODULE=off

    mkdir -p $GOPATH/src/github.com/ipfs
    cd $GOPATH/src/github.com/ipfs
    git clone https://github.com/ipfs/go-ipfs.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd go-ipfs
    
    #clean/reset
    git full
    make clean
    git reset --hard

    make install

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
    export GOPATH=~/dhnt/go
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
    export GOPATH=~/dhnt/go
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
    export GOPATH=~/dhnt/go
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

function install_no_m3() {
    install_ipfs
    install_gogs
    install_gotty
    install_traefik
}

function install_all() {
    install_no_m3
    install_m3
}

## setup

set_env

#mkdir -p ~/dhnt
mkdir -p ~/dhnt/etc

#mkdir -p ~/dhnt/go
mkdir -p ~/dhnt/go/bin

mkdir -p ~/dhnt/home

# export DHNT_BASE=~/dhnt

export GO111MODULE=auto
export GOPATH=~/dhnt/go
# export PATH=$GOPATH/bin:$PATH

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
        no_m3)
            install_no_m3
            ;;
        m3)
            install_m3
            ;;
        help)
            echo $"Usage: $0 {ipfs|gogs|gotty|traefik|m3|help|_all_}"
            exit 1
            ;;
        *)
            install_all
esac

chmod -R 755 $GOPATH/bin/

echo "Done!"

##
