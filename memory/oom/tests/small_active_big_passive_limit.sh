#!/bin/sh

just build-docker
just start big-1 1 1024000 75 1
just start big-2 1 1024000 75 1
just start small 100 10240 75

ITERS=10

for i in $(seq 1 "$ITERS");
do
    just usage
    sleep 1
done

just monitor