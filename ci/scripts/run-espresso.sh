#!/bin/bash
set -e

# docker run --network=host -e ESPRESSO_SEQUENCER_L1_PROVIDER=http://localhost:8545/ -e ESPRESSO_SEQUENCER_API_PORT=24000 ghcr.io/espressosystems/espresso-sequencer/espresso-dev-node:main

# docker run --network=host -e ESPRESSO_SEQUENCER_L1_PROVIDER=http://localhost:8545/ -e ESPRESSO_SEQUENCER_API_PORT=24000 ghcr.io/espressosystems/espresso-sequencer/espresso-dev-node:main
# docker network create my_bridge
docker rm espresso-dev-node || true
docker run --network my_bridge --rm --name devnet -p 8545:8545 -d cartesi/rollups-node-devnet:devel
sleep 10
docker run --network my_bridge -e ESPRESSO_SEQUENCER_L1_PROVIDER=http://devnet:8545/ -e ESPRESSO_SEQUENCER_API_PORT=24000 -p 24000:24000 --name espresso-dev-node ghcr.io/espressosystems/espresso-sequencer/espresso-dev-node:20241120-patch5