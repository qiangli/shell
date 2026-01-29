#!/bin/bash

set -x

echo "city $city"

/agent:search --message "weather in ${city} city"

@agent "tell me a joke"

ai @ask "what is fish"

echo "done"
