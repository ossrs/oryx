'use strict';

const keys = require('./keys');
const cloud = require('./cloud');

exports.createVodService = async (redis, VodClient, AbstractClient, region) => {
  const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
  const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');
  if (!secretId || !secretKey) return console.log(`COS: Ignore for no secret`);

  // Whethether VoD config changed.
  let vodConfigChanged = false;

  // Create services if not exists.
  const service = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'service');
  if (!service) {
    try {
      await cloud.tencent.vod(
        AbstractClient, secretId, secretKey, 'CreateService',
      );

      console.log(`VoD: CreateService ok`);
      vodConfigChanged = true;
    } catch (e) {
      if (e && e.code === 'FailedOperation.ServiceExist') {
        console.log(`VoD: CreateService exist`);
      } else {
        throw e;
      }
    }
  } else {
    console.log(`VoD: Already exists service`);
  }

  // Create tencent vod storage region.
  const storage = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'storage');
  if (!storage) {
    await cloud.tencent.vod(
      AbstractClient, secretId, secretKey, 'CreateStorageRegion', {StorageRegion: region},
    );

    console.log(`VoD: CreateStorageRegion ok, region=${region}`);
    vodConfigChanged = true;
  } else {
    console.log(`VoD: Already exists storage region ${storage}`);
  }

  // Query templates.
  const transcode = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'transcode');
  if (!transcode) {
    const {
      TotalCount: count,
      TranscodeTemplateSet,
    } = await new VodClient({
      credential: {secretId, secretKey},
      region,
      profile: {
        httpProfile: {
          endpoint: "vod.tencentcloudapi.com",
        },
      },
    }).DescribeTranscodeTemplates({
      "Type": "Preset",
      "ContainerType": "Video",
      "Limit": 100,
      "Offset": 0,
      "TEHDType": "Common"
    });
    if (count >= 100) console.warn(`VoD: Exceed templates count=${count}`);

    const templates = TranscodeTemplateSet.filter(e => {
      return (e.Name.indexOf('Deprecated') >= 0) ? null : e;
    });

    console.log(`VoD: Query templates nn=${templates.length}`);
    await redis.hset(keys.redis.SRS_TENCENT_VOD, 'transcode', JSON.stringify({nn: templates.length, templates}));
  } else {
    console.log(`VoD: Already exists transcode templates nn=${JSON.parse(transcode).nn}`);
  }

  // Filter the remux template, covert to MP4.
  const remux = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'remux');
  if (!remux) {
    const {templates} = JSON.parse(await redis.hget(keys.redis.SRS_TENCENT_VOD, 'transcode'));
    const remuxMp4 = templates.filter(e => {
      return (e.Container === 'mp4' && e.VideoTemplate.Codec === 'copy' && e.AudioTemplate.Codec === 'copy') ? e : null;
    });

    if (remuxMp4 || remuxMp4.length) {
      await redis.hset(keys.redis.SRS_TENCENT_VOD, 'remux', JSON.stringify({
        definition: parseInt(remuxMp4[0].Definition),
        name: remuxMp4[0].Name,
        comment: remuxMp4[0].Comment,
        container: remuxMp4[0].Container,
        update: remuxMp4[0].UpdateTime,
      }));
      console.log(`VoD: Set remux template definition=${remuxMp4[0].Definition}`);
    } else {
      console.log(`VoD: Remux template not found`);
    }
  } else {
    console.log(`VoD: Already exists remux template definition=${JSON.parse(remux).definition}`);
  }

  // Query the global vod domain.
  const domain = await redis.hget(keys.redis.SRS_TENCENT_VOD, 'domain');
  if (!domain) {
    // See https://cloud.tencent.com/document/product/266/54176
    const {
      DomainSet,
    } = await new VodClient({
      credential: {secretId, secretKey},
      region,
      profile: {
        httpProfile: {
          endpoint: "vod.tencentcloudapi.com",
        },
      },
    }).DescribeVodDomains({});

    const domains = DomainSet.filter(e => {
      return e.DeployStatus === 'Online' ? e : null;
    });
    if (domains?.length) {
      await redis.hset(keys.redis.SRS_TENCENT_VOD, 'domain', domains[0].Domain);
      console.log(`VoD: Setup domain ${domains[0].Domain}`);
    } else {
      console.log(`VoD: No vod domain`);
    }
  } else {
    console.log(`VoD: Already exists domain ${domain}`);
  }

  if (vodConfigChanged) {
    await redis.hset(keys.redis.SRS_TENCENT_VOD, 'service', 'ok');
    await redis.hset(keys.redis.SRS_TENCENT_VOD, 'storage', region);
    console.log(`Vod: Config service for region=${region}`);
  }
};

