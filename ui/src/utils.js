//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
const SRS_TERRAFORM_TOKEN = 'SRS_TERRAFORM_TOKEN';
const SRS_STACK_LOCALE = 'SRS_STACK_LOCALE';
const SRS_STREAM_NAME = 'SRS_STREAM_NAME';

export const Token = {
  save: (data) => {
    localStorage.setItem(SRS_TERRAFORM_TOKEN, JSON.stringify(data));
  },
  load: () => {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    if (!info) return null;

    const o = JSON.parse(info);
    return {token: o.token};
  },
  updateBearer: (bearer) => {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    const o = JSON.parse(info || '{}');
    o.bearer = bearer;
    localStorage.setItem(SRS_TERRAFORM_TOKEN, JSON.stringify(o));
  },
  loadBearer: () => {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    const o = JSON.parse(info || '{}');
    return {token: o.bearer};
  },
  loadBearerHeader: () => {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    const o = JSON.parse(info || '{}');
    return o?.bearer ? {'Authorization': `Bearer ${o?.bearer}`} : {};
  },
  remove: () => {
    localStorage.removeItem(SRS_TERRAFORM_TOKEN);
  },
};

export const Locale = {
  _cache: null,
  save: (data) => {
    Locale._cache = data;
    localStorage.setItem(SRS_STACK_LOCALE, JSON.stringify(data));
  },
  load: () => {
    const info = localStorage.getItem(SRS_STACK_LOCALE);
    Locale._cache = info ? JSON.parse(info) : null;
    return Locale._cache;
  },
  current: () => {
    return Locale._cache?.lang || process.env.REACT_APP_LOCALE || 'zh';
  }
};

export const StreamName = {
  save: (name) => {
    localStorage.setItem(SRS_STREAM_NAME, name);
  },
  load: () => {
    return localStorage.getItem(SRS_STREAM_NAME) || 'livestream';
  }
};

export const MediaSource = {
  exts: ['.mp4', '.flv', '.ts', '.m4a', '.mp3', '.aac']
};

export const Tools = {
  mask(data) {
    return JSON.stringify({
      ...data,
      token: `***${data?.token?.length}B***`,
      bearer: `***${data?.bearer?.length}B***`,
    });
  },

  // Copy object, with optional extras fields, for example:
  //    copy({id: 0}, ['msg': 'hi'])
  // Return an object:
  //    {id: 0, msg: 'hi'}
  copy(from, extras) {
    let cp = Tools.merge({}, from);

    for (let i = 0; i < extras?.length; i += 2) {
      const k = extras[i];
      const v = extras[i + 1];
      const ov = cp[k];

      const obj = {};
      obj[k] = Tools.merge(ov, v);
      cp = Tools.merge(cp, obj);
    }
    return cp;
  },
  // Merge two object, rewrite dst by src fields.
  merge(dst, src) {
    if (typeof dst !== 'object') return src;
    if (typeof src !== 'object') return src;

    const cp = {};
    for (const k in dst) {
      cp[k] = dst[k];
    }
    for (const k in src) {
      cp[k] = src[k];
    }
    return cp;
  }
};

export const Clipboard = {
  copy(text) {
    if (navigator.clipboard) {
      return navigator.clipboard.writeText(text);
    }

    const e = document.createElement("textarea");
    e.value = text;
    e.style.position = "fixed";

    document.body.appendChild(e);
    e.focus();
    e.select();

    return new Promise((resolve, reject) => {
      try {
        const success = document.execCommand('copy');
        if (success) {
          resolve(true);
        } else {
          reject(false);
        }
      } catch (e) {
        reject(e);
      } finally {
        document.body.removeChild(e);
      }
    });
  }
};

export const Errors = {
  redis: 1007, // Redis is not ready.
  auth: 2001, // Verify token failed.
  btHttps: 3001, // Please use BT to configure HTTPS.
};

function buildStreamURL (vhost, app, stream) {
  if (vhost === '__defaultVhost__') {
    return `${app}/${stream}`;
  } else {
    return `${vhost}/${app}/${stream}`;
  }
};

export const StreamURL = {
  build: (vhost, app, stream) => {
    return buildStreamURL(vhost, app, stream);
  },
};
