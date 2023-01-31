ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

ADD mgmt /g/mgmt
WORKDIR /g/mgmt

RUN make alpine

# To use if in RUN, see https://github.com/moby/moby/issues/7281#issuecomment-389440503
# Note that only exists issue like "/bin/sh: 1: [[: not found" for Ubuntu20, no such problem in CentOS7.
#SHELL ["/bin/bash", "-c"]
#RUN if [[ $TARGETPLATFORM != 'linux/arm/v7' && $TARGETPLATFORM != 'linux/arm64/v8' && $TARGETPLATFORM != 'linux/arm64' ]]; then \
#        make alpine; \
#    else \
#        make alpine-$(echo $TARGETPLATFORM |sed 's|/|-|g'); \
#    fi

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

