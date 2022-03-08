'use strict';

exports.inUpgradeWindow = (uwStart, uwDuration, now) => {
  if (uwStart === undefined || uwStart === null || !uwDuration || !now) {
    return true;
  }

  const [start, duration] = [parseInt(uwStart), parseInt(uwDuration)];
  const end = duration >= 24 ? start + 24 : (start + duration) % 24;

  let inWindow;
  if (start < end) {
    inWindow = (start <= now.hours() && now.hours() <= end);
  } else {
    inWindow = !(end < now.hours() && now.hours() < start);
  }
  console.log(`Upgrade window=${inWindow}, start=${start}, duration=${duration}, end=${end}, now=${now.format()}, hours=${now.hours()}`);

  return inWindow;
}

