
############################################################
# build
############################################################
FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs-terraform:nodejs-1 AS build

COPY js-core /usr/local/srs-terraform/js-core
COPY platform /usr/local/srs-terraform/platform

WORKDIR /usr/local/srs-terraform/platform
RUN npm install && (cd ui && npm install)
RUN npm run build

############################################################
# dist
############################################################
FROM node:slim as dist

COPY --from=build /usr/local/srs-terraform /usr/local/srs-terraform

ENV NODE_ENV=production
WORKDIR /usr/local/srs-terraform/platform
CMD ["node", "."]
