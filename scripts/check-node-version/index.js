const semver = require('semver');

if (semver.lt(process.version, 'v14.0.0')) {
  console.log(`Only support nodejs 14+`);
  process.exit(1);
}

console.log(`Your nodejs ${process.version} is ok`);
process.exit(0);

