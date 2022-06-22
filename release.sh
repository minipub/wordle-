#!/bin/bash

## Auto Tag & Release

if [ $# -ne 1 ]
then
	echo "Usage: bash release.sh <version>"
	exit 1
fi

version=$1
chksum=wordle-bundle_${version}_checksums.txt

gh release create ${version} -t ${version} -F - ./wordle-*_${version}_*.tar.gz ./${chksum}
