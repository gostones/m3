#!/usr/bin/env sh

##
# export GOOS=darwin
export GOOS=linux

export GOARCH=amd64
export CGO_ENABLED=0

# export SKIP_TEST=false
export SKIP_TEST=true

##
source env.sh
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
    if [ "x${SKIP_TEST}" != "xtrue" ]; then
        go test $FLAG ./...; if [ $? -ne 0 ]; then
            return 1
        fi
    fi

    echo "## Building ..."

    go build $FLAG -a -ldflags '-w -extldflags "-static"' ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    go install $FLAG ./cmd/...
    
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
