#!/bin/bash

kill $(ps -ax | grep grep-server | awk '{print $1}')
