#!/bin/sh

version=$(go version | awk '{ print $3 }' | awk -F. '{ print $2 }')

if [ "$version" != "5" ]
then
  echo "Installing golint into your GOPATH..."
  go get golang.org/x/lint/...
  echo "Checking with golint..."
  golint ./...
fi
