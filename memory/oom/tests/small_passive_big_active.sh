#!/bin/sh

just build-docker
just start small 1 1024 10000
just start big-1 10 102400 100
just start big-2 10 102400 100

ITERS=10

for i in $(seq 1 "$ITERS");
do
    just usage
    sleep 1
done

just monitor