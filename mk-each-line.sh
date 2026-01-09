#!/bin/bash

seq=1
for i in $( < $1 ) ; do
	~/bin/qrcode -o out__$seq.png "$i"
	seq=$(($seq+1))
done

