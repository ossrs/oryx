'use strict';

const keys = require('./keys');

exports.createCosBucket = async (redis, COS, region) => {
  const appId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'appId');
  const secretId = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretId');
  const secretKey = await redis.hget(keys.redis.SRS_TENCENT_CAM, 'secretKey');
  if (!appId || !secretId || !secretKey) return console.log(`COS: Ignore for no secret`);

  const cos = new COS({SecretId: secretId, SecretKey: secretKey});

  // Create COS bucket if not exists.
  // Note that should never use headBucket to query it, because it cause always fail if bucket not found, even though
  // the bucket is created.
  let bucket = await redis.hget(keys.redis.SRS_TENCENT_COS, 'bucket');
  if (!bucket) {
    // Add nonce to bucket name, to avoid conflict on different region as bellow:
    //    The requested bucket has already existed in other region.
    const nonce = Math.random().toString(16).slice(-4);
    bucket = `srs-lighthouse-${nonce}-${appId}`;

    // Avoid duplicated bucket.
    await redis.hset(keys.redis.SRS_TENCENT_COS, 'bucket', bucket);
    try {
      await putBucket(cos, bucket, region);
    } catch (e) {
      // Rollback when error.
      await redis.hdel(keys.redis.SRS_TENCENT_COS, 'bucket');
      throw e;
    }

    // See https://cloud.tencent.com/document/product/436/56556
    await redis.hset(keys.redis.SRS_TENCENT_COS, 'location', `${bucket}.cos.${region}.myqcloud.com`);
  } else {
    console.log(`COS: Already exists bucket=${bucket}`);
  }

  // Setup COS bucket if no policy.
  const policy = await redis.hget(keys.redis.SRS_TENCENT_COS, 'policy');
  if (!policy) {
    await putBucketPolicy(cos, bucket, region, appId);
    await redis.hset(keys.redis.SRS_TENCENT_COS, 'policy', 'read-without-list-files');
  } else {
    console.log(`COS: Already exists policy`);
  }

  // Setup the CORS of bucket.
  const cors = await redis.hget(keys.redis.SRS_TENCENT_COS, 'cors');
  if (!cors) {
    await putBucketCors(cos, bucket, region);
    await redis.hset(keys.redis.SRS_TENCENT_COS, 'cors', true);
  } else {
    console.log(`COS: Already exists CORS`);
  }
}

async function putBucket(cos, bucket, region) {
  const acl = 'private';
  await new Promise((resolve, reject) => {
    // See https://cloud.tencent.com/document/product/436/36118
    cos.putBucket({
      Bucket: bucket,
      Region: region,
      ACL: acl,
    }, function (err, data) {
      if (err) return reject(err);
      resolve(data);
    });
  });

  console.log(`COS: create cos bucket=${bucket}, acl=${acl}, region=${region}`);
}

async function putBucketPolicy(cos, bucket, region, appId) {
  await new Promise((resolve, reject) => {
    // Allow read without list files actions:
    //    cos:GetBucket
    //    cos:GetBucketObjectVersions
    // See https://cloud.tencent.com/document/product/436/43812
    cos.putBucketPolicy({
      Bucket: bucket,
      Region: region,
      Policy: {
        "version": "2.0",
        "Statement": [{
          "Effect": "Allow",
          "Principal": {
            "qcs": [
              "qcs::cam::anyone:anyone"
            ]
          },
          "Action": [
            "name/cos:HeadBucket",
            "name/cos:ListMultipartUploads",
            "name/cos:ListParts",
            "name/cos:GetObject",
            "name/cos:HeadObject",
            "name/cos:OptionsObject",
          ],
          "Resource": [
            `qcs::cos:${region}:uid/${appId}:${bucket}/*`,
          ]
        }],
      },
    }, function (err, data) {
      if (err) return reject(err);
      resolve(data);
    });
  });

  console.log(`COS: setup public read without list files ok, bucket=${bucket}`);
}

async function putBucketCors(cos, bucket, region) {
  await new Promise((resolve, reject) => {
    // See https://cloud.tencent.com/document/product/436/43811
    cos.putBucketCors({
      Bucket: bucket,
      Region: region,
      CORSRules: [{
        "AllowedOrigin": ["*"],
        "AllowedMethod": ["GET", "POST", "PUT", "DELETE", "HEAD"],
        "AllowedHeader": ["*"],
        "ExposeHeader": ["ETag", "Content-Length", "x-cos-request-id"],
        "MaxAgeSeconds": "0"
      }]
    }, function (err, data) {
      if (err) return reject(err);
      resolve(data);
    });
  });

  console.log(`COS: setup CORS ok, bucket=${bucket}`);
}

