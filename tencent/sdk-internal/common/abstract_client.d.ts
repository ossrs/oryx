import { ClientProfile, Credential, ClientConfig } from "./interface";
declare type ResponseCallback = (error: string, rep: any) => void;
interface RequestOptions {
    multipart: boolean;
}
declare type ResponseData = any;
/**
 * @inner
 */
export declare class AbstractClient {
    sdkVersion: string;
    path: string;
    credential: Credential;
    region: string;
    apiVersion: string;
    endpoint: string;
    profile: ClientProfile;
    /**
     * 实例化client对象
     * @param {string} endpoint 接入点域名
     * @param {string} version 产品版本
     * @param {Credential} credential 认证信息实例
     * @param {string} region 产品地域
     * @param {ClientProfile} profile 可选配置实例
     */
    constructor(endpoint: string, version: string, { credential, region, profile }: ClientConfig);
    /**
     * @inner
     */
    request(action: string, req: any, options?: ResponseCallback | RequestOptions, cb?: ResponseCallback): Promise<ResponseData>;
    /**
     * @inner
     */
    private doRequest;
    /**
     * @inner
     */
    private doRequestWithSign3;
    private parseResponse;
    /**
     * @inner
     */
    private mergeData;
    /**
     * @inner
     */
    private formatRequestData;
    /**
     * @inner
     */
    private formatSignString;
}
export {};
