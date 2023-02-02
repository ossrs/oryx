// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20180717

import (
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-07-17"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewApplyUploadRequest() (request *ApplyUploadRequest) {
    request = &ApplyUploadRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ApplyUpload")
    
    
    return
}

func NewApplyUploadResponse() (response *ApplyUploadResponse) {
    response = &ApplyUploadResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ApplyUpload
// * 我们强烈建议您使用云点播提供的 [服务端上传 SDK](/document/product/266/9759#1.-.E5.8F.91.E8.B5.B7.E4.B8.8A.E4.BC.A0) 来上传文件。直接调用 API 进行上传的难度和工作量都显著大于使用 SDK。
//
// * 该接口用于申请媒体文件（和封面文件）的上传，获取文件上传到云点播的元信息（包括上传路径、上传签名等），用于后续上传接口。
//
// * 上传流程请参考 [服务端上传综述](/document/product/266/9759)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXPIRETIME = "InvalidParameter.ExpireTime"
//  INVALIDPARAMETERVALUE_COVERTYPE = "InvalidParameterValue.CoverType"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_MEDIATYPE = "InvalidParameterValue.MediaType"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ApplyUpload(request *ApplyUploadRequest) (response *ApplyUploadResponse, err error) {
    return c.ApplyUploadWithContext(context.Background(), request)
}

// ApplyUpload
// * 我们强烈建议您使用云点播提供的 [服务端上传 SDK](/document/product/266/9759#1.-.E5.8F.91.E8.B5.B7.E4.B8.8A.E4.BC.A0) 来上传文件。直接调用 API 进行上传的难度和工作量都显著大于使用 SDK。
//
// * 该接口用于申请媒体文件（和封面文件）的上传，获取文件上传到云点播的元信息（包括上传路径、上传签名等），用于后续上传接口。
//
// * 上传流程请参考 [服务端上传综述](/document/product/266/9759)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXPIRETIME = "InvalidParameter.ExpireTime"
//  INVALIDPARAMETERVALUE_COVERTYPE = "InvalidParameterValue.CoverType"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_MEDIATYPE = "InvalidParameterValue.MediaType"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ApplyUploadWithContext(ctx context.Context, request *ApplyUploadRequest) (response *ApplyUploadResponse, err error) {
    if request == nil {
        request = NewApplyUploadRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ApplyUpload require credential")
    }

    request.SetContext(ctx)
    
    response = NewApplyUploadResponse()
    err = c.Send(request, response)
    return
}

func NewAttachMediaSubtitlesRequest() (request *AttachMediaSubtitlesRequest) {
    request = &AttachMediaSubtitlesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "AttachMediaSubtitles")
    
    
    return
}

func NewAttachMediaSubtitlesResponse() (response *AttachMediaSubtitlesResponse) {
    response = &AttachMediaSubtitlesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AttachMediaSubtitles
// 关联媒资字幕，将指定的字幕关联到转自适应码流模板号对应的媒体输出文件中（或解除关联）。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) AttachMediaSubtitles(request *AttachMediaSubtitlesRequest) (response *AttachMediaSubtitlesResponse, err error) {
    return c.AttachMediaSubtitlesWithContext(context.Background(), request)
}

// AttachMediaSubtitles
// 关联媒资字幕，将指定的字幕关联到转自适应码流模板号对应的媒体输出文件中（或解除关联）。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) AttachMediaSubtitlesWithContext(ctx context.Context, request *AttachMediaSubtitlesRequest) (response *AttachMediaSubtitlesResponse, err error) {
    if request == nil {
        request = NewAttachMediaSubtitlesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AttachMediaSubtitles require credential")
    }

    request.SetContext(ctx)
    
    response = NewAttachMediaSubtitlesResponse()
    err = c.Send(request, response)
    return
}

func NewCommitUploadRequest() (request *CommitUploadRequest) {
    request = &CommitUploadRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CommitUpload")
    
    
    return
}

func NewCommitUploadResponse() (response *CommitUploadResponse) {
    response = &CommitUploadResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CommitUpload
// 该接口用于确认媒体文件（和封面文件）上传到腾讯云点播的结果，并存储媒体信息，返回文件的播放地址和文件 ID。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_VODSESSIONKEY = "InvalidParameterValue.VodSessionKey"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CommitUpload(request *CommitUploadRequest) (response *CommitUploadResponse, err error) {
    return c.CommitUploadWithContext(context.Background(), request)
}

// CommitUpload
// 该接口用于确认媒体文件（和封面文件）上传到腾讯云点播的结果，并存储媒体信息，返回文件的播放地址和文件 ID。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_VODSESSIONKEY = "InvalidParameterValue.VodSessionKey"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CommitUploadWithContext(ctx context.Context, request *CommitUploadRequest) (response *CommitUploadResponse, err error) {
    if request == nil {
        request = NewCommitUploadRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CommitUpload require credential")
    }

    request.SetContext(ctx)
    
    response = NewCommitUploadResponse()
    err = c.Send(request, response)
    return
}

func NewComposeMediaRequest() (request *ComposeMediaRequest) {
    request = &ComposeMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ComposeMedia")
    
    
    return
}

func NewComposeMediaResponse() (response *ComposeMediaResponse) {
    response = &ComposeMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ComposeMedia
// 该接口用于合成媒体文件，可以达到以下效果：
//
// 
//
// 1. **画面旋转**：对视频、图片的画面旋转一定角度，或按照某个方向翻转。
//
// 2. **声音控制**：升高降低视频、音频中声音的音量，或者对视频静音。
//
// 3. **画面叠加**：将视频、图片中的画面依序叠加在一起，如实现“画中画”的效果。
//
// 4. **声音混合**：将视频、音频中的声音混合在一起（混音）。
//
// 5. **声音提取**：将视频中的音频提取出来（不保留画面）。
//
// 6. **裁剪**：对视频、音频裁剪出指定时间段。
//
// 7. **拼接**：对视频、音频、图片按时间顺序前后拼接。
//
// 8. **转场**：将多段视频或图片拼接时，可以在段落之间添加转场效果。
//
// 
//
// 合成后的媒体封装格式可以是 MP4（视频）或 MP3（音频）。如使用事件通知，事件通知的类型为 [视频合成完成](https://cloud.tencent.com/document/product/266/43000)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_UPLOADCOSFAIL = "FailedOperation.UploadCosFail"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLIPDURATION = "InvalidParameterValue.ClipDuration"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ComposeMedia(request *ComposeMediaRequest) (response *ComposeMediaResponse, err error) {
    return c.ComposeMediaWithContext(context.Background(), request)
}

// ComposeMedia
// 该接口用于合成媒体文件，可以达到以下效果：
//
// 
//
// 1. **画面旋转**：对视频、图片的画面旋转一定角度，或按照某个方向翻转。
//
// 2. **声音控制**：升高降低视频、音频中声音的音量，或者对视频静音。
//
// 3. **画面叠加**：将视频、图片中的画面依序叠加在一起，如实现“画中画”的效果。
//
// 4. **声音混合**：将视频、音频中的声音混合在一起（混音）。
//
// 5. **声音提取**：将视频中的音频提取出来（不保留画面）。
//
// 6. **裁剪**：对视频、音频裁剪出指定时间段。
//
// 7. **拼接**：对视频、音频、图片按时间顺序前后拼接。
//
// 8. **转场**：将多段视频或图片拼接时，可以在段落之间添加转场效果。
//
// 
//
// 合成后的媒体封装格式可以是 MP4（视频）或 MP3（音频）。如使用事件通知，事件通知的类型为 [视频合成完成](https://cloud.tencent.com/document/product/266/43000)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_UPLOADCOSFAIL = "FailedOperation.UploadCosFail"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLIPDURATION = "InvalidParameterValue.ClipDuration"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ComposeMediaWithContext(ctx context.Context, request *ComposeMediaRequest) (response *ComposeMediaResponse, err error) {
    if request == nil {
        request = NewComposeMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ComposeMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewComposeMediaResponse()
    err = c.Send(request, response)
    return
}

func NewConfirmEventsRequest() (request *ConfirmEventsRequest) {
    request = &ConfirmEventsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ConfirmEvents")
    
    
    return
}

func NewConfirmEventsResponse() (response *ConfirmEventsResponse) {
    response = &ConfirmEventsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ConfirmEvents
// * 开发者调用拉取事件通知，获取到事件后，必须调用该接口来确认消息已经收到；
//
// * 开发者获取到事件句柄后，等待确认的有效时间为 30 秒，超出 30 秒会报参数错误（4000）；
//
// * 更多参考事件通知的[可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ConfirmEvents(request *ConfirmEventsRequest) (response *ConfirmEventsResponse, err error) {
    return c.ConfirmEventsWithContext(context.Background(), request)
}

// ConfirmEvents
// * 开发者调用拉取事件通知，获取到事件后，必须调用该接口来确认消息已经收到；
//
// * 开发者获取到事件句柄后，等待确认的有效时间为 30 秒，超出 30 秒会报参数错误（4000）；
//
// * 更多参考事件通知的[可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ConfirmEventsWithContext(ctx context.Context, request *ConfirmEventsRequest) (response *ConfirmEventsResponse, err error) {
    if request == nil {
        request = NewConfirmEventsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ConfirmEvents require credential")
    }

    request.SetContext(ctx)
    
    response = NewConfirmEventsResponse()
    err = c.Send(request, response)
    return
}

func NewCreateAIAnalysisTemplateRequest() (request *CreateAIAnalysisTemplateRequest) {
    request = &CreateAIAnalysisTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateAIAnalysisTemplate")
    
    
    return
}

func NewCreateAIAnalysisTemplateResponse() (response *CreateAIAnalysisTemplateResponse) {
    response = &CreateAIAnalysisTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateAIAnalysisTemplate
// 创建用户自定义音视频内容分析模板，数量上限：50。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLASSIFCATIONCONFIGURE = "InvalidParameterValue.ClassifcationConfigure"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COVERCONFIGURE = "InvalidParameterValue.CoverConfigure"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  INVALIDPARAMETERVALUE_FRAMETAGCONFIGURE = "InvalidParameterValue.FrameTagConfigure"
//  INVALIDPARAMETERVALUE_HIGHLIGHTCONFIGURE = "InvalidParameterValue.HighlightConfigure"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGCONFIGURE = "InvalidParameterValue.TagConfigure"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAIAnalysisTemplate(request *CreateAIAnalysisTemplateRequest) (response *CreateAIAnalysisTemplateResponse, err error) {
    return c.CreateAIAnalysisTemplateWithContext(context.Background(), request)
}

