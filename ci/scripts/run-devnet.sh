#!/bin/bash
set -e

echo "Run Anvil"
cd rollups-node
make devnet
make run-devnet
cd -
