ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

FROM ${ARCH}redis:5.0 AS redis

FROM ${ARCH}ossrs/srs:4 AS srs

# http://releases.ubuntu.com/focal/
FROM ${ARCH}ubuntu:focal AS dist

# https://serverfault.com/questions/949991/how-to-install-tzdata-on-a-ubuntu-docker-image
ENV DEBIAN_FRONTEND=noninteractive
# See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get
RUN apt update && apt-get install -y certbot && rm -rf /var/lib/apt/lists/*

COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build /usr/local/bin/ffprobe /usr/local/bin/ffprobe
# For HTTPS client, or failed as: x509: certificate signed by unknown authority
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates
# For redis server and tools.
COPY --from=redis /usr/local/bin/redis-cli /usr/local/bin/redis-cli
COPY --from=redis /usr/local/bin/redis-server /usr/local/bin/redis-server
# For SRS server.
COPY --from=srs /usr/local/srs /usr/local/srs

# We write the common config for platform to base image, however you can
# overwrite it if not for platform, for example, used in releases.
ENV PORT=":2024"

# Setup the default node env, for both Go and nodejs env.
ENV NODE_ENV=production

# Run application in docker.
ENV SRS_DOCKERIZED=true
# Run mgmt in docker by default.
ENV MGMT_DOCKER=true
# Run platform in docker by default.
ENV PLATFORM_DOCKER=true
# Setup the .env to docker.
ENV CLOUD=DOCKER

# Create link for compatible reason.
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform

# Create directory for compatible reason.
# TODO: FIXME: Remove it when upgrade to new versions.
RUN mkdir -p /usr/local/lighthouse/softwares
RUN ln -sf /usr/local/srs-cloud /usr/local/lighthouse/softwares/srs-cloud

# Setup the work directory to platform.
WORKDIR /usr/local/srs-cloud/platform
