#!/usr/bin/env bash

set -e
export DEVELOPER_NODE=1
export RETAILCRM_URL=https://test.retailcrm.pro
export RETAILCRM_KEY=key
touch coverage.txt

go test ./errs/ -race -coverprofile=errs.out -covermode=atomic "$d"
go test ./v5/ -race -coverprofile=v5.out -covermode=atomic "$d"

cat errs.out >> coverage.txt
cat v5.out >> coverage.txt

rm errs.out
rm v5.out
