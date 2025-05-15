#!/bin/sh

exec anvil --base-fee 0 --block-time 1 --host 0.0.0.0 --chain-id 13370 \
    --load-state /usr/share/devnet/anvil_state.json \
    --dump-state /dump/anvil_state.json
