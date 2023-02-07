ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

# http://releases.ubuntu.com/focal/
FROM ${ARCH}ubuntu:focal AS dist

COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build /usr/local/bin/ffprobe /usr/local/bin/ffprobe
# For HTTPS client, or failed as: x509: certificate signed by unknown authority
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates

# We write the common config for platform to base image, however you can
# overwrite it if not for platform, for example, used in releases.
ENV PORT=":2024"
# Create link for compatible reason.
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
# Create directory for compatible reason.
RUN mkdir -p /usr/local/lighthouse/softwares
RUN ln -sf /usr/local/srs-cloud /usr/local/lighthouse/softwares/srs-cloud
# Setup the work directory to platform.
WORKDIR /usr/local/srs-cloud/platform
