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

echo "done testing!"

exit 0


