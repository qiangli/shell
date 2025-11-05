#!/bin/bash
set -x -e
date

DIR=testdata/.local
mkdir -p $DIR

touch $DIR/code_dict
echo /agent:gpte/gen_code hello there > $DIR/code_dict
touch $DIR/entry_dict
echo /agent:gpte/gen_entrypoint main > $DIR/entry_dict

cd $DIR
pwd

cat code_dict entry_dict | tac > file_dict
cat file_dict | shasum

exit 0