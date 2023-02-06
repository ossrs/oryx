ARG ARCH

FROM ${ARCH}alpine:3.16 AS dist

# See https://stackoverflow.com/a/40944512/17679565
RUN apk add bash

ENV PORT=":2022"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/mgmt
