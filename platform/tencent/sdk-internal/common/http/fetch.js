"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const node_fetch_1 = require("node-fetch");
const HttpsProxyAgent = require("https-proxy-agent");
function default_1(url, options) {
    const instanceOptions = options || {};
    if (!options.agent && process.env.http_proxy) {
        instanceOptions.agent = new HttpsProxyAgent(process.env.http_proxy);
    }
    return node_fetch_1.default(url, instanceOptions);
}
exports.default = default_1;
