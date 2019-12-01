#!/bin/sh

path="tmp"

if [ -d $path ]
then
        rm -rf $path
fi
mkdir $path

nohup sh ./run.sh > ./tmp/run.log.out 2>&1 &

tail -f ./tmp/run.log.out