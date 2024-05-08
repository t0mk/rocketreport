#!/bin/sh

set -x

VERSION=$1
    
[ -z "$VERSION" ] && echo "Usage: $0 <version>" && exit 1

# VERSION shouldnt have v prefix

VERSION=${VERSION#v}


make docker-image || exit 1

docker tag t0mk/rocketreport:latest t0mk/rocketreport:$VERSION || exit 1

docker push t0mk/rocketreport:$VERSION || exit 1
docker push t0mk/rocketreport:latest || exit 1

gh release create v${VERSION} 'rocketreport-amd64' || exit 1