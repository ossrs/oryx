FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs:node-av

COPY js-core /usr/local/srs-cloud/js-core
COPY ffmpeg /usr/local/srs-cloud/ffmpeg
RUN cd /usr/local/srs-cloud/ffmpeg && npm i

RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/ffmpeg
CMD ["node", "."]
