ARG ARCH

FROM ${ARCH}ossrs/node:18 AS node
FROM ${ARCH}ossrs/srs:5 AS srs

RUN mv /usr/local/srs/objs/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg && \
    ln -sf /usr/local/bin/ffmpeg /usr/local/srs/objs/ffmpeg/bin/ffmpeg

RUN rm -rf /usr/local/srs/objs/nginx/html/console \
    /usr/local/srs/objs/nginx/html/players

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
# We define SRS_NO_LINT to disable the lint check.
RUN export SRS_NO_LINT=1 && \
    make clean && make -j ${MAKEARGS} && make install

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/srs-stack:focal-1 AS dist

# Expose ports @see https://github.com/ossrs/srs-stack/blob/main/DEVELOPER.md#docker-allocated-ports
EXPOSE 2022 2443 1935 8080 5060 9000 8000/udp 10080/udp

# For srs-stack, build it.
COPY --from=build /usr/local/srs-stack /usr/local/srs-stack
# For SRS server, always use the latest release version.
COPY --from=srs /usr/local/srs /usr/local/srs

# Prepare data directory.
RUN mkdir -p /data && \
    cd /usr/local/srs-stack/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
