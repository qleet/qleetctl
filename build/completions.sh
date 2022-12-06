#!/bin/sh
set -e
rm -rf completions
mkdir completions
pwd
ls -al
ls -al ../
for sh in bash zsh fish; do
  ./qleetctl completion "$sh" >"completions/qleetctl.$sh"
#	go run ./main.go completion "$sh" >"completions/qleetctl.$sh"
done
chmod +x ./completions/
