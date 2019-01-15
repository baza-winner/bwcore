#!/bin/sh
for i in jlf.sublime-syntax jlf.tmPreferences syntax_test_jlf; do
  targetFileSpec="$HOME/Library/Application Support/Sublime Text 3/Packages/User/$i"
  rm -f "$targetFileSpec"
done
for i in edf.sublime-syntax edf.tmPreferences; do
  targetFileSpec="$HOME/Library/Application Support/Sublime Text 3/Packages/User/$i"
  rm -f "$targetFileSpec"
  ln "./$i" "$targetFileSpec"
done

