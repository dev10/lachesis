#!/bin/bash

for i in `seq 1 20`
do
    go test -run TestCatchUp  > ~/gossip.logs
    if grep "FAIL" ~/gossip.logs 
    then
        echo 'CHECK LOGS'
        exit
    else
        echo $i "OK"
    fi
done
echo 'NO ERRORS'