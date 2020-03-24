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


install-go-swagger() {
    install_path=$1
    mkdir -p $(dirname $install_path)

    VERSION="0.23.0"
    OS=$(uname -s  | tr '[:upper:]' '[:lower:]')
    FILE="swagger_${OS}_amd64"
    URL="https://github.com/go-swagger/go-swagger/releases/download/v${VERSION}/${FILE}"
    CACHE_FILE="${CACHE_DIR}/${FILE}"

    # Ensure the source file exists
    if [[ ! -f $CACHE_FILE ]]; then
	curl -fsSL $URL > $CACHE_FILE
    fi

    # Create the target
    if [[ ! -f $install_path ]]; then
	cp $CACHE_FILE $install_path
	chmod a+x $install_path
    fi
}

#####################################################################
# Run the main program
main "$@"
