#!/usr/bin/env bash

##
[[ $DEBUG ]] && FLAG="-x"

function build() {
    echo "## Cleaning ..."
    go clean $FLAG ./...

    echo "## Formatting ..."
    go fmt $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi
    
    echo "## Vetting ..."
    go vet $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    echo "## Testing ..."
    go test $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    echo "## Building ..."
#    go build $FLAG -buildmode=exe -o bin/mirr -ldflags '-extldflags "-static"'; if [ $? -ne 0 ]; then
#        return 1
#    fi

    go build $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    #bin
    go build $FLAG -buildmode=exe -o $GOPATH/bin/mirr ./cmd/mirr; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o $GOPATH/bin/hexid ./cmd/hexid; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o $GOPATH/bin/systray ./cmd/systray; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o $GOPATH/bin/pmd ./cmd/pmd; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o $GOPATH/bin/m3d ./cmd/m3d; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o $GOPATH/bin/pmctl ./cmd/pmctl; if [ $? -ne 0 ]; then
        return 1
    fi

    echo "## Tidying up modules ..."
    go mod tidy
}

echo "#### Building ..."

build; if [ $? -ne 0 ]; then
    echo "#### Build failure"
    exit 1
fi

echo "#### Build success"

exit 0
