#!/usr/bin/env bash
source env.sh
#
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

export SKIP_TEST=true

##
./build.sh; if [ $? -ne 0 ]; then
    exit $?
fi

#
./setup.sh no_m3

#workaround
mv -f $GOPATH/bin/linux_amd64/* $GOPATH/bin
rmdir $GOPATH/bin/linux_amd64/
rm -f $GOPATH/bin/go-bindata #not required

exit 0