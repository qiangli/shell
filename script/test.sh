#!/bin/bash
set -x 

echo "testing..."

cd tool/; echo $?
pwd

/bin/ls .

cd /tmp; echo $?
pwd

@agent --max-history 0 --max-span 0 what is weather today?
ai @agent --models default/any --message "tell me a joke"
ai /kit:func --arg name=value --arg name=v1,v2,v3 --arguments "{name:value,}"

# TODO
# bash $PWD/script/sub.sh

echo "done testing!"

exit 0


