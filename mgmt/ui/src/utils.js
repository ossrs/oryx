const SRS_TERRAFORM_TOKEN = 'SRS_TERRAFORM_TOKEN';

export const Token = {
  save(data) {
    localStorage.setItem(SRS_TERRAFORM_TOKEN, JSON.stringify(data));
  },
  load(callback) {
    const info = localStorage.getItem(SRS_TERRAFORM_TOKEN);
    callback && callback(JSON.parse(info));
  },
  remove(callback) {
    localStorage.removeItem(SRS_TERRAFORM_TOKEN);
    callback && callback();
  },
};

export const Tools = {
  mask(data) {
    const mask = `***${data.token.length}B***`;
    return JSON.stringify({...data, token: mask});
  }
};

