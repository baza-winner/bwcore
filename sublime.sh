#!/bin/sh
for i in jld jlf; do
  for j in sublime-syntax tmPreferences; do
    targetFileSpec="$HOME/Library/Application Support/Sublime Text 3/Packages/User/$i.$j"
    rm -f "$targetFileSpec"
    ln "./$i.$j" "$targetFileSpec"
  done
done
