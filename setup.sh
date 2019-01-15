#!/usr/bin/env bash

##
function set_env() {
    export DHNT_BASE=~/dhnt

    #
    export GOPATH=$DHNT_BASE/go
    export PATH=$GOPATH/bin:$DHNT_BASE/bin:$PATH

    #
    export IPFS_PATH=$DHNT_BASE/home/ipfs
    export GOGS_WORK_DIR=$DHNT_BASE/var/gogs
}
##
function install_m3() {
    export GO111MODULE=on

    mkdir -p ~/dhnt
    cd ~/dhnt
    git clone https://github.com/dhnt/m3.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd m3
    #no auto pull here
    ./build.sh
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
    git pull
    make clean
    git reset --hard

    make install

    ipfs init 
    #optional - change default ports
    #ipfs config Addresses
    ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/9001
    #ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001

    #$GOPATH/bin/ipfs
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
    git pull

    go install -tags "sqlite pam cert"

    #
    export GOGS_WORK_DIR=~/dhnt/var/gogs

    mkdir -p $GOGS_WORK_DIR
    cp -R $GOPATH/src/github.com/gogs/gogs/templates $GOGS_WORK_DIR
    #
}

# web terminal
function install_gotty() {
    export GOPATH=~/dhnt/go
    export GO111MODULE=on

    mkdir -p $GOPATH/src/github.com/yudai
    cd $GOPATH/src/github.com/yudai
    git clone https://github.com/yudai/gotty.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd gotty
    git pull

    go install
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

install_ipfs

install_gogs

install_gotty

install_m3

chmod -R 755 $GOPATH/bin/

echo "Done!"

##
