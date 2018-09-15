#!/bin/bash

kill $(ps -ax | grep server | awk '{print $1}')
