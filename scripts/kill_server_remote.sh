#!/bin/bash

echo Which VM Sever would you kill"(01-10)"?
read vmnum
ssh kechenl3@fa18-cs425-g29-$vmnum.cs.illinois.edu "~/go/src/dist-grep/scripts/kill_server_local.sh"

