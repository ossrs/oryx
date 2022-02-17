FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /srs-terraform/js-core
COPY tencent /srs-terraform/tencent
RUN cd /srs-terraform/tencent && npm i

WORKDIR /srs-terraform/tencent
CMD ["node", "."]
