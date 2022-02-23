FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /usr/local/srs-terraform/js-core
COPY tencent /usr/local/srs-terraform/tencent
RUN cd /usr/local/srs-terraform/tencent && npm i

WORKDIR /usr/local/srs-terraform/tencent
CMD ["node", "."]
