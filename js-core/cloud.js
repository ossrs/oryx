'use strict';

exports.tencent = {
  cam: async (AbstractClient, secretId, secretKey, action, req) => {
    return await new AbstractClient(
      null,
      '2019-01-16',
      {
        credential: {
          secretId: secretId,
          secretKey: secretKey,
        },
        profile: {
          httpProfile: {
            endpoint: 'cam.tencentcloudapi.com',
          },
        },
        region: null,
      },
    ).request(
      action,
      req || {},
    );
  },
  vod: async (AbstractClient, secretId, secretKey, action, req) => {
    return await new AbstractClient(
      null,
      '2018-07-17',
      {
        credential: {
          secretId: secretId,
          secretKey: secretKey,
        },
        profile: {
          httpProfile: {
            endpoint: 'vod.tencentcloudapi.com',
          },
        },
        region: null,
      },
    ).request(
      action,
      req || {},
    );
  },
};

