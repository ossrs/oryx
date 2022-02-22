FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /usr/local/srs-terraform/js-core
COPY hooks /usr/local/srs-terraform/hooks
RUN cd /usr/local/srs-terraform/hooks && npm i

WORKDIR /usr/local/srs-terraform/hooks
CMD ["node", "."]
