#!/bin/bash

cd `dirname $0`
./build_all.sh
cd ../server
nohup ./server > /dev/null 2>&1 &
cd -

