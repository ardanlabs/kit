#!/bin/bash

# list all the packges, trim out the vendor directory and any main packages,
# then strip off the package name
PACKAGES="$(go list -f "{{.Name}}:{{.ImportPath}}" ./... | grep -v -E "main:|vendor/|examples" | cut -d ":" -f 2)"

# loop over all packages generating all their documentation
for PACKAGE in $PACKAGES
do

  echo "godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md"

  godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md

done
