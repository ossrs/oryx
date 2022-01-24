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
  auth: 201, // Verify token failed.
};

