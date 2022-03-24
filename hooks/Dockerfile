FROM registry.cn-hangzhou.aliyuncs.com/ossrs/node:slim

COPY js-core /usr/local/srs-cloud/js-core
COPY hooks /usr/local/srs-cloud/hooks
RUN cd /usr/local/srs-cloud/hooks && npm i

RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/hooks
CMD ["node", "."]
