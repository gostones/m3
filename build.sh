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

    # echo "## Testing ..."
    # go test $FLAG ./...; if [ $? -ne 0 ]; then
    #     return 1
    # fi

    echo "## Building ..."
#    go build $FLAG -buildmode=exe -o bin/mirr -ldflags '-extldflags "-static"'; if [ $? -ne 0 ]; then
#        return 1
#    fi

    go build $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    #bin
    go build $FLAG -buildmode=exe -o bin/mirr ./cmd/mirr; if [ $? -ne 0 ]; then
        return 1
    fi
    go build $FLAG -buildmode=exe -o bin/hexid ./cmd/hexid; if [ $? -ne 0 ]; then
        return 1
    fi
}

echo "#### Building ..."

build; if [ $? -ne 0 ]; then
    echo "#### Build failure"
    exit 1
fi

echo "#### Build success"

exit 0