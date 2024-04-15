#!/bin/bash

# Print usage
usage() {
    echo "Usage: build-release.sh -v <version number>"
    echo "This script builds the Smartnode builder image used to build the daemon binaries."
    exit 0
}

# =================
# === Main Body ===
# =================

# Get the version
while getopts "admv:" FLAG; do
    case "$FLAG" in
        v) VERSION="$OPTARG" ;;
        *) usage ;;
    esac
done
if [ -z "$VERSION" ]; then
    usage
fi

echo -n "Building Docker image... "
docker build -t t0mk/rocketreport-builder:$VERSION -f docker/builder .
#docker tag t0mk/rocketreport-builder:$VERSION t0mk/rocketreport-builder:latest
#docker push t0mk/rocketreport-builder:$VERSION
#docker push t0mk/rocketreport-builder:latest
echo "done!"
