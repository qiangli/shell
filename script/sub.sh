#!/bin/bash

set -xue

pwd

echo "------"
echo "sub script called!"
echo "------"

ai @ --max-history=0 --max-span=0 \
    --max-time-10 --max-turns=3 \
    --log-level quiet \
    --models default/L1 \
    --instruction "you are a cool agent"\
    --message "anonymous agent"


exit 0
