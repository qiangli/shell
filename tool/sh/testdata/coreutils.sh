#!/bin/bash
set -x -e
date

DIR=testdata/.local
mkdir -p $DIR

touch $DIR/code_dict
echo /agent:gpte/gen_code hello there > $DIR/code_dict
touch $DIR/entry_dict
echo /agent:gpte/gen_entrypoint main > $DIR/entry_dict

# cd $DIR
# pwd

cat $DIR/code_dict $DIR/entry_dict  > $DIR/file_dict
head $DIR/file_dict 
shasum $DIR/file_dict 

exit 0