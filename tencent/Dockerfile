FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /usr/local/srs-cloud/js-core
COPY tencent /usr/local/srs-cloud/tencent
RUN cd /usr/local/srs-cloud/tencent && npm i

RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/tencent
CMD ["node", "."]
