#!/bin/bash

cd `dirname $0`
./build_all.sh
cd ../src/grep-server
./server
cd -

