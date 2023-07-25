ARG ARCH

FROM ${ARCH}ossrs/node:18 AS node

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

# For ui build.
COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /usr/local/lib /usr/local/lib

ADD releases /g/releases
ADD mgmt /g/mgmt
ADD platform /g/platform
ADD ui /g/ui

# Note that we only build the platform without ui, because already build ui for all OS.
# See platform.yml command:
#     make ui-build-cn && make ui-build-en
WORKDIR /g
RUN make clean && make -j && make install

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/srs-cloud:focal-1 AS dist

COPY --from=build /usr/local/srs-cloud /usr/local/srs-cloud

# Prepare data directory.
RUN mkdir -p /data && \
    cd /usr/local/srs-cloud/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
