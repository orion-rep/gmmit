#!/usr/bin/env bash
set -o errexit

set -o nounset
set -o pipefail
if [[ "${TRACE-0}" == "1" ]]; then
    set -o xtrace
fi

main() {
    build darwin arm64
    build darwin amd64
    build linux 386
    build linux amd64
    build windows 386
    build windows amd64
}

function build() {
    OS=$1
    ARCH=$2
    echo "[INFO] Building for ${OS} (${ARCH})"
    env GOOS=${OS} GOARCH=${ARCH} go build -o build/${OS}-${ARCH}/gmmit ./cmd/gmmit/
    tar -czf build/${OS}-${ARCH}.tgz --directory=build/${OS}-${ARCH} gmmit
    rm -rf build/${OS}-${ARCH}
}

main "$@"