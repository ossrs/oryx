ARG ARCH

FROM ${ARCH}alpine:3.16 AS dist

# See https://stackoverflow.com/a/40944512/17679565
RUN apk add bash

# Note that we do not install FFmpeg, because it make the alpine too large,
# so that we encourage to use Ubuntu focal instead.
#
# We write the common config for mgmt to base image, however you can
# overwrite it if not for mgmt.
ENV PORT=":2022"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/mgmt
