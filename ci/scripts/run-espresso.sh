#!/bin/bash
set -e

# docker run --network=host -e ESPRESSO_SEQUENCER_L1_PROVIDER=http://localhost:8545/ -e ESPRESSO_SEQUENCER_API_PORT=24000 ghcr.io/espressosystems/espresso-sequencer/espresso-dev-node:main

docker run -d --network=host -e ESPRESSO_SEQUENCER_L1_PROVIDER=http://localhost:8545/ -e ESPRESSO_SEQUENCER_API_PORT=24000 ghcr.io/espressosystems/espresso-sequencer/espresso-dev-node:20241120-patch4
