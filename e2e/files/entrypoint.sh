#!/bin/bash

for i in `seq 1 2000000`; do echo "$i.example.org A" >> n.txt ;done

dnsperf -d ./n.txt -s 10.0.0.10 -c 30 -Q 60 -v
