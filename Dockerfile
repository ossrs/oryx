ARG ARCH

FROM ${ARCH}ossrs/node:18 AS node
FROM ${ARCH}ossrs/srs:5 AS srs

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
ARG MAKEARGS
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH, MAKEARGS: $MAKEARGS"

# For ui build.
COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /usr/local/lib /usr/local/lib

ADD releases /g/releases
ADD mgmt /g/mgmt
ADD platform /g/platform
ADD ui /g/ui
ADD usr /g/usr
ADD test /g/test
ADD Makefile /g/Makefile

# By default, make all, including platform and ui, but it will take a long time,
# so there is a MAKEARGS to build without UI, see platform.yml.
WORKDIR /g
RUN make clean && make -j ${MAKEARGS} && make install

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/srs-cloud:focal-1 AS dist

# For SRS-Cloud, build it.
COPY --from=build /usr/local/srs-cloud /usr/local/srs-cloud
# For SRS server, always use the latest release version.
COPY --from=srs /usr/local/srs /usr/local/srs

# Prepare data directory.
RUN mkdir -p /data && \
    cd /usr/local/srs-cloud/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
