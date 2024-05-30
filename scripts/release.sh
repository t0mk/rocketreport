#!/bin/sh

set -x

VERSION=$1
    
[ -z "$VERSION" ] && echo "Usage: $0 <version>" && exit 1

# VERSION shouldnt have v prefix

VERSION=${VERSION#v}

make static-builds

gh release create v${VERSION} 'rocketreport-amd64' || exit 1
