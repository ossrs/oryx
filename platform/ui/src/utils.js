const SRS_TERRAFORM_TOKEN = 'SRS_TERRAFORM_TOKEN';
const SRS_CLOUD_LOCALE = 'SRS_CLOUD_LOCALE';

export const Token = {
  save(data) {
    localStorage.setItem(SRS_TERRAFORM_TOKEN, JSON.stringify(data));
  },
  load() {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    if (!info) return null;

    const o = JSON.parse(info);
    return {token: o.token};
  },
  remove() {
    localStorage.removeItem(SRS_TERRAFORM_TOKEN);
  },
};

export const Locale = {
  _cache: null,
  save(data) {
    this._cache = data;
    localStorage.setItem(SRS_CLOUD_LOCALE, JSON.stringify(data));
  },
  load() {
    const info = localStorage.getItem(SRS_CLOUD_LOCALE);
    this._cache = info ? JSON.parse(info) : null;
    return this._cache;
  },
  current() {
    return this._cache?.lang || process.env.REACT_APP_LOCALE || 'zh';
  }
};

export const Tools = {
  mask(data) {
    const mask = `***${data.token.length}B***`;
    return JSON.stringify({...data, token: mask});
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

export const PlatformPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1c+ZAfJ93/qJ3bYp3SuVaMqYMniMCjNnFf20asK+oM7HJqFbBe/VZM2/Z2UkDHZiOqiArb1RLjYQeUFbUmPj2A5cCE8IPaeu28thbOdEC6wTztnAOdVzQBkBPytZiVR8DUUAzgz0tLoXB4nXGXQDntTgXoL/Rzn59BQIa7BzLlLnKc4TCn+LPpsOOmDPnnvjjJXpBKTY/rRTYvvgCUCQ/clSfBsgfQgP1p0nVRlH3FoZaJS4QRdzFVRKJtCytC1NwtgVNwRxpqYsJos9YW+yw+X/K5w7JAjG0v+9TycIzl5/Wd7R3zHMENe2uYx7XayksLc1ZLfgBD1/gldYd6l5VCcgHZJWKVsur8dNwvs0yWj3y9iOi1Lx+J8gLkMSqNouHVV2nVvSILoeWHaadd1+3ghuXKmbvauYI6mYai/T12vnEcxZ1yc6rVah8oy+vNwmpcKj2lixExrNW8JrhjLUU/Rlzla89es8JAZNfQDy7+ZOU1UGt//QqGZaiC8VhtV0= video@MB0`;

