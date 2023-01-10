"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AbstractClient = void 0;
const sdk_version_1 = require("./sdk_version");
const interface_1 = require("./interface");
const sign_1 = require("./sign");
const http_connection_1 = require("./http/http_connection");
const tencent_cloud_sdk_exception_1 = require("./exception/tencent_cloud_sdk_exception");
/**
 * @inner
 */
class AbstractClient {
    /**
     * 实例化client对象
     * @param {string} endpoint 接入点域名
     * @param {string} version 产品版本
     * @param {Credential} credential 认证信息实例
     * @param {string} region 产品地域
     * @param {ClientProfile} profile 可选配置实例
     */
    constructor(endpoint, version, { credential, region, profile }) {
        this.path = "/";
        /**
         * 认证信息实例
         */
        this.credential = Object.assign({
            secretId: null,
            secretKey: null,
            token: null,
        }, credential);
        /**
         * 产品地域
         */
        this.region = region || null;
        this.sdkVersion = "SDK_NODEJS_" + sdk_version_1.sdkVersion;
        this.apiVersion = version;
        this.endpoint = (profile && profile.httpProfile && profile.httpProfile.endpoint) || endpoint;
        /**
         * 可选配置实例
         * @type {ClientProfile}
         */
        this.profile = {
            signMethod: (profile && profile.signMethod) || "TC3-HMAC-SHA256",
            httpProfile: Object.assign({
                reqMethod: "POST",
                endpoint: null,
                protocol: "https://",
                reqTimeout: 60,
            }, profile && profile.httpProfile),
            language: profile.language,
        };
        if (this.profile.language && !interface_1.SUPPORT_LANGUAGE_LIST.includes(this.profile.language)) {
            throw new tencent_cloud_sdk_exception_1.default(`Language invalid, choices: ${interface_1.SUPPORT_LANGUAGE_LIST.join("|")}`);
        }
    }
    /**
     * @inner
     */
    async request(action, req, options, cb) {
        if (typeof options === "function") {
            cb = options;
            options = {};
        }
        try {
            const result = await this.doRequest(action, req, options);
            cb && cb(null, result);
            return result;
        }
        catch (e) {
            cb && cb(e, null);
            throw e;
        }
    }
    /**
     * @inner
     */
    async doRequest(action, req, options) {
        if (this.profile.signMethod === "TC3-HMAC-SHA256") {
            return this.doRequestWithSign3(action, req, options);
        }
        let params = this.mergeData(req);
        params = this.formatRequestData(action, params);
        let res;
        try {
            res = await http_connection_1.HttpConnection.doRequest({
                method: this.profile.httpProfile.reqMethod,
                url: this.profile.httpProfile.protocol + this.endpoint + this.path,
                data: params,
                timeout: this.profile.httpProfile.reqTimeout * 1000,
            });
        }
        catch (error) {
            throw new tencent_cloud_sdk_exception_1.default(error.message);
        }
        return this.parseResponse(res);
    }
    /**
     * @inner
     */
    async doRequestWithSign3(action, params, options) {
        let res;
        try {
            res = await http_connection_1.HttpConnection.doRequestWithSign3({
                method: this.profile.httpProfile.reqMethod,
                url: this.profile.httpProfile.protocol + this.endpoint + this.path,
                secretId: this.credential.secretId,
                secretKey: this.credential.secretKey,
                region: this.region,
                data: params || "",
                service: this.endpoint.split(".")[0],
                action: action,
                version: this.apiVersion,
                multipart: options && options.multipart,
                timeout: this.profile.httpProfile.reqTimeout * 1000,
                token: this.credential.token,
                requestClient: this.sdkVersion,
                language: this.profile.language,
            });
        }
        catch (e) {
            throw new tencent_cloud_sdk_exception_1.default(e.message);
        }
        return this.parseResponse(res);
    }
    async parseResponse(res) {
        if (res.status !== 200) {
            const tcError = new tencent_cloud_sdk_exception_1.default(res.statusText);
            tcError.httpCode = res.status;
            throw tcError;
        }
        else {
            const data = await res.json();
            if (data.Response.Error) {
                const tcError = new tencent_cloud_sdk_exception_1.default(data.Response.Error.Message, data.Response.RequestId);
                tcError.code = data.Response.Error.Code;
                throw tcError;
            }
            else {
                return data.Response;
            }
        }
    }
    /**
     * @inner
     */
    mergeData(data, prefix = "") {
        const ret = {};
        for (const k in data) {
            if (data[k] === null) {
                continue;
            }
            if (data[k] instanceof Array || data[k] instanceof Object) {
                Object.assign(ret, this.mergeData(data[k], prefix + k + "."));
            }
            else {
                ret[prefix + k] = data[k];
            }
        }
        return ret;
    }
    /**
     * @inner
     */
    formatRequestData(action, params) {
        params.Action = action;
        params.RequestClient = this.sdkVersion;
        params.Nonce = Math.round(Math.random() * 65535);
        params.Timestamp = Math.round(Date.now() / 1000);
        params.Version = this.apiVersion;
        if (this.credential.secretId) {
            params.SecretId = this.credential.secretId;
        }
        if (this.region) {
            params.Region = this.region;
        }
        if (this.credential.token) {
            params.Token = this.credential.token;
        }
        if (this.profile.language) {
            params.Language = this.profile.language;
        }
        if (this.profile.signMethod) {
            params.SignatureMethod = this.profile.signMethod;
        }
        const signStr = this.formatSignString(params);
        params.Signature = sign_1.default.sign(this.credential.secretKey, signStr, this.profile.signMethod);
        return params;
    }
    /**
     * @inner
     */
    formatSignString(params) {
        let strParam = "";
        const keys = Object.keys(params);
        keys.sort();
        for (const k in keys) {
            if (!keys.hasOwnProperty(k)) {
                continue;
            }
            //k = k.replace(/_/g, '.');
            strParam += "&" + keys[k] + "=" + params[keys[k]];
        }
        const strSign = this.profile.httpProfile.reqMethod.toLocaleUpperCase() +
            this.endpoint +
            this.path +
            "?" +
            strParam.slice(1);
        return strSign;
    }
}
exports.AbstractClient = AbstractClient;
