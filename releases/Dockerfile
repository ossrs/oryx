ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ADD releases /g/releases
WORKDIR /g/releases
RUN make

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/srs-cloud:focal-1 AS dist

COPY --from=build /g/releases/releases /usr/local/srs-cloud/releases/releases

ENV PORT=":9000"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/releases
CMD ["./releases"]
