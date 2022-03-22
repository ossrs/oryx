FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs:node-av

COPY js-core /usr/local/srs-terraform/js-core
COPY platform /usr/local/srs-terraform/platform

WORKDIR /usr/local/srs-terraform/platform
RUN npm install && (cd ui && npm install) && npm run build

ENV NODE_ENV=production
WORKDIR /usr/local/srs-terraform/platform
CMD ["node", "."]
