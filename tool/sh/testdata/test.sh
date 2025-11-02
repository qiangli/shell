#!/bin/bash

set -x

/agent:search --message "weather in sfo"

@agent "tell me a joke"

ai @ask "what is fish"

echo "done"
