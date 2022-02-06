'use strict';

// Local config.
exports.config = {
  port: process.env.PORT || 2022,
};

// Redis key for upgrading.
exports.SRS_UPGRADING = 'SRS_UPGRADING';

