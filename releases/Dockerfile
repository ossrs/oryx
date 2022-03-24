FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /usr/local/srs-cloud/js-core
COPY releases /usr/local/srs-cloud/releases
RUN cd /usr/local/srs-cloud/releases && npm i

ENV PORT 9000
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/releases
CMD ["node", "."]
