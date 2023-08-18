#!/bin/bash

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <number_of_retries> <command> [args...]"
    exit 1
fi

num_retries=$1; shift
command="$@"
echo "Running with retries=$num_retries command: $command"

retry_count=0
while [ $retry_count -lt $num_retries ]; do
    $command
    if [ $? -eq 0 ]; then exit 0; fi

    retry_count=$((retry_count + 1))
    echo "Command failed. Retrying ($retry_count/$num_retries)..."
done

echo "Command failed after $num_retries retries."
exit 1