// CreateAIAnalysisTemplate
// 创建用户自定义音视频内容分析模板，数量上限：50。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLASSIFCATIONCONFIGURE = "InvalidParameterValue.ClassifcationConfigure"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COVERCONFIGURE = "InvalidParameterValue.CoverConfigure"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  INVALIDPARAMETERVALUE_FRAMETAGCONFIGURE = "InvalidParameterValue.FrameTagConfigure"
//  INVALIDPARAMETERVALUE_HIGHLIGHTCONFIGURE = "InvalidParameterValue.HighlightConfigure"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGCONFIGURE = "InvalidParameterValue.TagConfigure"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAIAnalysisTemplateWithContext(ctx context.Context, request *CreateAIAnalysisTemplateRequest) (response *CreateAIAnalysisTemplateResponse, err error) {
    if request == nil {
        request = NewCreateAIAnalysisTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateAIAnalysisTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateAIAnalysisTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateAIRecognitionTemplateRequest() (request *CreateAIRecognitionTemplateRequest) {
    request = &CreateAIRecognitionTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateAIRecognitionTemplate")
    
    
    return
}

func NewCreateAIRecognitionTemplateResponse() (response *CreateAIRecognitionTemplateResponse) {
    response = &CreateAIRecognitionTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateAIRecognitionTemplate
// 创建用户自定义音视频内容识别模板，数量上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFAULTLIBRARYLABELSET = "InvalidParameterValue.DefaultLibraryLabelSet"
//  INVALIDPARAMETERVALUE_FACELIBRARY = "InvalidParameterValue.FaceLibrary"
//  INVALIDPARAMETERVALUE_FACESCORE = "InvalidParameterValue.FaceScore"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_OBJECTLIBRARY = "InvalidParameterValue.ObjectLibrary"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SUBTITLEFORMAT = "InvalidParameterValue.SubtitleFormat"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  INVALIDPARAMETERVALUE_USERDEFINELIBRARYLABELSET = "InvalidParameterValue.UserDefineLibraryLabelSet"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAIRecognitionTemplate(request *CreateAIRecognitionTemplateRequest) (response *CreateAIRecognitionTemplateResponse, err error) {
    return c.CreateAIRecognitionTemplateWithContext(context.Background(), request)
}

// CreateAIRecognitionTemplate
// 创建用户自定义音视频内容识别模板，数量上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFAULTLIBRARYLABELSET = "InvalidParameterValue.DefaultLibraryLabelSet"
//  INVALIDPARAMETERVALUE_FACELIBRARY = "InvalidParameterValue.FaceLibrary"
//  INVALIDPARAMETERVALUE_FACESCORE = "InvalidParameterValue.FaceScore"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_OBJECTLIBRARY = "InvalidParameterValue.ObjectLibrary"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SUBTITLEFORMAT = "InvalidParameterValue.SubtitleFormat"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  INVALIDPARAMETERVALUE_USERDEFINELIBRARYLABELSET = "InvalidParameterValue.UserDefineLibraryLabelSet"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAIRecognitionTemplateWithContext(ctx context.Context, request *CreateAIRecognitionTemplateRequest) (response *CreateAIRecognitionTemplateResponse, err error) {
    if request == nil {
        request = NewCreateAIRecognitionTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateAIRecognitionTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateAIRecognitionTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateAdaptiveDynamicStreamingTemplateRequest() (request *CreateAdaptiveDynamicStreamingTemplateRequest) {
    request = &CreateAdaptiveDynamicStreamingTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateAdaptiveDynamicStreamingTemplate")
    
    
    return
}

func NewCreateAdaptiveDynamicStreamingTemplateResponse() (response *CreateAdaptiveDynamicStreamingTemplateResponse) {
    response = &CreateAdaptiveDynamicStreamingTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateAdaptiveDynamicStreamingTemplate
// 创建转自适应码流模板，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEOBITRATE = "InvalidParameterValue.DisableHigherVideoBitrate"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEORESOLUTION = "InvalidParameterValue.DisableHigherVideoResolution"
//  INVALIDPARAMETERVALUE_DRMTYPE = "InvalidParameterValue.DrmType"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_SOUNDSYSTEM = "InvalidParameterValue.SoundSystem"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
func (c *Client) CreateAdaptiveDynamicStreamingTemplate(request *CreateAdaptiveDynamicStreamingTemplateRequest) (response *CreateAdaptiveDynamicStreamingTemplateResponse, err error) {
    return c.CreateAdaptiveDynamicStreamingTemplateWithContext(context.Background(), request)
}

// CreateAdaptiveDynamicStreamingTemplate
// 创建转自适应码流模板，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEOBITRATE = "InvalidParameterValue.DisableHigherVideoBitrate"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEORESOLUTION = "InvalidParameterValue.DisableHigherVideoResolution"
//  INVALIDPARAMETERVALUE_DRMTYPE = "InvalidParameterValue.DrmType"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_SOUNDSYSTEM = "InvalidParameterValue.SoundSystem"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
func (c *Client) CreateAdaptiveDynamicStreamingTemplateWithContext(ctx context.Context, request *CreateAdaptiveDynamicStreamingTemplateRequest) (response *CreateAdaptiveDynamicStreamingTemplateResponse, err error) {
    if request == nil {
        request = NewCreateAdaptiveDynamicStreamingTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateAdaptiveDynamicStreamingTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateAdaptiveDynamicStreamingTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateAnimatedGraphicsTemplateRequest() (request *CreateAnimatedGraphicsTemplateRequest) {
    request = &CreateAnimatedGraphicsTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateAnimatedGraphicsTemplate")
    
    
    return
}

func NewCreateAnimatedGraphicsTemplateResponse() (response *CreateAnimatedGraphicsTemplateResponse) {
    response = &CreateAnimatedGraphicsTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateAnimatedGraphicsTemplate
// 创建用户自定义转动图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FORMATWEBPLACKWIDTHANDHEIGHT = "InvalidParameterValue.FormatWebpLackWidthAndHeight"
//  INVALIDPARAMETERVALUE_FORMATWEBPWIDTHANDHEIGHTBOTHZERO = "InvalidParameterValue.FormatWebpWidthAndHeightBothZero"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_QUALITY = "InvalidParameterValue.Quality"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAnimatedGraphicsTemplate(request *CreateAnimatedGraphicsTemplateRequest) (response *CreateAnimatedGraphicsTemplateResponse, err error) {
    return c.CreateAnimatedGraphicsTemplateWithContext(context.Background(), request)
}

// CreateAnimatedGraphicsTemplate
// 创建用户自定义转动图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FORMATWEBPLACKWIDTHANDHEIGHT = "InvalidParameterValue.FormatWebpLackWidthAndHeight"
//  INVALIDPARAMETERVALUE_FORMATWEBPWIDTHANDHEIGHTBOTHZERO = "InvalidParameterValue.FormatWebpWidthAndHeightBothZero"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_QUALITY = "InvalidParameterValue.Quality"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateAnimatedGraphicsTemplateWithContext(ctx context.Context, request *CreateAnimatedGraphicsTemplateRequest) (response *CreateAnimatedGraphicsTemplateResponse, err error) {
    if request == nil {
        request = NewCreateAnimatedGraphicsTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateAnimatedGraphicsTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateAnimatedGraphicsTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateClassRequest() (request *CreateClassRequest) {
    request = &CreateClassRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateClass")
    
    
    return
}

func NewCreateClassResponse() (response *CreateClassResponse) {
    response = &CreateClassResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateClass
// * 用于对媒体进行分类管理；
//
// * 该接口不影响既有媒体的分类，如需修改媒体分类，请调用[修改媒体文件属性](/document/product/266/31762)接口。
//
// * 分类层次不可超过 4 层。
//
// * 每个分类的子类数量不可超过 500 个。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_CLASSLEVELLIMITEXCEEDED = "FailedOperation.ClassLevelLimitExceeded"
//  FAILEDOPERATION_CLASSNAMEDUPLICATE = "FailedOperation.ClassNameDuplicate"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_PARENTIDNOFOUND = "FailedOperation.ParentIdNoFound"
//  FAILEDOPERATION_SUBCLASSLIMITEXCEEDED = "FailedOperation.SubclassLimitExceeded"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSNAME = "InvalidParameterValue.ClassName"
//  INVALIDPARAMETERVALUE_PARENTID = "InvalidParameterValue.ParentId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateClass(request *CreateClassRequest) (response *CreateClassResponse, err error) {
    return c.CreateClassWithContext(context.Background(), request)
}

// CreateClass
// * 用于对媒体进行分类管理；
//
// * 该接口不影响既有媒体的分类，如需修改媒体分类，请调用[修改媒体文件属性](/document/product/266/31762)接口。
//
// * 分类层次不可超过 4 层。
//
// * 每个分类的子类数量不可超过 500 个。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_CLASSLEVELLIMITEXCEEDED = "FailedOperation.ClassLevelLimitExceeded"
//  FAILEDOPERATION_CLASSNAMEDUPLICATE = "FailedOperation.ClassNameDuplicate"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_PARENTIDNOFOUND = "FailedOperation.ParentIdNoFound"
//  FAILEDOPERATION_SUBCLASSLIMITEXCEEDED = "FailedOperation.SubclassLimitExceeded"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSNAME = "InvalidParameterValue.ClassName"
//  INVALIDPARAMETERVALUE_PARENTID = "InvalidParameterValue.ParentId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateClassWithContext(ctx context.Context, request *CreateClassRequest) (response *CreateClassResponse, err error) {
    if request == nil {
        request = NewCreateClassRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateClass require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateClassResponse()
    err = c.Send(request, response)
    return
}

func NewCreateContentReviewTemplateRequest() (request *CreateContentReviewTemplateRequest) {
    request = &CreateContentReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateContentReviewTemplate")
    
    
    return
}

func NewCreateContentReviewTemplateResponse() (response *CreateContentReviewTemplateResponse) {
    response = &CreateContentReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [创建审核模板](https://cloud.tencent.com/document/api/266/84391)。
//
// 创建用户自定义音视频内容审核模板，数量上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETERVALUE_BLOCKCONFIDENCE = "InvalidParameterValue.BlockConfidence"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REVIEWCONFIDENCE = "InvalidParameterValue.ReviewConfidence"
//  INVALIDPARAMETERVALUE_REVIEWWALLSWITCH = "InvalidParameterValue.ReviewWallSwitch"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateContentReviewTemplate(request *CreateContentReviewTemplateRequest) (response *CreateContentReviewTemplateResponse, err error) {
    return c.CreateContentReviewTemplateWithContext(context.Background(), request)
}

// CreateContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [创建审核模板](https://cloud.tencent.com/document/api/266/84391)。
//
// 创建用户自定义音视频内容审核模板，数量上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETERVALUE_BLOCKCONFIDENCE = "InvalidParameterValue.BlockConfidence"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REVIEWCONFIDENCE = "InvalidParameterValue.ReviewConfidence"
//  INVALIDPARAMETERVALUE_REVIEWWALLSWITCH = "InvalidParameterValue.ReviewWallSwitch"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateContentReviewTemplateWithContext(ctx context.Context, request *CreateContentReviewTemplateRequest) (response *CreateContentReviewTemplateResponse, err error) {
    if request == nil {
        request = NewCreateContentReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateContentReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateContentReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateHeadTailTemplateRequest() (request *CreateHeadTailTemplateRequest) {
    request = &CreateHeadTailTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateHeadTailTemplate")
    
    
    return
}

func NewCreateHeadTailTemplateResponse() (response *CreateHeadTailTemplateResponse) {
    response = &CreateHeadTailTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateHeadTailTemplate
// 创建片头片尾模板。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateHeadTailTemplate(request *CreateHeadTailTemplateRequest) (response *CreateHeadTailTemplateResponse, err error) {
    return c.CreateHeadTailTemplateWithContext(context.Background(), request)
}

// CreateHeadTailTemplate
// 创建片头片尾模板。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateHeadTailTemplateWithContext(ctx context.Context, request *CreateHeadTailTemplateRequest) (response *CreateHeadTailTemplateResponse, err error) {
    if request == nil {
        request = NewCreateHeadTailTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateHeadTailTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateHeadTailTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateImageProcessingTemplateRequest() (request *CreateImageProcessingTemplateRequest) {
    request = &CreateImageProcessingTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateImageProcessingTemplate")
    
    
    return
}

func NewCreateImageProcessingTemplateResponse() (response *CreateImageProcessingTemplateResponse) {
    response = &CreateImageProcessingTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateImageProcessingTemplate
// 创建一个用户自定义的图片处理模板，数量上限：16。最多支持三次操作，例如：裁剪-缩略-裁剪。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_CUTANDCROPS = "InvalidParameterValue.CutAndCrops"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_THUMBNAILS = "InvalidParameterValue.Thumbnails"
//  INVALIDPARAMETERVALUE_WATERMARKS = "InvalidParameterValue.Watermarks"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateImageProcessingTemplate(request *CreateImageProcessingTemplateRequest) (response *CreateImageProcessingTemplateResponse, err error) {
    return c.CreateImageProcessingTemplateWithContext(context.Background(), request)
}

// CreateImageProcessingTemplate
// 创建一个用户自定义的图片处理模板，数量上限：16。最多支持三次操作，例如：裁剪-缩略-裁剪。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_CUTANDCROPS = "InvalidParameterValue.CutAndCrops"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_THUMBNAILS = "InvalidParameterValue.Thumbnails"
//  INVALIDPARAMETERVALUE_WATERMARKS = "InvalidParameterValue.Watermarks"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateImageProcessingTemplateWithContext(ctx context.Context, request *CreateImageProcessingTemplateRequest) (response *CreateImageProcessingTemplateResponse, err error) {
    if request == nil {
        request = NewCreateImageProcessingTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateImageProcessingTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateImageProcessingTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateImageSpriteTemplateRequest() (request *CreateImageSpriteTemplateRequest) {
    request = &CreateImageSpriteTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateImageSpriteTemplate")
    
    
    return
}

func NewCreateImageSpriteTemplateResponse() (response *CreateImageSpriteTemplateResponse) {
    response = &CreateImageSpriteTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateImageSpriteTemplate
// 创建用户自定义雪碧图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COLUMNCOUNT = "InvalidParameterValue.ColumnCount"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_ROWCOUNT = "InvalidParameterValue.RowCount"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateImageSpriteTemplate(request *CreateImageSpriteTemplateRequest) (response *CreateImageSpriteTemplateResponse, err error) {
    return c.CreateImageSpriteTemplateWithContext(context.Background(), request)
}

// CreateImageSpriteTemplate
// 创建用户自定义雪碧图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COLUMNCOUNT = "InvalidParameterValue.ColumnCount"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_ROWCOUNT = "InvalidParameterValue.RowCount"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateImageSpriteTemplateWithContext(ctx context.Context, request *CreateImageSpriteTemplateRequest) (response *CreateImageSpriteTemplateResponse, err error) {
    if request == nil {
        request = NewCreateImageSpriteTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateImageSpriteTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateImageSpriteTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreatePersonSampleRequest() (request *CreatePersonSampleRequest) {
    request = &CreatePersonSampleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreatePersonSample")
    
    
    return
}

func NewCreatePersonSampleResponse() (response *CreatePersonSampleResponse) {
    response = &CreatePersonSampleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreatePersonSample
// 该接口用于创建素材样本，用于通过五官定位等技术，进行内容识别、不适宜视频识别等视频处理。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FACEDUPLICATE = "InvalidParameterValue.FaceDuplicate"
//  INVALIDPARAMETERVALUE_PICFORMATERROR = "InvalidParameterValue.PicFormatError"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreatePersonSample(request *CreatePersonSampleRequest) (response *CreatePersonSampleResponse, err error) {
    return c.CreatePersonSampleWithContext(context.Background(), request)
}

// CreatePersonSample
// 该接口用于创建素材样本，用于通过五官定位等技术，进行内容识别、不适宜视频识别等视频处理。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FACEDUPLICATE = "InvalidParameterValue.FaceDuplicate"
//  INVALIDPARAMETERVALUE_PICFORMATERROR = "InvalidParameterValue.PicFormatError"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreatePersonSampleWithContext(ctx context.Context, request *CreatePersonSampleRequest) (response *CreatePersonSampleResponse, err error) {
    if request == nil {
        request = NewCreatePersonSampleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreatePersonSample require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreatePersonSampleResponse()
    err = c.Send(request, response)
    return
}

func NewCreateProcedureTemplateRequest() (request *CreateProcedureTemplateRequest) {
    request = &CreateProcedureTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateProcedureTemplate")
    
    
    return
}

func NewCreateProcedureTemplateResponse() (response *CreateProcedureTemplateResponse) {
    response = &CreateProcedureTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateProcedureTemplate
// 创建用户自定义的任务流模板，模板上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXISTEDPROCEDURENAME = "InvalidParameter.ExistedProcedureName"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateProcedureTemplate(request *CreateProcedureTemplateRequest) (response *CreateProcedureTemplateResponse, err error) {
    return c.CreateProcedureTemplateWithContext(context.Background(), request)
}

// CreateProcedureTemplate
// 创建用户自定义的任务流模板，模板上限：50。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXISTEDPROCEDURENAME = "InvalidParameter.ExistedProcedureName"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateProcedureTemplateWithContext(ctx context.Context, request *CreateProcedureTemplateRequest) (response *CreateProcedureTemplateResponse, err error) {
    if request == nil {
        request = NewCreateProcedureTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateProcedureTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateProcedureTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateReviewTemplateRequest() (request *CreateReviewTemplateRequest) {
    request = &CreateReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateReviewTemplate")
    
    
    return
}

func NewCreateReviewTemplateResponse() (response *CreateReviewTemplateResponse) {
    response = &CreateReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateReviewTemplate
// 创建用户自定义审核模板，数量上限：50。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETER_LABELS = "InvalidParameter.Labels"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELS = "InvalidParameterValue.Labels"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateReviewTemplate(request *CreateReviewTemplateRequest) (response *CreateReviewTemplateResponse, err error) {
    return c.CreateReviewTemplateWithContext(context.Background(), request)
}

// CreateReviewTemplate
// 创建用户自定义审核模板，数量上限：50。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INVALIDPARAMETER_LABELS = "InvalidParameter.Labels"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELS = "InvalidParameterValue.Labels"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateReviewTemplateWithContext(ctx context.Context, request *CreateReviewTemplateRequest) (response *CreateReviewTemplateResponse, err error) {
    if request == nil {
        request = NewCreateReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSampleSnapshotTemplateRequest() (request *CreateSampleSnapshotTemplateRequest) {
    request = &CreateSampleSnapshotTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateSampleSnapshotTemplate")
    
    
    return
}

func NewCreateSampleSnapshotTemplateResponse() (response *CreateSampleSnapshotTemplateResponse) {
    response = &CreateSampleSnapshotTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateSampleSnapshotTemplate
// 创建用户自定义采样截图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSampleSnapshotTemplate(request *CreateSampleSnapshotTemplateRequest) (response *CreateSampleSnapshotTemplateResponse, err error) {
    return c.CreateSampleSnapshotTemplateWithContext(context.Background(), request)
}

// CreateSampleSnapshotTemplate
// 创建用户自定义采样截图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSampleSnapshotTemplateWithContext(ctx context.Context, request *CreateSampleSnapshotTemplateRequest) (response *CreateSampleSnapshotTemplateResponse, err error) {
    if request == nil {
        request = NewCreateSampleSnapshotTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateSampleSnapshotTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateSampleSnapshotTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSnapshotByTimeOffsetTemplateRequest() (request *CreateSnapshotByTimeOffsetTemplateRequest) {
    request = &CreateSnapshotByTimeOffsetTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateSnapshotByTimeOffsetTemplate")
    
    
    return
}

func NewCreateSnapshotByTimeOffsetTemplateResponse() (response *CreateSnapshotByTimeOffsetTemplateResponse) {
    response = &CreateSnapshotByTimeOffsetTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateSnapshotByTimeOffsetTemplate
// 创建用户自定义指定时间点截图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSnapshotByTimeOffsetTemplate(request *CreateSnapshotByTimeOffsetTemplateRequest) (response *CreateSnapshotByTimeOffsetTemplateResponse, err error) {
    return c.CreateSnapshotByTimeOffsetTemplateWithContext(context.Background(), request)
}

// CreateSnapshotByTimeOffsetTemplate
// 创建用户自定义指定时间点截图模板，数量上限：16。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSnapshotByTimeOffsetTemplateWithContext(ctx context.Context, request *CreateSnapshotByTimeOffsetTemplateRequest) (response *CreateSnapshotByTimeOffsetTemplateResponse, err error) {
    if request == nil {
        request = NewCreateSnapshotByTimeOffsetTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateSnapshotByTimeOffsetTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateSnapshotByTimeOffsetTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateStorageRegionRequest() (request *CreateStorageRegionRequest) {
    request = &CreateStorageRegionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateStorageRegion")
    
    
    return
}

func NewCreateStorageRegionResponse() (response *CreateStorageRegionResponse) {
    response = &CreateStorageRegionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStorageRegion
// 该接口用于开通某地域的存储。
//
//   1. 用户开通点播业务时，系统默认为用户开通了部分地域的存储，用户如果需要开通其它地域的存储，可以通过该接口进行开通。
//
//   2. 通过 DescribeStorageRegions 接口可以查询到所有存储地域及已经开通的地域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_USERSTATUSINAVLID = "FailedOperation.UserStatusInavlid"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateStorageRegion(request *CreateStorageRegionRequest) (response *CreateStorageRegionResponse, err error) {
    return c.CreateStorageRegionWithContext(context.Background(), request)
}

// CreateStorageRegion
// 该接口用于开通某地域的存储。
//
//   1. 用户开通点播业务时，系统默认为用户开通了部分地域的存储，用户如果需要开通其它地域的存储，可以通过该接口进行开通。
//
//   2. 通过 DescribeStorageRegions 接口可以查询到所有存储地域及已经开通的地域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_USERSTATUSINAVLID = "FailedOperation.UserStatusInavlid"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateStorageRegionWithContext(ctx context.Context, request *CreateStorageRegionRequest) (response *CreateStorageRegionResponse, err error) {
    if request == nil {
        request = NewCreateStorageRegionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStorageRegion require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStorageRegionResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSubAppIdRequest() (request *CreateSubAppIdRequest) {
    request = &CreateSubAppIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateSubAppId")
    
    
    return
}

func NewCreateSubAppIdResponse() (response *CreateSubAppIdResponse) {
    response = &CreateSubAppIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateSubAppId
// 该接口用于创建点播子应用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DESCRIPTION = "InvalidParameterValue.Description"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSubAppId(request *CreateSubAppIdRequest) (response *CreateSubAppIdResponse, err error) {
    return c.CreateSubAppIdWithContext(context.Background(), request)
}

// CreateSubAppId
// 该接口用于创建点播子应用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DESCRIPTION = "InvalidParameterValue.Description"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSubAppIdWithContext(ctx context.Context, request *CreateSubAppIdRequest) (response *CreateSubAppIdResponse, err error) {
    if request == nil {
        request = NewCreateSubAppIdRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateSubAppId require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateSubAppIdResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSuperPlayerConfigRequest() (request *CreateSuperPlayerConfigRequest) {
    request = &CreateSuperPlayerConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateSuperPlayerConfig")
    
    
    return
}

func NewCreateSuperPlayerConfigResponse() (response *CreateSuperPlayerConfigResponse) {
    response = &CreateSuperPlayerConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateSuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 创建播放器配置，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSuperPlayerConfig(request *CreateSuperPlayerConfigRequest) (response *CreateSuperPlayerConfigResponse, err error) {
    return c.CreateSuperPlayerConfigWithContext(context.Background(), request)
}

// CreateSuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 创建播放器配置，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateSuperPlayerConfigWithContext(ctx context.Context, request *CreateSuperPlayerConfigRequest) (response *CreateSuperPlayerConfigResponse, err error) {
    if request == nil {
        request = NewCreateSuperPlayerConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateSuperPlayerConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateSuperPlayerConfigResponse()
    err = c.Send(request, response)
    return
}

func NewCreateTranscodeTemplateRequest() (request *CreateTranscodeTemplateRequest) {
    request = &CreateTranscodeTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateTranscodeTemplate")
    
    
    return
}

func NewCreateTranscodeTemplateResponse() (response *CreateTranscodeTemplateResponse) {
    response = &CreateTranscodeTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateTranscodeTemplate
// 创建用户自定义转码模板，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_AUDIOBITRATE = "InvalidParameterValue.AudioBitrate"
//  INVALIDPARAMETERVALUE_AUDIOCHANNEL = "InvalidParameterValue.AudioChannel"
//  INVALIDPARAMETERVALUE_AUDIOCODEC = "InvalidParameterValue.AudioCodec"
//  INVALIDPARAMETERVALUE_AUDIOSAMPLERATE = "InvalidParameterValue.AudioSampleRate"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_CONTAINER = "InvalidParameterValue.Container"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_REMOVEVIDEO = "InvalidParameterValue.RemoveVideo"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_RESOLUTIONADAPTIVE = "InvalidParameterValue.ResolutionAdaptive"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_VCRF = "InvalidParameterValue.Vcrf"
//  INVALIDPARAMETERVALUE_VIDEOBITRATE = "InvalidParameterValue.VideoBitrate"
//  INVALIDPARAMETERVALUE_VIDEOCODEC = "InvalidParameterValue.VideoCodec"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateTranscodeTemplate(request *CreateTranscodeTemplateRequest) (response *CreateTranscodeTemplateResponse, err error) {
    return c.CreateTranscodeTemplateWithContext(context.Background(), request)
}

// CreateTranscodeTemplate
// 创建用户自定义转码模板，数量上限：100。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_AUDIOBITRATE = "InvalidParameterValue.AudioBitrate"
//  INVALIDPARAMETERVALUE_AUDIOCHANNEL = "InvalidParameterValue.AudioChannel"
//  INVALIDPARAMETERVALUE_AUDIOCODEC = "InvalidParameterValue.AudioCodec"
//  INVALIDPARAMETERVALUE_AUDIOSAMPLERATE = "InvalidParameterValue.AudioSampleRate"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_CONTAINER = "InvalidParameterValue.Container"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_REMOVEVIDEO = "InvalidParameterValue.RemoveVideo"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_RESOLUTIONADAPTIVE = "InvalidParameterValue.ResolutionAdaptive"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_VCRF = "InvalidParameterValue.Vcrf"
//  INVALIDPARAMETERVALUE_VIDEOBITRATE = "InvalidParameterValue.VideoBitrate"
//  INVALIDPARAMETERVALUE_VIDEOCODEC = "InvalidParameterValue.VideoCodec"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateTranscodeTemplateWithContext(ctx context.Context, request *CreateTranscodeTemplateRequest) (response *CreateTranscodeTemplateResponse, err error) {
    if request == nil {
        request = NewCreateTranscodeTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateTranscodeTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateTranscodeTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateVodDomainRequest() (request *CreateVodDomainRequest) {
    request = &CreateVodDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateVodDomain")
    
    
    return
}

func NewCreateVodDomainResponse() (response *CreateVodDomainResponse) {
    response = &CreateVodDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateVodDomain
// 该接口用于将加速域名添加到点播，一个用户最多添加20个加速域名。
//
// 1.域名添加成功后点播会进行域名的部署，域名由部署状态变为在线状态大概需要2分钟的时间。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDACCOUNT = "FailedOperation.InvalidAccount"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DOMAINNAMEINBLACKLIST = "InvalidParameterValue.DomainNameInBlackList"
func (c *Client) CreateVodDomain(request *CreateVodDomainRequest) (response *CreateVodDomainResponse, err error) {
    return c.CreateVodDomainWithContext(context.Background(), request)
}

// CreateVodDomain
// 该接口用于将加速域名添加到点播，一个用户最多添加20个加速域名。
//
// 1.域名添加成功后点播会进行域名的部署，域名由部署状态变为在线状态大概需要2分钟的时间。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDACCOUNT = "FailedOperation.InvalidAccount"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DOMAINNAMEINBLACKLIST = "InvalidParameterValue.DomainNameInBlackList"
func (c *Client) CreateVodDomainWithContext(ctx context.Context, request *CreateVodDomainRequest) (response *CreateVodDomainResponse, err error) {
    if request == nil {
        request = NewCreateVodDomainRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateVodDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateVodDomainResponse()
    err = c.Send(request, response)
    return
}

func NewCreateWatermarkTemplateRequest() (request *CreateWatermarkTemplateRequest) {
    request = &CreateWatermarkTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateWatermarkTemplate")
    
    
    return
}

func NewCreateWatermarkTemplateResponse() (response *CreateWatermarkTemplateResponse) {
    response = &CreateWatermarkTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateWatermarkTemplate
// 创建用户自定义水印模板，数量上限：1000。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INTERNALERROR_UPLOADWATERMARKERROR = "InternalError.UploadWatermarkError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COORDINATEORIGIN = "InvalidParameterValue.CoordinateOrigin"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_IMAGECONTENT = "InvalidParameterValue.ImageContent"
//  INVALIDPARAMETERVALUE_IMAGETEMPLATE = "InvalidParameterValue.ImageTemplate"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REPEATTYPE = "InvalidParameterValue.RepeatType"
//  INVALIDPARAMETERVALUE_SVGTEMPLATE = "InvalidParameterValue.SvgTemplate"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEHEIGHT = "InvalidParameterValue.SvgTemplateHeight"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEWIDTH = "InvalidParameterValue.SvgTemplateWidth"
//  INVALIDPARAMETERVALUE_TEXTALPHA = "InvalidParameterValue.TextAlpha"
//  INVALIDPARAMETERVALUE_TEXTTEMPLATE = "InvalidParameterValue.TextTemplate"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  INVALIDPARAMETERVALUE_XPOS = "InvalidParameterValue.XPos"
//  INVALIDPARAMETERVALUE_YPOS = "InvalidParameterValue.YPos"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateWatermarkTemplate(request *CreateWatermarkTemplateRequest) (response *CreateWatermarkTemplateResponse, err error) {
    return c.CreateWatermarkTemplateWithContext(context.Background(), request)
}

// CreateWatermarkTemplate
// 创建用户自定义水印模板，数量上限：1000。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"
//  INTERNALERROR_UPLOADWATERMARKERROR = "InternalError.UploadWatermarkError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COORDINATEORIGIN = "InvalidParameterValue.CoordinateOrigin"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_IMAGECONTENT = "InvalidParameterValue.ImageContent"
//  INVALIDPARAMETERVALUE_IMAGETEMPLATE = "InvalidParameterValue.ImageTemplate"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REPEATTYPE = "InvalidParameterValue.RepeatType"
//  INVALIDPARAMETERVALUE_SVGTEMPLATE = "InvalidParameterValue.SvgTemplate"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEHEIGHT = "InvalidParameterValue.SvgTemplateHeight"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEWIDTH = "InvalidParameterValue.SvgTemplateWidth"
//  INVALIDPARAMETERVALUE_TEXTALPHA = "InvalidParameterValue.TextAlpha"
//  INVALIDPARAMETERVALUE_TEXTTEMPLATE = "InvalidParameterValue.TextTemplate"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  INVALIDPARAMETERVALUE_XPOS = "InvalidParameterValue.XPos"
//  INVALIDPARAMETERVALUE_YPOS = "InvalidParameterValue.YPos"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateWatermarkTemplateWithContext(ctx context.Context, request *CreateWatermarkTemplateRequest) (response *CreateWatermarkTemplateResponse, err error) {
    if request == nil {
        request = NewCreateWatermarkTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateWatermarkTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateWatermarkTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateWordSamplesRequest() (request *CreateWordSamplesRequest) {
    request = &CreateWordSamplesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "CreateWordSamples")
    
    
    return
}

func NewCreateWordSamplesResponse() (response *CreateWordSamplesResponse) {
    response = &CreateWordSamplesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateWordSamples
// 该接口用于批量创建关键词样本，样本用于通过OCR、ASR技术，进行不适宜内容识别、内容识别等视频处理。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateWordSamples(request *CreateWordSamplesRequest) (response *CreateWordSamplesResponse, err error) {
    return c.CreateWordSamplesWithContext(context.Background(), request)
}

// CreateWordSamples
// 该接口用于批量创建关键词样本，样本用于通过OCR、ASR技术，进行不适宜内容识别、内容识别等视频处理。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateWordSamplesWithContext(ctx context.Context, request *CreateWordSamplesRequest) (response *CreateWordSamplesResponse, err error) {
    if request == nil {
        request = NewCreateWordSamplesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateWordSamples require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateWordSamplesResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAIAnalysisTemplateRequest() (request *DeleteAIAnalysisTemplateRequest) {
    request = &DeleteAIAnalysisTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteAIAnalysisTemplate")
    
    
    return
}

func NewDeleteAIAnalysisTemplateResponse() (response *DeleteAIAnalysisTemplateResponse) {
    response = &DeleteAIAnalysisTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteAIAnalysisTemplate
// 删除用户自定义音视频内容分析模板。
//
// 
//
// 注意：模板 ID 为 10000 以下的为系统预置模板，不允许删除。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAIAnalysisTemplate(request *DeleteAIAnalysisTemplateRequest) (response *DeleteAIAnalysisTemplateResponse, err error) {
    return c.DeleteAIAnalysisTemplateWithContext(context.Background(), request)
}

// DeleteAIAnalysisTemplate
// 删除用户自定义音视频内容分析模板。
//
// 
//
// 注意：模板 ID 为 10000 以下的为系统预置模板，不允许删除。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAIAnalysisTemplateWithContext(ctx context.Context, request *DeleteAIAnalysisTemplateRequest) (response *DeleteAIAnalysisTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteAIAnalysisTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteAIAnalysisTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteAIAnalysisTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAIRecognitionTemplateRequest() (request *DeleteAIRecognitionTemplateRequest) {
    request = &DeleteAIRecognitionTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteAIRecognitionTemplate")
    
    
    return
}

func NewDeleteAIRecognitionTemplateResponse() (response *DeleteAIRecognitionTemplateResponse) {
    response = &DeleteAIRecognitionTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteAIRecognitionTemplate
// 删除用户自定义音视频内容识别模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAIRecognitionTemplate(request *DeleteAIRecognitionTemplateRequest) (response *DeleteAIRecognitionTemplateResponse, err error) {
    return c.DeleteAIRecognitionTemplateWithContext(context.Background(), request)
}

// DeleteAIRecognitionTemplate
// 删除用户自定义音视频内容识别模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAIRecognitionTemplateWithContext(ctx context.Context, request *DeleteAIRecognitionTemplateRequest) (response *DeleteAIRecognitionTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteAIRecognitionTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteAIRecognitionTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteAIRecognitionTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAdaptiveDynamicStreamingTemplateRequest() (request *DeleteAdaptiveDynamicStreamingTemplateRequest) {
    request = &DeleteAdaptiveDynamicStreamingTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteAdaptiveDynamicStreamingTemplate")
    
    
    return
}

func NewDeleteAdaptiveDynamicStreamingTemplateResponse() (response *DeleteAdaptiveDynamicStreamingTemplateResponse) {
    response = &DeleteAdaptiveDynamicStreamingTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteAdaptiveDynamicStreamingTemplate
// 删除转自适应码流模板
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
func (c *Client) DeleteAdaptiveDynamicStreamingTemplate(request *DeleteAdaptiveDynamicStreamingTemplateRequest) (response *DeleteAdaptiveDynamicStreamingTemplateResponse, err error) {
    return c.DeleteAdaptiveDynamicStreamingTemplateWithContext(context.Background(), request)
}

// DeleteAdaptiveDynamicStreamingTemplate
// 删除转自适应码流模板
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
func (c *Client) DeleteAdaptiveDynamicStreamingTemplateWithContext(ctx context.Context, request *DeleteAdaptiveDynamicStreamingTemplateRequest) (response *DeleteAdaptiveDynamicStreamingTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteAdaptiveDynamicStreamingTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteAdaptiveDynamicStreamingTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteAdaptiveDynamicStreamingTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAnimatedGraphicsTemplateRequest() (request *DeleteAnimatedGraphicsTemplateRequest) {
    request = &DeleteAnimatedGraphicsTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteAnimatedGraphicsTemplate")
    
    
    return
}

func NewDeleteAnimatedGraphicsTemplateResponse() (response *DeleteAnimatedGraphicsTemplateResponse) {
    response = &DeleteAnimatedGraphicsTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteAnimatedGraphicsTemplate
// 删除用户自定义转动图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAnimatedGraphicsTemplate(request *DeleteAnimatedGraphicsTemplateRequest) (response *DeleteAnimatedGraphicsTemplateResponse, err error) {
    return c.DeleteAnimatedGraphicsTemplateWithContext(context.Background(), request)
}

// DeleteAnimatedGraphicsTemplate
// 删除用户自定义转动图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteAnimatedGraphicsTemplateWithContext(ctx context.Context, request *DeleteAnimatedGraphicsTemplateRequest) (response *DeleteAnimatedGraphicsTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteAnimatedGraphicsTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteAnimatedGraphicsTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteAnimatedGraphicsTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteClassRequest() (request *DeleteClassRequest) {
    request = &DeleteClassRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteClass")
    
    
    return
}

func NewDeleteClassResponse() (response *DeleteClassResponse) {
    response = &DeleteClassResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteClass
// * 仅当待删分类无子分类且无媒体关联情况下，可删除分类；
//
// * 否则，请先执行[删除媒体](/document/product/266/31764)及子分类，再删除该分类；
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSID = "InvalidParameterValue.ClassId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION_CLASSNOTEMPTY = "UnsupportedOperation.ClassNotEmpty"
func (c *Client) DeleteClass(request *DeleteClassRequest) (response *DeleteClassResponse, err error) {
    return c.DeleteClassWithContext(context.Background(), request)
}

// DeleteClass
// * 仅当待删分类无子分类且无媒体关联情况下，可删除分类；
//
// * 否则，请先执行[删除媒体](/document/product/266/31764)及子分类，再删除该分类；
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSID = "InvalidParameterValue.ClassId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION_CLASSNOTEMPTY = "UnsupportedOperation.ClassNotEmpty"
func (c *Client) DeleteClassWithContext(ctx context.Context, request *DeleteClassRequest) (response *DeleteClassResponse, err error) {
    if request == nil {
        request = NewDeleteClassRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteClass require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteClassResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteContentReviewTemplateRequest() (request *DeleteContentReviewTemplateRequest) {
    request = &DeleteContentReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteContentReviewTemplate")
    
    
    return
}

func NewDeleteContentReviewTemplateResponse() (response *DeleteContentReviewTemplateResponse) {
    response = &DeleteContentReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [删除审核模板](https://cloud.tencent.com/document/api/266/84390)。
//
// 删除用户自定义音视频内容审核模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteContentReviewTemplate(request *DeleteContentReviewTemplateRequest) (response *DeleteContentReviewTemplateResponse, err error) {
    return c.DeleteContentReviewTemplateWithContext(context.Background(), request)
}

// DeleteContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [删除审核模板](https://cloud.tencent.com/document/api/266/84390)。
//
// 删除用户自定义音视频内容审核模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteContentReviewTemplateWithContext(ctx context.Context, request *DeleteContentReviewTemplateRequest) (response *DeleteContentReviewTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteContentReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteContentReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteContentReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteHeadTailTemplateRequest() (request *DeleteHeadTailTemplateRequest) {
    request = &DeleteHeadTailTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteHeadTailTemplate")
    
    
    return
}

func NewDeleteHeadTailTemplateResponse() (response *DeleteHeadTailTemplateResponse) {
    response = &DeleteHeadTailTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteHeadTailTemplate
// 删除片头片尾模板。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteHeadTailTemplate(request *DeleteHeadTailTemplateRequest) (response *DeleteHeadTailTemplateResponse, err error) {
    return c.DeleteHeadTailTemplateWithContext(context.Background(), request)
}

// DeleteHeadTailTemplate
// 删除片头片尾模板。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteHeadTailTemplateWithContext(ctx context.Context, request *DeleteHeadTailTemplateRequest) (response *DeleteHeadTailTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteHeadTailTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteHeadTailTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteHeadTailTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteImageProcessingTemplateRequest() (request *DeleteImageProcessingTemplateRequest) {
    request = &DeleteImageProcessingTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteImageProcessingTemplate")
    
    
    return
}

func NewDeleteImageProcessingTemplateResponse() (response *DeleteImageProcessingTemplateResponse) {
    response = &DeleteImageProcessingTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteImageProcessingTemplate
// 删除用户自定义图片处理模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteImageProcessingTemplate(request *DeleteImageProcessingTemplateRequest) (response *DeleteImageProcessingTemplateResponse, err error) {
    return c.DeleteImageProcessingTemplateWithContext(context.Background(), request)
}

// DeleteImageProcessingTemplate
// 删除用户自定义图片处理模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteImageProcessingTemplateWithContext(ctx context.Context, request *DeleteImageProcessingTemplateRequest) (response *DeleteImageProcessingTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteImageProcessingTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteImageProcessingTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteImageProcessingTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteImageSpriteTemplateRequest() (request *DeleteImageSpriteTemplateRequest) {
    request = &DeleteImageSpriteTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteImageSpriteTemplate")
    
    
    return
}

func NewDeleteImageSpriteTemplateResponse() (response *DeleteImageSpriteTemplateResponse) {
    response = &DeleteImageSpriteTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteImageSpriteTemplate
// 删除雪碧图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteImageSpriteTemplate(request *DeleteImageSpriteTemplateRequest) (response *DeleteImageSpriteTemplateResponse, err error) {
    return c.DeleteImageSpriteTemplateWithContext(context.Background(), request)
}

// DeleteImageSpriteTemplate
// 删除雪碧图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteImageSpriteTemplateWithContext(ctx context.Context, request *DeleteImageSpriteTemplateRequest) (response *DeleteImageSpriteTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteImageSpriteTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteImageSpriteTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteImageSpriteTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteMediaRequest() (request *DeleteMediaRequest) {
    request = &DeleteMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteMedia")
    
    
    return
}

func NewDeleteMediaResponse() (response *DeleteMediaResponse) {
    response = &DeleteMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteMedia
// * 删除媒体及其对应的视频处理文件（原始文件、如转码视频、雪碧图、截图、微信发布视频等）；
//
// * 可单独删除指定 ID 的视频文件下的原文件、转码视频、微信发布视频等；
//
// * 注意：原文件删除后，无法发起转码、微信发布等任何视频处理操作。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteMedia(request *DeleteMediaRequest) (response *DeleteMediaResponse, err error) {
    return c.DeleteMediaWithContext(context.Background(), request)
}

// DeleteMedia
// * 删除媒体及其对应的视频处理文件（原始文件、如转码视频、雪碧图、截图、微信发布视频等）；
//
// * 可单独删除指定 ID 的视频文件下的原文件、转码视频、微信发布视频等；
//
// * 注意：原文件删除后，无法发起转码、微信发布等任何视频处理操作。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteMediaWithContext(ctx context.Context, request *DeleteMediaRequest) (response *DeleteMediaResponse, err error) {
    if request == nil {
        request = NewDeleteMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteMediaResponse()
    err = c.Send(request, response)
    return
}

func NewDeletePersonSampleRequest() (request *DeletePersonSampleRequest) {
    request = &DeletePersonSampleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeletePersonSample")
    
    
    return
}

func NewDeletePersonSampleResponse() (response *DeletePersonSampleResponse) {
    response = &DeletePersonSampleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeletePersonSample
// 该接口用于根据人物 ID，删除素材样本。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_PERSON = "ResourceNotFound.Person"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeletePersonSample(request *DeletePersonSampleRequest) (response *DeletePersonSampleResponse, err error) {
    return c.DeletePersonSampleWithContext(context.Background(), request)
}

// DeletePersonSample
// 该接口用于根据人物 ID，删除素材样本。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_PERSON = "ResourceNotFound.Person"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeletePersonSampleWithContext(ctx context.Context, request *DeletePersonSampleRequest) (response *DeletePersonSampleResponse, err error) {
    if request == nil {
        request = NewDeletePersonSampleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeletePersonSample require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeletePersonSampleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteProcedureTemplateRequest() (request *DeleteProcedureTemplateRequest) {
    request = &DeleteProcedureTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteProcedureTemplate")
    
    
    return
}

func NewDeleteProcedureTemplateResponse() (response *DeleteProcedureTemplateResponse) {
    response = &DeleteProcedureTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteProcedureTemplate
// 删除用户自定义的任务流模板。  
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteProcedureTemplate(request *DeleteProcedureTemplateRequest) (response *DeleteProcedureTemplateResponse, err error) {
    return c.DeleteProcedureTemplateWithContext(context.Background(), request)
}

// DeleteProcedureTemplate
// 删除用户自定义的任务流模板。  
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteProcedureTemplateWithContext(ctx context.Context, request *DeleteProcedureTemplateRequest) (response *DeleteProcedureTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteProcedureTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteProcedureTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteProcedureTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteReviewTemplateRequest() (request *DeleteReviewTemplateRequest) {
    request = &DeleteReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteReviewTemplate")
    
    
    return
}

func NewDeleteReviewTemplateResponse() (response *DeleteReviewTemplateResponse) {
    response = &DeleteReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteReviewTemplate
// 删除用户自定义审核模板。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteReviewTemplate(request *DeleteReviewTemplateRequest) (response *DeleteReviewTemplateResponse, err error) {
    return c.DeleteReviewTemplateWithContext(context.Background(), request)
}

// DeleteReviewTemplate
// 删除用户自定义审核模板。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteReviewTemplateWithContext(ctx context.Context, request *DeleteReviewTemplateRequest) (response *DeleteReviewTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSampleSnapshotTemplateRequest() (request *DeleteSampleSnapshotTemplateRequest) {
    request = &DeleteSampleSnapshotTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteSampleSnapshotTemplate")
    
    
    return
}

func NewDeleteSampleSnapshotTemplateResponse() (response *DeleteSampleSnapshotTemplateResponse) {
    response = &DeleteSampleSnapshotTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteSampleSnapshotTemplate
// 删除用户自定义采样截图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSampleSnapshotTemplate(request *DeleteSampleSnapshotTemplateRequest) (response *DeleteSampleSnapshotTemplateResponse, err error) {
    return c.DeleteSampleSnapshotTemplateWithContext(context.Background(), request)
}

// DeleteSampleSnapshotTemplate
// 删除用户自定义采样截图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSampleSnapshotTemplateWithContext(ctx context.Context, request *DeleteSampleSnapshotTemplateRequest) (response *DeleteSampleSnapshotTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteSampleSnapshotTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteSampleSnapshotTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteSampleSnapshotTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSnapshotByTimeOffsetTemplateRequest() (request *DeleteSnapshotByTimeOffsetTemplateRequest) {
    request = &DeleteSnapshotByTimeOffsetTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteSnapshotByTimeOffsetTemplate")
    
    
    return
}

func NewDeleteSnapshotByTimeOffsetTemplateResponse() (response *DeleteSnapshotByTimeOffsetTemplateResponse) {
    response = &DeleteSnapshotByTimeOffsetTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteSnapshotByTimeOffsetTemplate
// 删除用户自定义指定时间点截图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSnapshotByTimeOffsetTemplate(request *DeleteSnapshotByTimeOffsetTemplateRequest) (response *DeleteSnapshotByTimeOffsetTemplateResponse, err error) {
    return c.DeleteSnapshotByTimeOffsetTemplateWithContext(context.Background(), request)
}

// DeleteSnapshotByTimeOffsetTemplate
// 删除用户自定义指定时间点截图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSnapshotByTimeOffsetTemplateWithContext(ctx context.Context, request *DeleteSnapshotByTimeOffsetTemplateRequest) (response *DeleteSnapshotByTimeOffsetTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteSnapshotByTimeOffsetTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteSnapshotByTimeOffsetTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteSnapshotByTimeOffsetTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSuperPlayerConfigRequest() (request *DeleteSuperPlayerConfigRequest) {
    request = &DeleteSuperPlayerConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteSuperPlayerConfig")
    
    
    return
}

func NewDeleteSuperPlayerConfigResponse() (response *DeleteSuperPlayerConfigResponse) {
    response = &DeleteSuperPlayerConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteSuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 删除播放器配置。  
//
// *注：系统预置播放器配置不允许删除。*
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSuperPlayerConfig(request *DeleteSuperPlayerConfigRequest) (response *DeleteSuperPlayerConfigResponse, err error) {
    return c.DeleteSuperPlayerConfigWithContext(context.Background(), request)
}

// DeleteSuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 删除播放器配置。  
//
// *注：系统预置播放器配置不允许删除。*
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteSuperPlayerConfigWithContext(ctx context.Context, request *DeleteSuperPlayerConfigRequest) (response *DeleteSuperPlayerConfigResponse, err error) {
    if request == nil {
        request = NewDeleteSuperPlayerConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteSuperPlayerConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteSuperPlayerConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteTranscodeTemplateRequest() (request *DeleteTranscodeTemplateRequest) {
    request = &DeleteTranscodeTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteTranscodeTemplate")
    
    
    return
}

func NewDeleteTranscodeTemplateResponse() (response *DeleteTranscodeTemplateResponse) {
    response = &DeleteTranscodeTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteTranscodeTemplate
// 删除用户自定义转码模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteTranscodeTemplate(request *DeleteTranscodeTemplateRequest) (response *DeleteTranscodeTemplateResponse, err error) {
    return c.DeleteTranscodeTemplateWithContext(context.Background(), request)
}

// DeleteTranscodeTemplate
// 删除用户自定义转码模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteTranscodeTemplateWithContext(ctx context.Context, request *DeleteTranscodeTemplateRequest) (response *DeleteTranscodeTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteTranscodeTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteTranscodeTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteTranscodeTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteVodDomainRequest() (request *DeleteVodDomainRequest) {
    request = &DeleteVodDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteVodDomain")
    
    
    return
}

func NewDeleteVodDomainResponse() (response *DeleteVodDomainResponse) {
    response = &DeleteVodDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteVodDomain
// 该接口用于删除点播加速域名。
//
// 1、域名删除前需要先关闭所有区域的加速。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DeleteVodDomain(request *DeleteVodDomainRequest) (response *DeleteVodDomainResponse, err error) {
    return c.DeleteVodDomainWithContext(context.Background(), request)
}

// DeleteVodDomain
// 该接口用于删除点播加速域名。
//
// 1、域名删除前需要先关闭所有区域的加速。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DeleteVodDomainWithContext(ctx context.Context, request *DeleteVodDomainRequest) (response *DeleteVodDomainResponse, err error) {
    if request == nil {
        request = NewDeleteVodDomainRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteVodDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteVodDomainResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteWatermarkTemplateRequest() (request *DeleteWatermarkTemplateRequest) {
    request = &DeleteWatermarkTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteWatermarkTemplate")
    
    
    return
}

func NewDeleteWatermarkTemplateResponse() (response *DeleteWatermarkTemplateResponse) {
    response = &DeleteWatermarkTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteWatermarkTemplate
// 删除用户自定义水印模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteWatermarkTemplate(request *DeleteWatermarkTemplateRequest) (response *DeleteWatermarkTemplateResponse, err error) {
    return c.DeleteWatermarkTemplateWithContext(context.Background(), request)
}

// DeleteWatermarkTemplate
// 删除用户自定义水印模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteWatermarkTemplateWithContext(ctx context.Context, request *DeleteWatermarkTemplateRequest) (response *DeleteWatermarkTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteWatermarkTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteWatermarkTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteWatermarkTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteWordSamplesRequest() (request *DeleteWordSamplesRequest) {
    request = &DeleteWordSamplesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DeleteWordSamples")
    
    
    return
}

func NewDeleteWordSamplesResponse() (response *DeleteWordSamplesResponse) {
    response = &DeleteWordSamplesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteWordSamples
// 该接口用于批量删除关键词样本。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteWordSamples(request *DeleteWordSamplesRequest) (response *DeleteWordSamplesResponse, err error) {
    return c.DeleteWordSamplesWithContext(context.Background(), request)
}

// DeleteWordSamples
// 该接口用于批量删除关键词样本。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteWordSamplesWithContext(ctx context.Context, request *DeleteWordSamplesRequest) (response *DeleteWordSamplesResponse, err error) {
    if request == nil {
        request = NewDeleteWordSamplesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteWordSamples require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteWordSamplesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAIAnalysisTemplatesRequest() (request *DescribeAIAnalysisTemplatesRequest) {
    request = &DescribeAIAnalysisTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeAIAnalysisTemplates")
    
    
    return
}

func NewDescribeAIAnalysisTemplatesResponse() (response *DescribeAIAnalysisTemplatesResponse) {
    response = &DescribeAIAnalysisTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeAIAnalysisTemplates
// 根据音视频内容分析模板唯一标识，获取音视频内容分析模板详情列表。返回结果包含符合条件的所有用户自定义音视频内容分析模板及[系统预置音视频内容分析模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E5.88.86.E6.9E.90.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAIAnalysisTemplates(request *DescribeAIAnalysisTemplatesRequest) (response *DescribeAIAnalysisTemplatesResponse, err error) {
    return c.DescribeAIAnalysisTemplatesWithContext(context.Background(), request)
}

// DescribeAIAnalysisTemplates
// 根据音视频内容分析模板唯一标识，获取音视频内容分析模板详情列表。返回结果包含符合条件的所有用户自定义音视频内容分析模板及[系统预置音视频内容分析模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E5.88.86.E6.9E.90.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAIAnalysisTemplatesWithContext(ctx context.Context, request *DescribeAIAnalysisTemplatesRequest) (response *DescribeAIAnalysisTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeAIAnalysisTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeAIAnalysisTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeAIAnalysisTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAIRecognitionTemplatesRequest() (request *DescribeAIRecognitionTemplatesRequest) {
    request = &DescribeAIRecognitionTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeAIRecognitionTemplates")
    
    
    return
}

func NewDescribeAIRecognitionTemplatesResponse() (response *DescribeAIRecognitionTemplatesResponse) {
    response = &DescribeAIRecognitionTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeAIRecognitionTemplates
// 根据音视频内容识别模板唯一标识，获取音视频内容识别模板详情列表。返回结果包含符合条件的所有用户自定义音视频内容识别模板及[系统预置音视频内容识别模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E8.AF.86.E5.88.AB.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAIRecognitionTemplates(request *DescribeAIRecognitionTemplatesRequest) (response *DescribeAIRecognitionTemplatesResponse, err error) {
    return c.DescribeAIRecognitionTemplatesWithContext(context.Background(), request)
}

// DescribeAIRecognitionTemplates
// 根据音视频内容识别模板唯一标识，获取音视频内容识别模板详情列表。返回结果包含符合条件的所有用户自定义音视频内容识别模板及[系统预置音视频内容识别模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E8.AF.86.E5.88.AB.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAIRecognitionTemplatesWithContext(ctx context.Context, request *DescribeAIRecognitionTemplatesRequest) (response *DescribeAIRecognitionTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeAIRecognitionTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeAIRecognitionTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeAIRecognitionTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAdaptiveDynamicStreamingTemplatesRequest() (request *DescribeAdaptiveDynamicStreamingTemplatesRequest) {
    request = &DescribeAdaptiveDynamicStreamingTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeAdaptiveDynamicStreamingTemplates")
    
    
    return
}

func NewDescribeAdaptiveDynamicStreamingTemplatesResponse() (response *DescribeAdaptiveDynamicStreamingTemplatesResponse) {
    response = &DescribeAdaptiveDynamicStreamingTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeAdaptiveDynamicStreamingTemplates
// 查询转自适应码流模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAdaptiveDynamicStreamingTemplates(request *DescribeAdaptiveDynamicStreamingTemplatesRequest) (response *DescribeAdaptiveDynamicStreamingTemplatesResponse, err error) {
    return c.DescribeAdaptiveDynamicStreamingTemplatesWithContext(context.Background(), request)
}

// DescribeAdaptiveDynamicStreamingTemplates
// 查询转自适应码流模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAdaptiveDynamicStreamingTemplatesWithContext(ctx context.Context, request *DescribeAdaptiveDynamicStreamingTemplatesRequest) (response *DescribeAdaptiveDynamicStreamingTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeAdaptiveDynamicStreamingTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeAdaptiveDynamicStreamingTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeAdaptiveDynamicStreamingTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAllClassRequest() (request *DescribeAllClassRequest) {
    request = &DescribeAllClassRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeAllClass")
    
    
    return
}

func NewDescribeAllClassResponse() (response *DescribeAllClassResponse) {
    response = &DescribeAllClassResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeAllClass
// * 获得用户的所有分类信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAllClass(request *DescribeAllClassRequest) (response *DescribeAllClassResponse, err error) {
    return c.DescribeAllClassWithContext(context.Background(), request)
}

// DescribeAllClass
// * 获得用户的所有分类信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAllClassWithContext(ctx context.Context, request *DescribeAllClassRequest) (response *DescribeAllClassResponse, err error) {
    if request == nil {
        request = NewDescribeAllClassRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeAllClass require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeAllClassResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAnimatedGraphicsTemplatesRequest() (request *DescribeAnimatedGraphicsTemplatesRequest) {
    request = &DescribeAnimatedGraphicsTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeAnimatedGraphicsTemplates")
    
    
    return
}

func NewDescribeAnimatedGraphicsTemplatesResponse() (response *DescribeAnimatedGraphicsTemplatesResponse) {
    response = &DescribeAnimatedGraphicsTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeAnimatedGraphicsTemplates
// 查询转动图模板列表，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAnimatedGraphicsTemplates(request *DescribeAnimatedGraphicsTemplatesRequest) (response *DescribeAnimatedGraphicsTemplatesResponse, err error) {
    return c.DescribeAnimatedGraphicsTemplatesWithContext(context.Background(), request)
}

// DescribeAnimatedGraphicsTemplates
// 查询转动图模板列表，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeAnimatedGraphicsTemplatesWithContext(ctx context.Context, request *DescribeAnimatedGraphicsTemplatesRequest) (response *DescribeAnimatedGraphicsTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeAnimatedGraphicsTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeAnimatedGraphicsTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeAnimatedGraphicsTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCDNStatDetailsRequest() (request *DescribeCDNStatDetailsRequest) {
    request = &DescribeCDNStatDetailsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeCDNStatDetails")
    
    
    return
}

func NewDescribeCDNStatDetailsResponse() (response *DescribeCDNStatDetailsResponse) {
    response = &DescribeCDNStatDetailsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeCDNStatDetails
// 该接口用于查询点播域名的 CDN 带宽、流量等统计数据。
//
// * 查询的起始时间和结束时间跨度不超过90天。
//
// * 可以查询不同服务区域的数据。
//
// * 中国境内的数据支持查询指定地区、运营商的统计数据。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDACCOUNT = "FailedOperation.InvalidAccount"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETERVALUE_AREA = "InvalidParameterValue.Area"
//  INVALIDPARAMETERVALUE_DISTRICTS = "InvalidParameterValue.Districts"
//  INVALIDPARAMETERVALUE_DOMAINNAMES = "InvalidParameterValue.DomainNames"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_ISPS = "InvalidParameterValue.Isps"
//  INVALIDPARAMETERVALUE_METRIC = "InvalidParameterValue.Metric"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeCDNStatDetails(request *DescribeCDNStatDetailsRequest) (response *DescribeCDNStatDetailsResponse, err error) {
    return c.DescribeCDNStatDetailsWithContext(context.Background(), request)
}

// DescribeCDNStatDetails
// 该接口用于查询点播域名的 CDN 带宽、流量等统计数据。
//
// * 查询的起始时间和结束时间跨度不超过90天。
//
// * 可以查询不同服务区域的数据。
//
// * 中国境内的数据支持查询指定地区、运营商的统计数据。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDACCOUNT = "FailedOperation.InvalidAccount"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETERVALUE_AREA = "InvalidParameterValue.Area"
//  INVALIDPARAMETERVALUE_DISTRICTS = "InvalidParameterValue.Districts"
//  INVALIDPARAMETERVALUE_DOMAINNAMES = "InvalidParameterValue.DomainNames"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_ISPS = "InvalidParameterValue.Isps"
//  INVALIDPARAMETERVALUE_METRIC = "InvalidParameterValue.Metric"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeCDNStatDetailsWithContext(ctx context.Context, request *DescribeCDNStatDetailsRequest) (response *DescribeCDNStatDetailsResponse, err error) {
    if request == nil {
        request = NewDescribeCDNStatDetailsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeCDNStatDetails require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeCDNStatDetailsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCDNUsageDataRequest() (request *DescribeCDNUsageDataRequest) {
    request = &DescribeCDNUsageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeCDNUsageData")
    
    
    return
}

func NewDescribeCDNUsageDataResponse() (response *DescribeCDNUsageDataResponse) {
    response = &DescribeCDNUsageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeCDNUsageData
// 该接口用于查询点播 CDN 的流量、带宽等统计数据。
//
//    1. 可以查询最近365天内的 CDN 用量数据。
//
//    2.  查询时间跨度不超过90天。
//
//    3. 可以指定用量数据的时间粒度，支持5分钟、1小时、1天的时间粒度。
//
//    4.  流量为查询时间粒度内的总流量，带宽为查询时间粒度内的峰值带宽。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_DATAINTERVAL = "InvalidParameterValue.DataInterval"
//  INVALIDPARAMETERVALUE_DATATYPE = "InvalidParameterValue.DataType"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_DOMAINNAMES = "InvalidParameterValue.DomainNames"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeCDNUsageData(request *DescribeCDNUsageDataRequest) (response *DescribeCDNUsageDataResponse, err error) {
    return c.DescribeCDNUsageDataWithContext(context.Background(), request)
}

// DescribeCDNUsageData
// 该接口用于查询点播 CDN 的流量、带宽等统计数据。
//
//    1. 可以查询最近365天内的 CDN 用量数据。
//
//    2.  查询时间跨度不超过90天。
//
//    3. 可以指定用量数据的时间粒度，支持5分钟、1小时、1天的时间粒度。
//
//    4.  流量为查询时间粒度内的总流量，带宽为查询时间粒度内的峰值带宽。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_DATAINTERVAL = "InvalidParameterValue.DataInterval"
//  INVALIDPARAMETERVALUE_DATATYPE = "InvalidParameterValue.DataType"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_DOMAINNAMES = "InvalidParameterValue.DomainNames"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeCDNUsageDataWithContext(ctx context.Context, request *DescribeCDNUsageDataRequest) (response *DescribeCDNUsageDataResponse, err error) {
    if request == nil {
        request = NewDescribeCDNUsageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeCDNUsageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeCDNUsageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCdnLogsRequest() (request *DescribeCdnLogsRequest) {
    request = &DescribeCdnLogsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeCdnLogs")
    
    
    return
}

func NewDescribeCdnLogsResponse() (response *DescribeCdnLogsResponse) {
    response = &DescribeCdnLogsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeCdnLogs
// 查询点播域名的 CDN 访问日志的下载链接。
//
//     1. 可以查询最近30天内的 CDN 日志下载链接。
//
//     2. 默认情况下 CDN 每小时生成一个日志文件，如果某一个小时没有 CDN 访问，不会生成日志文件。    
//
//     3. CDN 日志下载链接的有效期为24小时。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeCdnLogs(request *DescribeCdnLogsRequest) (response *DescribeCdnLogsResponse, err error) {
    return c.DescribeCdnLogsWithContext(context.Background(), request)
}

// DescribeCdnLogs
// 查询点播域名的 CDN 访问日志的下载链接。
//
//     1. 可以查询最近30天内的 CDN 日志下载链接。
//
//     2. 默认情况下 CDN 每小时生成一个日志文件，如果某一个小时没有 CDN 访问，不会生成日志文件。    
//
//     3. CDN 日志下载链接的有效期为24小时。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeCdnLogsWithContext(ctx context.Context, request *DescribeCdnLogsRequest) (response *DescribeCdnLogsResponse, err error) {
    if request == nil {
        request = NewDescribeCdnLogsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeCdnLogs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeCdnLogsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeClientUploadAccelerationUsageDataRequest() (request *DescribeClientUploadAccelerationUsageDataRequest) {
    request = &DescribeClientUploadAccelerationUsageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeClientUploadAccelerationUsageData")
    
    
    return
}

func NewDescribeClientUploadAccelerationUsageDataResponse() (response *DescribeClientUploadAccelerationUsageDataResponse) {
    response = &DescribeClientUploadAccelerationUsageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeClientUploadAccelerationUsageData
// 该接口返回查询时间范围内客户端上传加速统计信息。
//
//    1. 可以查询最近365天内的客户端上传加速统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeClientUploadAccelerationUsageData(request *DescribeClientUploadAccelerationUsageDataRequest) (response *DescribeClientUploadAccelerationUsageDataResponse, err error) {
    return c.DescribeClientUploadAccelerationUsageDataWithContext(context.Background(), request)
}

// DescribeClientUploadAccelerationUsageData
// 该接口返回查询时间范围内客户端上传加速统计信息。
//
//    1. 可以查询最近365天内的客户端上传加速统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeClientUploadAccelerationUsageDataWithContext(ctx context.Context, request *DescribeClientUploadAccelerationUsageDataRequest) (response *DescribeClientUploadAccelerationUsageDataResponse, err error) {
    if request == nil {
        request = NewDescribeClientUploadAccelerationUsageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeClientUploadAccelerationUsageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeClientUploadAccelerationUsageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeContentReviewTemplatesRequest() (request *DescribeContentReviewTemplatesRequest) {
    request = &DescribeContentReviewTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeContentReviewTemplates")
    
    
    return
}

func NewDescribeContentReviewTemplatesResponse() (response *DescribeContentReviewTemplatesResponse) {
    response = &DescribeContentReviewTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeContentReviewTemplates
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [获取审核模板列表](https://cloud.tencent.com/document/api/266/84389)。
//
// 根据音视频内容审核模板唯一标识，获取音视频内容审核模板详情列表。返回结果包含符合条件的所有用户自定义模板及[系统预置内容审核模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E5.AE.A1.E6.A0.B8.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeContentReviewTemplates(request *DescribeContentReviewTemplatesRequest) (response *DescribeContentReviewTemplatesResponse, err error) {
    return c.DescribeContentReviewTemplatesWithContext(context.Background(), request)
}

// DescribeContentReviewTemplates
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [获取审核模板列表](https://cloud.tencent.com/document/api/266/84389)。
//
// 根据音视频内容审核模板唯一标识，获取音视频内容审核模板详情列表。返回结果包含符合条件的所有用户自定义模板及[系统预置内容审核模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.A7.86.E9.A2.91.E5.86.85.E5.AE.B9.E5.AE.A1.E6.A0.B8.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeContentReviewTemplatesWithContext(ctx context.Context, request *DescribeContentReviewTemplatesRequest) (response *DescribeContentReviewTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeContentReviewTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeContentReviewTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeContentReviewTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDailyMediaPlayStatRequest() (request *DescribeDailyMediaPlayStatRequest) {
    request = &DescribeDailyMediaPlayStatRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeDailyMediaPlayStat")
    
    
    return
}

func NewDescribeDailyMediaPlayStatResponse() (response *DescribeDailyMediaPlayStatResponse) {
    response = &DescribeDailyMediaPlayStatResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeDailyMediaPlayStat
// 该接口用于查询指定日期范围内每天的播放统计数据。
//
// * 可以查询最近一年的播放统计数据。
//
// * 结束日期和起始日期的时间跨度最大为90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_ENDDATE = "InvalidParameterValue.EndDate"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_STARTDATE = "InvalidParameterValue.StartDate"
func (c *Client) DescribeDailyMediaPlayStat(request *DescribeDailyMediaPlayStatRequest) (response *DescribeDailyMediaPlayStatResponse, err error) {
    return c.DescribeDailyMediaPlayStatWithContext(context.Background(), request)
}

// DescribeDailyMediaPlayStat
// 该接口用于查询指定日期范围内每天的播放统计数据。
//
// * 可以查询最近一年的播放统计数据。
//
// * 结束日期和起始日期的时间跨度最大为90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_ENDDATE = "InvalidParameterValue.EndDate"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_STARTDATE = "InvalidParameterValue.StartDate"
func (c *Client) DescribeDailyMediaPlayStatWithContext(ctx context.Context, request *DescribeDailyMediaPlayStatRequest) (response *DescribeDailyMediaPlayStatResponse, err error) {
    if request == nil {
        request = NewDescribeDailyMediaPlayStatRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDailyMediaPlayStat require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDailyMediaPlayStatResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDailyMostPlayedStatRequest() (request *DescribeDailyMostPlayedStatRequest) {
    request = &DescribeDailyMostPlayedStatRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeDailyMostPlayedStat")
    
    
    return
}

func NewDescribeDailyMostPlayedStatResponse() (response *DescribeDailyMostPlayedStatResponse) {
    response = &DescribeDailyMostPlayedStatResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeDailyMostPlayedStat
// 该接口用于查询每日播放Top100 的媒体文件的播放统计数据。
//
// * 可以查询最近一年的播放统计数据。
//
// * 可以按播放次数或者播放流量查询。
//
// * 播放次数统计说明：
//
//     1. HLS 文件：访问 M3U8 文件时统计播放次数；访问 TS 文件不统计播放次数。
//
//     2. 其它文件（如 MP4 文件）：播放请求带有 range 参数且 range 的 start 参数不等于0时不统计播放次数，其它情况统计播放次数。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INVALIDPARAMETERVALUE_DATE = "InvalidParameterValue.Date"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_METRIC = "InvalidParameterValue.Metric"
func (c *Client) DescribeDailyMostPlayedStat(request *DescribeDailyMostPlayedStatRequest) (response *DescribeDailyMostPlayedStatResponse, err error) {
    return c.DescribeDailyMostPlayedStatWithContext(context.Background(), request)
}

// DescribeDailyMostPlayedStat
// 该接口用于查询每日播放Top100 的媒体文件的播放统计数据。
//
// * 可以查询最近一年的播放统计数据。
//
// * 可以按播放次数或者播放流量查询。
//
// * 播放次数统计说明：
//
//     1. HLS 文件：访问 M3U8 文件时统计播放次数；访问 TS 文件不统计播放次数。
//
//     2. 其它文件（如 MP4 文件）：播放请求带有 range 参数且 range 的 start 参数不等于0时不统计播放次数，其它情况统计播放次数。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"
//  INVALIDPARAMETERVALUE_DATE = "InvalidParameterValue.Date"
//  INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"
//  INVALIDPARAMETERVALUE_METRIC = "InvalidParameterValue.Metric"
func (c *Client) DescribeDailyMostPlayedStatWithContext(ctx context.Context, request *DescribeDailyMostPlayedStatRequest) (response *DescribeDailyMostPlayedStatResponse, err error) {
    if request == nil {
        request = NewDescribeDailyMostPlayedStatRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDailyMostPlayedStat require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDailyMostPlayedStatResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDailyPlayStatFileListRequest() (request *DescribeDailyPlayStatFileListRequest) {
    request = &DescribeDailyPlayStatFileListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeDailyPlayStatFileList")
    
    
    return
}

func NewDescribeDailyPlayStatFileListResponse() (response *DescribeDailyPlayStatFileListResponse) {
    response = &DescribeDailyPlayStatFileListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeDailyPlayStatFileList
// 该接口用于查询播放统计文件的下载地址。
//
// * 可以查询最近一年的播放统计文件下载地址，查询的起始日期和结束日期的时间跨度不超过90天。
//
// * 云点播每天对前一天的 CDN 请求日志进行分析处理，生成播放统计文件。
//
// * 播放统计文件内容包含媒体文件的播放次数、播放流量等统计信息。
//
// * 播放次数统计说明：
//
//     1. HLS 文件：访问M3U8 文件时统计播放次数；访问TS 文件不统计播放次数。
//
//     2. 其它文件（如 MP4 文件）：播放请求带有 range 参数且 range 的 start 参数不等于0时不统计播放次数，其它情况统计播放次数。
//
// * 播放设备的统计：播放请求带了 UserAgent 参数，并且 UserAgent 包含 Android 或者 iPhone 等标识，会统计为移动端播放次数，否则统计为 PC 端播放次数。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeDailyPlayStatFileList(request *DescribeDailyPlayStatFileListRequest) (response *DescribeDailyPlayStatFileListResponse, err error) {
    return c.DescribeDailyPlayStatFileListWithContext(context.Background(), request)
}

// DescribeDailyPlayStatFileList
// 该接口用于查询播放统计文件的下载地址。
//
// * 可以查询最近一年的播放统计文件下载地址，查询的起始日期和结束日期的时间跨度不超过90天。
//
// * 云点播每天对前一天的 CDN 请求日志进行分析处理，生成播放统计文件。
//
// * 播放统计文件内容包含媒体文件的播放次数、播放流量等统计信息。
//
// * 播放次数统计说明：
//
//     1. HLS 文件：访问M3U8 文件时统计播放次数；访问TS 文件不统计播放次数。
//
//     2. 其它文件（如 MP4 文件）：播放请求带有 range 参数且 range 的 start 参数不等于0时不统计播放次数，其它情况统计播放次数。
//
// * 播放设备的统计：播放请求带了 UserAgent 参数，并且 UserAgent 包含 Android 或者 iPhone 等标识，会统计为移动端播放次数，否则统计为 PC 端播放次数。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_DBERROR = "InternalError.DBError"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeDailyPlayStatFileListWithContext(ctx context.Context, request *DescribeDailyPlayStatFileListRequest) (response *DescribeDailyPlayStatFileListResponse, err error) {
    if request == nil {
        request = NewDescribeDailyPlayStatFileListRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDailyPlayStatFileList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDailyPlayStatFileListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDrmDataKeyRequest() (request *DescribeDrmDataKeyRequest) {
    request = &DescribeDrmDataKeyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeDrmDataKey")
    
    
    return
}

func NewDescribeDrmDataKeyResponse() (response *DescribeDrmDataKeyResponse) {
    response = &DescribeDrmDataKeyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeDrmDataKey
// 本 API 是 [旧版本加密](https://cloud.tencent.com/document/product/266/9638) 中 [DescribeDrmDataKey 的 API 2017 接口](https://cloud.tencent.com/document/product/266/9643) 的升级版本。
//
// 
//
// 如果您是新接入点播加密的用户，不要使用该 API，请参考 [视频加密综述](https://cloud.tencent.com/document/product/266/45552) 使用推荐的加密方式。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DescribeDrmDataKey(request *DescribeDrmDataKeyRequest) (response *DescribeDrmDataKeyResponse, err error) {
    return c.DescribeDrmDataKeyWithContext(context.Background(), request)
}

// DescribeDrmDataKey
// 本 API 是 [旧版本加密](https://cloud.tencent.com/document/product/266/9638) 中 [DescribeDrmDataKey 的 API 2017 接口](https://cloud.tencent.com/document/product/266/9643) 的升级版本。
//
// 
//
// 如果您是新接入点播加密的用户，不要使用该 API，请参考 [视频加密综述](https://cloud.tencent.com/document/product/266/45552) 使用推荐的加密方式。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DescribeDrmDataKeyWithContext(ctx context.Context, request *DescribeDrmDataKeyRequest) (response *DescribeDrmDataKeyResponse, err error) {
    if request == nil {
        request = NewDescribeDrmDataKeyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDrmDataKey require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDrmDataKeyResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDrmKeyProviderInfoRequest() (request *DescribeDrmKeyProviderInfoRequest) {
    request = &DescribeDrmKeyProviderInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeDrmKeyProviderInfo")
    
    
    return
}

func NewDescribeDrmKeyProviderInfoResponse() (response *DescribeDrmKeyProviderInfoResponse) {
    response = &DescribeDrmKeyProviderInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeDrmKeyProviderInfo
// 查询 DRM 密钥提供商信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeDrmKeyProviderInfo(request *DescribeDrmKeyProviderInfoRequest) (response *DescribeDrmKeyProviderInfoResponse, err error) {
    return c.DescribeDrmKeyProviderInfoWithContext(context.Background(), request)
}

// DescribeDrmKeyProviderInfo
// 查询 DRM 密钥提供商信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeDrmKeyProviderInfoWithContext(ctx context.Context, request *DescribeDrmKeyProviderInfoRequest) (response *DescribeDrmKeyProviderInfoResponse, err error) {
    if request == nil {
        request = NewDescribeDrmKeyProviderInfoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDrmKeyProviderInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDrmKeyProviderInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeEventConfigRequest() (request *DescribeEventConfigRequest) {
    request = &DescribeEventConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeEventConfig")
    
    
    return
}

func NewDescribeEventConfigResponse() (response *DescribeEventConfigResponse) {
    response = &DescribeEventConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeEventConfig
// 腾讯云点播为客户提供了媒体上传、媒体管理、媒体处理等等服务，在这些服务执行过程或执行结束时，腾讯云点播也提供各种对应的事件通知，方便开发者感知服务处理状态，并做下一步的业务操作。
//
// 
//
// 开发者可以通过本接口来查询当前配置事件通知的接收方式、接收地址以及哪些事件开启了接收回调通知。
//
// 
//
// 默认接口请求频率限制：100次/秒。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeEventConfig(request *DescribeEventConfigRequest) (response *DescribeEventConfigResponse, err error) {
    return c.DescribeEventConfigWithContext(context.Background(), request)
}

// DescribeEventConfig
// 腾讯云点播为客户提供了媒体上传、媒体管理、媒体处理等等服务，在这些服务执行过程或执行结束时，腾讯云点播也提供各种对应的事件通知，方便开发者感知服务处理状态，并做下一步的业务操作。
//
// 
//
// 开发者可以通过本接口来查询当前配置事件通知的接收方式、接收地址以及哪些事件开启了接收回调通知。
//
// 
//
// 默认接口请求频率限制：100次/秒。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeEventConfigWithContext(ctx context.Context, request *DescribeEventConfigRequest) (response *DescribeEventConfigResponse, err error) {
    if request == nil {
        request = NewDescribeEventConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeEventConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeEventConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeEventsStateRequest() (request *DescribeEventsStateRequest) {
    request = &DescribeEventsStateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeEventsState")
    
    
    return
}

func NewDescribeEventsStateResponse() (response *DescribeEventsStateResponse) {
    response = &DescribeEventsStateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeEventsState
// * 该接口用于业务服务器获取 [可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83) 事件通知的状态。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeEventsState(request *DescribeEventsStateRequest) (response *DescribeEventsStateResponse, err error) {
    return c.DescribeEventsStateWithContext(context.Background(), request)
}

// DescribeEventsState
// * 该接口用于业务服务器获取 [可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83) 事件通知的状态。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeEventsStateWithContext(ctx context.Context, request *DescribeEventsStateRequest) (response *DescribeEventsStateResponse, err error) {
    if request == nil {
        request = NewDescribeEventsStateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeEventsState require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeEventsStateResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeFileAttributesRequest() (request *DescribeFileAttributesRequest) {
    request = &DescribeFileAttributesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeFileAttributes")
    
    
    return
}

func NewDescribeFileAttributesResponse() (response *DescribeFileAttributesResponse) {
    response = &DescribeFileAttributesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeFileAttributes
// 用于异步获取文件属性。
//
// - 当前仅支持获取源文件的 Md5。
//
// - 对输入文件为 HLS 或 DASH 的情况，仅获取索引文件的属性。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeFileAttributes(request *DescribeFileAttributesRequest) (response *DescribeFileAttributesResponse, err error) {
    return c.DescribeFileAttributesWithContext(context.Background(), request)
}

// DescribeFileAttributes
// 用于异步获取文件属性。
//
// - 当前仅支持获取源文件的 Md5。
//
// - 对输入文件为 HLS 或 DASH 的情况，仅获取索引文件的属性。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeFileAttributesWithContext(ctx context.Context, request *DescribeFileAttributesRequest) (response *DescribeFileAttributesResponse, err error) {
    if request == nil {
        request = NewDescribeFileAttributesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeFileAttributes require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeFileAttributesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeHeadTailTemplatesRequest() (request *DescribeHeadTailTemplatesRequest) {
    request = &DescribeHeadTailTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeHeadTailTemplates")
    
    
    return
}

func NewDescribeHeadTailTemplatesResponse() (response *DescribeHeadTailTemplatesResponse) {
    response = &DescribeHeadTailTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeHeadTailTemplates
// 获取片头片尾模板列表。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeHeadTailTemplates(request *DescribeHeadTailTemplatesRequest) (response *DescribeHeadTailTemplatesResponse, err error) {
    return c.DescribeHeadTailTemplatesWithContext(context.Background(), request)
}

// DescribeHeadTailTemplates
// 获取片头片尾模板列表。
//
// 可能返回的错误码:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeHeadTailTemplatesWithContext(ctx context.Context, request *DescribeHeadTailTemplatesRequest) (response *DescribeHeadTailTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeHeadTailTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeHeadTailTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeHeadTailTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageProcessingTemplatesRequest() (request *DescribeImageProcessingTemplatesRequest) {
    request = &DescribeImageProcessingTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeImageProcessingTemplates")
    
    
    return
}

func NewDescribeImageProcessingTemplatesResponse() (response *DescribeImageProcessingTemplatesResponse) {
    response = &DescribeImageProcessingTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeImageProcessingTemplates
// 获取图片处理模板列表，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeImageProcessingTemplates(request *DescribeImageProcessingTemplatesRequest) (response *DescribeImageProcessingTemplatesResponse, err error) {
    return c.DescribeImageProcessingTemplatesWithContext(context.Background(), request)
}

// DescribeImageProcessingTemplates
// 获取图片处理模板列表，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeImageProcessingTemplatesWithContext(ctx context.Context, request *DescribeImageProcessingTemplatesRequest) (response *DescribeImageProcessingTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeImageProcessingTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageProcessingTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageProcessingTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageReviewUsageDataRequest() (request *DescribeImageReviewUsageDataRequest) {
    request = &DescribeImageReviewUsageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeImageReviewUsageData")
    
    
    return
}

func NewDescribeImageReviewUsageDataResponse() (response *DescribeImageReviewUsageDataResponse) {
    response = &DescribeImageReviewUsageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeImageReviewUsageData
// 该接口返回查询时间范围内每天使用的图片审核用量信息。
//
//    1. 可以查询最近365天内的图片审核统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeImageReviewUsageData(request *DescribeImageReviewUsageDataRequest) (response *DescribeImageReviewUsageDataResponse, err error) {
    return c.DescribeImageReviewUsageDataWithContext(context.Background(), request)
}

// DescribeImageReviewUsageData
// 该接口返回查询时间范围内每天使用的图片审核用量信息。
//
//    1. 可以查询最近365天内的图片审核统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeImageReviewUsageDataWithContext(ctx context.Context, request *DescribeImageReviewUsageDataRequest) (response *DescribeImageReviewUsageDataResponse, err error) {
    if request == nil {
        request = NewDescribeImageReviewUsageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageReviewUsageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageReviewUsageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageSpriteTemplatesRequest() (request *DescribeImageSpriteTemplatesRequest) {
    request = &DescribeImageSpriteTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeImageSpriteTemplates")
    
    
    return
}

func NewDescribeImageSpriteTemplatesResponse() (response *DescribeImageSpriteTemplatesResponse) {
    response = &DescribeImageSpriteTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeImageSpriteTemplates
// 查询雪碧图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeImageSpriteTemplates(request *DescribeImageSpriteTemplatesRequest) (response *DescribeImageSpriteTemplatesResponse, err error) {
    return c.DescribeImageSpriteTemplatesWithContext(context.Background(), request)
}

// DescribeImageSpriteTemplates
// 查询雪碧图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeImageSpriteTemplatesWithContext(ctx context.Context, request *DescribeImageSpriteTemplatesRequest) (response *DescribeImageSpriteTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeImageSpriteTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageSpriteTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageSpriteTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeLicenseUsageDataRequest() (request *DescribeLicenseUsageDataRequest) {
    request = &DescribeLicenseUsageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeLicenseUsageData")
    
    
    return
}

func NewDescribeLicenseUsageDataResponse() (response *DescribeLicenseUsageDataResponse) {
    response = &DescribeLicenseUsageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeLicenseUsageData
// 该接口返回查询时间范围内每天 License 请求次数信息。
//
//    1. 可以查询最近365天内的 License 请求次数统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeLicenseUsageData(request *DescribeLicenseUsageDataRequest) (response *DescribeLicenseUsageDataResponse, err error) {
    return c.DescribeLicenseUsageDataWithContext(context.Background(), request)
}

// DescribeLicenseUsageData
// 该接口返回查询时间范围内每天 License 请求次数信息。
//
//    1. 可以查询最近365天内的 License 请求次数统计数据。
//
//    2. 查询时间跨度不超过90天。
//
//    3. 查询时间跨度超过1天的，返回以天为粒度的数据，否则，返回以5分钟为粒度的数据。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeLicenseUsageDataWithContext(ctx context.Context, request *DescribeLicenseUsageDataRequest) (response *DescribeLicenseUsageDataResponse, err error) {
    if request == nil {
        request = NewDescribeLicenseUsageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeLicenseUsageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeLicenseUsageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMediaInfosRequest() (request *DescribeMediaInfosRequest) {
    request = &DescribeMediaInfosRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeMediaInfos")
    
    
    return
}

func NewDescribeMediaInfosResponse() (response *DescribeMediaInfosResponse) {
    response = &DescribeMediaInfosResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeMediaInfos
// 1. 该接口可以获取多个媒体文件的多种信息，包括：
//
//     1. 基础信息（basicInfo）：包括媒体名称、分类、播放地址、封面图片等。
//
//     2. 元信息（metaData）：包括大小、时长、视频流信息、音频流信息等。
//
//     3. 转码结果信息（transcodeInfo）：包括该媒体转码生成的各种规格的媒体地址、视频流参数、音频流参数等。
//
//     4. 转动图结果信息（animatedGraphicsInfo）：对视频转动图（如 gif）后的动图信息。
//
//     5. 采样截图信息（sampleSnapshotInfo）：对视频采样截图后的截图信息。
//
//     6. 雪碧图信息（imageSpriteInfo）：对视频截取雪碧图后的雪碧图信息。
//
//     7. 指定时间点截图信息（snapshotByTimeOffsetInfo）：对视频依照指定时间点截图后，的截图信息。
//
//     8. 视频打点信息（keyFrameDescInfo）：对视频设置的打点信息。
//
//     9. 转自适应码流信息（adaptiveDynamicStreamingInfo）：包括规格、加密类型、打包格式等相关信息。
//
//     10. 审核信息（reviewInfo）：包括媒体审核及媒体封面审核信息。
//
// 2. 可以指定回包只返回部分信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_FILEIDSEMPTY = "InvalidParameterValue.FileIdsEmpty"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeMediaInfos(request *DescribeMediaInfosRequest) (response *DescribeMediaInfosResponse, err error) {
    return c.DescribeMediaInfosWithContext(context.Background(), request)
}

// DescribeMediaInfos
// 1. 该接口可以获取多个媒体文件的多种信息，包括：
//
//     1. 基础信息（basicInfo）：包括媒体名称、分类、播放地址、封面图片等。
//
//     2. 元信息（metaData）：包括大小、时长、视频流信息、音频流信息等。
//
//     3. 转码结果信息（transcodeInfo）：包括该媒体转码生成的各种规格的媒体地址、视频流参数、音频流参数等。
//
//     4. 转动图结果信息（animatedGraphicsInfo）：对视频转动图（如 gif）后的动图信息。
//
//     5. 采样截图信息（sampleSnapshotInfo）：对视频采样截图后的截图信息。
//
//     6. 雪碧图信息（imageSpriteInfo）：对视频截取雪碧图后的雪碧图信息。
//
//     7. 指定时间点截图信息（snapshotByTimeOffsetInfo）：对视频依照指定时间点截图后，的截图信息。
//
//     8. 视频打点信息（keyFrameDescInfo）：对视频设置的打点信息。
//
//     9. 转自适应码流信息（adaptiveDynamicStreamingInfo）：包括规格、加密类型、打包格式等相关信息。
//
//     10. 审核信息（reviewInfo）：包括媒体审核及媒体封面审核信息。
//
// 2. 可以指定回包只返回部分信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_FILEIDSEMPTY = "InvalidParameterValue.FileIdsEmpty"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeMediaInfosWithContext(ctx context.Context, request *DescribeMediaInfosRequest) (response *DescribeMediaInfosResponse, err error) {
    if request == nil {
        request = NewDescribeMediaInfosRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeMediaInfos require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeMediaInfosResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMediaPlayStatDetailsRequest() (request *DescribeMediaPlayStatDetailsRequest) {
    request = &DescribeMediaPlayStatDetailsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeMediaPlayStatDetails")
    
    
    return
}

func NewDescribeMediaPlayStatDetailsResponse() (response *DescribeMediaPlayStatDetailsResponse) {
    response = &DescribeMediaPlayStatDetailsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeMediaPlayStatDetails
// 该接口用于查询媒体文件按指定时间粒度统计的播放数据
//
// * 可以查询最近一年的播放统计数据。
//
// * 时间粒度为小时，结束时间和起始时间的跨度最大为7天。
//
// * 时间粒度为天，结束时间和起始时间的跨度最大为90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_INTERVAL = "InvalidParameterValue.Interval"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeMediaPlayStatDetails(request *DescribeMediaPlayStatDetailsRequest) (response *DescribeMediaPlayStatDetailsResponse, err error) {
    return c.DescribeMediaPlayStatDetailsWithContext(context.Background(), request)
}

// DescribeMediaPlayStatDetails
// 该接口用于查询媒体文件按指定时间粒度统计的播放数据
//
// * 可以查询最近一年的播放统计数据。
//
// * 时间粒度为小时，结束时间和起始时间的跨度最大为7天。
//
// * 时间粒度为天，结束时间和起始时间的跨度最大为90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_INTERVAL = "InvalidParameterValue.Interval"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
func (c *Client) DescribeMediaPlayStatDetailsWithContext(ctx context.Context, request *DescribeMediaPlayStatDetailsRequest) (response *DescribeMediaPlayStatDetailsResponse, err error) {
    if request == nil {
        request = NewDescribeMediaPlayStatDetailsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeMediaPlayStatDetails require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeMediaPlayStatDetailsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMediaProcessUsageDataRequest() (request *DescribeMediaProcessUsageDataRequest) {
    request = &DescribeMediaProcessUsageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeMediaProcessUsageData")
    
    
    return
}

func NewDescribeMediaProcessUsageDataResponse() (response *DescribeMediaProcessUsageDataResponse) {
    response = &DescribeMediaProcessUsageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeMediaProcessUsageData
// 该接口返回查询时间范围内每天使用的视频处理用量信息。
//
//    1. 可以查询最近365天内的视频处理统计数据。
//
//    2. 查询时间跨度不超过90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeMediaProcessUsageData(request *DescribeMediaProcessUsageDataRequest) (response *DescribeMediaProcessUsageDataResponse, err error) {
    return c.DescribeMediaProcessUsageDataWithContext(context.Background(), request)
}

// DescribeMediaProcessUsageData
// 该接口返回查询时间范围内每天使用的视频处理用量信息。
//
//    1. 可以查询最近365天内的视频处理统计数据。
//
//    2. 查询时间跨度不超过90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) DescribeMediaProcessUsageDataWithContext(ctx context.Context, request *DescribeMediaProcessUsageDataRequest) (response *DescribeMediaProcessUsageDataResponse, err error) {
    if request == nil {
        request = NewDescribeMediaProcessUsageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeMediaProcessUsageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeMediaProcessUsageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonSamplesRequest() (request *DescribePersonSamplesRequest) {
    request = &DescribePersonSamplesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribePersonSamples")
    
    
    return
}

func NewDescribePersonSamplesResponse() (response *DescribePersonSamplesResponse) {
    response = &DescribePersonSamplesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribePersonSamples
// 该接口用于查询素材样本信息，支持根据素材 ID、名称、标签，分页查询。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribePersonSamples(request *DescribePersonSamplesRequest) (response *DescribePersonSamplesResponse, err error) {
    return c.DescribePersonSamplesWithContext(context.Background(), request)
}

// DescribePersonSamples
// 该接口用于查询素材样本信息，支持根据素材 ID、名称、标签，分页查询。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribePersonSamplesWithContext(ctx context.Context, request *DescribePersonSamplesRequest) (response *DescribePersonSamplesResponse, err error) {
    if request == nil {
        request = NewDescribePersonSamplesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribePersonSamples require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribePersonSamplesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePrepaidProductsRequest() (request *DescribePrepaidProductsRequest) {
    request = &DescribePrepaidProductsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribePrepaidProducts")
    
    
    return
}

func NewDescribePrepaidProductsResponse() (response *DescribePrepaidProductsResponse) {
    response = &DescribePrepaidProductsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribePrepaidProducts
// 该接口可以查询用户已经购买的预付费商品的信息，包括：
//
//     1. 商品的类型、生效和失效日期。
//
//     2. 商品中每种资源的额度和剩余额度。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
func (c *Client) DescribePrepaidProducts(request *DescribePrepaidProductsRequest) (response *DescribePrepaidProductsResponse, err error) {
    return c.DescribePrepaidProductsWithContext(context.Background(), request)
}

// DescribePrepaidProducts
// 该接口可以查询用户已经购买的预付费商品的信息，包括：
//
//     1. 商品的类型、生效和失效日期。
//
//     2. 商品中每种资源的额度和剩余额度。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
func (c *Client) DescribePrepaidProductsWithContext(ctx context.Context, request *DescribePrepaidProductsRequest) (response *DescribePrepaidProductsResponse, err error) {
    if request == nil {
        request = NewDescribePrepaidProductsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribePrepaidProducts require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribePrepaidProductsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeProcedureTemplatesRequest() (request *DescribeProcedureTemplatesRequest) {
    request = &DescribeProcedureTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeProcedureTemplates")
    
    
    return
}

func NewDescribeProcedureTemplatesResponse() (response *DescribeProcedureTemplatesResponse) {
    response = &DescribeProcedureTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeProcedureTemplates
// 根据任务流模板名字，获取任务流模板详情列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_NAMES = "InvalidParameterValue.Names"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeProcedureTemplates(request *DescribeProcedureTemplatesRequest) (response *DescribeProcedureTemplatesResponse, err error) {
    return c.DescribeProcedureTemplatesWithContext(context.Background(), request)
}

// DescribeProcedureTemplates
// 根据任务流模板名字，获取任务流模板详情列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_NAMES = "InvalidParameterValue.Names"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeProcedureTemplatesWithContext(ctx context.Context, request *DescribeProcedureTemplatesRequest) (response *DescribeProcedureTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeProcedureTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeProcedureTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeProcedureTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeReviewDetailsRequest() (request *DescribeReviewDetailsRequest) {
    request = &DescribeReviewDetailsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeReviewDetails")
    
    
    return
}

func NewDescribeReviewDetailsResponse() (response *DescribeReviewDetailsResponse) {
    response = &DescribeReviewDetailsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeReviewDetails
// <b>本接口已不推荐使用，用 [DescribeMediaProcessUsageData](/document/product/266/41464) 替代</b>
//
// 
//
// 该接口返回查询时间范围内每天使用的视频内容智能识别时长数据，单位： 秒。
//
// 
//
// 1. 可以查询最近365天内的视频内容智能识别时长统计数据。
//
// 2. 查询时间跨度不超过90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeReviewDetails(request *DescribeReviewDetailsRequest) (response *DescribeReviewDetailsResponse, err error) {
    return c.DescribeReviewDetailsWithContext(context.Background(), request)
}

// DescribeReviewDetails
// <b>本接口已不推荐使用，用 [DescribeMediaProcessUsageData](/document/product/266/41464) 替代</b>
//
// 
//
// 该接口返回查询时间范围内每天使用的视频内容智能识别时长数据，单位： 秒。
//
// 
//
// 1. 可以查询最近365天内的视频内容智能识别时长统计数据。
//
// 2. 查询时间跨度不超过90天。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeReviewDetailsWithContext(ctx context.Context, request *DescribeReviewDetailsRequest) (response *DescribeReviewDetailsResponse, err error) {
    if request == nil {
        request = NewDescribeReviewDetailsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeReviewDetails require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeReviewDetailsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeReviewTemplatesRequest() (request *DescribeReviewTemplatesRequest) {
    request = &DescribeReviewTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeReviewTemplates")
    
    
    return
}

func NewDescribeReviewTemplatesResponse() (response *DescribeReviewTemplatesResponse) {
    response = &DescribeReviewTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeReviewTemplates
// 获取审核模板列表。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeReviewTemplates(request *DescribeReviewTemplatesRequest) (response *DescribeReviewTemplatesResponse, err error) {
    return c.DescribeReviewTemplatesWithContext(context.Background(), request)
}

// DescribeReviewTemplates
// 获取审核模板列表。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeReviewTemplatesWithContext(ctx context.Context, request *DescribeReviewTemplatesRequest) (response *DescribeReviewTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeReviewTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeReviewTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeReviewTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSampleSnapshotTemplatesRequest() (request *DescribeSampleSnapshotTemplatesRequest) {
    request = &DescribeSampleSnapshotTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeSampleSnapshotTemplates")
    
    
    return
}

func NewDescribeSampleSnapshotTemplatesResponse() (response *DescribeSampleSnapshotTemplatesResponse) {
    response = &DescribeSampleSnapshotTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSampleSnapshotTemplates
// 查询采样截图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSampleSnapshotTemplates(request *DescribeSampleSnapshotTemplatesRequest) (response *DescribeSampleSnapshotTemplatesResponse, err error) {
    return c.DescribeSampleSnapshotTemplatesWithContext(context.Background(), request)
}

// DescribeSampleSnapshotTemplates
// 查询采样截图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSampleSnapshotTemplatesWithContext(ctx context.Context, request *DescribeSampleSnapshotTemplatesRequest) (response *DescribeSampleSnapshotTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeSampleSnapshotTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSampleSnapshotTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSampleSnapshotTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSnapshotByTimeOffsetTemplatesRequest() (request *DescribeSnapshotByTimeOffsetTemplatesRequest) {
    request = &DescribeSnapshotByTimeOffsetTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeSnapshotByTimeOffsetTemplates")
    
    
    return
}

func NewDescribeSnapshotByTimeOffsetTemplatesResponse() (response *DescribeSnapshotByTimeOffsetTemplatesResponse) {
    response = &DescribeSnapshotByTimeOffsetTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSnapshotByTimeOffsetTemplates
// 查询指定时间点截图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSnapshotByTimeOffsetTemplates(request *DescribeSnapshotByTimeOffsetTemplatesRequest) (response *DescribeSnapshotByTimeOffsetTemplatesResponse, err error) {
    return c.DescribeSnapshotByTimeOffsetTemplatesWithContext(context.Background(), request)
}

// DescribeSnapshotByTimeOffsetTemplates
// 查询指定时间点截图模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSnapshotByTimeOffsetTemplatesWithContext(ctx context.Context, request *DescribeSnapshotByTimeOffsetTemplatesRequest) (response *DescribeSnapshotByTimeOffsetTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeSnapshotByTimeOffsetTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSnapshotByTimeOffsetTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSnapshotByTimeOffsetTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStorageDataRequest() (request *DescribeStorageDataRequest) {
    request = &DescribeStorageDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeStorageData")
    
    
    return
}

func NewDescribeStorageDataResponse() (response *DescribeStorageDataResponse) {
    response = &DescribeStorageDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStorageData
// 查询存储空间使用情况和文件数量。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageData(request *DescribeStorageDataRequest) (response *DescribeStorageDataResponse, err error) {
    return c.DescribeStorageDataWithContext(context.Background(), request)
}

// DescribeStorageData
// 查询存储空间使用情况和文件数量。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageDataWithContext(ctx context.Context, request *DescribeStorageDataRequest) (response *DescribeStorageDataResponse, err error) {
    if request == nil {
        request = NewDescribeStorageDataRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStorageData require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStorageDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStorageDetailsRequest() (request *DescribeStorageDetailsRequest) {
    request = &DescribeStorageDetailsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeStorageDetails")
    
    
    return
}

func NewDescribeStorageDetailsResponse() (response *DescribeStorageDetailsResponse) {
    response = &DescribeStorageDetailsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStorageDetails
// 该接口返回查询时间范围内使用的点播存储空间，单位：字节。
//
//     1. 可以查询最近365天内的存储空间数据；
//
//     2. 查询时间跨度不超过90天；
//
//     3. 分钟粒度查询跨度不超过7天；
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_TIMEPARSEERROR = "InternalError.TimeParseError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AREA = "InvalidParameterValue.Area"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STORAGETYPE = "InvalidParameterValue.StorageType"
//  INVALIDPARAMETERVALUE_TIMETYPE = "InvalidParameterValue.TimeType"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageDetails(request *DescribeStorageDetailsRequest) (response *DescribeStorageDetailsResponse, err error) {
    return c.DescribeStorageDetailsWithContext(context.Background(), request)
}

// DescribeStorageDetails
// 该接口返回查询时间范围内使用的点播存储空间，单位：字节。
//
//     1. 可以查询最近365天内的存储空间数据；
//
//     2. 查询时间跨度不超过90天；
//
//     3. 分钟粒度查询跨度不超过7天；
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_TIMEPARSEERROR = "InternalError.TimeParseError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AREA = "InvalidParameterValue.Area"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STORAGETYPE = "InvalidParameterValue.StorageType"
//  INVALIDPARAMETERVALUE_TIMETYPE = "InvalidParameterValue.TimeType"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageDetailsWithContext(ctx context.Context, request *DescribeStorageDetailsRequest) (response *DescribeStorageDetailsResponse, err error) {
    if request == nil {
        request = NewDescribeStorageDetailsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStorageDetails require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStorageDetailsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStorageRegionsRequest() (request *DescribeStorageRegionsRequest) {
    request = &DescribeStorageRegionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeStorageRegions")
    
    
    return
}

func NewDescribeStorageRegionsResponse() (response *DescribeStorageRegionsResponse) {
    response = &DescribeStorageRegionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStorageRegions
// 该接口用于：
//
//   1. 查询点播可开通的所有存储园区列表。
//
//   2. 查询已经开通的园区列表。
//
//   3. 查询默认使用的存储园区。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageRegions(request *DescribeStorageRegionsRequest) (response *DescribeStorageRegionsResponse, err error) {
    return c.DescribeStorageRegionsWithContext(context.Background(), request)
}

// DescribeStorageRegions
// 该接口用于：
//
//   1. 查询点播可开通的所有存储园区列表。
//
//   2. 查询已经开通的园区列表。
//
//   3. 查询默认使用的存储园区。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeStorageRegionsWithContext(ctx context.Context, request *DescribeStorageRegionsRequest) (response *DescribeStorageRegionsResponse, err error) {
    if request == nil {
        request = NewDescribeStorageRegionsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStorageRegions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStorageRegionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSubAppIdsRequest() (request *DescribeSubAppIdsRequest) {
    request = &DescribeSubAppIdsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeSubAppIds")
    
    
    return
}

func NewDescribeSubAppIdsResponse() (response *DescribeSubAppIdsResponse) {
    response = &DescribeSubAppIdsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSubAppIds
// 该接口用于获取当前账号的子应用列表，包含主应用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_SERVICENOTEXIST = "ResourceNotFound.ServiceNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSubAppIds(request *DescribeSubAppIdsRequest) (response *DescribeSubAppIdsResponse, err error) {
    return c.DescribeSubAppIdsWithContext(context.Background(), request)
}

// DescribeSubAppIds
// 该接口用于获取当前账号的子应用列表，包含主应用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_SERVICENOTEXIST = "ResourceNotFound.ServiceNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSubAppIdsWithContext(ctx context.Context, request *DescribeSubAppIdsRequest) (response *DescribeSubAppIdsResponse, err error) {
    if request == nil {
        request = NewDescribeSubAppIdsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSubAppIds require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSubAppIdsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSuperPlayerConfigsRequest() (request *DescribeSuperPlayerConfigsRequest) {
    request = &DescribeSuperPlayerConfigsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeSuperPlayerConfigs")
    
    
    return
}

func NewDescribeSuperPlayerConfigsResponse() (response *DescribeSuperPlayerConfigsResponse) {
    response = &DescribeSuperPlayerConfigsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSuperPlayerConfigs
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 查询播放器配置，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSuperPlayerConfigs(request *DescribeSuperPlayerConfigsRequest) (response *DescribeSuperPlayerConfigsResponse, err error) {
    return c.DescribeSuperPlayerConfigsWithContext(context.Background(), request)
}

// DescribeSuperPlayerConfigs
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 查询播放器配置，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeSuperPlayerConfigsWithContext(ctx context.Context, request *DescribeSuperPlayerConfigsRequest) (response *DescribeSuperPlayerConfigsResponse, err error) {
    if request == nil {
        request = NewDescribeSuperPlayerConfigsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSuperPlayerConfigs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSuperPlayerConfigsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskDetailRequest() (request *DescribeTaskDetailRequest) {
    request = &DescribeTaskDetailRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeTaskDetail")
    
    
    return
}

func NewDescribeTaskDetailResponse() (response *DescribeTaskDetailResponse) {
    response = &DescribeTaskDetailResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTaskDetail
// 通过任务 ID 查询任务的执行状态和结果的详细信息（最多可以查询3天之内提交的任务）。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_TASKID = "InvalidParameterValue.TaskId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTaskDetail(request *DescribeTaskDetailRequest) (response *DescribeTaskDetailResponse, err error) {
    return c.DescribeTaskDetailWithContext(context.Background(), request)
}

// DescribeTaskDetail
// 通过任务 ID 查询任务的执行状态和结果的详细信息（最多可以查询3天之内提交的任务）。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_TASKID = "InvalidParameterValue.TaskId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTaskDetailWithContext(ctx context.Context, request *DescribeTaskDetailRequest) (response *DescribeTaskDetailResponse, err error) {
    if request == nil {
        request = NewDescribeTaskDetailRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTaskDetail require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTaskDetailResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTasksRequest() (request *DescribeTasksRequest) {
    request = &DescribeTasksRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeTasks")
    
    
    return
}

func NewDescribeTasksResponse() (response *DescribeTasksResponse) {
    response = &DescribeTasksResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTasks
// * 该接口用于查询任务列表；
//
// * 当列表数据比较多时，单次接口调用无法拉取整个列表，可通过 ScrollToken 参数，分批拉取；
//
// * 只能查询到最近三天（72 小时）内的任务。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_STATUS = "InvalidParameterValue.Status"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTasks(request *DescribeTasksRequest) (response *DescribeTasksResponse, err error) {
    return c.DescribeTasksWithContext(context.Background(), request)
}

// DescribeTasks
// * 该接口用于查询任务列表；
//
// * 当列表数据比较多时，单次接口调用无法拉取整个列表，可通过 ScrollToken 参数，分批拉取；
//
// * 只能查询到最近三天（72 小时）内的任务。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_STATUS = "InvalidParameterValue.Status"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTasksWithContext(ctx context.Context, request *DescribeTasksRequest) (response *DescribeTasksResponse, err error) {
    if request == nil {
        request = NewDescribeTasksRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTasks require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTasksResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTranscodeTemplatesRequest() (request *DescribeTranscodeTemplatesRequest) {
    request = &DescribeTranscodeTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeTranscodeTemplates")
    
    
    return
}

func NewDescribeTranscodeTemplatesResponse() (response *DescribeTranscodeTemplatesResponse) {
    response = &DescribeTranscodeTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeTranscodeTemplates
// 根据转码模板唯一标识，获取转码模板详情列表。返回结果包含符合条件的所有用户自定义模板及[系统预置转码模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.BD.AC.E7.A0.81.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CONTAINERTYPE = "InvalidParameterValue.ContainerType"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTranscodeTemplates(request *DescribeTranscodeTemplatesRequest) (response *DescribeTranscodeTemplatesResponse, err error) {
    return c.DescribeTranscodeTemplatesWithContext(context.Background(), request)
}

// DescribeTranscodeTemplates
// 根据转码模板唯一标识，获取转码模板详情列表。返回结果包含符合条件的所有用户自定义模板及[系统预置转码模板](https://cloud.tencent.com/document/product/266/33476#.E9.A2.84.E7.BD.AE.E8.BD.AC.E7.A0.81.E6.A8.A1.E6.9D.BF)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CONTAINERTYPE = "InvalidParameterValue.ContainerType"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTranscodeTemplatesWithContext(ctx context.Context, request *DescribeTranscodeTemplatesRequest) (response *DescribeTranscodeTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeTranscodeTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTranscodeTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTranscodeTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeVodDomainsRequest() (request *DescribeVodDomainsRequest) {
    request = &DescribeVodDomainsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeVodDomains")
    
    
    return
}

func NewDescribeVodDomainsResponse() (response *DescribeVodDomainsResponse) {
    response = &DescribeVodDomainsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeVodDomains
// 该接口用于查询点播域名信息列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_LIMITTOOLARGE = "InvalidParameterValue.LimitTooLarge"
//  INVALIDPARAMETERVALUE_OFFSET = "InvalidParameterValue.Offset"
//  INVALIDPARAMETERVALUE_OFFSETTOOLARGE = "InvalidParameterValue.OffsetTooLarge"
func (c *Client) DescribeVodDomains(request *DescribeVodDomainsRequest) (response *DescribeVodDomainsResponse, err error) {
    return c.DescribeVodDomainsWithContext(context.Background(), request)
}

// DescribeVodDomains
// 该接口用于查询点播域名信息列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_LIMITTOOLARGE = "InvalidParameterValue.LimitTooLarge"
//  INVALIDPARAMETERVALUE_OFFSET = "InvalidParameterValue.Offset"
//  INVALIDPARAMETERVALUE_OFFSETTOOLARGE = "InvalidParameterValue.OffsetTooLarge"
func (c *Client) DescribeVodDomainsWithContext(ctx context.Context, request *DescribeVodDomainsRequest) (response *DescribeVodDomainsResponse, err error) {
    if request == nil {
        request = NewDescribeVodDomainsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeVodDomains require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeVodDomainsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeWatermarkTemplatesRequest() (request *DescribeWatermarkTemplatesRequest) {
    request = &DescribeWatermarkTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeWatermarkTemplates")
    
    
    return
}

func NewDescribeWatermarkTemplatesResponse() (response *DescribeWatermarkTemplatesResponse) {
    response = &DescribeWatermarkTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeWatermarkTemplates
// 查询用户自定义水印模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeWatermarkTemplates(request *DescribeWatermarkTemplatesRequest) (response *DescribeWatermarkTemplatesResponse, err error) {
    return c.DescribeWatermarkTemplatesWithContext(context.Background(), request)
}

// DescribeWatermarkTemplates
// 查询用户自定义水印模板，支持根据条件，分页查询。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"
//  INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeWatermarkTemplatesWithContext(ctx context.Context, request *DescribeWatermarkTemplatesRequest) (response *DescribeWatermarkTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeWatermarkTemplatesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeWatermarkTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeWatermarkTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeWordSamplesRequest() (request *DescribeWordSamplesRequest) {
    request = &DescribeWordSamplesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "DescribeWordSamples")
    
    
    return
}

func NewDescribeWordSamplesResponse() (response *DescribeWordSamplesResponse) {
    response = &DescribeWordSamplesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeWordSamples
// 该接口用于根据应用场景、关键词、标签，分页查询关键词样本信息。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeWordSamples(request *DescribeWordSamplesRequest) (response *DescribeWordSamplesResponse, err error) {
    return c.DescribeWordSamplesWithContext(context.Background(), request)
}

// DescribeWordSamples
// 该接口用于根据应用场景、关键词、标签，分页查询关键词样本信息。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeWordSamplesWithContext(ctx context.Context, request *DescribeWordSamplesRequest) (response *DescribeWordSamplesResponse, err error) {
    if request == nil {
        request = NewDescribeWordSamplesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeWordSamples require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeWordSamplesResponse()
    err = c.Send(request, response)
    return
}

func NewEditMediaRequest() (request *EditMediaRequest) {
    request = &EditMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "EditMedia")
    
    
    return
}

func NewEditMediaResponse() (response *EditMediaResponse) {
    response = &EditMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// EditMedia
// 对视频进行编辑（剪辑、拼接等），生成一个新的点播视频。编辑的功能包括：
//
// 
//
// 1. 对点播中的一个文件进行剪辑，生成一个新的视频；
//
// 2. 对点播中的多个文件进行拼接，生成一个新的视频；
//
// 3. 对点播中的多个文件进行剪辑，然后再拼接，生成一个新的视频；
//
// 4. 对点播中的一个流，直接生成一个新的视频；
//
// 5. 对点播中的一个流进行剪辑，生成一个新的视频；
//
// 6. 对点播中的多个流进行拼接，生成一个新的视频；
//
// 7. 对点播中的多个流进行剪辑，然后拼接，生成一个新的视频。
//
// 
//
// 对于生成的新视频，还可以指定生成后的视频是否要执行任务流。
//
// 
//
// >当对直播流做剪辑、拼接等操作时，请确保流结束后再操作。否则生成的视频可能不完整。
//
// 
//
// 如使用事件通知，事件通知的类型为 [视频编辑完成](https://cloud.tencent.com/document/product/266/33794)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) EditMedia(request *EditMediaRequest) (response *EditMediaResponse, err error) {
    return c.EditMediaWithContext(context.Background(), request)
}

// EditMedia
// 对视频进行编辑（剪辑、拼接等），生成一个新的点播视频。编辑的功能包括：
//
// 
//
// 1. 对点播中的一个文件进行剪辑，生成一个新的视频；
//
// 2. 对点播中的多个文件进行拼接，生成一个新的视频；
//
// 3. 对点播中的多个文件进行剪辑，然后再拼接，生成一个新的视频；
//
// 4. 对点播中的一个流，直接生成一个新的视频；
//
// 5. 对点播中的一个流进行剪辑，生成一个新的视频；
//
// 6. 对点播中的多个流进行拼接，生成一个新的视频；
//
// 7. 对点播中的多个流进行剪辑，然后拼接，生成一个新的视频。
//
// 
//
// 对于生成的新视频，还可以指定生成后的视频是否要执行任务流。
//
// 
//
// >当对直播流做剪辑、拼接等操作时，请确保流结束后再操作。否则生成的视频可能不完整。
//
// 
//
// 如使用事件通知，事件通知的类型为 [视频编辑完成](https://cloud.tencent.com/document/product/266/33794)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) EditMediaWithContext(ctx context.Context, request *EditMediaRequest) (response *EditMediaResponse, err error) {
    if request == nil {
        request = NewEditMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("EditMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewEditMediaResponse()
    err = c.Send(request, response)
    return
}

func NewExecuteFunctionRequest() (request *ExecuteFunctionRequest) {
    request = &ExecuteFunctionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ExecuteFunction")
    
    
    return
}

func NewExecuteFunctionResponse() (response *ExecuteFunctionResponse) {
    response = &ExecuteFunctionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ExecuteFunction
// 本接口仅用于定制开发的特殊场景，除非云点播客服人员主动告知您需要使用本接口，其它情况请勿调用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_FUNCTIONARG = "InvalidParameterValue.FunctionArg"
//  INVALIDPARAMETERVALUE_FUNCTIONNAME = "InvalidParameterValue.FunctionName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ExecuteFunction(request *ExecuteFunctionRequest) (response *ExecuteFunctionResponse, err error) {
    return c.ExecuteFunctionWithContext(context.Background(), request)
}

// ExecuteFunction
// 本接口仅用于定制开发的特殊场景，除非云点播客服人员主动告知您需要使用本接口，其它情况请勿调用。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_FUNCTIONARG = "InvalidParameterValue.FunctionArg"
//  INVALIDPARAMETERVALUE_FUNCTIONNAME = "InvalidParameterValue.FunctionName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ExecuteFunctionWithContext(ctx context.Context, request *ExecuteFunctionRequest) (response *ExecuteFunctionResponse, err error) {
    if request == nil {
        request = NewExecuteFunctionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ExecuteFunction require credential")
    }

    request.SetContext(ctx)
    
    response = NewExecuteFunctionResponse()
    err = c.Send(request, response)
    return
}

func NewExtractTraceWatermarkRequest() (request *ExtractTraceWatermarkRequest) {
    request = &ExtractTraceWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ExtractTraceWatermark")
    
    
    return
}

func NewExtractTraceWatermarkResponse() (response *ExtractTraceWatermarkResponse) {
    response = &ExtractTraceWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ExtractTraceWatermark
// 用于提取溯源水印。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_FUNCTIONARG = "InvalidParameterValue.FunctionArg"
//  INVALIDPARAMETERVALUE_FUNCTIONNAME = "InvalidParameterValue.FunctionName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ExtractTraceWatermark(request *ExtractTraceWatermarkRequest) (response *ExtractTraceWatermarkResponse, err error) {
    return c.ExtractTraceWatermarkWithContext(context.Background(), request)
}

// ExtractTraceWatermark
// 用于提取溯源水印。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_FUNCTIONARG = "InvalidParameterValue.FunctionArg"
//  INVALIDPARAMETERVALUE_FUNCTIONNAME = "InvalidParameterValue.FunctionName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ExtractTraceWatermarkWithContext(ctx context.Context, request *ExtractTraceWatermarkRequest) (response *ExtractTraceWatermarkResponse, err error) {
    if request == nil {
        request = NewExtractTraceWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ExtractTraceWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewExtractTraceWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewForbidMediaDistributionRequest() (request *ForbidMediaDistributionRequest) {
    request = &ForbidMediaDistributionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ForbidMediaDistribution")
    
    
    return
}

func NewForbidMediaDistributionResponse() (response *ForbidMediaDistributionResponse) {
    response = &ForbidMediaDistributionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ForbidMediaDistribution
// * 对媒体禁播后，除了点播控制台预览，其他场景访问视频各种资源的 URL（原始文件、转码输出文件、截图等）均会返回 403。
//
//   禁播/解禁操作全网生效时间约 5~10 分钟。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIAFORBIDEDBYSYSTEM = "FailedOperation.MediaForbidedBySystem"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_FILEIDSTOOMANY = "InvalidParameterValue.FileIdsTooMany"
//  INVALIDPARAMETERVALUE_OPERATION = "InvalidParameterValue.Operation"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ForbidMediaDistribution(request *ForbidMediaDistributionRequest) (response *ForbidMediaDistributionResponse, err error) {
    return c.ForbidMediaDistributionWithContext(context.Background(), request)
}

// ForbidMediaDistribution
// * 对媒体禁播后，除了点播控制台预览，其他场景访问视频各种资源的 URL（原始文件、转码输出文件、截图等）均会返回 403。
//
//   禁播/解禁操作全网生效时间约 5~10 分钟。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIAFORBIDEDBYSYSTEM = "FailedOperation.MediaForbidedBySystem"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_FILEIDSTOOMANY = "InvalidParameterValue.FileIdsTooMany"
//  INVALIDPARAMETERVALUE_OPERATION = "InvalidParameterValue.Operation"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ForbidMediaDistributionWithContext(ctx context.Context, request *ForbidMediaDistributionRequest) (response *ForbidMediaDistributionResponse, err error) {
    if request == nil {
        request = NewForbidMediaDistributionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ForbidMediaDistribution require credential")
    }

    request.SetContext(ctx)
    
    response = NewForbidMediaDistributionResponse()
    err = c.Send(request, response)
    return
}

func NewLiveRealTimeClipRequest() (request *LiveRealTimeClipRequest) {
    request = &LiveRealTimeClipRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "LiveRealTimeClip")
    
    
    return
}

func NewLiveRealTimeClipResponse() (response *LiveRealTimeClipResponse) {
    response = &LiveRealTimeClipResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// LiveRealTimeClip
// 直播即时剪辑，是指在直播过程中（即直播尚未结束时），客户可以在过往直播内容中选择一段，实时生成一个新的视频（HLS 格式），开发者可以将其立即分享出去，或者长久保存起来。
//
// 
//
// 腾讯云点播支持两种即时剪辑模式：
//
// - 剪辑固化：将剪辑出来的视频保存成独立的视频，拥有独立 FileId；适用于将精彩片段**长久保存**的场景；
//
// - 剪辑不固化：剪辑得到的视频附属于直播录制文件，没有独立 FileId；适用于将精彩片段**临时分享**的场景。
//
// 
//
// 注意：
//
// - 使用直播即时剪辑功能的前提是：目标直播流开启了[时移回看](https://cloud.tencent.com/document/product/267/32742)功能。
//
// - 直播即时剪辑是基于直播录制生成的 m3u8 文件进行的，故而其最小剪辑精度为一个 ts 切片，无法实现秒级或者更为精确的剪辑精度。
//
// - 由于直播过程中可能存在断流的情况，所以有可能导致剪辑生成的实际视频时长与期望不一致。例如剪辑某个直播流的时间区间为 2018-09-20T10:30:00Z 到 2018-09-20T10:40:00Z ，如果在该时间区间中发生过断流，那么返回的媒资文件的时长将少于 10 分钟，在这种情况下，可以通过输出参数 <a href="#p_segmentset">SegmentSet</a> 感知到。
//
// 
//
// ### 剪辑固化
//
// 所谓剪辑固化，是指将剪辑出来的视频是保存成一个独立的视频（拥有独立的 FileId）。其生命周期不受原始直播录制视频影响（即使原始录制视频被删除，剪辑结果也不会受到任何影响）；也可以对其进行转码、微信发布等二次处理。
//
// 
//
// 举例如下：一场完整的足球比赛，直播录制出来的原始视频可能长达 2 个小时，客户出于节省成本的目的可以对这个视频存储 2 个月，但对于直播即时剪辑的「精彩时刻」视频却可以指定存储更长时间，同时可以单独对「精彩时刻」视频进行转码、微信发布等额外的点播操作，这时候可以选择直播即时剪辑并且固化的方案。
//
// 
//
// 剪辑固化的优势在于其生命周期与原始录制视频相互独立，可以独立管理、长久保存。
//
// 
//
// ### 剪辑不固化
//
// 所谓剪辑不固化，是指剪辑所得到的结果（m3u8 文件）与直播录制视频共享相同的 ts 分片，新生成的视频不是一个独立完整的视频（没有独立 FileId，只有播放 URL），其有效期与直播录制的完整视频有效期是一致的。一旦直播录制出来的视频被删除，也会导致该片段无法播放。
//
// 
//
// 剪辑不固化，由于其剪辑结果不是一个独立的视频，因而也不会纳入点播媒资视频管理（例如控制台的视频总数不会统计这一片段）中，也无法单独针对这个片段做转码、微信发布等任何视频处理操作。
//
// 
//
// 剪辑不固化的优势在于其剪辑操作十分“轻量化”，不会产生额外的存储开销。但其不足之处在于生命周期与原始录制视频相同，且无法进一步进行转码等视频处理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLIPDURATION = "InvalidParameterValue.ClipDuration"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STREAMIDINVALID = "InvalidParameterValue.StreamIdInvalid"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) LiveRealTimeClip(request *LiveRealTimeClipRequest) (response *LiveRealTimeClipResponse, err error) {
    return c.LiveRealTimeClipWithContext(context.Background(), request)
}

// LiveRealTimeClip
// 直播即时剪辑，是指在直播过程中（即直播尚未结束时），客户可以在过往直播内容中选择一段，实时生成一个新的视频（HLS 格式），开发者可以将其立即分享出去，或者长久保存起来。
//
// 
//
// 腾讯云点播支持两种即时剪辑模式：
//
// - 剪辑固化：将剪辑出来的视频保存成独立的视频，拥有独立 FileId；适用于将精彩片段**长久保存**的场景；
//
// - 剪辑不固化：剪辑得到的视频附属于直播录制文件，没有独立 FileId；适用于将精彩片段**临时分享**的场景。
//
// 
//
// 注意：
//
// - 使用直播即时剪辑功能的前提是：目标直播流开启了[时移回看](https://cloud.tencent.com/document/product/267/32742)功能。
//
// - 直播即时剪辑是基于直播录制生成的 m3u8 文件进行的，故而其最小剪辑精度为一个 ts 切片，无法实现秒级或者更为精确的剪辑精度。
//
// - 由于直播过程中可能存在断流的情况，所以有可能导致剪辑生成的实际视频时长与期望不一致。例如剪辑某个直播流的时间区间为 2018-09-20T10:30:00Z 到 2018-09-20T10:40:00Z ，如果在该时间区间中发生过断流，那么返回的媒资文件的时长将少于 10 分钟，在这种情况下，可以通过输出参数 <a href="#p_segmentset">SegmentSet</a> 感知到。
//
// 
//
// ### 剪辑固化
//
// 所谓剪辑固化，是指将剪辑出来的视频是保存成一个独立的视频（拥有独立的 FileId）。其生命周期不受原始直播录制视频影响（即使原始录制视频被删除，剪辑结果也不会受到任何影响）；也可以对其进行转码、微信发布等二次处理。
//
// 
//
// 举例如下：一场完整的足球比赛，直播录制出来的原始视频可能长达 2 个小时，客户出于节省成本的目的可以对这个视频存储 2 个月，但对于直播即时剪辑的「精彩时刻」视频却可以指定存储更长时间，同时可以单独对「精彩时刻」视频进行转码、微信发布等额外的点播操作，这时候可以选择直播即时剪辑并且固化的方案。
//
// 
//
// 剪辑固化的优势在于其生命周期与原始录制视频相互独立，可以独立管理、长久保存。
//
// 
//
// ### 剪辑不固化
//
// 所谓剪辑不固化，是指剪辑所得到的结果（m3u8 文件）与直播录制视频共享相同的 ts 分片，新生成的视频不是一个独立完整的视频（没有独立 FileId，只有播放 URL），其有效期与直播录制的完整视频有效期是一致的。一旦直播录制出来的视频被删除，也会导致该片段无法播放。
//
// 
//
// 剪辑不固化，由于其剪辑结果不是一个独立的视频，因而也不会纳入点播媒资视频管理（例如控制台的视频总数不会统计这一片段）中，也无法单独针对这个片段做转码、微信发布等任何视频处理操作。
//
// 
//
// 剪辑不固化的优势在于其剪辑操作十分“轻量化”，不会产生额外的存储开销。但其不足之处在于生命周期与原始录制视频相同，且无法进一步进行转码等视频处理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLIPDURATION = "InvalidParameterValue.ClipDuration"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STREAMIDINVALID = "InvalidParameterValue.StreamIdInvalid"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) LiveRealTimeClipWithContext(ctx context.Context, request *LiveRealTimeClipRequest) (response *LiveRealTimeClipResponse, err error) {
    if request == nil {
        request = NewLiveRealTimeClipRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("LiveRealTimeClip require credential")
    }

    request.SetContext(ctx)
    
    response = NewLiveRealTimeClipResponse()
    err = c.Send(request, response)
    return
}

func NewManageTaskRequest() (request *ManageTaskRequest) {
    request = &ManageTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ManageTask")
    
    
    return
}

func NewManageTaskResponse() (response *ManageTaskResponse) {
    response = &ManageTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ManageTask
// 对已发起的任务进行管理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INVALIDOPERATIONTYPE = "InvalidParameterValue.InvalidOperationType"
//  INVALIDPARAMETERVALUE_TASKID = "InvalidParameterValue.TaskId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ManageTask(request *ManageTaskRequest) (response *ManageTaskResponse, err error) {
    return c.ManageTaskWithContext(context.Background(), request)
}

// ManageTask
// 对已发起的任务进行管理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INVALIDOPERATIONTYPE = "InvalidParameterValue.InvalidOperationType"
//  INVALIDPARAMETERVALUE_TASKID = "InvalidParameterValue.TaskId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ManageTaskWithContext(ctx context.Context, request *ManageTaskRequest) (response *ManageTaskResponse, err error) {
    if request == nil {
        request = NewManageTaskRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ManageTask require credential")
    }

    request.SetContext(ctx)
    
    response = NewManageTaskResponse()
    err = c.Send(request, response)
    return
}

func NewModifyAIAnalysisTemplateRequest() (request *ModifyAIAnalysisTemplateRequest) {
    request = &ModifyAIAnalysisTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyAIAnalysisTemplate")
    
    
    return
}

func NewModifyAIAnalysisTemplateResponse() (response *ModifyAIAnalysisTemplateResponse) {
    response = &ModifyAIAnalysisTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyAIAnalysisTemplate
// 修改用户自定义音视频内容分析模板。
//
// 
//
// 注意：模板 ID 10000 以下的为系统预置模板，不允许修改。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLASSIFCATIONCONFIGURE = "InvalidParameterValue.ClassifcationConfigure"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COVERCONFIGURE = "InvalidParameterValue.CoverConfigure"
//  INVALIDPARAMETERVALUE_FRAMETAGCONFIGURE = "InvalidParameterValue.FrameTagConfigure"
//  INVALIDPARAMETERVALUE_HIGHLIGHTCONFIGURE = "InvalidParameterValue.HighlightConfigure"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGCONFIGURE = "InvalidParameterValue.TagConfigure"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAIAnalysisTemplate(request *ModifyAIAnalysisTemplateRequest) (response *ModifyAIAnalysisTemplateResponse, err error) {
    return c.ModifyAIAnalysisTemplateWithContext(context.Background(), request)
}

// ModifyAIAnalysisTemplate
// 修改用户自定义音视频内容分析模板。
//
// 
//
// 注意：模板 ID 10000 以下的为系统预置模板，不允许修改。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CLASSIFCATIONCONFIGURE = "InvalidParameterValue.ClassifcationConfigure"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COVERCONFIGURE = "InvalidParameterValue.CoverConfigure"
//  INVALIDPARAMETERVALUE_FRAMETAGCONFIGURE = "InvalidParameterValue.FrameTagConfigure"
//  INVALIDPARAMETERVALUE_HIGHLIGHTCONFIGURE = "InvalidParameterValue.HighlightConfigure"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGCONFIGURE = "InvalidParameterValue.TagConfigure"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAIAnalysisTemplateWithContext(ctx context.Context, request *ModifyAIAnalysisTemplateRequest) (response *ModifyAIAnalysisTemplateResponse, err error) {
    if request == nil {
        request = NewModifyAIAnalysisTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyAIAnalysisTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyAIAnalysisTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyAIRecognitionTemplateRequest() (request *ModifyAIRecognitionTemplateRequest) {
    request = &ModifyAIRecognitionTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyAIRecognitionTemplate")
    
    
    return
}

func NewModifyAIRecognitionTemplateResponse() (response *ModifyAIRecognitionTemplateResponse) {
    response = &ModifyAIRecognitionTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyAIRecognitionTemplate
// 修改用户自定义音视频内容识别模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFAULTLIBRARYLABELSET = "InvalidParameterValue.DefaultLibraryLabelSet"
//  INVALIDPARAMETERVALUE_FACELIBRARY = "InvalidParameterValue.FaceLibrary"
//  INVALIDPARAMETERVALUE_FACESCORE = "InvalidParameterValue.FaceScore"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_MODIFYDEFAULTTEMPLATE = "InvalidParameterValue.ModifyDefaultTemplate"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_OBJECTLIBRARY = "InvalidParameterValue.ObjectLibrary"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SUBTITLEFORMAT = "InvalidParameterValue.SubtitleFormat"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  INVALIDPARAMETERVALUE_USERDEFINELIBRARYLABELSET = "InvalidParameterValue.UserDefineLibraryLabelSet"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAIRecognitionTemplate(request *ModifyAIRecognitionTemplateRequest) (response *ModifyAIRecognitionTemplateResponse, err error) {
    return c.ModifyAIRecognitionTemplateWithContext(context.Background(), request)
}

// ModifyAIRecognitionTemplate
// 修改用户自定义音视频内容识别模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFAULTLIBRARYLABELSET = "InvalidParameterValue.DefaultLibraryLabelSet"
//  INVALIDPARAMETERVALUE_FACELIBRARY = "InvalidParameterValue.FaceLibrary"
//  INVALIDPARAMETERVALUE_FACESCORE = "InvalidParameterValue.FaceScore"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_MODIFYDEFAULTTEMPLATE = "InvalidParameterValue.ModifyDefaultTemplate"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_OBJECTLIBRARY = "InvalidParameterValue.ObjectLibrary"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SUBTITLEFORMAT = "InvalidParameterValue.SubtitleFormat"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  INVALIDPARAMETERVALUE_USERDEFINELIBRARYLABELSET = "InvalidParameterValue.UserDefineLibraryLabelSet"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAIRecognitionTemplateWithContext(ctx context.Context, request *ModifyAIRecognitionTemplateRequest) (response *ModifyAIRecognitionTemplateResponse, err error) {
    if request == nil {
        request = NewModifyAIRecognitionTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyAIRecognitionTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyAIRecognitionTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyAdaptiveDynamicStreamingTemplateRequest() (request *ModifyAdaptiveDynamicStreamingTemplateRequest) {
    request = &ModifyAdaptiveDynamicStreamingTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyAdaptiveDynamicStreamingTemplate")
    
    
    return
}

func NewModifyAdaptiveDynamicStreamingTemplateResponse() (response *ModifyAdaptiveDynamicStreamingTemplateResponse) {
    response = &ModifyAdaptiveDynamicStreamingTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyAdaptiveDynamicStreamingTemplate
// 修改转自适应码流模板
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEOBITRATE = "InvalidParameterValue.DisableHigherVideoBitrate"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEORESOLUTION = "InvalidParameterValue.DisableHigherVideoResolution"
//  INVALIDPARAMETERVALUE_DRMTYPE = "InvalidParameterValue.DrmType"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_SOUNDSYSTEM = "InvalidParameterValue.SoundSystem"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
func (c *Client) ModifyAdaptiveDynamicStreamingTemplate(request *ModifyAdaptiveDynamicStreamingTemplateRequest) (response *ModifyAdaptiveDynamicStreamingTemplateResponse, err error) {
    return c.ModifyAdaptiveDynamicStreamingTemplateWithContext(context.Background(), request)
}

// ModifyAdaptiveDynamicStreamingTemplate
// 修改转自适应码流模板
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEOBITRATE = "InvalidParameterValue.DisableHigherVideoBitrate"
//  INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEORESOLUTION = "InvalidParameterValue.DisableHigherVideoResolution"
//  INVALIDPARAMETERVALUE_DRMTYPE = "InvalidParameterValue.DrmType"
//  INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_SOUNDSYSTEM = "InvalidParameterValue.SoundSystem"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
func (c *Client) ModifyAdaptiveDynamicStreamingTemplateWithContext(ctx context.Context, request *ModifyAdaptiveDynamicStreamingTemplateRequest) (response *ModifyAdaptiveDynamicStreamingTemplateResponse, err error) {
    if request == nil {
        request = NewModifyAdaptiveDynamicStreamingTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyAdaptiveDynamicStreamingTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyAdaptiveDynamicStreamingTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyAnimatedGraphicsTemplateRequest() (request *ModifyAnimatedGraphicsTemplateRequest) {
    request = &ModifyAnimatedGraphicsTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyAnimatedGraphicsTemplate")
    
    
    return
}

func NewModifyAnimatedGraphicsTemplateResponse() (response *ModifyAnimatedGraphicsTemplateResponse) {
    response = &ModifyAnimatedGraphicsTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyAnimatedGraphicsTemplate
// 修改用户自定义转动图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FORMATWEBPLACKWIDTHANDHEIGHT = "InvalidParameterValue.FormatWebpLackWidthAndHeight"
//  INVALIDPARAMETERVALUE_FORMATWEBPWIDTHANDHEIGHTBOTHZERO = "InvalidParameterValue.FormatWebpWidthAndHeightBothZero"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_QUALITY = "InvalidParameterValue.Quality"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAnimatedGraphicsTemplate(request *ModifyAnimatedGraphicsTemplateRequest) (response *ModifyAnimatedGraphicsTemplateResponse, err error) {
    return c.ModifyAnimatedGraphicsTemplateWithContext(context.Background(), request)
}

// ModifyAnimatedGraphicsTemplate
// 修改用户自定义转动图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_FORMATWEBPLACKWIDTHANDHEIGHT = "InvalidParameterValue.FormatWebpLackWidthAndHeight"
//  INVALIDPARAMETERVALUE_FORMATWEBPWIDTHANDHEIGHTBOTHZERO = "InvalidParameterValue.FormatWebpWidthAndHeightBothZero"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_QUALITY = "InvalidParameterValue.Quality"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyAnimatedGraphicsTemplateWithContext(ctx context.Context, request *ModifyAnimatedGraphicsTemplateRequest) (response *ModifyAnimatedGraphicsTemplateResponse, err error) {
    if request == nil {
        request = NewModifyAnimatedGraphicsTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyAnimatedGraphicsTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyAnimatedGraphicsTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyClassRequest() (request *ModifyClassRequest) {
    request = &ModifyClassRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyClass")
    
    
    return
}

func NewModifyClassResponse() (response *ModifyClassResponse) {
    response = &ModifyClassResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyClass
// 修改媒体分类属性。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_CLASSNAMEDUPLICATE = "FailedOperation.ClassNameDuplicate"
//  FAILEDOPERATION_CLASSNOFOUND = "FailedOperation.ClassNoFound"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSID = "InvalidParameterValue.ClassId"
//  INVALIDPARAMETERVALUE_CLASSNAME = "InvalidParameterValue.ClassName"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyClass(request *ModifyClassRequest) (response *ModifyClassResponse, err error) {
    return c.ModifyClassWithContext(context.Background(), request)
}

// ModifyClass
// 修改媒体分类属性。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_CLASSNAMEDUPLICATE = "FailedOperation.ClassNameDuplicate"
//  FAILEDOPERATION_CLASSNOFOUND = "FailedOperation.ClassNoFound"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_CLASSID = "InvalidParameterValue.ClassId"
//  INVALIDPARAMETERVALUE_CLASSNAME = "InvalidParameterValue.ClassName"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyClassWithContext(ctx context.Context, request *ModifyClassRequest) (response *ModifyClassResponse, err error) {
    if request == nil {
        request = NewModifyClassRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyClass require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyClassResponse()
    err = c.Send(request, response)
    return
}

func NewModifyContentReviewTemplateRequest() (request *ModifyContentReviewTemplateRequest) {
    request = &ModifyContentReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyContentReviewTemplate")
    
    
    return
}

func NewModifyContentReviewTemplateResponse() (response *ModifyContentReviewTemplateResponse) {
    response = &ModifyContentReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [修改审核模板](https://cloud.tencent.com/document/api/266/84388)。
//
// 修改用户自定义音视频内容审核模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_BLOCKCONFIDENCE = "InvalidParameterValue.BlockConfidence"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REVIEWCONFIDENCE = "InvalidParameterValue.ReviewConfidence"
//  INVALIDPARAMETERVALUE_REVIEWWALLSWITCH = "InvalidParameterValue.ReviewWallSwitch"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyContentReviewTemplate(request *ModifyContentReviewTemplateRequest) (response *ModifyContentReviewTemplateResponse, err error) {
    return c.ModifyContentReviewTemplateWithContext(context.Background(), request)
}

// ModifyContentReviewTemplate
// 该 API 已经<font color=red>不再维护</font>，新版审核模板支持音视频审核和图片审核，详细请参考 [修改审核模板](https://cloud.tencent.com/document/api/266/84388)。
//
// 修改用户自定义音视频内容审核模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_BLOCKCONFIDENCE = "InvalidParameterValue.BlockConfidence"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REVIEWCONFIDENCE = "InvalidParameterValue.ReviewConfidence"
//  INVALIDPARAMETERVALUE_REVIEWWALLSWITCH = "InvalidParameterValue.ReviewWallSwitch"
//  INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"
//  INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyContentReviewTemplateWithContext(ctx context.Context, request *ModifyContentReviewTemplateRequest) (response *ModifyContentReviewTemplateResponse, err error) {
    if request == nil {
        request = NewModifyContentReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyContentReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyContentReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDefaultStorageRegionRequest() (request *ModifyDefaultStorageRegionRequest) {
    request = &ModifyDefaultStorageRegionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyDefaultStorageRegion")
    
    
    return
}

func NewModifyDefaultStorageRegionResponse() (response *ModifyDefaultStorageRegionResponse) {
    response = &ModifyDefaultStorageRegionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyDefaultStorageRegion
// 该接口用于设置默认的存储地域。上传文件时如果没有指定地域，将上传到默认地域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyDefaultStorageRegion(request *ModifyDefaultStorageRegionRequest) (response *ModifyDefaultStorageRegionResponse, err error) {
    return c.ModifyDefaultStorageRegionWithContext(context.Background(), request)
}

// ModifyDefaultStorageRegion
// 该接口用于设置默认的存储地域。上传文件时如果没有指定地域，将上传到默认地域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyDefaultStorageRegionWithContext(ctx context.Context, request *ModifyDefaultStorageRegionRequest) (response *ModifyDefaultStorageRegionResponse, err error) {
    if request == nil {
        request = NewModifyDefaultStorageRegionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDefaultStorageRegion require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDefaultStorageRegionResponse()
    err = c.Send(request, response)
    return
}

func NewModifyEventConfigRequest() (request *ModifyEventConfigRequest) {
    request = &ModifyEventConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyEventConfig")
    
    
    return
}

func NewModifyEventConfigResponse() (response *ModifyEventConfigResponse) {
    response = &ModifyEventConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyEventConfig
// 腾讯云点播为客户提供了媒体上传、媒体管理、媒体处理等等服务，在这些服务执行过程或执行结束时，腾讯云点播也提供各种对应的事件通知，方便开发者感知服务处理状态，并做下一步的业务操作。
//
// 
//
// 开发者可以通过调用本接口来实现：
//
// - 设置接收回调通知的类型，目前有[ HTTP 回调通知](https://cloud.tencent.com/document/product/266/33779) 和 [基于消息队列的可靠通知](https://cloud.tencent.com/document/product/266/33779) 两种类型。
//
// - 对于[ HTTP 回调通知](https://cloud.tencent.com/document/product/266/33779)，可设置 3.0 格式回调的地址。3.0 格式回调的说明参见 [历史格式回调](https://cloud.tencent.com/document/product/266/33796)。
//
// - 对具体事件服务的通知事件选择设置接收或者忽略。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyEventConfig(request *ModifyEventConfigRequest) (response *ModifyEventConfigResponse, err error) {
    return c.ModifyEventConfigWithContext(context.Background(), request)
}

// ModifyEventConfig
// 腾讯云点播为客户提供了媒体上传、媒体管理、媒体处理等等服务，在这些服务执行过程或执行结束时，腾讯云点播也提供各种对应的事件通知，方便开发者感知服务处理状态，并做下一步的业务操作。
//
// 
//
// 开发者可以通过调用本接口来实现：
//
// - 设置接收回调通知的类型，目前有[ HTTP 回调通知](https://cloud.tencent.com/document/product/266/33779) 和 [基于消息队列的可靠通知](https://cloud.tencent.com/document/product/266/33779) 两种类型。
//
// - 对于[ HTTP 回调通知](https://cloud.tencent.com/document/product/266/33779)，可设置 3.0 格式回调的地址。3.0 格式回调的说明参见 [历史格式回调](https://cloud.tencent.com/document/product/266/33796)。
//
// - 对具体事件服务的通知事件选择设置接收或者忽略。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyEventConfigWithContext(ctx context.Context, request *ModifyEventConfigRequest) (response *ModifyEventConfigResponse, err error) {
    if request == nil {
        request = NewModifyEventConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyEventConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyEventConfigResponse()
    err = c.Send(request, response)
    return
}

func NewModifyHeadTailTemplateRequest() (request *ModifyHeadTailTemplateRequest) {
    request = &ModifyHeadTailTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyHeadTailTemplate")
    
    
    return
}

func NewModifyHeadTailTemplateResponse() (response *ModifyHeadTailTemplateResponse) {
    response = &ModifyHeadTailTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyHeadTailTemplate
// 修改片头片尾模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_STORAGEREGIONS = "InvalidParameterValue.StorageRegions"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyHeadTailTemplate(request *ModifyHeadTailTemplateRequest) (response *ModifyHeadTailTemplateResponse, err error) {
    return c.ModifyHeadTailTemplateWithContext(context.Background(), request)
}

// ModifyHeadTailTemplate
// 修改片头片尾模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_STORAGEREGIONS = "InvalidParameterValue.StorageRegions"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyHeadTailTemplateWithContext(ctx context.Context, request *ModifyHeadTailTemplateRequest) (response *ModifyHeadTailTemplateResponse, err error) {
    if request == nil {
        request = NewModifyHeadTailTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyHeadTailTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyHeadTailTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyImageSpriteTemplateRequest() (request *ModifyImageSpriteTemplateRequest) {
    request = &ModifyImageSpriteTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyImageSpriteTemplate")
    
    
    return
}

func NewModifyImageSpriteTemplateResponse() (response *ModifyImageSpriteTemplateResponse) {
    response = &ModifyImageSpriteTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyImageSpriteTemplate
// 修改用户自定义雪碧图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COLUMNCOUNT = "InvalidParameterValue.ColumnCount"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_ROWCOUNT = "InvalidParameterValue.RowCount"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyImageSpriteTemplate(request *ModifyImageSpriteTemplateRequest) (response *ModifyImageSpriteTemplateResponse, err error) {
    return c.ModifyImageSpriteTemplateWithContext(context.Background(), request)
}

// ModifyImageSpriteTemplate
// 修改用户自定义雪碧图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COLUMNCOUNT = "InvalidParameterValue.ColumnCount"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_ROWCOUNT = "InvalidParameterValue.RowCount"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyImageSpriteTemplateWithContext(ctx context.Context, request *ModifyImageSpriteTemplateRequest) (response *ModifyImageSpriteTemplateResponse, err error) {
    if request == nil {
        request = NewModifyImageSpriteTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyImageSpriteTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyImageSpriteTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyMediaInfoRequest() (request *ModifyMediaInfoRequest) {
    request = &ModifyMediaInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyMediaInfo")
    
    
    return
}

func NewModifyMediaInfoResponse() (response *ModifyMediaInfoResponse) {
    response = &ModifyMediaInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyMediaInfo
// 修改媒体文件的属性，包括分类、名称、描述、标签、过期时间、打点信息、视频封面、字幕信息等。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INTERNALERROR_UPDATEMEDIAERROR = "InternalError.UpdateMediaError"
//  INTERNALERROR_UPLOADCOVERIMAGEERROR = "InternalError.UploadCoverImageError"
//  INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDCLEARKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndClearKeyFrameDescsConflict"
//  INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDDELETEKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndDeleteKeyFrameDescsConflict"
//  INVALIDPARAMETERVALUE_ADDTAGSANDCLEARTAGSCONFLICT = "InvalidParameterValue.AddTagsAndClearTagsConflict"
//  INVALIDPARAMETERVALUE_ADDTAGSANDDELETETAGSCONFLICT = "InvalidParameterValue.AddTagsAndDeleteTagsConflict"
//  INVALIDPARAMETERVALUE_DESCRIPTION = "InvalidParameterValue.Description"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_IMAGEDECODEERROR = "InvalidParameterValue.ImageDecodeError"
//  INVALIDPARAMETERVALUE_KEYFRAMEDESCCONTENTTOOLONG = "InvalidParameterValue.KeyFrameDescContentTooLong"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGTOOLONG = "InvalidParameterValue.TagTooLong"
//  LIMITEXCEEDED_KEYFRAMEDESCCOUNTREACHMAX = "LimitExceeded.KeyFrameDescCountReachMax"
//  LIMITEXCEEDED_TAGCOUNTREACHMAX = "LimitExceeded.TagCountReachMax"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyMediaInfo(request *ModifyMediaInfoRequest) (response *ModifyMediaInfoResponse, err error) {
    return c.ModifyMediaInfoWithContext(context.Background(), request)
}

// ModifyMediaInfo
// 修改媒体文件的属性，包括分类、名称、描述、标签、过期时间、打点信息、视频封面、字幕信息等。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INTERNALERROR_UPDATEMEDIAERROR = "InternalError.UpdateMediaError"
//  INTERNALERROR_UPLOADCOVERIMAGEERROR = "InternalError.UploadCoverImageError"
//  INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDCLEARKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndClearKeyFrameDescsConflict"
//  INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDDELETEKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndDeleteKeyFrameDescsConflict"
//  INVALIDPARAMETERVALUE_ADDTAGSANDCLEARTAGSCONFLICT = "InvalidParameterValue.AddTagsAndClearTagsConflict"
//  INVALIDPARAMETERVALUE_ADDTAGSANDDELETETAGSCONFLICT = "InvalidParameterValue.AddTagsAndDeleteTagsConflict"
//  INVALIDPARAMETERVALUE_DESCRIPTION = "InvalidParameterValue.Description"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_IMAGEDECODEERROR = "InvalidParameterValue.ImageDecodeError"
//  INVALIDPARAMETERVALUE_KEYFRAMEDESCCONTENTTOOLONG = "InvalidParameterValue.KeyFrameDescContentTooLong"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_TAGTOOLONG = "InvalidParameterValue.TagTooLong"
//  LIMITEXCEEDED_KEYFRAMEDESCCOUNTREACHMAX = "LimitExceeded.KeyFrameDescCountReachMax"
//  LIMITEXCEEDED_TAGCOUNTREACHMAX = "LimitExceeded.TagCountReachMax"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyMediaInfoWithContext(ctx context.Context, request *ModifyMediaInfoRequest) (response *ModifyMediaInfoResponse, err error) {
    if request == nil {
        request = NewModifyMediaInfoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyMediaInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyMediaInfoResponse()
    err = c.Send(request, response)
    return
}

func NewModifyMediaStorageClassRequest() (request *ModifyMediaStorageClassRequest) {
    request = &ModifyMediaStorageClassRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyMediaStorageClass")
    
    
    return
}

func NewModifyMediaStorageClassResponse() (response *ModifyMediaStorageClassResponse) {
    response = &ModifyMediaStorageClassResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyMediaStorageClass
// 修改媒体文件的存储类型。
//
// 当媒体文件的存储类型为标准存储时，可以修改为以下类型：
//
// <li>低频存储</li>
//
// <li>归档存储</li>
//
// <li>深度归档存储</li>
//
// 当媒体文件的当前存储类型为低频存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// <li>归档存储</li>
//
// <li>深度归档存储</li>
//
// 当媒体文件的当前存储类型为归档存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// 当媒体文件的当前存储类型为深度归档存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_ORIGINALSTORAGECLASS = "InvalidParameterValue.OriginalStorageClass"
//  INVALIDPARAMETERVALUE_STORAGECLASS = "InvalidParameterValue.StorageClass"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDRESTORETIER = "InvalidParameterValue.UnsupportedRestoreTier"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDSTORAGECLASS = "InvalidParameterValue.UnsupportedStorageClass"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDTRANSITION = "InvalidParameterValue.UnsupportedTransition"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) ModifyMediaStorageClass(request *ModifyMediaStorageClassRequest) (response *ModifyMediaStorageClassResponse, err error) {
    return c.ModifyMediaStorageClassWithContext(context.Background(), request)
}

// ModifyMediaStorageClass
// 修改媒体文件的存储类型。
//
// 当媒体文件的存储类型为标准存储时，可以修改为以下类型：
//
// <li>低频存储</li>
//
// <li>归档存储</li>
//
// <li>深度归档存储</li>
//
// 当媒体文件的当前存储类型为低频存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// <li>归档存储</li>
//
// <li>深度归档存储</li>
//
// 当媒体文件的当前存储类型为归档存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// 当媒体文件的当前存储类型为深度归档存储时，可以修改为以下类型：
//
// <li>标准存储</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_ORIGINALSTORAGECLASS = "InvalidParameterValue.OriginalStorageClass"
//  INVALIDPARAMETERVALUE_STORAGECLASS = "InvalidParameterValue.StorageClass"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDRESTORETIER = "InvalidParameterValue.UnsupportedRestoreTier"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDSTORAGECLASS = "InvalidParameterValue.UnsupportedStorageClass"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDTRANSITION = "InvalidParameterValue.UnsupportedTransition"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) ModifyMediaStorageClassWithContext(ctx context.Context, request *ModifyMediaStorageClassRequest) (response *ModifyMediaStorageClassResponse, err error) {
    if request == nil {
        request = NewModifyMediaStorageClassRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyMediaStorageClass require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyMediaStorageClassResponse()
    err = c.Send(request, response)
    return
}

func NewModifyPersonSampleRequest() (request *ModifyPersonSampleRequest) {
    request = &ModifyPersonSampleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyPersonSample")
    
    
    return
}

func NewModifyPersonSampleResponse() (response *ModifyPersonSampleResponse) {
    response = &ModifyPersonSampleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyPersonSample
// 该接口用于根据素材 ID，修改素材样本信息，包括名称、描述的修改，以及五官、标签的添加、删除、重置操作。五官删除操作需保证至少剩余 1 张图片，否则，请使用重置操作。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FACEDUPLICATE = "InvalidParameterValue.FaceDuplicate"
//  INVALIDPARAMETERVALUE_PICFORMATERROR = "InvalidParameterValue.PicFormatError"
//  RESOURCENOTFOUND_PERSON = "ResourceNotFound.Person"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyPersonSample(request *ModifyPersonSampleRequest) (response *ModifyPersonSampleResponse, err error) {
    return c.ModifyPersonSampleWithContext(context.Background(), request)
}

// ModifyPersonSample
// 该接口用于根据素材 ID，修改素材样本信息，包括名称、描述的修改，以及五官、标签的添加、删除、重置操作。五官删除操作需保证至少剩余 1 张图片，否则，请使用重置操作。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FACEDUPLICATE = "InvalidParameterValue.FaceDuplicate"
//  INVALIDPARAMETERVALUE_PICFORMATERROR = "InvalidParameterValue.PicFormatError"
//  RESOURCENOTFOUND_PERSON = "ResourceNotFound.Person"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyPersonSampleWithContext(ctx context.Context, request *ModifyPersonSampleRequest) (response *ModifyPersonSampleResponse, err error) {
    if request == nil {
        request = NewModifyPersonSampleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyPersonSample require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyPersonSampleResponse()
    err = c.Send(request, response)
    return
}

func NewModifyReviewTemplateRequest() (request *ModifyReviewTemplateRequest) {
    request = &ModifyReviewTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyReviewTemplate")
    
    
    return
}

func NewModifyReviewTemplateResponse() (response *ModifyReviewTemplateResponse) {
    response = &ModifyReviewTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyReviewTemplate
// 修改用户自定义审核模板。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELS = "InvalidParameterValue.Labels"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyReviewTemplate(request *ModifyReviewTemplateRequest) (response *ModifyReviewTemplateResponse, err error) {
    return c.ModifyReviewTemplateWithContext(context.Background(), request)
}

// ModifyReviewTemplate
// 修改用户自定义审核模板。
//
// >模板仅适用于 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 和 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_LABELS = "InvalidParameterValue.Labels"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyReviewTemplateWithContext(ctx context.Context, request *ModifyReviewTemplateRequest) (response *ModifyReviewTemplateResponse, err error) {
    if request == nil {
        request = NewModifyReviewTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyReviewTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyReviewTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifySampleSnapshotTemplateRequest() (request *ModifySampleSnapshotTemplateRequest) {
    request = &ModifySampleSnapshotTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifySampleSnapshotTemplate")
    
    
    return
}

func NewModifySampleSnapshotTemplateResponse() (response *ModifySampleSnapshotTemplateResponse) {
    response = &ModifySampleSnapshotTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifySampleSnapshotTemplate
// 修改用户自定义采样截图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySampleSnapshotTemplate(request *ModifySampleSnapshotTemplateRequest) (response *ModifySampleSnapshotTemplateResponse, err error) {
    return c.ModifySampleSnapshotTemplateWithContext(context.Background(), request)
}

// ModifySampleSnapshotTemplate
// 修改用户自定义采样截图模板。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"
//  INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySampleSnapshotTemplateWithContext(ctx context.Context, request *ModifySampleSnapshotTemplateRequest) (response *ModifySampleSnapshotTemplateResponse, err error) {
    if request == nil {
        request = NewModifySampleSnapshotTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifySampleSnapshotTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifySampleSnapshotTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifySnapshotByTimeOffsetTemplateRequest() (request *ModifySnapshotByTimeOffsetTemplateRequest) {
    request = &ModifySnapshotByTimeOffsetTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifySnapshotByTimeOffsetTemplate")
    
    
    return
}

func NewModifySnapshotByTimeOffsetTemplateResponse() (response *ModifySnapshotByTimeOffsetTemplateResponse) {
    response = &ModifySnapshotByTimeOffsetTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifySnapshotByTimeOffsetTemplate
// 修改用户自定义指定时间点截图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySnapshotByTimeOffsetTemplate(request *ModifySnapshotByTimeOffsetTemplateRequest) (response *ModifySnapshotByTimeOffsetTemplateResponse, err error) {
    return c.ModifySnapshotByTimeOffsetTemplateWithContext(context.Background(), request)
}

// ModifySnapshotByTimeOffsetTemplate
// 修改用户自定义指定时间点截图模板。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySnapshotByTimeOffsetTemplateWithContext(ctx context.Context, request *ModifySnapshotByTimeOffsetTemplateRequest) (response *ModifySnapshotByTimeOffsetTemplateResponse, err error) {
    if request == nil {
        request = NewModifySnapshotByTimeOffsetTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifySnapshotByTimeOffsetTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifySnapshotByTimeOffsetTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifySubAppIdInfoRequest() (request *ModifySubAppIdInfoRequest) {
    request = &ModifySubAppIdInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifySubAppIdInfo")
    
    
    return
}

func NewModifySubAppIdInfoResponse() (response *ModifySubAppIdInfoResponse) {
    response = &ModifySubAppIdInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifySubAppIdInfo
// 该接口用于修改子应用信息，但不允许修改主应用信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySubAppIdInfo(request *ModifySubAppIdInfoRequest) (response *ModifySubAppIdInfoResponse, err error) {
    return c.ModifySubAppIdInfoWithContext(context.Background(), request)
}

// ModifySubAppIdInfo
// 该接口用于修改子应用信息，但不允许修改主应用信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySubAppIdInfoWithContext(ctx context.Context, request *ModifySubAppIdInfoRequest) (response *ModifySubAppIdInfoResponse, err error) {
    if request == nil {
        request = NewModifySubAppIdInfoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifySubAppIdInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifySubAppIdInfoResponse()
    err = c.Send(request, response)
    return
}

func NewModifySubAppIdStatusRequest() (request *ModifySubAppIdStatusRequest) {
    request = &ModifySubAppIdStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifySubAppIdStatus")
    
    
    return
}

func NewModifySubAppIdStatusResponse() (response *ModifySubAppIdStatusResponse) {
    response = &ModifySubAppIdStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifySubAppIdStatus
// 该接口用于启用、停用子应用。被停用的子应用将封停对应域名，并限制控制台访问。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySubAppIdStatus(request *ModifySubAppIdStatusRequest) (response *ModifySubAppIdStatusResponse, err error) {
    return c.ModifySubAppIdStatusWithContext(context.Background(), request)
}

// ModifySubAppIdStatus
// 该接口用于启用、停用子应用。被停用的子应用将封停对应域名，并限制控制台访问。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySubAppIdStatusWithContext(ctx context.Context, request *ModifySubAppIdStatusRequest) (response *ModifySubAppIdStatusResponse, err error) {
    if request == nil {
        request = NewModifySubAppIdStatusRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifySubAppIdStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifySubAppIdStatusResponse()
    err = c.Send(request, response)
    return
}

func NewModifySuperPlayerConfigRequest() (request *ModifySuperPlayerConfigRequest) {
    request = &ModifySuperPlayerConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifySuperPlayerConfig")
    
    
    return
}

func NewModifySuperPlayerConfigResponse() (response *ModifySuperPlayerConfigResponse) {
    response = &ModifySuperPlayerConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifySuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 修改播放器配置。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySuperPlayerConfig(request *ModifySuperPlayerConfigRequest) (response *ModifySuperPlayerConfigResponse, err error) {
    return c.ModifySuperPlayerConfigWithContext(context.Background(), request)
}

// ModifySuperPlayerConfig
// 该 API 已经<font color='red'>不再维护</font>，新版播放器签名不再使用播放器配置模板，详细请参考 [播放器签名](https://cloud.tencent.com/document/product/266/45554)。
//
// 修改播放器配置。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifySuperPlayerConfigWithContext(ctx context.Context, request *ModifySuperPlayerConfigRequest) (response *ModifySuperPlayerConfigResponse, err error) {
    if request == nil {
        request = NewModifySuperPlayerConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifySuperPlayerConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifySuperPlayerConfigResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTranscodeTemplateRequest() (request *ModifyTranscodeTemplateRequest) {
    request = &ModifyTranscodeTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyTranscodeTemplate")
    
    
    return
}

func NewModifyTranscodeTemplateResponse() (response *ModifyTranscodeTemplateResponse) {
    response = &ModifyTranscodeTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyTranscodeTemplate
// 修改用户自定义转码模板信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_AUDIOBITRATE = "InvalidParameterValue.AudioBitrate"
//  INVALIDPARAMETERVALUE_AUDIOCHANNEL = "InvalidParameterValue.AudioChannel"
//  INVALIDPARAMETERVALUE_AUDIOCODEC = "InvalidParameterValue.AudioCodec"
//  INVALIDPARAMETERVALUE_AUDIOSAMPLERATE = "InvalidParameterValue.AudioSampleRate"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_CONTAINER = "InvalidParameterValue.Container"
//  INVALIDPARAMETERVALUE_FILTRATEAUDIO = "InvalidParameterValue.FiltrateAudio"
//  INVALIDPARAMETERVALUE_FILTRATEVIDEO = "InvalidParameterValue.FiltrateVideo"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_REMOVEVIDEO = "InvalidParameterValue.RemoveVideo"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_RESOLUTIONADAPTIVE = "InvalidParameterValue.ResolutionAdaptive"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_VCRF = "InvalidParameterValue.Vcrf"
//  INVALIDPARAMETERVALUE_VIDEOBITRATE = "InvalidParameterValue.VideoBitrate"
//  INVALIDPARAMETERVALUE_VIDEOCODEC = "InvalidParameterValue.VideoCodec"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTranscodeTemplate(request *ModifyTranscodeTemplateRequest) (response *ModifyTranscodeTemplateResponse, err error) {
    return c.ModifyTranscodeTemplateWithContext(context.Background(), request)
}

// ModifyTranscodeTemplate
// 修改用户自定义转码模板信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_AUDIOBITRATE = "InvalidParameterValue.AudioBitrate"
//  INVALIDPARAMETERVALUE_AUDIOCHANNEL = "InvalidParameterValue.AudioChannel"
//  INVALIDPARAMETERVALUE_AUDIOCODEC = "InvalidParameterValue.AudioCodec"
//  INVALIDPARAMETERVALUE_AUDIOSAMPLERATE = "InvalidParameterValue.AudioSampleRate"
//  INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"
//  INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"
//  INVALIDPARAMETERVALUE_CONTAINER = "InvalidParameterValue.Container"
//  INVALIDPARAMETERVALUE_FILTRATEAUDIO = "InvalidParameterValue.FiltrateAudio"
//  INVALIDPARAMETERVALUE_FILTRATEVIDEO = "InvalidParameterValue.FiltrateVideo"
//  INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"
//  INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"
//  INVALIDPARAMETERVALUE_REMOVEVIDEO = "InvalidParameterValue.RemoveVideo"
//  INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"
//  INVALIDPARAMETERVALUE_RESOLUTIONADAPTIVE = "InvalidParameterValue.ResolutionAdaptive"
//  INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"
//  INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_VCRF = "InvalidParameterValue.Vcrf"
//  INVALIDPARAMETERVALUE_VIDEOBITRATE = "InvalidParameterValue.VideoBitrate"
//  INVALIDPARAMETERVALUE_VIDEOCODEC = "InvalidParameterValue.VideoCodec"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTranscodeTemplateWithContext(ctx context.Context, request *ModifyTranscodeTemplateRequest) (response *ModifyTranscodeTemplateResponse, err error) {
    if request == nil {
        request = NewModifyTranscodeTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyTranscodeTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyTranscodeTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyVodDomainAccelerateConfigRequest() (request *ModifyVodDomainAccelerateConfigRequest) {
    request = &ModifyVodDomainAccelerateConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyVodDomainAccelerateConfig")
    
    
    return
}

func NewModifyVodDomainAccelerateConfigResponse() (response *ModifyVodDomainAccelerateConfigResponse) {
    response = &ModifyVodDomainAccelerateConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyVodDomainAccelerateConfig
// 该接口用于修改点播域名的加速区域。
//
// 1、域名部署状态为 Online 状态时才允许修改加速区域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINDEPLOYING = "FailedOperation.DomainDeploying"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) ModifyVodDomainAccelerateConfig(request *ModifyVodDomainAccelerateConfigRequest) (response *ModifyVodDomainAccelerateConfigResponse, err error) {
    return c.ModifyVodDomainAccelerateConfigWithContext(context.Background(), request)
}

// ModifyVodDomainAccelerateConfig
// 该接口用于修改点播域名的加速区域。
//
// 1、域名部署状态为 Online 状态时才允许修改加速区域。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINDEPLOYING = "FailedOperation.DomainDeploying"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) ModifyVodDomainAccelerateConfigWithContext(ctx context.Context, request *ModifyVodDomainAccelerateConfigRequest) (response *ModifyVodDomainAccelerateConfigResponse, err error) {
    if request == nil {
        request = NewModifyVodDomainAccelerateConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyVodDomainAccelerateConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyVodDomainAccelerateConfigResponse()
    err = c.Send(request, response)
    return
}

func NewModifyVodDomainConfigRequest() (request *ModifyVodDomainConfigRequest) {
    request = &ModifyVodDomainConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyVodDomainConfig")
    
    
    return
}

func NewModifyVodDomainConfigResponse() (response *ModifyVodDomainConfigResponse) {
    response = &ModifyVodDomainConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyVodDomainConfig
// 该接口用于修改域名配置，可以修改域名的防盗链配置。
//
// 1、域名部署状态为 Online 状态时才允许修改域名的配置。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINDEPLOYING = "FailedOperation.DomainDeploying"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) ModifyVodDomainConfig(request *ModifyVodDomainConfigRequest) (response *ModifyVodDomainConfigResponse, err error) {
    return c.ModifyVodDomainConfigWithContext(context.Background(), request)
}

// ModifyVodDomainConfig
// 该接口用于修改域名配置，可以修改域名的防盗链配置。
//
// 1、域名部署状态为 Online 状态时才允许修改域名的配置。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINDEPLOYING = "FailedOperation.DomainDeploying"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) ModifyVodDomainConfigWithContext(ctx context.Context, request *ModifyVodDomainConfigRequest) (response *ModifyVodDomainConfigResponse, err error) {
    if request == nil {
        request = NewModifyVodDomainConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyVodDomainConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyVodDomainConfigResponse()
    err = c.Send(request, response)
    return
}

func NewModifyWatermarkTemplateRequest() (request *ModifyWatermarkTemplateRequest) {
    request = &ModifyWatermarkTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyWatermarkTemplate")
    
    
    return
}

func NewModifyWatermarkTemplateResponse() (response *ModifyWatermarkTemplateResponse) {
    response = &ModifyWatermarkTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyWatermarkTemplate
// 修改用户自定义水印模板，水印类型不允许修改。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_UPLOADWATERMARKERROR = "InternalError.UploadWatermarkError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COORDINATEORIGIN = "InvalidParameterValue.CoordinateOrigin"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_IMAGECONTENT = "InvalidParameterValue.ImageContent"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REPEATTYPE = "InvalidParameterValue.RepeatType"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEHEIGHT = "InvalidParameterValue.SvgTemplateHeight"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEWIDTH = "InvalidParameterValue.SvgTemplateWidth"
//  INVALIDPARAMETERVALUE_TEXTALPHA = "InvalidParameterValue.TextAlpha"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  INVALIDPARAMETERVALUE_XPOS = "InvalidParameterValue.XPos"
//  INVALIDPARAMETERVALUE_YPOS = "InvalidParameterValue.YPos"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyWatermarkTemplate(request *ModifyWatermarkTemplateRequest) (response *ModifyWatermarkTemplateResponse, err error) {
    return c.ModifyWatermarkTemplateWithContext(context.Background(), request)
}

// ModifyWatermarkTemplate
// 修改用户自定义水印模板，水印类型不允许修改。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_UPLOADWATERMARKERROR = "InternalError.UploadWatermarkError"
//  INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"
//  INVALIDPARAMETERVALUE_COORDINATEORIGIN = "InvalidParameterValue.CoordinateOrigin"
//  INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"
//  INVALIDPARAMETERVALUE_IMAGECONTENT = "InvalidParameterValue.ImageContent"
//  INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"
//  INVALIDPARAMETERVALUE_REPEATTYPE = "InvalidParameterValue.RepeatType"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEHEIGHT = "InvalidParameterValue.SvgTemplateHeight"
//  INVALIDPARAMETERVALUE_SVGTEMPLATEWIDTH = "InvalidParameterValue.SvgTemplateWidth"
//  INVALIDPARAMETERVALUE_TEXTALPHA = "InvalidParameterValue.TextAlpha"
//  INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"
//  INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"
//  INVALIDPARAMETERVALUE_XPOS = "InvalidParameterValue.XPos"
//  INVALIDPARAMETERVALUE_YPOS = "InvalidParameterValue.YPos"
//  RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyWatermarkTemplateWithContext(ctx context.Context, request *ModifyWatermarkTemplateRequest) (response *ModifyWatermarkTemplateResponse, err error) {
    if request == nil {
        request = NewModifyWatermarkTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyWatermarkTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyWatermarkTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewModifyWordSampleRequest() (request *ModifyWordSampleRequest) {
    request = &ModifyWordSampleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ModifyWordSample")
    
    
    return
}

func NewModifyWordSampleResponse() (response *ModifyWordSampleResponse) {
    response = &ModifyWordSampleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyWordSample
// 该接口用于修改关键词的应用场景、标签，关键词本身不可修改，如需修改，可删除重建。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_WORD = "ResourceNotFound.Word"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyWordSample(request *ModifyWordSampleRequest) (response *ModifyWordSampleResponse, err error) {
    return c.ModifyWordSampleWithContext(context.Background(), request)
}

// ModifyWordSample
// 该接口用于修改关键词的应用场景、标签，关键词本身不可修改，如需修改，可删除重建。
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND_WORD = "ResourceNotFound.Word"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyWordSampleWithContext(ctx context.Context, request *ModifyWordSampleRequest) (response *ModifyWordSampleResponse, err error) {
    if request == nil {
        request = NewModifyWordSampleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyWordSample require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyWordSampleResponse()
    err = c.Send(request, response)
    return
}

func NewParseStreamingManifestRequest() (request *ParseStreamingManifestRequest) {
    request = &ParseStreamingManifestRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ParseStreamingManifest")
    
    
    return
}

func NewParseStreamingManifestResponse() (response *ParseStreamingManifestResponse) {
    response = &ParseStreamingManifestResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ParseStreamingManifest
// 上传 HLS 视频时，解析索引文件内容，返回待上传的分片文件列表。分片文件路径必须是当前目录或子目录的相对路径，不能是 URL，不能是绝对路径。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_MEDIAMANIFESTCONTENT = "InvalidParameterValue.MediaManifestContent"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ParseStreamingManifest(request *ParseStreamingManifestRequest) (response *ParseStreamingManifestResponse, err error) {
    return c.ParseStreamingManifestWithContext(context.Background(), request)
}

// ParseStreamingManifest
// 上传 HLS 视频时，解析索引文件内容，返回待上传的分片文件列表。分片文件路径必须是当前目录或子目录的相对路径，不能是 URL，不能是绝对路径。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_MEDIAMANIFESTCONTENT = "InvalidParameterValue.MediaManifestContent"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ParseStreamingManifestWithContext(ctx context.Context, request *ParseStreamingManifestRequest) (response *ParseStreamingManifestResponse, err error) {
    if request == nil {
        request = NewParseStreamingManifestRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ParseStreamingManifest require credential")
    }

    request.SetContext(ctx)
    
    response = NewParseStreamingManifestResponse()
    err = c.Send(request, response)
    return
}

func NewProcessImageRequest() (request *ProcessImageRequest) {
    request = &ProcessImageRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ProcessImage")
    
    
    return
}

func NewProcessImageResponse() (response *ProcessImageResponse) {
    response = &ProcessImageResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ProcessImage
// 该 API 已经<font color='red'>不再维护</font>，智能识别任务请使用图片智能识别 [ReviewImage](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 
//
// 对点播中的图片文件发起处理任务，功能包括：
//
// 
//
// 1. 智能识别（令人反感的信息、不安全的信息、不适宜的信息）;
//
// 
//
// ><li>图片文件大小支持：文件 < 5M；</li>
//
// ><li>图片文件分辨率支持：建议分辨率大于256x256，否则可能会影响识别效果；</li>
//
// ><li>图片文件支持格式：PNG、JPG、JPEG、BMP、GIF、WEBP格式。</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED = "LimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ProcessImage(request *ProcessImageRequest) (response *ProcessImageResponse, err error) {
    return c.ProcessImageWithContext(context.Background(), request)
}

// ProcessImage
// 该 API 已经<font color='red'>不再维护</font>，智能识别任务请使用图片智能识别 [ReviewImage](https://cloud.tencent.com/document/api/266/73217) 接口。
//
// 
//
// 对点播中的图片文件发起处理任务，功能包括：
//
// 
//
// 1. 智能识别（令人反感的信息、不安全的信息、不适宜的信息）;
//
// 
//
// ><li>图片文件大小支持：文件 < 5M；</li>
//
// ><li>图片文件分辨率支持：建议分辨率大于256x256，否则可能会影响识别效果；</li>
//
// ><li>图片文件支持格式：PNG、JPG、JPEG、BMP、GIF、WEBP格式。</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED = "LimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ProcessImageWithContext(ctx context.Context, request *ProcessImageRequest) (response *ProcessImageResponse, err error) {
    if request == nil {
        request = NewProcessImageRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ProcessImage require credential")
    }

    request.SetContext(ctx)
    
    response = NewProcessImageResponse()
    err = c.Send(request, response)
    return
}

func NewProcessMediaRequest() (request *ProcessMediaRequest) {
    request = &ProcessMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ProcessMedia")
    
    
    return
}

func NewProcessMediaResponse() (response *ProcessMediaResponse) {
    response = &ProcessMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ProcessMedia
// 对点播中的音视频媒体发起处理任务，功能包括：
//
// 1. 视频转码（带水印）；
//
// 2. 视频转动图；
//
// 3. 对视频按指定时间点截图；
//
// 4. 对视频采样截图；
//
// 5. 对视频截图雪碧图；
//
// 6. 对视频截取一张图做封面；
//
// 7. 对视频转自适应码流（并加密）；
//
// 8. 内容审核（令人反感的信息、不安全的信息、不适宜的信息），<font color=red>不建议</font> 使用该接口发起，推荐使用 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 或 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217)；
//
// 9. 内容分析（标签、分类、封面、按帧标签）；
//
// 10. 内容识别（视频片头片尾、人脸、文本全文、文本关键词、语音全文、语音关键词、物体）。
//
// 
//
// 如使用事件通知，事件通知的类型为 [任务流状态变更](https://cloud.tencent.com/document/product/266/9636)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NONEEDTOREDUCEMEDIABITRATE = "FailedOperation.NoNeedToReduceMediaBitrate"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_AIANALYSISTASKDEFINITION = "InvalidParameterValue.AiAnalysisTaskDefinition"
//  INVALIDPARAMETERVALUE_AICONTENTREVIEWTASKDEFINITION = "InvalidParameterValue.AiContentReviewTaskDefinition"
//  INVALIDPARAMETERVALUE_AIRECOGNITIONTASKDEFINITION = "InvalidParameterValue.AiRecognitionTaskDefinition"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMedia(request *ProcessMediaRequest) (response *ProcessMediaResponse, err error) {
    return c.ProcessMediaWithContext(context.Background(), request)
}

// ProcessMedia
// 对点播中的音视频媒体发起处理任务，功能包括：
//
// 1. 视频转码（带水印）；
//
// 2. 视频转动图；
//
// 3. 对视频按指定时间点截图；
//
// 4. 对视频采样截图；
//
// 5. 对视频截图雪碧图；
//
// 6. 对视频截取一张图做封面；
//
// 7. 对视频转自适应码流（并加密）；
//
// 8. 内容审核（令人反感的信息、不安全的信息、不适宜的信息），<font color=red>不建议</font> 使用该接口发起，推荐使用 [音视频审核(ReviewAudioVideo)](https://cloud.tencent.com/document/api/266/80283) 或 [图片审核(ReviewImage)](https://cloud.tencent.com/document/api/266/73217)；
//
// 9. 内容分析（标签、分类、封面、按帧标签）；
//
// 10. 内容识别（视频片头片尾、人脸、文本全文、文本关键词、语音全文、语音关键词、物体）。
//
// 
//
// 如使用事件通知，事件通知的类型为 [任务流状态变更](https://cloud.tencent.com/document/product/266/9636)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NONEEDTOREDUCEMEDIABITRATE = "FailedOperation.NoNeedToReduceMediaBitrate"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_AIANALYSISTASKDEFINITION = "InvalidParameterValue.AiAnalysisTaskDefinition"
//  INVALIDPARAMETERVALUE_AICONTENTREVIEWTASKDEFINITION = "InvalidParameterValue.AiContentReviewTaskDefinition"
//  INVALIDPARAMETERVALUE_AIRECOGNITIONTASKDEFINITION = "InvalidParameterValue.AiRecognitionTaskDefinition"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMediaWithContext(ctx context.Context, request *ProcessMediaRequest) (response *ProcessMediaResponse, err error) {
    if request == nil {
        request = NewProcessMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ProcessMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewProcessMediaResponse()
    err = c.Send(request, response)
    return
}

func NewProcessMediaByProcedureRequest() (request *ProcessMediaByProcedureRequest) {
    request = &ProcessMediaByProcedureRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ProcessMediaByProcedure")
    
    
    return
}

func NewProcessMediaByProcedureResponse() (response *ProcessMediaByProcedureResponse) {
    response = &ProcessMediaByProcedureResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ProcessMediaByProcedure
// 使用任务流模板，对点播中的视频发起处理任务。
//
// 有两种方式创建任务流模板：
//
// 1. 在控制台上创建和修改任务流模板；
//
// 2. 通过任务流模板接口创建任务流模板。
//
// 
//
// 如使用事件通知，除音视频审核任务外的事件通知的类型为 [任务流状态变更](https://cloud.tencent.com/document/product/266/9636)；音视频审核任务事件通知的类型为 [音视频审核完成](https://cloud.tencent.com/document/product/266/81258)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_PROCEDURENAME = "InvalidParameterValue.ProcedureName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_UNIQUEIDENTIFIER = "InvalidParameterValue.UniqueIdentifier"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMediaByProcedure(request *ProcessMediaByProcedureRequest) (response *ProcessMediaByProcedureResponse, err error) {
    return c.ProcessMediaByProcedureWithContext(context.Background(), request)
}

// ProcessMediaByProcedure
// 使用任务流模板，对点播中的视频发起处理任务。
//
// 有两种方式创建任务流模板：
//
// 1. 在控制台上创建和修改任务流模板；
//
// 2. 通过任务流模板接口创建任务流模板。
//
// 
//
// 如使用事件通知，除音视频审核任务外的事件通知的类型为 [任务流状态变更](https://cloud.tencent.com/document/product/266/9636)；音视频审核任务事件通知的类型为 [音视频审核完成](https://cloud.tencent.com/document/product/266/81258)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_PROCEDURENAME = "InvalidParameterValue.ProcedureName"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_UNIQUEIDENTIFIER = "InvalidParameterValue.UniqueIdentifier"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMediaByProcedureWithContext(ctx context.Context, request *ProcessMediaByProcedureRequest) (response *ProcessMediaByProcedureResponse, err error) {
    if request == nil {
        request = NewProcessMediaByProcedureRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ProcessMediaByProcedure require credential")
    }

    request.SetContext(ctx)
    
    response = NewProcessMediaByProcedureResponse()
    err = c.Send(request, response)
    return
}

func NewProcessMediaByUrlRequest() (request *ProcessMediaByUrlRequest) {
    request = &ProcessMediaByUrlRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ProcessMediaByUrl")
    
    
    return
}

func NewProcessMediaByUrlResponse() (response *ProcessMediaByUrlResponse) {
    response = &ProcessMediaByUrlResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ProcessMediaByUrl
// 该 API 已经<font color='red'>不再维护</font>，请使用 MPS 产品的 [ProcessMedia](https://cloud.tencent.com/document/product/862/37578) 接口，在入参 InputInfo.UrlInputInfo.Url 中指定视频 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AIANALYSISTASKDEFINITION = "InvalidParameterValue.AiAnalysisTaskDefinition"
//  INVALIDPARAMETERVALUE_AICONTENTREVIEWTASKDEFINITION = "InvalidParameterValue.AiContentReviewTaskDefinition"
//  INVALIDPARAMETERVALUE_AIRECOGNITIONTASKDEFINITION = "InvalidParameterValue.AiRecognitionTaskDefinition"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMediaByUrl(request *ProcessMediaByUrlRequest) (response *ProcessMediaByUrlResponse, err error) {
    return c.ProcessMediaByUrlWithContext(context.Background(), request)
}

// ProcessMediaByUrl
// 该 API 已经<font color='red'>不再维护</font>，请使用 MPS 产品的 [ProcessMedia](https://cloud.tencent.com/document/product/862/37578) 接口，在入参 InputInfo.UrlInputInfo.Url 中指定视频 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AIANALYSISTASKDEFINITION = "InvalidParameterValue.AiAnalysisTaskDefinition"
//  INVALIDPARAMETERVALUE_AICONTENTREVIEWTASKDEFINITION = "InvalidParameterValue.AiContentReviewTaskDefinition"
//  INVALIDPARAMETERVALUE_AIRECOGNITIONTASKDEFINITION = "InvalidParameterValue.AiRecognitionTaskDefinition"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ProcessMediaByUrlWithContext(ctx context.Context, request *ProcessMediaByUrlRequest) (response *ProcessMediaByUrlResponse, err error) {
    if request == nil {
        request = NewProcessMediaByUrlRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ProcessMediaByUrl require credential")
    }

    request.SetContext(ctx)
    
    response = NewProcessMediaByUrlResponse()
    err = c.Send(request, response)
    return
}

func NewPullEventsRequest() (request *PullEventsRequest) {
    request = &PullEventsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "PullEvents")
    
    
    return
}

func NewPullEventsResponse() (response *PullEventsResponse) {
    response = &PullEventsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// PullEvents
// * 该接口用于业务服务器以 [可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83) 的方式获取事件通知；
//
// * 接口为长轮询模式，即：如果服务端存在未消费事件，则立即返回给请求方；如果服务端没有未消费事件，则后台会将请求挂起，直到有新的事件产生为止；
//
// * 请求最多挂起5秒，建议请求方将超时时间设置为10秒；
//
// * 未被拉取的事件通知最多保留4天，超过该时限的事件通知可能会被清除；
//
// * 若该接口有事件返回，调用方必须在<font color="red">30秒</font>内调用 [确认事件通知](https://cloud.tencent.com/document/product/266/33434) 接口，确认事件通知已经处理，否则该事件通知在<font color="red">30秒</font>后会再次被拉取到。
//
// * 当前，API 每次最多可以获取16个事件通知。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PullEvents(request *PullEventsRequest) (response *PullEventsResponse, err error) {
    return c.PullEventsWithContext(context.Background(), request)
}

// PullEvents
// * 该接口用于业务服务器以 [可靠回调](https://cloud.tencent.com/document/product/266/33779#.E5.8F.AF.E9.9D.A0.E5.9B.9E.E8.B0.83) 的方式获取事件通知；
//
// * 接口为长轮询模式，即：如果服务端存在未消费事件，则立即返回给请求方；如果服务端没有未消费事件，则后台会将请求挂起，直到有新的事件产生为止；
//
// * 请求最多挂起5秒，建议请求方将超时时间设置为10秒；
//
// * 未被拉取的事件通知最多保留4天，超过该时限的事件通知可能会被清除；
//
// * 若该接口有事件返回，调用方必须在<font color="red">30秒</font>内调用 [确认事件通知](https://cloud.tencent.com/document/product/266/33434) 接口，确认事件通知已经处理，否则该事件通知在<font color="red">30秒</font>后会再次被拉取到。
//
// * 当前，API 每次最多可以获取16个事件通知。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PullEventsWithContext(ctx context.Context, request *PullEventsRequest) (response *PullEventsResponse, err error) {
    if request == nil {
        request = NewPullEventsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("PullEvents require credential")
    }

    request.SetContext(ctx)
    
    response = NewPullEventsResponse()
    err = c.Send(request, response)
    return
}

func NewPullUploadRequest() (request *PullUploadRequest) {
    request = &PullUploadRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "PullUpload")
    
    
    return
}

func NewPullUploadResponse() (response *PullUploadResponse) {
    response = &PullUploadResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// PullUpload
// 该接口用于将一个网络上的视频拉取到云点播平台。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_COVERTYPE = "FailedOperation.CoverType"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXPIRETIME = "InvalidParameter.ExpireTime"
//  INVALIDPARAMETER_STORAGEREGION = "InvalidParameter.StorageRegion"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COVERURL = "InvalidParameterValue.CoverUrl"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_MEDIAURL = "InvalidParameterValue.MediaUrl"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND_COVERURL = "ResourceNotFound.CoverUrl"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PullUpload(request *PullUploadRequest) (response *PullUploadResponse, err error) {
    return c.PullUploadWithContext(context.Background(), request)
}

// PullUpload
// 该接口用于将一个网络上的视频拉取到云点播平台。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_COVERTYPE = "FailedOperation.CoverType"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXPIRETIME = "InvalidParameter.ExpireTime"
//  INVALIDPARAMETER_STORAGEREGION = "InvalidParameter.StorageRegion"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_COVERURL = "InvalidParameterValue.CoverUrl"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_MEDIAURL = "InvalidParameterValue.MediaUrl"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"
//  RESOURCENOTFOUND_COVERURL = "ResourceNotFound.CoverUrl"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PullUploadWithContext(ctx context.Context, request *PullUploadRequest) (response *PullUploadResponse, err error) {
    if request == nil {
        request = NewPullUploadRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("PullUpload require credential")
    }

    request.SetContext(ctx)
    
    response = NewPullUploadResponse()
    err = c.Send(request, response)
    return
}

func NewPushUrlCacheRequest() (request *PushUrlCacheRequest) {
    request = &PushUrlCacheRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "PushUrlCache")
    
    
    return
}

func NewPushUrlCacheResponse() (response *PushUrlCacheResponse) {
    response = &PushUrlCacheResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// PushUrlCache
// 1. 预热指定的 URL 列表。
//
// 2. URL 的域名必须已在云点播中注册。
//
// 3. 单次请求最多指定20个 URL。
//
// 4. 默认预热配额为每天10000个 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PushUrlCache(request *PushUrlCacheRequest) (response *PushUrlCacheResponse, err error) {
    return c.PushUrlCacheWithContext(context.Background(), request)
}

// PushUrlCache
// 1. 预热指定的 URL 列表。
//
// 2. URL 的域名必须已在云点播中注册。
//
// 3. 单次请求最多指定20个 URL。
//
// 4. 默认预热配额为每天10000个 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PushUrlCacheWithContext(ctx context.Context, request *PushUrlCacheRequest) (response *PushUrlCacheResponse, err error) {
    if request == nil {
        request = NewPushUrlCacheRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("PushUrlCache require credential")
    }

    request.SetContext(ctx)
    
    response = NewPushUrlCacheResponse()
    err = c.Send(request, response)
    return
}

func NewRefreshUrlCacheRequest() (request *RefreshUrlCacheRequest) {
    request = &RefreshUrlCacheRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "RefreshUrlCache")
    
    
    return
}

func NewRefreshUrlCacheResponse() (response *RefreshUrlCacheResponse) {
    response = &RefreshUrlCacheResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// RefreshUrlCache
// 1. 刷新指定的 URL 列表。
//
// 2. URL 的域名必须已在云点播中注册。
//
// 3. 单次请求最多指定20个 URL。
//
// 4. 默认刷新配额为每天100000个 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RefreshUrlCache(request *RefreshUrlCacheRequest) (response *RefreshUrlCacheResponse, err error) {
    return c.RefreshUrlCacheWithContext(context.Background(), request)
}

// RefreshUrlCache
// 1. 刷新指定的 URL 列表。
//
// 2. URL 的域名必须已在云点播中注册。
//
// 3. 单次请求最多指定20个 URL。
//
// 4. 默认刷新配额为每天100000个 URL。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RefreshUrlCacheWithContext(ctx context.Context, request *RefreshUrlCacheRequest) (response *RefreshUrlCacheResponse, err error) {
    if request == nil {
        request = NewRefreshUrlCacheRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("RefreshUrlCache require credential")
    }

    request.SetContext(ctx)
    
    response = NewRefreshUrlCacheResponse()
    err = c.Send(request, response)
    return
}

func NewRemoveWatermarkRequest() (request *RemoveWatermarkRequest) {
    request = &RemoveWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "RemoveWatermark")
    
    
    return
}

func NewRemoveWatermarkResponse() (response *RemoveWatermarkResponse) {
    response = &RemoveWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// RemoveWatermark
// 智能去除水印
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RemoveWatermark(request *RemoveWatermarkRequest) (response *RemoveWatermarkResponse, err error) {
    return c.RemoveWatermarkWithContext(context.Background(), request)
}

// RemoveWatermark
// 智能去除水印
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RemoveWatermarkWithContext(ctx context.Context, request *RemoveWatermarkRequest) (response *RemoveWatermarkResponse, err error) {
    if request == nil {
        request = NewRemoveWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("RemoveWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewRemoveWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewResetProcedureTemplateRequest() (request *ResetProcedureTemplateRequest) {
    request = &ResetProcedureTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ResetProcedureTemplate")
    
    
    return
}

func NewResetProcedureTemplateResponse() (response *ResetProcedureTemplateResponse) {
    response = &ResetProcedureTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ResetProcedureTemplate
// 重新设置用户自定义任务流模板的内容。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PROCEDURENAMENOTEXIST = "InvalidParameter.ProcedureNameNotExist"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ResetProcedureTemplate(request *ResetProcedureTemplateRequest) (response *ResetProcedureTemplateResponse, err error) {
    return c.ResetProcedureTemplateWithContext(context.Background(), request)
}

// ResetProcedureTemplate
// 重新设置用户自定义任务流模板的内容。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PROCEDURENAMENOTEXIST = "InvalidParameter.ProcedureNameNotExist"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ResetProcedureTemplateWithContext(ctx context.Context, request *ResetProcedureTemplateRequest) (response *ResetProcedureTemplateResponse, err error) {
    if request == nil {
        request = NewResetProcedureTemplateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResetProcedureTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewResetProcedureTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewRestoreMediaRequest() (request *RestoreMediaRequest) {
    request = &RestoreMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "RestoreMedia")
    
    
    return
}

func NewRestoreMediaResponse() (response *RestoreMediaResponse) {
    response = &RestoreMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// RestoreMedia
// 当媒体文件的存储类型是归档存储或深度归档存储时，是不可访问的。如需访问，则需要调用本接口进行解冻，解冻后可访问的媒体文件是临时的，在有效期过后，则不可访问。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_NOTRESTORABLE = "InvalidParameterValue.NotRestorable"
//  INVALIDPARAMETERVALUE_ORIGINALSTORAGECLASS = "InvalidParameterValue.OriginalStorageClass"
//  INVALIDPARAMETERVALUE_RESTOREDAY = "InvalidParameterValue.RestoreDay"
//  INVALIDPARAMETERVALUE_RESTORETIER = "InvalidParameterValue.RestoreTier"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDRESTORETIER = "InvalidParameterValue.UnsupportedRestoreTier"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) RestoreMedia(request *RestoreMediaRequest) (response *RestoreMediaResponse, err error) {
    return c.RestoreMediaWithContext(context.Background(), request)
}

// RestoreMedia
// 当媒体文件的存储类型是归档存储或深度归档存储时，是不可访问的。如需访问，则需要调用本接口进行解冻，解冻后可访问的媒体文件是临时的，在有效期过后，则不可访问。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_NOTRESTORABLE = "InvalidParameterValue.NotRestorable"
//  INVALIDPARAMETERVALUE_ORIGINALSTORAGECLASS = "InvalidParameterValue.OriginalStorageClass"
//  INVALIDPARAMETERVALUE_RESTOREDAY = "InvalidParameterValue.RestoreDay"
//  INVALIDPARAMETERVALUE_RESTORETIER = "InvalidParameterValue.RestoreTier"
//  INVALIDPARAMETERVALUE_UNSUPPORTEDRESTORETIER = "InvalidParameterValue.UnsupportedRestoreTier"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) RestoreMediaWithContext(ctx context.Context, request *RestoreMediaRequest) (response *RestoreMediaResponse, err error) {
    if request == nil {
        request = NewRestoreMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("RestoreMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewRestoreMediaResponse()
    err = c.Send(request, response)
    return
}

func NewReviewAudioVideoRequest() (request *ReviewAudioVideoRequest) {
    request = &ReviewAudioVideoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ReviewAudioVideo")
    
    
    return
}

func NewReviewAudioVideoResponse() (response *ReviewAudioVideoResponse) {
    response = &ReviewAudioVideoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ReviewAudioVideo
// 对点播中的音视频媒体发起审核任务，智能检测视频画面、画面中的文字、语音中的文字、声音出现的违规内容。
//
// 
//
// 如使用事件通知，事件通知的类型为 [音视频审核完成](https://cloud.tencent.com/document/product/266/81258)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  LIMITEXCEEDED = "LimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ReviewAudioVideo(request *ReviewAudioVideoRequest) (response *ReviewAudioVideoResponse, err error) {
    return c.ReviewAudioVideoWithContext(context.Background(), request)
}

// ReviewAudioVideo
// 对点播中的音视频媒体发起审核任务，智能检测视频画面、画面中的文字、语音中的文字、声音出现的违规内容。
//
// 
//
// 如使用事件通知，事件通知的类型为 [音视频审核完成](https://cloud.tencent.com/document/product/266/81258)。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  LIMITEXCEEDED = "LimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ReviewAudioVideoWithContext(ctx context.Context, request *ReviewAudioVideoRequest) (response *ReviewAudioVideoResponse, err error) {
    if request == nil {
        request = NewReviewAudioVideoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ReviewAudioVideo require credential")
    }

    request.SetContext(ctx)
    
    response = NewReviewAudioVideoResponse()
    err = c.Send(request, response)
    return
}

func NewReviewImageRequest() (request *ReviewImageRequest) {
    request = &ReviewImageRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "ReviewImage")
    
    
    return
}

func NewReviewImageResponse() (response *ReviewImageResponse) {
    response = &ReviewImageResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ReviewImage
// 对点播中的图片文件发起审核（令人反感的信息、不安全的信息、不适宜的信息）任务。
//
// 
//
// ><li>图片文件大小支持：文件 < 5M；</li>
//
// ><li>图片文件分辨率支持：建议分辨率大于256x256，否则可能会影响审核效果；</li>
//
// ><li>图片文件支持格式：PNG、JPG、JPEG、BMP、GIF、WEBP格式。</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ReviewImage(request *ReviewImageRequest) (response *ReviewImageResponse, err error) {
    return c.ReviewImageWithContext(context.Background(), request)
}

// ReviewImage
// 对点播中的图片文件发起审核（令人反感的信息、不安全的信息、不适宜的信息）任务。
//
// 
//
// ><li>图片文件大小支持：文件 < 5M；</li>
//
// ><li>图片文件分辨率支持：建议分辨率大于256x256，否则可能会影响审核效果；</li>
//
// ><li>图片文件支持格式：PNG、JPG、JPEG、BMP、GIF、WEBP格式。</li>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) ReviewImageWithContext(ctx context.Context, request *ReviewImageRequest) (response *ReviewImageResponse, err error) {
    if request == nil {
        request = NewReviewImageRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ReviewImage require credential")
    }

    request.SetContext(ctx)
    
    response = NewReviewImageResponse()
    err = c.Send(request, response)
    return
}

func NewSearchMediaRequest() (request *SearchMediaRequest) {
    request = &SearchMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "SearchMedia")
    
    
    return
}

func NewSearchMediaResponse() (response *SearchMediaResponse) {
    response = &SearchMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SearchMedia
// 搜索媒体信息，支持多种条件筛选，以及支持对返回结果排序、过滤等功能，具体包括：
//
// - 指定文件 ID 集合 FileIds ，返回匹配集合中任意 ID 的媒体。
//
// - 根据多个媒体文件名 Names 或描述信息 Descriptions 进行模糊搜索。
//
// - 根据多个文件名前缀 NamePrefixes 进行搜索。
//
// - 指定分类集合 ClassIds（见输入参数），返回满足集合中任意分类的媒体。例如：媒体分类有电影、电视剧、综艺等，其中电影分类下又有子分类历史片、动作片、言情片。如果 ClassIds 指定了电影、电视剧，那么电影和电视剧下的所有子分类都会返回；而如果 ClassIds 指定的是历史片、动作片，那么只有这2个子分类下的媒体才会返回。
//
// - 指定标签集合 Tags（见输入参数），返回满足集合中任意标签的媒体。例如：媒体标签有二次元、宫斗、鬼畜，如果 Tags 指定了二次元、鬼畜2个标签，那么只要符合这2个标签中任意一个的媒体都会被检索出来。
//
// - 指定文件类型集合 Categories（见输入参数），返回满足集合中任意类型的媒体。例如：文件类型有 Video（视频）、 Audio （音频）、 Image （图片）。如果Categories指定了 Video 和 Audio 2个文件类型，那么符合这些类型的媒体都会被检索出来。
//
// - 指定来源集合 SourceTypes（见输入参数），返回满足集合中任意来源的媒体。例如：媒体来源有 Record (直播录制)、Upload （上传）等。如果 SourceTypes 指定了 Record 和 Upload ，那么符合这些来源的媒体都会被检索出来。
//
// - 指定直播推流码集合 StreamIds（见输入参数）筛选直播录制的媒体。
//
// - 指定媒体的创建时间范围筛选媒体。
//
// - 指定 TRTC 应用 ID 集合筛选媒体。
//
// - 指定 TRTC 房间 ID 集合筛选媒体。
//
// - （不推荐：应使用 Names、NamePrefixes 或 Descriptions 替代）指定单个文本 Text 对媒体文件名或描述信息进行模糊搜索。
//
// - （不推荐：应使用 SourceTypes 替代）指定单个媒体文件来源 SourceType 进行搜索。
//
// - （不推荐：应使用 StreamIds 替代）指定单个推流直播码 StreamId 进行搜索。
//
// - （不推荐：应使用 CreateTime 替代）指定单个起始创建时间 StartTime 进行搜索。
//
// - （不推荐：应使用 CreateTime 替代）指定单个结尾创建时间 EndTime 进行搜索。
//
// - 以上参数之间可以任意组合进行检索。例如：筛选创建时间在2018年12月1日12:00:00到2018年12月8日12:00:00之间、分类为电影或电视剧、带有宫斗和悬疑标签的媒体。注意，任何支持数组输入的参数，其元素之间的搜索逻辑为‘或’。所有参数之间的逻辑关系为‘与’。
//
// 
//
// - 允许对结果根据创建时间进行排序并分页返回，通过 Offset 和 Limit （见输入参数）来控制分页。
//
// - 允许通过 Filters 控制返回的媒体信息种类（默认返回所有信息）。可选输入包括：
//
//     1. 基础信息（basicInfo）：包括媒体名称、分类、播放地址、封面图片等。
//
//     2. 元信息（metaData）：包括大小、时长、视频流信息、音频流信息等。
//
//     3. 转码结果信息（transcodeInfo）：包括该媒体转码生成的各种规格的媒体地址、视频流参数、音频流参数等。
//
//     4. 转动图结果信息（animatedGraphicsInfo）：对视频转动图（如 gif）后的动图信息。
//
//     5. 采样截图信息（sampleSnapshotInfo）：对视频采样截图后的截图信息。
//
//     6. 雪碧图信息（imageSpriteInfo）：对视频截取雪碧图后的雪碧图信息。
//
//     7. 指定时间点截图信息（snapshotByTimeOffsetInfo）：对视频依照指定时间点截图后，的截图信息。
//
//     8. 视频打点信息（keyFrameDescInfo）：对视频设置的打点信息。
//
//     9. 转自适应码流信息（adaptiveDynamicStreamingInfo）：包括规格、加密类型、打包格式等相关信息。
//
// 
//
// <div id="maxResultsDesc">接口返回结果数限制：</div>
//
// - <b><a href="#p_offset">Offset</a> 和 <a href="#p_limit">Limit</a> 两个参数影响单次分页查询结果数。特别注意：当这2个值都缺省时，本接口最多只返回10条查询结果。</b>
//
// - <b>最大支持返回5000条搜索结果，超出部分不再支持查询。如果搜索结果量太大，建议使用更精细的筛选条件来减少搜索结果。</b>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETERVALUE_CATEGORIES = "InvalidParameterValue.Categories"
//  INVALIDPARAMETERVALUE_CLASSIDS = "InvalidParameterValue.ClassIds"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_NAMEPREFIXES = "InvalidParameterValue.NamePrefixes"
//  INVALIDPARAMETERVALUE_NAMES = "InvalidParameterValue.Names"
//  INVALIDPARAMETERVALUE_OFFSET = "InvalidParameterValue.Offset"
//  INVALIDPARAMETERVALUE_SORT = "InvalidParameterValue.Sort"
//  INVALIDPARAMETERVALUE_SOURCETYPE = "InvalidParameterValue.SourceType"
//  INVALIDPARAMETERVALUE_SOURCETYPES = "InvalidParameterValue.SourceTypes"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STORAGEREGIONS = "InvalidParameterValue.StorageRegions"
//  INVALIDPARAMETERVALUE_STREAMIDS = "InvalidParameterValue.StreamIds"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_TAGS = "InvalidParameterValue.Tags"
//  INVALIDPARAMETERVALUE_TEXT = "InvalidParameterValue.Text"
//  INVALIDPARAMETERVALUE_TYPES = "InvalidParameterValue.Types"
//  INVALIDPARAMETERVALUE_VIDS = "InvalidParameterValue.Vids"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SearchMedia(request *SearchMediaRequest) (response *SearchMediaResponse, err error) {
    return c.SearchMediaWithContext(context.Background(), request)
}

// SearchMedia
// 搜索媒体信息，支持多种条件筛选，以及支持对返回结果排序、过滤等功能，具体包括：
//
// - 指定文件 ID 集合 FileIds ，返回匹配集合中任意 ID 的媒体。
//
// - 根据多个媒体文件名 Names 或描述信息 Descriptions 进行模糊搜索。
//
// - 根据多个文件名前缀 NamePrefixes 进行搜索。
//
// - 指定分类集合 ClassIds（见输入参数），返回满足集合中任意分类的媒体。例如：媒体分类有电影、电视剧、综艺等，其中电影分类下又有子分类历史片、动作片、言情片。如果 ClassIds 指定了电影、电视剧，那么电影和电视剧下的所有子分类都会返回；而如果 ClassIds 指定的是历史片、动作片，那么只有这2个子分类下的媒体才会返回。
//
// - 指定标签集合 Tags（见输入参数），返回满足集合中任意标签的媒体。例如：媒体标签有二次元、宫斗、鬼畜，如果 Tags 指定了二次元、鬼畜2个标签，那么只要符合这2个标签中任意一个的媒体都会被检索出来。
//
// - 指定文件类型集合 Categories（见输入参数），返回满足集合中任意类型的媒体。例如：文件类型有 Video（视频）、 Audio （音频）、 Image （图片）。如果Categories指定了 Video 和 Audio 2个文件类型，那么符合这些类型的媒体都会被检索出来。
//
// - 指定来源集合 SourceTypes（见输入参数），返回满足集合中任意来源的媒体。例如：媒体来源有 Record (直播录制)、Upload （上传）等。如果 SourceTypes 指定了 Record 和 Upload ，那么符合这些来源的媒体都会被检索出来。
//
// - 指定直播推流码集合 StreamIds（见输入参数）筛选直播录制的媒体。
//
// - 指定媒体的创建时间范围筛选媒体。
//
// - 指定 TRTC 应用 ID 集合筛选媒体。
//
// - 指定 TRTC 房间 ID 集合筛选媒体。
//
// - （不推荐：应使用 Names、NamePrefixes 或 Descriptions 替代）指定单个文本 Text 对媒体文件名或描述信息进行模糊搜索。
//
// - （不推荐：应使用 SourceTypes 替代）指定单个媒体文件来源 SourceType 进行搜索。
//
// - （不推荐：应使用 StreamIds 替代）指定单个推流直播码 StreamId 进行搜索。
//
// - （不推荐：应使用 CreateTime 替代）指定单个起始创建时间 StartTime 进行搜索。
//
// - （不推荐：应使用 CreateTime 替代）指定单个结尾创建时间 EndTime 进行搜索。
//
// - 以上参数之间可以任意组合进行检索。例如：筛选创建时间在2018年12月1日12:00:00到2018年12月8日12:00:00之间、分类为电影或电视剧、带有宫斗和悬疑标签的媒体。注意，任何支持数组输入的参数，其元素之间的搜索逻辑为‘或’。所有参数之间的逻辑关系为‘与’。
//
// 
//
// - 允许对结果根据创建时间进行排序并分页返回，通过 Offset 和 Limit （见输入参数）来控制分页。
//
// - 允许通过 Filters 控制返回的媒体信息种类（默认返回所有信息）。可选输入包括：
//
//     1. 基础信息（basicInfo）：包括媒体名称、分类、播放地址、封面图片等。
//
//     2. 元信息（metaData）：包括大小、时长、视频流信息、音频流信息等。
//
//     3. 转码结果信息（transcodeInfo）：包括该媒体转码生成的各种规格的媒体地址、视频流参数、音频流参数等。
//
//     4. 转动图结果信息（animatedGraphicsInfo）：对视频转动图（如 gif）后的动图信息。
//
//     5. 采样截图信息（sampleSnapshotInfo）：对视频采样截图后的截图信息。
//
//     6. 雪碧图信息（imageSpriteInfo）：对视频截取雪碧图后的雪碧图信息。
//
//     7. 指定时间点截图信息（snapshotByTimeOffsetInfo）：对视频依照指定时间点截图后，的截图信息。
//
//     8. 视频打点信息（keyFrameDescInfo）：对视频设置的打点信息。
//
//     9. 转自适应码流信息（adaptiveDynamicStreamingInfo）：包括规格、加密类型、打包格式等相关信息。
//
// 
//
// <div id="maxResultsDesc">接口返回结果数限制：</div>
//
// - <b><a href="#p_offset">Offset</a> 和 <a href="#p_limit">Limit</a> 两个参数影响单次分页查询结果数。特别注意：当这2个值都缺省时，本接口最多只返回10条查询结果。</b>
//
// - <b>最大支持返回5000条搜索结果，超出部分不再支持查询。如果搜索结果量太大，建议使用更精细的筛选条件来减少搜索结果。</b>
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"
//  INVALIDPARAMETERVALUE_CATEGORIES = "InvalidParameterValue.Categories"
//  INVALIDPARAMETERVALUE_CLASSIDS = "InvalidParameterValue.ClassIds"
//  INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"
//  INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"
//  INVALIDPARAMETERVALUE_NAMEPREFIXES = "InvalidParameterValue.NamePrefixes"
//  INVALIDPARAMETERVALUE_NAMES = "InvalidParameterValue.Names"
//  INVALIDPARAMETERVALUE_OFFSET = "InvalidParameterValue.Offset"
//  INVALIDPARAMETERVALUE_SORT = "InvalidParameterValue.Sort"
//  INVALIDPARAMETERVALUE_SOURCETYPE = "InvalidParameterValue.SourceType"
//  INVALIDPARAMETERVALUE_SOURCETYPES = "InvalidParameterValue.SourceTypes"
//  INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"
//  INVALIDPARAMETERVALUE_STORAGEREGIONS = "InvalidParameterValue.StorageRegions"
//  INVALIDPARAMETERVALUE_STREAMIDS = "InvalidParameterValue.StreamIds"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  INVALIDPARAMETERVALUE_TAGS = "InvalidParameterValue.Tags"
//  INVALIDPARAMETERVALUE_TEXT = "InvalidParameterValue.Text"
//  INVALIDPARAMETERVALUE_TYPES = "InvalidParameterValue.Types"
//  INVALIDPARAMETERVALUE_VIDS = "InvalidParameterValue.Vids"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SearchMediaWithContext(ctx context.Context, request *SearchMediaRequest) (response *SearchMediaResponse, err error) {
    if request == nil {
        request = NewSearchMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SearchMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewSearchMediaResponse()
    err = c.Send(request, response)
    return
}

func NewSetDrmKeyProviderInfoRequest() (request *SetDrmKeyProviderInfoRequest) {
    request = &SetDrmKeyProviderInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "SetDrmKeyProviderInfo")
    
    
    return
}

func NewSetDrmKeyProviderInfoResponse() (response *SetDrmKeyProviderInfoResponse) {
    response = &SetDrmKeyProviderInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SetDrmKeyProviderInfo
// 设置 DRM 密钥提供商信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) SetDrmKeyProviderInfo(request *SetDrmKeyProviderInfoRequest) (response *SetDrmKeyProviderInfoResponse, err error) {
    return c.SetDrmKeyProviderInfoWithContext(context.Background(), request)
}

// SetDrmKeyProviderInfo
// 设置 DRM 密钥提供商信息。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
func (c *Client) SetDrmKeyProviderInfoWithContext(ctx context.Context, request *SetDrmKeyProviderInfoRequest) (response *SetDrmKeyProviderInfoResponse, err error) {
    if request == nil {
        request = NewSetDrmKeyProviderInfoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SetDrmKeyProviderInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewSetDrmKeyProviderInfoResponse()
    err = c.Send(request, response)
    return
}

func NewSimpleHlsClipRequest() (request *SimpleHlsClipRequest) {
    request = &SimpleHlsClipRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "SimpleHlsClip")
    
    
    return
}

func NewSimpleHlsClipResponse() (response *SimpleHlsClipResponse) {
    response = &SimpleHlsClipResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SimpleHlsClip
// 对 HLS 视频进行按时间段裁剪，实时生成一个新的视频（HLS 格式），开发者可以将其立即分享出去，或者长久保存起来。
//
// 
//
// 腾讯云点播支持两种剪辑模式：
//
// - 剪辑固化：将剪辑出来的视频保存成独立的视频，拥有独立 FileId；适用于将精彩片段长久保存的场景；
//
// - 剪辑不固化：剪辑得到的视频附属于输入文件，没有独立 FileId；适用于将精彩片段临时分享的场景。
//
// 
//
// 注意：
//
// - 剪辑是基于输入 m3u8 文件进行的，故而其最小剪辑精度为一个 ts 切片，无法实现秒级或者更为精确的剪辑精度。
//
// 
//
// 
//
// ### 剪辑固化
//
// 所谓剪辑固化，是指将剪辑出来的视频保存成一个独立的视频（拥有独立的 FileId）。其生命周期不受原始输入视频影响（即使原始输入视频被删除，剪辑结果也不会受到任何影响）；也可以对其进行转码、微信发布等二次处理。
//
// 
//
// 举例如下：一场完整的足球比赛，原始视频可能长达 2 个小时，客户出于节省成本的目的可以对这个视频存储 2 个月，但对于剪辑的「精彩时刻」视频却可以指定存储更长时间，同时可以单独对「精彩时刻」视频进行转码、微信发布等额外的点播操作，这时候可以选择剪辑并且固化的方案。
//
// 
//
// 剪辑固化的优势在于其生命周期与原始输入视频相互独立，可以独立管理、长久保存。
//
// 
//
// ### 剪辑不固化
//
// 所谓剪辑不固化，是指剪辑所得到的结果（m3u8 文件）与原始输入视频共享相同的 ts 分片，新生成的视频不是一个独立完整的视频（没有独立 FileId，只有播放 URL），其有效期与原始输入的完整视频有效期是一致的。一旦原始输入的视频被删除，也会导致该片段无法播放。
//
// 
//
// 剪辑不固化，由于其剪辑结果不是一个独立的视频，因而也不会纳入点播媒资视频管理（例如控制台的视频总数不会统计这一片段）中，也无法单独针对这个片段做转码、微信发布等任何视频处理操作。
//
// 
//
// 剪辑不固化的优势在于其剪辑操作十分“轻量化”，不会产生额外的存储开销。但其不足之处在于生命周期与原始录制视频相同，且无法进一步进行转码等视频处理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_ENDTIMEOFFSET = "InvalidParameterValue.EndTimeOffset"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_STARTTIMEOFFSET = "InvalidParameterValue.StartTimeOffset"
//  INVALIDPARAMETERVALUE_URL = "InvalidParameterValue.Url"
//  RESOURCEUNAVAILABLE_MASTERPLAYLIST = "ResourceUnavailable.MasterPlaylist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SimpleHlsClip(request *SimpleHlsClipRequest) (response *SimpleHlsClipResponse, err error) {
    return c.SimpleHlsClipWithContext(context.Background(), request)
}

// SimpleHlsClip
// 对 HLS 视频进行按时间段裁剪，实时生成一个新的视频（HLS 格式），开发者可以将其立即分享出去，或者长久保存起来。
//
// 
//
// 腾讯云点播支持两种剪辑模式：
//
// - 剪辑固化：将剪辑出来的视频保存成独立的视频，拥有独立 FileId；适用于将精彩片段长久保存的场景；
//
// - 剪辑不固化：剪辑得到的视频附属于输入文件，没有独立 FileId；适用于将精彩片段临时分享的场景。
//
// 
//
// 注意：
//
// - 剪辑是基于输入 m3u8 文件进行的，故而其最小剪辑精度为一个 ts 切片，无法实现秒级或者更为精确的剪辑精度。
//
// 
//
// 
//
// ### 剪辑固化
//
// 所谓剪辑固化，是指将剪辑出来的视频保存成一个独立的视频（拥有独立的 FileId）。其生命周期不受原始输入视频影响（即使原始输入视频被删除，剪辑结果也不会受到任何影响）；也可以对其进行转码、微信发布等二次处理。
//
// 
//
// 举例如下：一场完整的足球比赛，原始视频可能长达 2 个小时，客户出于节省成本的目的可以对这个视频存储 2 个月，但对于剪辑的「精彩时刻」视频却可以指定存储更长时间，同时可以单独对「精彩时刻」视频进行转码、微信发布等额外的点播操作，这时候可以选择剪辑并且固化的方案。
//
// 
//
// 剪辑固化的优势在于其生命周期与原始输入视频相互独立，可以独立管理、长久保存。
//
// 
//
// ### 剪辑不固化
//
// 所谓剪辑不固化，是指剪辑所得到的结果（m3u8 文件）与原始输入视频共享相同的 ts 分片，新生成的视频不是一个独立完整的视频（没有独立 FileId，只有播放 URL），其有效期与原始输入的完整视频有效期是一致的。一旦原始输入的视频被删除，也会导致该片段无法播放。
//
// 
//
// 剪辑不固化，由于其剪辑结果不是一个独立的视频，因而也不会纳入点播媒资视频管理（例如控制台的视频总数不会统计这一片段）中，也无法单独针对这个片段做转码、微信发布等任何视频处理操作。
//
// 
//
// 剪辑不固化的优势在于其剪辑操作十分“轻量化”，不会产生额外的存储开销。但其不足之处在于生命周期与原始录制视频相同，且无法进一步进行转码等视频处理。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_ENDTIMEOFFSET = "InvalidParameterValue.EndTimeOffset"
//  INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"
//  INVALIDPARAMETERVALUE_STARTTIMEOFFSET = "InvalidParameterValue.StartTimeOffset"
//  INVALIDPARAMETERVALUE_URL = "InvalidParameterValue.Url"
//  RESOURCEUNAVAILABLE_MASTERPLAYLIST = "ResourceUnavailable.MasterPlaylist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SimpleHlsClipWithContext(ctx context.Context, request *SimpleHlsClipRequest) (response *SimpleHlsClipResponse, err error) {
    if request == nil {
        request = NewSimpleHlsClipRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SimpleHlsClip require credential")
    }

    request.SetContext(ctx)
    
    response = NewSimpleHlsClipResponse()
    err = c.Send(request, response)
    return
}

func NewSplitMediaRequest() (request *SplitMediaRequest) {
    request = &SplitMediaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "SplitMedia")
    
    
    return
}

func NewSplitMediaResponse() (response *SplitMediaResponse) {
    response = &SplitMediaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SplitMedia
// 对点播视频进行拆条，生成多个新的点播视频。
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SplitMedia(request *SplitMediaRequest) (response *SplitMediaResponse, err error) {
    return c.SplitMediaWithContext(context.Background(), request)
}

// SplitMedia
// 对点播视频进行拆条，生成多个新的点播视频。
//
// 可能返回的错误码:
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"
//  INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"
//  INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) SplitMediaWithContext(ctx context.Context, request *SplitMediaRequest) (response *SplitMediaResponse, err error) {
    if request == nil {
        request = NewSplitMediaRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SplitMedia require credential")
    }

    request.SetContext(ctx)
    
    response = NewSplitMediaResponse()
    err = c.Send(request, response)
    return
}

func NewWeChatMiniProgramPublishRequest() (request *WeChatMiniProgramPublishRequest) {
    request = &WeChatMiniProgramPublishRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("vod", APIVersion, "WeChatMiniProgramPublish")
    
    
    return
}

func NewWeChatMiniProgramPublishResponse() (response *WeChatMiniProgramPublishResponse) {
    response = &WeChatMiniProgramPublishResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// WeChatMiniProgramPublish
// 将点播视频发布到微信小程序，供微信小程序播放器播放。
//
// 本接口支持发布原始视频和转码后视频，暂不支持发布自适应码流。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NOPRIVILEGES = "FailedOperation.NoPrivileges"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_FILETYPE = "InvalidParameterValue.FileType"
//  INVALIDPARAMETERVALUE_SOURCEDEFINITION = "InvalidParameterValue.SourceDefinition"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) WeChatMiniProgramPublish(request *WeChatMiniProgramPublishRequest) (response *WeChatMiniProgramPublishResponse, err error) {
    return c.WeChatMiniProgramPublishWithContext(context.Background(), request)
}

// WeChatMiniProgramPublish
// 将点播视频发布到微信小程序，供微信小程序播放器播放。
//
// 本接口支持发布原始视频和转码后视频，暂不支持发布自适应码流。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"
//  FAILEDOPERATION_NOPRIVILEGES = "FailedOperation.NoPrivileges"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"
//  INVALIDPARAMETERVALUE_FILETYPE = "InvalidParameterValue.FileType"
//  INVALIDPARAMETERVALUE_SOURCEDEFINITION = "InvalidParameterValue.SourceDefinition"
//  INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"
//  RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) WeChatMiniProgramPublishWithContext(ctx context.Context, request *WeChatMiniProgramPublishRequest) (response *WeChatMiniProgramPublishResponse, err error) {
    if request == nil {
        request = NewWeChatMiniProgramPublishRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("WeChatMiniProgramPublish require credential")
    }

    request.SetContext(ctx)
    
    response = NewWeChatMiniProgramPublishResponse()
    err = c.Send(request, response)
    return
}
