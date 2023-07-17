ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

ADD platform /g/platform

# Note that we only build the platform without ui, because already build ui for all OS.
# See platform.yml command:
#     make ui-build-cn && make ui-build-en
WORKDIR /g/platform
RUN make clean && make
RUN cd /g/platform/ui && rm -rf js-core node_modules package* public src

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/srs-cloud:focal-1 AS dist

COPY --from=build /g/platform /usr/local/srs-cloud/platform

# Prepare data directory.
RUN mkdir -p /data && \
    cd /usr/local/srs-cloud/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
