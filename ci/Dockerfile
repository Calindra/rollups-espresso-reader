# syntax=docker.io/docker/dockerfile:1
# (c) Cartesi and individual authors (see AUTHORS)
# SPDX-License-Identifier: Apache-2.0 (see LICENSE)

ARG EMULATOR_VERSION=0.18.1

# Build directories.
ARG GO_BUILD_PATH=/build/cartesi/go

FROM golang:1.23.2 AS go-builder

ARG GO_BUILD_PATH

WORKDIR ${GO_BUILD_PATH}
COPY . .

# Build espresso reader.
RUN go build -o cartesi-rollups-espresso-reader

# =============================================================================
# STAGE: rollups-node
#
# This stage prepares the final Docker image that will be used in the production
# environment. It installs in /usr/bin all the binaries necessary to run the
# node.
#
# (This stage copies the binaries from previous stages.)
# =============================================================================

FROM debian:bookworm-20250113-slim AS rollups-node

ARG NODE_RUNTIME_DIR=/var/lib/cartesi-rollups-node

USER root

# Download system dependencies required at runtime.
ARG DEBIAN_FRONTEND=noninteractive
RUN <<EOF
    set -e
    apt-get update
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        libboost-all-dev \
        libslirp0 \
        lua5.4 \
        procps \
        tini
    rm -rf /var/lib/apt/lists/*
    mkdir -p ${NODE_RUNTIME_DIR}/snapshots ${NODE_RUNTIME_DIR}/data
EOF

# Copy Go binary.
ARG GO_BUILD_PATH

# Install cartesi machine
WORKDIR /tmp
ENV CARTESI_MACHINE_VERSION=0.19.0
ENV CARTESI_MACHINE_SUFFIX="-alpha3"
RUN <<EOF
    set -e
    TARGETARCH=$(dpkg --print-architecture)
    curl -LO https://github.com/cartesi/machine-emulator/releases/download/v${CARTESI_MACHINE_VERSION}${CARTESI_MACHINE_SUFFIX}/cartesi-machine-v${CARTESI_MACHINE_VERSION}_${TARGETARCH}.deb
    dpkg -i ./cartesi-machine-v${CARTESI_MACHINE_VERSION}_${TARGETARCH}.deb
    rm ./cartesi-machine-v${CARTESI_MACHINE_VERSION}_${TARGETARCH}.deb
EOF

# Copy kernel and root filesystem images with their SHA512 checksums
ADD --chmod=644 https://github.com/cartesi/image-kernel/releases/download/v0.20.0/linux-6.5.13-ctsi-1-v0.20.0.bin /tmp/artifacts/linux-6.5.13-ctsi-1-v0.20.0.bin
ADD --chmod=644 https://github.com/cartesi/image-kernel/releases/download/v0.20.0/linux-6.5.13-ctsi-1-v0.20.0.bin.sha512 /tmp/linux-6.5.13-ctsi-1-v0.20.0.bin.sha512
ADD --chmod=644 https://github.com/cartesi/machine-emulator-tools/releases/download/v0.16.1/rootfs-tools-v0.16.1.ext2 /tmp/rootfs-tools-v0.16.1.ext2
ADD --chmod=644 https://github.com/cartesi/machine-emulator-tools/releases/download/v0.16.1/rootfs-tools-v0.16.1.ext2.sha512 /tmp/rootfs-tools-v0.16.1.ext2.sha512

# Verify the downloaded files
WORKDIR /tmp
RUN <<EOF
    set -e
    sha512sum --check \
        rootfs-tools-v0.16.1.ext2.sha512 \
        linux-6.5.13-ctsi-1-v0.20.0.bin.sha512
EOF

# Copy kernel and root filesystem images to the appropriate directory
RUN <<EOF
    set -e
    mkdir -pv /usr/share/cartesi-machine/images
    mv -v /tmp/artifacts/linux-6.5.13-ctsi-1-v0.20.0.bin /usr/share/cartesi-machine/images/linux.bin
    mv -v /tmp/rootfs-tools-v0.16.1.ext2 /usr/share/cartesi-machine/images/rootfs.ext2
EOF

# Install rollups node
ARG ROLLUPS_NODE_VERSION
RUN <<EOF
    set -e
    TARGETARCH=$(dpkg --print-architecture)
    curl -LO https://github.com/cartesi/rollups-node/releases/download/v${ROLLUPS_NODE_VERSION}/cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
    dpkg -i ./cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
    rm ./cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
EOF

WORKDIR ${NODE_RUNTIME_DIR}

# HEALTHCHECK --interval=1s --timeout=1s --retries=5 \
    # CMD curl -G -f -H 'Content-Type: application/json' http://127.0.0.1:10000/healthz

RUN mkdir applications
RUN cartesi-machine --ram-length=128Mi --store=applications/echo-dapp --final-hash -- ioctl-echo-loop --vouchers=1 --notices=1 --reports=1 --verbose=1

COPY --from=go-builder ${GO_BUILD_PATH}/cartesi-rollups-* /usr/bin

COPY --chown=cartesi:cartesi ./ci/entrypoint.sh /entrypoint.sh
COPY --chown=cartesi:cartesi ./ci/scripts/migrate-db.sh /migrate-db.sh
COPY --chown=cartesi:cartesi ./ci/scripts/espresso-dev-node-health.sh /espresso-dev-node-health.sh
COPY --chown=cartesi:cartesi ./internal/repository/postgres/schema/migrations /internal/repository/postgres/schema/migrations
ADD https://raw.githubusercontent.com/cartesi/rollups-node/refs/tags/v2.0.0-alpha.4/internal/repository/postgres/schema/migrations/000001_create_initial_schema.up.sql /rollups-node/migrations/

RUN chmod +x /entrypoint.sh
CMD ["tini", "--", "/entrypoint.sh"]
