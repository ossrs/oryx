ARG ARCH

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ADD . /g
WORKDIR /g/releases

# Build for alpine, see https://www.cloudbees.com/blog/building-minimal-docker-containers-for-go-applications
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo  -o releases .

FROM ${ARCH}alpine:3.16 AS dist

COPY --from=build /g/releases/releases /usr/local/srs-cloud/releases/releases

ENV PORT=":9000"
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/releases
CMD ["./releases"]

