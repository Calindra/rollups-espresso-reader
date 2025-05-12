ARG FOUNDRY_VERSION=1.0.0
ARG CONTRACTS_VERSION=2.0.0-rc.17
ARG DEVNET_BUILD_PATH=/opt/cartesi/rollups-contracts

FROM ghcr.io/foundry-rs/foundry:${FOUNDRY_VERSION}

ARG FOUNDRY_VERSION
ARG CONTRACTS_VERSION
ARG DEVNET_BUILD_PATH

USER root

# Install ca-certificates, curl, and git (setup).
ENV DEBIAN_FRONTEND=noninteractive
RUN <<EOF
    set -e
    apt-get update
    apt-get install -y --no-install-recommends ca-certificates curl git jq
EOF

# Install anvil state file and devnet deployment info.
RUN <<EOF
    set -e
    mkdir -p ${DEVNET_BUILD_PATH}
    URL=https://github.com/cartesi/rollups-contracts/releases/download
    VERSION=v${CONTRACTS_VERSION}
    ARTIFACT=rollups-contracts-${CONTRACTS_VERSION}-anvil-v${FOUNDRY_VERSION}.tar.gz
    curl -sSL ${URL}/${VERSION}/${ARTIFACT} -o /tmp/contracts.tar.gz
    echo "3c19e44868db7b8bf2edd55a089cda39787dcfb2715d4a2d877644d68c40141b  /tmp/contracts.tar.gz" | sha256sum --check
    tar -zxf /tmp/contracts.tar.gz -C ${DEVNET_BUILD_PATH}/
    mkdir -p /usr/share/devnet
    cat ${DEVNET_BUILD_PATH}/deployments/*.json | jq -s 'map({ (.contractName): .address }) | add' > /usr/share/devnet/deployment.json
    mv ${DEVNET_BUILD_PATH}/state.json /usr/share/devnet/anvil_state.json
EOF

HEALTHCHECK --interval=1s --timeout=1s --retries=5 \
	CMD curl \
	-X \
	POST \
	-s \
	-H 'Content-Type: application/json' \
	-d '{"jsonrpc":"2.0","id":"1","method":"net_listening","params":[]}' \
	http://127.0.0.1:8545

CMD ["anvil", "--base-fee", "0", "--block-time", "1", "--silent", "--host", "0.0.0.0", "--chain-id", "13370", "--load-state", "/usr/share/devnet/anvil_state.json"]
