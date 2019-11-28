#!/bin/sh

mkdir tmp
mkdir log.out
nohup sh ./run.sh >> ./tmp/run.log.out 2>&1 &
