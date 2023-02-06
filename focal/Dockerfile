ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

# http://releases.ubuntu.com/focal/
FROM ${ARCH}ubuntu:focal AS dist

COPY --from=build /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build /usr/local/bin/ffprobe /usr/local/bin/ffprobe
# For HTTPS client, or failed as: x509: certificate signed by unknown authority
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates

ENV PORT=":2024"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/platform
