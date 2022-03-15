'use strict';

const helper = require('./helper');
const moment = require('moment');

// See https://github.com/facebook/react/issues/11098#issuecomment-523977830
function disableConsoleLog(pfn) {
  const spy = jest.spyOn(console, 'log');
  spy.mockImplementation(() => {});

  try {
    pfn();
  } finally {
    spy.mockRestore();
  }
}

test('upgrade window not set', () => {
  disableConsoleLog(() => {
    expect(helper.inUpgradeWindow(null, null, null)).toBe(true);
    expect(helper.inUpgradeWindow(undefined, null, null)).toBe(true);
    expect(helper.inUpgradeWindow(undefined, 0, null)).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, null)).toBe(true);
  });
});

test('upgrade normal set', () => {
  disableConsoleLog(() => {
    expect(helper.inUpgradeWindow(0, 24, moment())).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 0}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 1}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 2}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 4}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 8}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 16}))).toBe(true);
    expect(helper.inUpgradeWindow(0, 24, moment({hour: 24}))).toBe(true);
  });
});

test('in window', () => {
  disableConsoleLog(() => {
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 1}))).toBe(true);
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 9}))).toBe(true);
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 2}))).toBe(true);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 20}))).toBe(true);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 2}))).toBe(true);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 23}))).toBe(true);
  });
});

test('not in window', () => {
  disableConsoleLog(() => {
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 0}))).toBe(false);
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 10}))).toBe(false);
    expect(helper.inUpgradeWindow(1, 8, moment({hour: 14}))).toBe(false);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 19}))).toBe(false);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 3}))).toBe(false);
    expect(helper.inUpgradeWindow(20, 6, moment({hour: 12}))).toBe(false);
  });
});

