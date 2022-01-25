FROM node:slim

COPY releases /releases
RUN cd /releases && npm i

WORKDIR /releases
CMD ["node", "."]
