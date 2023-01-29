ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ADD . /g
WORKDIR /g/releases
RUN make alpine

FROM ${ARCH}alpine:3.16 AS dist

COPY --from=build /g/releases/releases /usr/local/srs-cloud/releases/releases

ENV PORT=":9000"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/releases
CMD ["./releases"]

