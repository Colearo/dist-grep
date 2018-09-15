#!/bin/bash

for val in {1..9}
do
    echo VM$val
    ssh kechenl3@fa18-cs425-g29-0$val.cs.illinois.edu "cd ~/go/src/dist-grep; git pull; exit"
done
ssh kechenl3@fa18-cs425-g29-10.cs.illinois.edu "cd ~/go/src/dist-grep; git pull; exit"
echo 'Git Updated!'

