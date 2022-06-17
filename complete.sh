#!/bin/sh

path=$(sergeant query paths | fzf)

printf "Start (HH:MM): "
read -r start

printf "End (HH:MM): "
read -r end

echo "Score:"
score=$(printf "perfect\nmajor\nminor" | fzf)

time=$(echo $(( ($(gdate --date "$(date +%Y-%m-%d) $end" +%s) - $(gdate --date "$(date +%Y-%m-%d) $start" +%s)) ))s)

~/go/bin/sergeant complete --path $path --user olly --time $time $score
