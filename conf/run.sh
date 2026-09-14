#!/bin/bash
runapp="$1"
echo "准备重启 $runapp"
kill $runapp
sleep 2
nohup ./${runapp} > /dev/null 2>&1 &
