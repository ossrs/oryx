'use strict';

// MySQL日期字段格式化字符串 @see https://stackoverflow.com/a/27381633
exports.MYSQL_DATETIME = 'YYYY-MM-DD HH:mm:ss';

// The redis key.
exports.SRS_SECRET_PUBLISH = 'SRS_SECRET_PUBLISH';
exports.SRS_FIRST_BOOT_DONE = 'SRS_FIRST_BOOT_DONE';

// Local redis config.
exports.redis = {
  host: 'localhost',
  port: 6379,
  password: '',
};

