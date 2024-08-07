ARG ARCH

FROM ${ARCH}ossrs/node:18 AS node
FROM ${ARCH}ossrs/srs:5 AS srs

RUN mv /usr/local/srs/objs/ffmpeg/bin/ffmpeg /usr/local/bin/ffmpeg && \
    ln -sf /usr/local/bin/ffmpeg /usr/local/srs/objs/ffmpeg/bin/ffmpeg

RUN rm -rf /usr/local/srs/objs/nginx/html/console \
    /usr/local/srs/objs/nginx/html/players

FROM ${ARCH}ossrs/srs:ubuntu20 AS build

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
ARG MAKEARGS
RUN echo "BUILDPLATFORM: $BUILDPLATFORM, TARGETPLATFORM: $TARGETPLATFORM, TARGETARCH: $TARGETARCH, MAKEARGS: $MAKEARGS"

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

# Use UPX to compress the binary.
# https://serverfault.com/questions/949991/how-to-install-tzdata-on-a-ubuntu-docker-image
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y upx

RUN echo "Before UPX for $TARGETARCH" && \
    ls -lh /usr/local/srs/objs/srs /usr/local/oryx/platform/platform && \
    upx --best --lzma /usr/local/srs/objs/srs && \
    upx --best --lzma /usr/local/oryx/platform/platform && \
    echo "After UPX for $TARGETARCH" && \
    ls -lh /usr/local/srs/objs/srs /usr/local/oryx/platform/platform

# For youtube-dl, see https://github.com/ytdl-org/ytdl-nightly
FROM ${ARCH}python:3.9-slim-bullseye AS ytdl

RUN apt-get update -y && apt-get install -y binutils curl unzip && \
    pip install pyinstaller

WORKDIR /g
RUN curl -O -L https://github.com/ytdl-org/youtube-dl/archive/refs/heads/master.zip && \
    unzip -q master.zip && cd youtube-dl-master && \
    pyinstaller --onefile --clean --noconfirm --name youtube-dl youtube_dl/__main__.py && \
    cp dist/youtube-dl /usr/local/bin/ && \
    ldd /usr/local/bin/youtube-dl

# http://releases.ubuntu.com/focal/
#FROM ${ARCH}ubuntu:focal AS dist
FROM ${ARCH}ossrs/oryx:focal-1 AS dist

# Expose ports @see https://github.com/ossrs/oryx/blob/main/DEVELOPER.md#docker-allocated-ports
EXPOSE 2022 2443 1935 8080 5060 9000 8000/udp 10080/udp

# Copy files from build.
COPY --from=build /usr/local/oryx /usr/local/oryx
COPY --from=build /usr/local/srs /usr/local/srs
COPY --from=ytdl /usr/local/bin/youtube-dl /usr/local/bin/

# Prepare data directory.
RUN mkdir -p /data && \
    cd /usr/local/oryx/platform/containers && \
    rm -rf data && ln -sf /data .

CMD ["./bootstrap"]
