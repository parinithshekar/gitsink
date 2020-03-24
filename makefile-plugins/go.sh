#!/usr/bin/env bash
set -euo pipefail

CACHE_DIR=.cache

function main() {
    # Ensure cache directory exists
    if [[ ! -d ${CACHE_DIR} ]]; then
	rm -rf ${CACHE_DIR} && mkdir -p ${CACHE_DIR}
    fi

    # Just launch the parameters passed as: function var-args...
    $@
}

__install-go-package() {
    install_path=$1
    mkdir -p $(dirname $install_path)

    package=$2
    binname=$3

    if [[ ! -f ${CACHE_DIR}/bin/$binname ]]; then
	# Do the install from /tmp to prevent modification of go.mod and go.sum
	cwd=$(pwd)
        (cd /tmp && GO111MODULE=on GOPATH=$cwd/.cache go get $package)

	# Make files writeable so make clean change remove them
	chmod -R +w $cwd/.cache/pkg/mod
    fi

    if [[ ! -f $install_path ]]; then
        cp ${CACHE_DIR}/bin/$binname $install_path
    fi
}

install-dlv() {
    install_path=$1
    VERSION="master"
    __install-go-package $install_path github.com/go-delve/delve/cmd/dlv@$VERSION dlv
}

install-gox() {
    install_path=$1
    VERSION="v1.0.1"
    __install-go-package $install_path github.com/mitchellh/gox@$VERSION gox
}

install-ttyrec2gif() {
    install_path=$1
    VERSION="master"
    __install-go-package $install_path github.com/sugyan/ttyrec2gif@$VERSION ttyrec2gif
}

install-golangci-lint() {
    install_path=$1
    mkdir -p $(dirname $install_path)

    VERSION="1.23.8"
    OS=$(uname -s  | tr '[:upper:]' '[:lower:]')
    FILE="golangci-lint-${VERSION}-${OS}-amd64.tar.gz"
    URL="https://github.com/golangci/golangci-lint/releases/download/v${VERSION}/${FILE}"
    CACHE_FILE="${CACHE_DIR}/${FILE}"

    # Ensure the source file exists
    if [[ ! -f $CACHE_FILE ]]; then
	curl -fsSL $URL > $CACHE_FILE
    fi

    # Create the target
    if [[ ! -f $install_path ]]; then
	tar xzf $CACHE_FILE -C $CACHE_DIR
	mv $CACHE_DIR/golangci-lint-${VERSION}-${OS}-amd64/golangci-lint $install_path
	rm -rf $CACHE_DIR/golangci-lint-${VERSION}-${OS}-amd64
    fi
}

package-deps-unchanged?() {
    # Ensure cache directory exists
    if [[ ! -d .cache ]]; then
	rm -rf .cache && mkdir -p .cache
    fi

    # Creates a list of all dependencies
    rm -f .cache/package-deps.cur
    for arg in "$@" ; do
	if [[ -f $arg ]]; then
	    go list -deps -f "{{.ImportPath}} {{.Imports}}" "./$arg" >> .cache/package-deps.cur
	elif [[ -d $arg ]]; then
	    go list -deps -f "{{.ImportPath}} {{.Imports}}" "./$arg/..." >> .cache/package-deps.cur
	else
	    echo "Dependency target is not a file or directory: $arg"
	    return 2
	fi
    done

    # Only update the ./cache/package-deps file if deps have changed
    if [[ -f .cache/package-deps ]]; then
	if ! diff -q .cache/package-deps.cur .cache/package-deps 2>&1; then
	    echo "Package Deps Updated"
	    mv .cache/package-deps.cur .cache/package-deps
	    return 1
	else
	    echo "Package Deps Unchanged"
	    return 0
	fi
    else
	echo "Package Deps Created"
	mv .cache/package-deps.cur .cache/package-deps
	return 1
    fi
}

#####################################################################
# Run the main program
main "$@"
