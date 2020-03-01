#!/bin/bash -eu

tools=($(sed -En 's/[[:space:]]+_ "(.*)"/\1/p' ./tools/tools.go))

echo "install tools"
for tool in ${tools[@]}
do
    echo " - $tool"
    go install -mod="" "$tool"
done
