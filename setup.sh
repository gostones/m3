#!/usr/bin/env bash

##
function install_m3() {
    export GO111MODULE=on

    cd ~/dhnt
    git clone https://github.com/dhnt/m3.git; if [ $? -ne 0 ]; then
        echo "Git repo exists?"
    fi
    cd m3
    #no auto pull here
    ./build.sh

    # cp ~/dhnt/m3/bin/mirr $GOPATH/bin/
}

function install_ipfs() {
    export GOPATH=~/dhnt/go
    export GO111MODULE=off

    mkdir $GOPATH/src/github.com/ipfs
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

    #optional - change default ports
    #ipfs config Addresses
    ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/9001
    #ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001

    #$GOPATH/bin/ipfs
}

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

    go build -tags "sqlite pam cert"

    #./gogs
    #export GOGS_WORK_DIR=./
}


## setup

#mkdir -p ~/dhnt
#mkdir -p ~/dhnt/go
mkdir -p ~/dhnt/go/bin

export GO111MODULE=auto
export GOPATH=~/dhnt/go
export PATH=$GOPATH/bin:$PATH

#

install_ipfs

install_gogs

install_m3

chmod -R 755 $GOPATH/bin/

echo "Done!"

##
