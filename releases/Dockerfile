ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ADD releases /g/releases
WORKDIR /g/releases
RUN make

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ubuntu:focal AS dist

COPY --from=build /g/releases/releases /usr/local/oryx/releases/releases

ENV PORT=":9000"
WORKDIR /usr/local/oryx/releases
CMD ["./releases"]
