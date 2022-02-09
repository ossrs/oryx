'use strict';

// Local config.
exports.config = {
  port: process.env.PORT || 2022,
};

// Redis key for upgrading.
exports.SRS_UPGRADING = 'SRS_UPGRADING';
exports.SRS_UPGRADE_STRATEGY = 'SRS_UPGRADE_STRATEGY';

