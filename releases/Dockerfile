FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /srs-terraform/js-core
COPY releases /srs-terraform/releases
RUN cd /srs-terraform/releases && npm i

ENV PORT 9000
WORKDIR /srs-terraform/releases
CMD ["node", "."]
