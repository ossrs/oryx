FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /srs-terraform/js-core
COPY hooks /srs-terraform/hooks
RUN cd /srs-terraform/hooks && npm i

WORKDIR /srs-terraform/hooks
CMD ["node", "."]
