'use strict';

// The redis key.
exports.redis = {
  // For LightHouse information, like region or source.
  SRS_TENCENT_LH: 'SRS_TENCENT_LH',
  // For account or auth, like secretId and secretKey.
  SRS_TENCENT_CAM: 'SRS_TENCENT_CAM',
  // For COS bucket and region.
  SRS_TENCENT_COS: 'SRS_TENCENT_COS',
  // For DVR HLS to COS, the stream or m3u8 list status.
  SRS_DVR_M3U8_ACTIVE: 'SRS_DVR_M3U8_ACTIVE',
  // Fresh ts files of m3u8, store in local disk, not uploaded yet.
  SRS_DVR_M3U8_LOCAL: 'SRS_DVR_M3U8_LOCAL',
  // Uploaded ts files of m3u8, store in COS.
  SRS_DVR_M3U8_UPLOADED: 'SRS_DVR_M3U8_UPLOADED',
  // The metadata of m3u8, generate the final result.
  SRS_DVR_M3U8_METADATA: 'SRS_DVR_M3U8_METADATA',
  // The patterns apply to DVR.
  SRS_DVR_PATTERNS: 'SRS_DVR_PATTERNS',
};

