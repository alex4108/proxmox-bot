#!/usr/bin/env bash

set -ex

# Define .secrets like
# PROXMOX_BOT_DISCORD_TOKEN=val
# Then for every PVE host...
# HOST1_USERNAME=foo@pam
# HOST1_PASSWORD=bar

for line in $(cat .secrets); do
    export $line
done

make clean
make build

printenv
echo $PVE1_PASSWORD
echo $PROXMOX_BOT_DISCORD_TOKEN
export CONFIG_PATH=$(pwd)/.secret-config.yml
./bin/proxmox-bot

