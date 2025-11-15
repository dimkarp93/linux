#!/bin/sh

just build-docker
just start small 20 256 10
just start big-1 3 10240000 500
just start big-2 3 10240000 500

ITERS=10

for i in $(seq 1 "$ITERS");
do
    just usage
    sleep 1
done

just monitor