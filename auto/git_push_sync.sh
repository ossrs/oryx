#!/bin/bash

for file in $(git remote -v |grep -v tmp |grep push |awk '{print $1}'); do
    echo ""; 
    echo "git push $file $@"; 
    git push $file $@; 
done

