#!/bin/sh
script_file=$(readlink -e "$0")
script_dir=$(dirname ${script_file})
pushd ${script_dir}
    go build -o . ./..
popd
