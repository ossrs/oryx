ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build
FROM ${ARCH}redis:5.0 AS redis
FROM ${ARCH}goacme/lego AS lego


FROM ${ARCH}ubuntu:focal AS files

 # certs: For HTTPS client, or failed as: x509: certificate signed by unknown authority
 # redis: For redis server and tools.
 # lego: For ACME client, request and renew the HTTPS certificate.
COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /usr/local/bin/
COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /usr/local/bin/
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=redis /usr/local/bin/redis-cli /usr/local/bin/redis-server /usr/local/bin/
COPY --from=lego /lego /usr/local/bin/

# http://releases.ubuntu.com/focal/
FROM ${ARCH}ubuntu:focal AS dist

COPY --from=files /usr/local /usr/local
COPY --from=files /usr/share /usr/share
COPY --from=files /etc/ssl /etc/ssl

# We write the common config for platform to base image, however you can
# overwrite it if not for platform, for example, used in releases.
ENV PORT=":2024" NODE_ENV=production SRS_DOCKERIZED=true MGMT_DOCKER=true PLATFORM_DOCKER=true CLOUD=DOCKER

# Setup the work directory to platform.
WORKDIR /usr/local/srs-cloud/platform
