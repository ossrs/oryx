FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs:node-av

COPY js-core /usr/local/srs-terraform/js-core
COPY ffmpeg /usr/local/srs-terraform/ffmpeg
RUN cd /usr/local/srs-terraform/ffmpeg && npm i

WORKDIR /usr/local/srs-terraform/ffmpeg
CMD ["node", "."]
