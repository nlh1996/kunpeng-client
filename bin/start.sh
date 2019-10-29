#!/bin/bash
nohup ./client -team=$1 -ip="$2" -port=$3 >../battle.out 2>&1 &
