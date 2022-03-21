FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs:node-av

COPY js-core /usr/local/srs-terraform/js-core
COPY platform /usr/local/srs-terraform/platform
RUN cd /usr/local/srs-terraform/platform && npm i

WORKDIR /usr/local/srs-terraform/platform
CMD ["node", "."]
