#!/usr/bin/env bash

# Load the common demo lib
demo_dir="$(dirname "$0")"
. "$demo_dir/demo_magic_common.sh"
. "$demo_dir/common.sh"

#####################################################################
# Start the good stuff

SLEEPTIME=2

pe "# Show CLI usage."
pe "${CLI} --help"
sleep ${SLEEPTIME}

pe "# Show CLI version."
pe "${CLI} version"
sleep ${SLEEPTIME}

pe "# Export your MERAKI_AUTH_TOKEN as an environment variable."
p  "export MERAKI_AUTH_TOKEN=<secret_token>"
sleep ${SLEEPTIME}

pe "# List your organizations."
pe "${CLI} organizations list"
sleep ${SLEEPTIME}

ORGANIZATION=$(${CLI} organizations list | jq -r .[0].id)
pe  "# Subsequent calls will use the organization ${ORGANIZATION}"
sleep ${SLEEPTIME}

pe "# List the devices within the organization."
pe "${CLI} devices list --organization-id=${ORGANIZATION}"
sleep ${SLEEPTIME}

pe "# List the networks within the organization."
pe "${CLI} networks list --organization-id=${ORGANIZATION}"
sleep ${SLEEPTIME}

NETWORK=$(${CLI} networks list --organization-id=${ORGANIZATION} | jq -r .[0].id)
pe  "# Subsequent calls will use the network ${NETWORK}"
sleep ${SLEEPTIME}

pe "# List the wireless SSIDs within the network."
pe "${CLI} networks ssid list --network-id=${NETWORK}"
sleep ${SLEEPTIME}

pe "# List the events of type 'appliance' within the network."
pe "${CLI} events list --network-id=${NETWORK} --product-type=appliance"
sleep ${SLEEPTIME}

pe "# That's it for the demo.  Thanks for watching!"
sleep ${SLEEPTIME}

clear
