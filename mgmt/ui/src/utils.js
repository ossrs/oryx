const SRS_TERRAFORM_TOKEN = 'SRS_TERRAFORM_TOKEN';

export const Token = {
  save(data) {
    localStorage.setItem(SRS_TERRAFORM_TOKEN, JSON.stringify(data));
  },
  load() {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    return JSON.parse(info);
  },
  remove() {
    localStorage.removeItem(SRS_TERRAFORM_TOKEN);
  },
};

export const Tools = {
  mask(data) {
    const mask = `***${data.token.length}B***`;
    return JSON.stringify({...data, token: mask});
  }
};

export const Errors = {
  auth: 2001, // Verify token failed.
};

export const PlatformPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1c+ZAfJ93/qJ3bYp3SuVaMqYMniMCjNnFf20asK+oM7HJqFbBe/VZM2/Z2UkDHZiOqiArb1RLjYQeUFbUmPj2A5cCE8IPaeu28thbOdEC6wTztnAOdVzQBkBPytZiVR8DUUAzgz0tLoXB4nXGXQDntTgXoL/Rzn59BQIa7BzLlLnKc4TCn+LPpsOOmDPnnvjjJXpBKTY/rRTYvvgCUCQ/clSfBsgfQgP1p0nVRlH3FoZaJS4QRdzFVRKJtCytC1NwtgVNwRxpqYsJos9YW+yw+X/K5w7JAjG0v+9TycIzl5/Wd7R3zHMENe2uYx7XayksLc1ZLfgBD1/gldYd6l5VCcgHZJWKVsur8dNwvs0yWj3y9iOi1Lx+J8gLkMSqNouHVV2nVvSILoeWHaadd1+3ghuXKmbvauYI6mYai/T12vnEcxZ1yc6rVah8oy+vNwmpcKj2lixExrNW8JrhjLUU/Rlzla89es8JAZNfQDy7+ZOU1UGt//QqGZaiC8VhtV0= video@MB0`;


