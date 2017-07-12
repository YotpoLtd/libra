#!/bin/bash

BASE_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD_DIR=${BASE_DIR}/tmp

echo "$BUILD_DIR"

mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

: ${CI_VERSION_TAG:=$SHORT_GIT_COMMIT}
: ${CI_VERSION_TAG:=CIVERSION}

if [ "$PREBUILD_SCRIPT" = "" ]; then
	pushd ../..
	docker build -t docker.uacf.io/infra/libra/build:${CI_VERSION_TAG} -f build/Dockerfile.build .
	popd
fi

docker run --rm docker.uacf.io/infra/libra/build:${CI_VERSION_TAG} bash -c "cd /go/bin; tar -cv ." | tar -x