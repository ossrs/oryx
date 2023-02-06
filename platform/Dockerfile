ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH"

ADD mgmt /g/mgmt
ADD platform /g/platform
WORKDIR /g/platform

RUN make alpine
RUN cd /g/platform/ui && rm -rf js-core node_modules package* public src

#FROM ${ARCH}alpine:3.16 AS dist
FROM ${ARCH}ossrs/srs-cloud:alpine-1 AS dist

COPY --from=build /g/mgmt /usr/local/srs-cloud/mgmt
COPY --from=build /g/platform /usr/local/srs-cloud/platform

## See https://stackoverflow.com/a/40944512/17679565
#RUN apk add bash
#RUN apk add ffmpeg

ENV PORT=":2024"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/platform
CMD ["./platform"]
