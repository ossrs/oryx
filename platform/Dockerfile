
############################################################
# build
############################################################
ARG ARCH
FROM ${ARCH}ossrs/srs:node-av AS build

COPY js-core /usr/local/srs-cloud/js-core
COPY platform /usr/local/srs-cloud/platform

WORKDIR /usr/local/srs-cloud/platform
RUN npm install && (cd ui && npm install)

RUN env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh npm run buildDocker
RUN env PUBLIC_URL=/mgmt REACT_APP_LOCALE=en BUILD_PATH=build/en npm run buildDocker

############################################################
# dist
############################################################
FROM ${ARCH}node:slim as dist

COPY --from=build /usr/local/bin /usr/local/bin
COPY --from=build /usr/local/srs-cloud /usr/local/srs-cloud

ENV NODE_ENV=production
RUN ln -sf /usr/local/srs-cloud /usr/local/srs-terraform
WORKDIR /usr/local/srs-cloud/platform
CMD ["node", "."]
