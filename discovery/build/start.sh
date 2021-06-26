#!/bin/sh
set -e 

cd /home

i=1
while [ "$i" -le "200" ]
do
    ./discovery &
done

while true; do
    sleep 60
done

