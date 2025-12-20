#!/bin/bash
set -x 

echo "testing..."

cd tool/; echo $?
pwd

/bin/ls .

cd /tmp; echo $?
pwd

@agent --max-history 0 --max-span 0 what is weather today? --adapter echo
ai @ed --model default/any --message "tell me a joke" --adapter echo
ai /kit:func --arg name=value --arg name=v1,v2,v3 --arguments "{name:value,}" --adapter echo

# TODO
# bash $PWD/script/sub.sh

# sleep 1

time sleep 2

echo "done testing!"

exit 0


