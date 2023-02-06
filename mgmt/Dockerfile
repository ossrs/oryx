ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

ADD mgmt /g/mgmt
WORKDIR /g/mgmt

RUN make alpine

FROM ${ARCH}alpine:3.16 AS dist

COPY --from=build /g/mgmt /usr/local/srs-cloud/mgmt

# See https://stackoverflow.com/a/40944512/17679565
RUN apk add --no-cache bash

ENV PORT=":2022"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/mgmt
# TODO: FIXME: Implements it.
#CMD ["./bootstrap"]
CMD ["bash"]

