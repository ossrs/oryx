FROM registry.cn-hangzhou.aliyuncs.com/ossrs/srs:node-av

# Cache packages for koa.
RUN npm i -g fs uuid jest moment dotenv koa koa-bodyparser koa-router koa-static koa2-cors \
  axios ioredis jsonwebtoken koa-mount koa-static-cache semver \
  cos-nodejs-sdk-v5 tencentcloud-sdk-nodejs

# Cache packages for react.
RUN npm i -g bootstrap http-proxy-middleware querystring react react-bootstrap react-bootstrap-icons \
  react-dom react-error-boundary react-qr-code react-router-dom react-scripts recharts web-vitals \
  @testing-library/jest-dom @testing-library/react @testing-library/user-event

