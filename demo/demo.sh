#!/usr/bin/env bash

# Clean up from the prior run
SILENCED=`killall -9 meraki-cli 2>&1 > /dev/null`

# Do some checks
if [[ -z ${MERAKI_AUTH_TOKEN:-""} ]]; then
    echo "Please set environment variable for MERAKI_AUTH_TOKEN"
    exit 0
fi

tmux new-session -n goapi -s test "tmux set-option -t test status off; yes | MERAKI_AUTH_TOKEN=${MERAKI_AUTH_TOKEN} ./demo/demo_client.sh; tmux kill-session -t test"
