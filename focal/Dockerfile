ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build
FROM ${ARCH}redis:5.0 AS redis

FROM ${ARCH}goacme/lego AS lego

# https://github.com/upx/upx
FROM ${ARCH}ubuntu:focal AS upx

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

# Use UPX to compress the binary.
# https://serverfault.com/questions/949991/how-to-install-tzdata-on-a-ubuntu-docker-image
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y upx

COPY --from=lego /lego /usr/local/bin/
COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /usr/local/bin/
COPY --from=redis /usr/local/bin/redis-cli /usr/local/bin/redis-server /usr/local/bin/

RUN echo "Before UPX for $TARGETARCH" && \
    ls -lh /usr/local/bin/* && \
    upx --best --lzma /usr/local/bin/lego && \
    upx --best --lzma /usr/local/bin/ffmpeg && \
    upx --best --lzma /usr/local/bin/ffprobe && \
    upx --best --lzma /usr/local/bin/redis-cli && \
    upx --best --lzma /usr/local/bin/redis-server && \
    echo "After UPX for $TARGETARCH" && \
    ls -lh /usr/local/bin/*

# http://releases.ubuntu.com/focal/
FROM ${ARCH}ubuntu:focal AS dist

 # certs: For HTTPS client, or failed as: x509: certificate signed by unknown authority
 # redis: For redis server and tools.
 # lego: For ACME client, request and renew the HTTPS certificate.
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=upx /usr/local/bin/lego \
    /usr/local/bin/ffmpeg /usr/local/bin/ffprobe \
    /usr/local/bin/redis-cli /usr/local/bin/redis-server \
    /usr/local/bin/

# We write the common config for platform to base image, however you can
# overwrite it if not for platform, for example, used in releases.
ENV PORT=":2024" NODE_ENV=production CLOUD=DOCKER PLATFORM_DOCKER=on

# Setup the work directory to platform.
WORKDIR /usr/local/srs-stack/platform
