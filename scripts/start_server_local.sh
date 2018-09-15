#!/bin/bash

cd `dirname $0`
./build_all.sh
cd ../server
./server
cd -

