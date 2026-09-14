ARG ARCH

FROM ${ARCH}node:22-bookworm-slim AS node
FROM ${ARCH}golang:1.25-bookworm AS golang
FROM ${ARCH}ossrs/srs:7 AS srs
FROM ${ARCH}ossrs/srs:ubuntu20 AS tools

RUN rm -rf /usr/local/srs/objs/nginx/html/console \
    /usr/local/srs/objs/nginx/html/players

FROM ${ARCH}goacme/lego AS lego

FROM ${ARCH}ubuntu:jammy AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
ARG MAKEARGS
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH, MAKEARGS: $MAKEARGS"

# Use a supported build toolchain while keeping the produced binaries compatible
# with the Ubuntu Jammy runtime image.
ENV PATH="/usr/local/go/bin:${PATH}"
COPY --from=golang /usr/local/go /usr/local/go
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends ca-certificates gcc g++ git make && \
    rm -rf /var/lib/apt/lists/*

# For ui build.
COPY --from=node /usr/local/bin /usr/local/bin
COPY --from=node /usr/local/lib /usr/local/lib
# For SRS server, always use the latest release version.
COPY --from=srs /usr/local/srs /usr/local/srs

ADD releases /g/releases
ADD mgmt /g/mgmt
ADD platform /g/platform
ADD ui /g/ui
ADD usr /g/usr
ADD test /g/test
ADD Makefile /g/Makefile

# For node to use more memory to fix: JavaScript heap out of memory
ENV NODE_OPTIONS="--max-old-space-size=4096"

# By default, make all, including platform and ui, but it will take a long time,
# so there is a MAKEARGS to build without UI, see platform.yml.
WORKDIR /g
# We define SRS_NO_LINT to disable the lint check.
RUN export SRS_NO_LINT=1 && \
    make clean && make -j ${MAKEARGS} && make install

# For youtube-dl, see https://github.com/ytdl-org/ytdl-nightly
FROM ${ARCH}python:3.11-slim-bookworm AS ytdl

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends binutils curl unzip && \
    pip install --no-cache-dir pyinstaller && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /g
RUN curl -O -L https://github.com/ytdl-org/youtube-dl/archive/refs/heads/master.zip && \
    unzip -q master.zip && cd youtube-dl-master && \
    pyinstaller --onefile --clean --noconfirm --name youtube-dl youtube_dl/__main__.py && \
    cp dist/youtube-dl /usr/local/bin/ && \
    ldd /usr/local/bin/youtube-dl

# Keep the runtime new enough for binaries produced by the Bookworm builders.
FROM ${ARCH}ubuntu:jammy AS dist

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends ca-certificates fontconfig fonts-noto-cjk \
        fonts-roboto-unhinted libbz2-1.0 liblzma5 libstdc++6 redis-server zlib1g && \
    rm -rf /var/lib/apt/lists/*

ENV PORT=":2024" NODE_ENV=production CLOUD=DOCKER PLATFORM_DOCKER=on \
    FONTCONFIG_PATH=/etc/fonts
WORKDIR /usr/local/oryx/platform

# Expose ports @see https://github.com/ossrs/oryx/blob/main/DEVELOPER.md#docker-allocated-ports
EXPOSE 2022 2443 1935 8080 5060 9000 8000/udp 10080/udp

# Copy files from build.
COPY --from=build /usr/local/oryx /usr/local/oryx
COPY --from=build /usr/local/srs /usr/local/srs
COPY --from=lego /lego /usr/local/bin/lego
COPY --from=srs /usr/local/srs/objs/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=tools /usr/local/bin/ffprobe /usr/local/bin/ffprobe
COPY --from=ytdl /usr/local/bin/youtube-dl /usr/local/bin/

# Prepare data directory.
RUN ln -sf /usr/local/bin/ffmpeg /usr/local/srs/objs/ffmpeg/bin/ffmpeg && \
    mkdir -p /data && \
    cd /usr/local/oryx/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
