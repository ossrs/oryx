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

const (
	// 此产品的特有错误码

	// CAM签名/鉴权错误。
	AUTHFAILURE = "AuthFailure"

	// 操作失败。
	FAILEDOPERATION = "FailedOperation"

	// 操作失败：超过分类层数限制。
	FAILEDOPERATION_CLASSLEVELLIMITEXCEEDED = "FailedOperation.ClassLevelLimitExceeded"

	// 操作失败：分类名称重复。
	FAILEDOPERATION_CLASSNAMEDUPLICATE = "FailedOperation.ClassNameDuplicate"

	// 操作失败：分类不存在。
	FAILEDOPERATION_CLASSNOFOUND = "FailedOperation.ClassNoFound"

	// 操作失败：不支持的封面类型。
	FAILEDOPERATION_COVERTYPE = "FailedOperation.CoverType"

	// 域名部署中，不能变更配置。
	FAILEDOPERATION_DOMAINDEPLOYING = "FailedOperation.DomainDeploying"

	// 用户账户异常。
	FAILEDOPERATION_INVALIDACCOUNT = "FailedOperation.InvalidAccount"

	// 没有开通点播业务。
	FAILEDOPERATION_INVALIDVODUSER = "FailedOperation.InvalidVodUser"

	// 媒体被系统封禁。
	FAILEDOPERATION_MEDIAFORBIDEDBYSYSTEM = "FailedOperation.MediaForbidedBySystem"

	// 操作失败：不支持的媒体类型。
	FAILEDOPERATION_MEDIATYPE = "FailedOperation.MediaType"

	// 网络错误。
	FAILEDOPERATION_NETWORKERROR = "FailedOperation.NetWorkError"

	// 无需做降码率处理。
	FAILEDOPERATION_NONEEDTOREDUCEMEDIABITRATE = "FailedOperation.NoNeedToReduceMediaBitrate"

	// 没有开通该接口使用权限。
	FAILEDOPERATION_NOPRIVILEGES = "FailedOperation.NoPrivileges"

	// 操作失败：父类 ID 不存在。
	FAILEDOPERATION_PARENTIDNOFOUND = "FailedOperation.ParentIdNoFound"

	// 操作失败：子类数量超过限制。
	FAILEDOPERATION_SUBCLASSLIMITEXCEEDED = "FailedOperation.SubclassLimitExceeded"

	// 操作失败：任务重复。
	FAILEDOPERATION_TASKDUPLICATE = "FailedOperation.TaskDuplicate"

	// 操作失败：上传文件到 cos 失败。
	FAILEDOPERATION_UPLOADCOSFAIL = "FailedOperation.UploadCosFail"

	// 用户已经停服。
	FAILEDOPERATION_USERSTATUSINAVLID = "FailedOperation.UserStatusInavlid"

	// 内部错误。
	INTERNALERROR = "InternalError"

	// 内部错误，访问DB失败。
	INTERNALERROR_DBERROR = "InternalError.DBError"

	// 内部错误：生成模板 ID 失败。
	INTERNALERROR_GENDEFINITION = "InternalError.GenDefinition"

	// 内部错误：获取媒体文件信息错误。
	INTERNALERROR_GETFILEINFOERROR = "InternalError.GetFileInfoError"

	// 内部错误：获取媒体列表错误。
	INTERNALERROR_GETMEDIALISTERROR = "InternalError.GetMediaListError"

	// 时间解析错误。
	INTERNALERROR_TIMEPARSEERROR = "InternalError.TimeParseError"

	// 内部错误：更新媒体文件信息错误。
	INTERNALERROR_UPDATEMEDIAERROR = "InternalError.UpdateMediaError"

	// 内部错误：上传封面图片错误。
	INTERNALERROR_UPLOADCOVERIMAGEERROR = "InternalError.UploadCoverImageError"

	// 内部错误：上传水印图片失败。
	INTERNALERROR_UPLOADWATERMARKERROR = "InternalError.UploadWatermarkError"

	// 参数错误。
	INVALIDPARAMETER = "InvalidParameter"

	// 任务流模板名字已存在。
	INVALIDPARAMETER_EXISTEDPROCEDURENAME = "InvalidParameter.ExistedProcedureName"

	// 参数值错误：过期时间。
	INVALIDPARAMETER_EXPIRETIME = "InvalidParameter.ExpireTime"

	// 参数错误：标签错误。
	INVALIDPARAMETER_LABELS = "InvalidParameter.Labels"

	// 任务流模板名字不存在。
	INVALIDPARAMETER_PROCEDURENAMENOTEXIST = "InvalidParameter.ProcedureNameNotExist"

	// 参数值错误：存储地域。
	INVALIDPARAMETER_STORAGEREGION = "InvalidParameter.StorageRegion"

	// 参数取值错误。
	INVALIDPARAMETERVALUE = "InvalidParameterValue"

	// 参数值错误：AddKeyFrameDescs 与 ClearKeyFrameDescs 参数冲突。
	INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDCLEARKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndClearKeyFrameDescsConflict"

	// 参数值错误：AddKeyFrameDescs 与 DeleteKeyFrameDescs 参数冲突。
	INVALIDPARAMETERVALUE_ADDKEYFRAMEDESCSANDDELETEKEYFRAMEDESCSCONFLICT = "InvalidParameterValue.AddKeyFrameDescsAndDeleteKeyFrameDescsConflict"

	// 参数值错误：AddTags 与 ClearTags 参数冲突。
	INVALIDPARAMETERVALUE_ADDTAGSANDCLEARTAGSCONFLICT = "InvalidParameterValue.AddTagsAndClearTagsConflict"

	// 参数值错误：AddTags 与 DeleteTags 参数冲突。
	INVALIDPARAMETERVALUE_ADDTAGSANDDELETETAGSCONFLICT = "InvalidParameterValue.AddTagsAndDeleteTagsConflict"

	// 参数值错误：AI 分析 Definition。
	INVALIDPARAMETERVALUE_AIANALYSISTASKDEFINITION = "InvalidParameterValue.AiAnalysisTaskDefinition"

	// 参数值错误：AI 内容审核 Definition。
	INVALIDPARAMETERVALUE_AICONTENTREVIEWTASKDEFINITION = "InvalidParameterValue.AiContentReviewTaskDefinition"

	// 参数值错误：AI 识别 Definition。
	INVALIDPARAMETERVALUE_AIRECOGNITIONTASKDEFINITION = "InvalidParameterValue.AiRecognitionTaskDefinition"

	// Area 参数错误。
	INVALIDPARAMETERVALUE_AREA = "InvalidParameterValue.Area"

	// 参数错误：音频流码率。
	INVALIDPARAMETERVALUE_AUDIOBITRATE = "InvalidParameterValue.AudioBitrate"

	// 参数值错误：AudioChannel。
	INVALIDPARAMETERVALUE_AUDIOCHANNEL = "InvalidParameterValue.AudioChannel"

	// 参数错误：音频流编码格式。
	INVALIDPARAMETERVALUE_AUDIOCODEC = "InvalidParameterValue.AudioCodec"

	// 参数错误：音频流采样率。
	INVALIDPARAMETERVALUE_AUDIOSAMPLERATE = "InvalidParameterValue.AudioSampleRate"

	// 无效的音频/视频码率。
	INVALIDPARAMETERVALUE_BITRATE = "InvalidParameterValue.Bitrate"

	// 参数值错误：BlockConfidence 参数取值非法。
	INVALIDPARAMETERVALUE_BLOCKCONFIDENCE = "InvalidParameterValue.BlockConfidence"

	// 无效的文件类型。
	INVALIDPARAMETERVALUE_CATEGORIES = "InvalidParameterValue.Categories"

	// 参数值错误：分类 ID。
	INVALIDPARAMETERVALUE_CLASSID = "InvalidParameterValue.ClassId"

	// 参数值错误：ClassIds 无效。
	INVALIDPARAMETERVALUE_CLASSIDS = "InvalidParameterValue.ClassIds"

	// 参数值错误：ClassName 无效。
	INVALIDPARAMETERVALUE_CLASSNAME = "InvalidParameterValue.ClassName"

	// 智能分类控制字段参数错误。
	INVALIDPARAMETERVALUE_CLASSIFCATIONCONFIGURE = "InvalidParameterValue.ClassifcationConfigure"

	// 参数值错误：裁剪时间段太长。
	INVALIDPARAMETERVALUE_CLIPDURATION = "InvalidParameterValue.ClipDuration"

	// 无效的音频/视频编码格式。
	INVALIDPARAMETERVALUE_CODEC = "InvalidParameterValue.Codec"

	// 参数值错误：ColumnCount。
	INVALIDPARAMETERVALUE_COLUMNCOUNT = "InvalidParameterValue.ColumnCount"

	// 参数错误：对该模板的描述。
	INVALIDPARAMETERVALUE_COMMENT = "InvalidParameterValue.Comment"

	// 参数错误：封装格式。
	INVALIDPARAMETERVALUE_CONTAINER = "InvalidParameterValue.Container"

	// 参数值错误：ContainerType。
	INVALIDPARAMETERVALUE_CONTAINERTYPE = "InvalidParameterValue.ContainerType"

	// 参数值错误：CoordinateOrigin。
	INVALIDPARAMETERVALUE_COORDINATEORIGIN = "InvalidParameterValue.CoordinateOrigin"

	// 智能封面控制字段参数错误。
	INVALIDPARAMETERVALUE_COVERCONFIGURE = "InvalidParameterValue.CoverConfigure"

	// 参数值错误：封面类型。
	INVALIDPARAMETERVALUE_COVERTYPE = "InvalidParameterValue.CoverType"

	// 参数值错误：封面 URL。
	INVALIDPARAMETERVALUE_COVERURL = "InvalidParameterValue.CoverUrl"

	// 参数值错误：CutAndCrops 参数取值非法。
	INVALIDPARAMETERVALUE_CUTANDCROPS = "InvalidParameterValue.CutAndCrops"

	// 参数值错误，时间粒度。
	INVALIDPARAMETERVALUE_DATAINTERVAL = "InvalidParameterValue.DataInterval"

	// 参数值错误，数据类型。
	INVALIDPARAMETERVALUE_DATATYPE = "InvalidParameterValue.DataType"

	// 参数值错误：Date。
	INVALIDPARAMETERVALUE_DATE = "InvalidParameterValue.Date"

	// 参数值错误：人脸默认库过滤标签非法。
	INVALIDPARAMETERVALUE_DEFAULTLIBRARYLABELSET = "InvalidParameterValue.DefaultLibraryLabelSet"

	// 参数错误：Definition。
	INVALIDPARAMETERVALUE_DEFINITION = "InvalidParameterValue.Definition"

	// 参数错误：Definitions。
	INVALIDPARAMETERVALUE_DEFINITIONS = "InvalidParameterValue.Definitions"

	// 参数值错误：不允许删除默认模板。
	INVALIDPARAMETERVALUE_DELETEDEFAULTTEMPLATE = "InvalidParameterValue.DeleteDefaultTemplate"

	// 参数值错误：Description 超过长度限制。
	INVALIDPARAMETERVALUE_DESCRIPTION = "InvalidParameterValue.Description"

	// 无效的禁止码率低转高开关值。
	INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEOBITRATE = "InvalidParameterValue.DisableHigherVideoBitrate"

	// 无效的禁止分辨率低转高开关值。
	INVALIDPARAMETERVALUE_DISABLEHIGHERVIDEORESOLUTION = "InvalidParameterValue.DisableHigherVideoResolution"

	// Districts 参数值错误。
	INVALIDPARAMETERVALUE_DISTRICTS = "InvalidParameterValue.Districts"

	// 参数错误：不存在的域名。
	INVALIDPARAMETERVALUE_DOMAINNAME = "InvalidParameterValue.DomainName"

	// 恶意域名，无法添加。
	INVALIDPARAMETERVALUE_DOMAINNAMEINBLACKLIST = "InvalidParameterValue.DomainNameInBlackList"

	// 参数值，域名列表太大。
	INVALIDPARAMETERVALUE_DOMAINNAMES = "InvalidParameterValue.DomainNames"

	// 无效的DRM类型。
	INVALIDPARAMETERVALUE_DRMTYPE = "InvalidParameterValue.DrmType"

	// 参数值错误：EndDate 无效。
	INVALIDPARAMETERVALUE_ENDDATE = "InvalidParameterValue.EndDate"

	// 参数值错误：EndTime 无效。
	INVALIDPARAMETERVALUE_ENDTIME = "InvalidParameterValue.EndTime"

	// 参数错误：无效的结束时间。
	INVALIDPARAMETERVALUE_ENDTIMEOFFSET = "InvalidParameterValue.EndTimeOffset"

	// 参数值错误：ExpireTime 格式错误。
	INVALIDPARAMETERVALUE_EXPIRETIME = "InvalidParameterValue.ExpireTime"

	// 参数值错误：人脸重复。
	INVALIDPARAMETERVALUE_FACEDUPLICATE = "InvalidParameterValue.FaceDuplicate"

	// 参数值错误：人脸库参数非法。
	INVALIDPARAMETERVALUE_FACELIBRARY = "InvalidParameterValue.FaceLibrary"

	// 参数值错误：人脸分数参数取值非法。
	INVALIDPARAMETERVALUE_FACESCORE = "InvalidParameterValue.FaceScore"

	// FileId 不存在。
	INVALIDPARAMETERVALUE_FILEID = "InvalidParameterValue.FileId"

	// FileIds 参数错误。
	INVALIDPARAMETERVALUE_FILEIDS = "InvalidParameterValue.FileIds"

	// FileIds 数组为空。
	INVALIDPARAMETERVALUE_FILEIDSEMPTY = "InvalidParameterValue.FileIdsEmpty"

	// 参数值错误：FileId 过多。
	INVALIDPARAMETERVALUE_FILEIDSTOOMANY = "InvalidParameterValue.FileIdsTooMany"

	// 错误的视频类型。
	INVALIDPARAMETERVALUE_FILETYPE = "InvalidParameterValue.FileType"

	// 参数错误：填充方式错误。
	INVALIDPARAMETERVALUE_FILLTYPE = "InvalidParameterValue.FillType"

	// 参数错误：是否去除视频，应为0或1。
	INVALIDPARAMETERVALUE_FILTRATEAUDIO = "InvalidParameterValue.FiltrateAudio"

	// 参数错误：去除视频。
	INVALIDPARAMETERVALUE_FILTRATEVIDEO = "InvalidParameterValue.FiltrateVideo"

	// 参数值错误：Format。
	INVALIDPARAMETERVALUE_FORMAT = "InvalidParameterValue.Format"

	// 参数值错误：Format 为 webp 时，Width、Height 均为空。
	INVALIDPARAMETERVALUE_FORMATWEBPLACKWIDTHANDHEIGHT = "InvalidParameterValue.FormatWebpLackWidthAndHeight"

	// 参数值错误：Format 为 webp 时，不允许 Width、Height 都为 0。
	INVALIDPARAMETERVALUE_FORMATWEBPWIDTHANDHEIGHTBOTHZERO = "InvalidParameterValue.FormatWebpWidthAndHeightBothZero"

	// 参数错误：视频帧率。
	INVALIDPARAMETERVALUE_FPS = "InvalidParameterValue.Fps"

	// 智能按帧标签控制字段参数错误。
	INVALIDPARAMETERVALUE_FRAMETAGCONFIGURE = "InvalidParameterValue.FrameTagConfigure"

	// 参数值错误：FunctionArg。
	INVALIDPARAMETERVALUE_FUNCTIONARG = "InvalidParameterValue.FunctionArg"

	// 参数值错误：FunctionName。
	INVALIDPARAMETERVALUE_FUNCTIONNAME = "InvalidParameterValue.FunctionName"

	// 参数错误：高度。
	INVALIDPARAMETERVALUE_HEIGHT = "InvalidParameterValue.Height"

	// 智能精彩集锦控制参数错误。
	INVALIDPARAMETERVALUE_HIGHLIGHTCONFIGURE = "InvalidParameterValue.HighlightConfigure"

	// ImageContent参数值无效。
	INVALIDPARAMETERVALUE_IMAGECONTENT = "InvalidParameterValue.ImageContent"

	// 图片解 Base64 编码失败。
	INVALIDPARAMETERVALUE_IMAGEDECODEERROR = "InvalidParameterValue.ImageDecodeError"

	// 参数错误：图片水印模板。
	INVALIDPARAMETERVALUE_IMAGETEMPLATE = "InvalidParameterValue.ImageTemplate"

	// 参数值错误：Interval 无效。
	INVALIDPARAMETERVALUE_INTERVAL = "InvalidParameterValue.Interval"

	// 参数错误：无效的操作类型。
	INVALIDPARAMETERVALUE_INVALIDOPERATIONTYPE = "InvalidParameterValue.InvalidOperationType"

	// Isps 参数错误。
	INVALIDPARAMETERVALUE_ISPS = "InvalidParameterValue.Isps"

	// 参数值错误：打点信息内容过长。
	INVALIDPARAMETERVALUE_KEYFRAMEDESCCONTENTTOOLONG = "InvalidParameterValue.KeyFrameDescContentTooLong"

	// 参数值错误：LabelSet 参数取值非法。
	INVALIDPARAMETERVALUE_LABELSET = "InvalidParameterValue.LabelSet"

	// 参数错误：标签错误。
	INVALIDPARAMETERVALUE_LABELS = "InvalidParameterValue.Labels"

	// 参数错误：Limit。
	INVALIDPARAMETERVALUE_LIMIT = "InvalidParameterValue.Limit"

	// 参数值错误：Limit 过大。
	INVALIDPARAMETERVALUE_LIMITTOOLARGE = "InvalidParameterValue.LimitTooLarge"

	// 参数取值错误：MediaManifestContent。
	INVALIDPARAMETERVALUE_MEDIAMANIFESTCONTENT = "InvalidParameterValue.MediaManifestContent"

	// 参数值错误：媒体类型。
	INVALIDPARAMETERVALUE_MEDIATYPE = "InvalidParameterValue.MediaType"

	// 参数值错误：媒体文件 URL。
	INVALIDPARAMETERVALUE_MEDIAURL = "InvalidParameterValue.MediaUrl"

	// 参数错误：Metric 无效。
	INVALIDPARAMETERVALUE_METRIC = "InvalidParameterValue.Metric"

	// 参数值错误：不允许修改默认模板。
	INVALIDPARAMETERVALUE_MODIFYDEFAULTTEMPLATE = "InvalidParameterValue.ModifyDefaultTemplate"

	// 参数值错误：Name 超过长度限制。
	INVALIDPARAMETERVALUE_NAME = "InvalidParameterValue.Name"

	// 无效的文件名前缀。
	INVALIDPARAMETERVALUE_NAMEPREFIXES = "InvalidParameterValue.NamePrefixes"

	// Names数组中元素过多。
	INVALIDPARAMETERVALUE_NAMES = "InvalidParameterValue.Names"

	// 参数错误：文件不支持解冻。
	INVALIDPARAMETERVALUE_NOTRESTORABLE = "InvalidParameterValue.NotRestorable"

	// 参数值错误：物体库参数非法。
	INVALIDPARAMETERVALUE_OBJECTLIBRARY = "InvalidParameterValue.ObjectLibrary"

	// 参数值错误：Offset 无效。
	INVALIDPARAMETERVALUE_OFFSET = "InvalidParameterValue.Offset"

	// 参数值错误：Offset 过大。
	INVALIDPARAMETERVALUE_OFFSETTOOLARGE = "InvalidParameterValue.OffsetTooLarge"

	// 参数值错误：Operation 无效。
	INVALIDPARAMETERVALUE_OPERATION = "InvalidParameterValue.Operation"

	// 参数值错误：文件原存储类型异常。
	INVALIDPARAMETERVALUE_ORIGINALSTORAGECLASS = "InvalidParameterValue.OriginalStorageClass"

	// 参数值错误：ParentId 无效。
	INVALIDPARAMETERVALUE_PARENTID = "InvalidParameterValue.ParentId"

	// 参数值错误：人脸图片格式错误。
	INVALIDPARAMETERVALUE_PICFORMATERROR = "InvalidParameterValue.PicFormatError"

	// 任务流模板名无效。
	INVALIDPARAMETERVALUE_PROCEDURENAME = "InvalidParameterValue.ProcedureName"

	// 参数值错误：Quality。
	INVALIDPARAMETERVALUE_QUALITY = "InvalidParameterValue.Quality"

	// 参数值错误：RemoveAudio。
	INVALIDPARAMETERVALUE_REMOVEAUDIO = "InvalidParameterValue.RemoveAudio"

	// 参数值错误：RemoveVideo。
	INVALIDPARAMETERVALUE_REMOVEVIDEO = "InvalidParameterValue.RemoveVideo"

	// 参数错误：RepeatType 无效。
	INVALIDPARAMETERVALUE_REPEATTYPE = "InvalidParameterValue.RepeatType"

	// 参数错误：分辨率错误。
	INVALIDPARAMETERVALUE_RESOLUTION = "InvalidParameterValue.Resolution"

	// 无效的ResolutionAdaptive。
	INVALIDPARAMETERVALUE_RESOLUTIONADAPTIVE = "InvalidParameterValue.ResolutionAdaptive"

	// 参数错误：解冻天数错误。
	INVALIDPARAMETERVALUE_RESTOREDAY = "InvalidParameterValue.RestoreDay"

	// 参数错误：不支持的解冻模式。
	INVALIDPARAMETERVALUE_RESTORETIER = "InvalidParameterValue.RestoreTier"

	// 参数值错误：ReviewConfidence 参数取值非法。
	INVALIDPARAMETERVALUE_REVIEWCONFIDENCE = "InvalidParameterValue.ReviewConfidence"

	// 参数值错误：ReviewWallSwitch 参数取值非法。
	INVALIDPARAMETERVALUE_REVIEWWALLSWITCH = "InvalidParameterValue.ReviewWallSwitch"

	// 参数值错误：RowCount。
	INVALIDPARAMETERVALUE_ROWCOUNT = "InvalidParameterValue.RowCount"

	// 参数值错误：SampleInterval。
	INVALIDPARAMETERVALUE_SAMPLEINTERVAL = "InvalidParameterValue.SampleInterval"

	// 无效的音频采样率。
	INVALIDPARAMETERVALUE_SAMPLERATE = "InvalidParameterValue.SampleRate"

	// 参数值错误：SampleType。
	INVALIDPARAMETERVALUE_SAMPLETYPE = "InvalidParameterValue.SampleType"

	// 参数值错误：ScreenshotInterval 参数取值非法。
	INVALIDPARAMETERVALUE_SCREENSHOTINTERVAL = "InvalidParameterValue.ScreenshotInterval"

	// SessionContext 过长。
	INVALIDPARAMETERVALUE_SESSIONCONTEXTTOOLONG = "InvalidParameterValue.SessionContextTooLong"

	// 去重识别码重复，请求被去重。
	INVALIDPARAMETERVALUE_SESSIONID = "InvalidParameterValue.SessionId"

	// SessionId 过长。
	INVALIDPARAMETERVALUE_SESSIONIDTOOLONG = "InvalidParameterValue.SessionIdTooLong"

	// 参数值错误：Sort 无效。
	INVALIDPARAMETERVALUE_SORT = "InvalidParameterValue.Sort"

	// 参数错误：音频通道方式。
	INVALIDPARAMETERVALUE_SOUNDSYSTEM = "InvalidParameterValue.SoundSystem"

	// SourceDefinition 错误，请检查媒体文件是否有对应的转码。
	INVALIDPARAMETERVALUE_SOURCEDEFINITION = "InvalidParameterValue.SourceDefinition"

	// 参数值错误：SourceType 无效。
	INVALIDPARAMETERVALUE_SOURCETYPE = "InvalidParameterValue.SourceType"

	// 未知的媒体文件来源。
	INVALIDPARAMETERVALUE_SOURCETYPES = "InvalidParameterValue.SourceTypes"

	// 参数值错误：StartDate 无效。
	INVALIDPARAMETERVALUE_STARTDATE = "InvalidParameterValue.StartDate"

	// 参数值错误：StartTime 无效。
	INVALIDPARAMETERVALUE_STARTTIME = "InvalidParameterValue.StartTime"

	// 参数错误：无效的起始时间。
	INVALIDPARAMETERVALUE_STARTTIMEOFFSET = "InvalidParameterValue.StartTimeOffset"

	// 参数值错误：人工确认结果取值非法。
	INVALIDPARAMETERVALUE_STATUS = "InvalidParameterValue.Status"

	// 参数值错误：目标存储类型不在可选值中。
	INVALIDPARAMETERVALUE_STORAGECLASS = "InvalidParameterValue.StorageClass"

	// 参数值错误：存储地域。
	INVALIDPARAMETERVALUE_STORAGEREGION = "InvalidParameterValue.StorageRegion"

	// 参数值错误：StorageRegions 无效。
	INVALIDPARAMETERVALUE_STORAGEREGIONS = "InvalidParameterValue.StorageRegions"

	// 参数值错误：StorageType。
	INVALIDPARAMETERVALUE_STORAGETYPE = "InvalidParameterValue.StorageType"

	// 参数值错误：StreamId无效。
	INVALIDPARAMETERVALUE_STREAMIDINVALID = "InvalidParameterValue.StreamIdInvalid"

	// 无效的流ID参数。
	INVALIDPARAMETERVALUE_STREAMIDS = "InvalidParameterValue.StreamIds"

	// 参数值错误：子应用 ID。
	INVALIDPARAMETERVALUE_SUBAPPID = "InvalidParameterValue.SubAppId"

	// 参数值错误：SubtitleFormat 参数非法。
	INVALIDPARAMETERVALUE_SUBTITLEFORMAT = "InvalidParameterValue.SubtitleFormat"

	// 参数值错误：SVG 为空。
	INVALIDPARAMETERVALUE_SVGTEMPLATE = "InvalidParameterValue.SvgTemplate"

	// 参数值错误：SVG 高度。
	INVALIDPARAMETERVALUE_SVGTEMPLATEHEIGHT = "InvalidParameterValue.SvgTemplateHeight"

	// 参数值错误：SVG 宽度。
	INVALIDPARAMETERVALUE_SVGTEMPLATEWIDTH = "InvalidParameterValue.SvgTemplateWidth"

	// 参数值错误：Switch 参数取值非法。
	INVALIDPARAMETERVALUE_SWITCH = "InvalidParameterValue.Switch"

	// 参数值错误：TEHD Type 无效。
	INVALIDPARAMETERVALUE_TEHDTYPE = "InvalidParameterValue.TEHDType"

	// 智能标签控制字段参数错误。
	INVALIDPARAMETERVALUE_TAGCONFIGURE = "InvalidParameterValue.TagConfigure"

	// 参数值错误：标签过长。
	INVALIDPARAMETERVALUE_TAGTOOLONG = "InvalidParameterValue.TagTooLong"

	// 参数值错误：Tags 无效。
	INVALIDPARAMETERVALUE_TAGS = "InvalidParameterValue.Tags"

	// 任务 ID 不存在。
	INVALIDPARAMETERVALUE_TASKID = "InvalidParameterValue.TaskId"

	// 参数值错误：搜索文本。
	INVALIDPARAMETERVALUE_TEXT = "InvalidParameterValue.Text"

	// 参数错误：文字透明度。
	INVALIDPARAMETERVALUE_TEXTALPHA = "InvalidParameterValue.TextAlpha"

	// 参数错误：文字模板。
	INVALIDPARAMETERVALUE_TEXTTEMPLATE = "InvalidParameterValue.TextTemplate"

	// 参数值错误：Thumbnail 参数取值非法。
	INVALIDPARAMETERVALUE_THUMBNAILS = "InvalidParameterValue.Thumbnails"

	// 参数值错误：TimeType。
	INVALIDPARAMETERVALUE_TIMETYPE = "InvalidParameterValue.TimeType"

	// Type 参数值错误。
	INVALIDPARAMETERVALUE_TYPE = "InvalidParameterValue.Type"

	// 无效的 Types 参数。
	INVALIDPARAMETERVALUE_TYPES = "InvalidParameterValue.Types"

	// 去重识别码一天内重复，请求被去重。
	INVALIDPARAMETERVALUE_UNIQUEIDENTIFIER = "InvalidParameterValue.UniqueIdentifier"

	// 参数值错误：文件不支持该解冻方式。
	INVALIDPARAMETERVALUE_UNSUPPORTEDRESTORETIER = "InvalidParameterValue.UnsupportedRestoreTier"

	// 参数值错误：不支持的目标存储类型，深度归档只支持部分园区。
	INVALIDPARAMETERVALUE_UNSUPPORTEDSTORAGECLASS = "InvalidParameterValue.UnsupportedStorageClass"

	// 参数值错误：不支持的存储类型转换。
	INVALIDPARAMETERVALUE_UNSUPPORTEDTRANSITION = "InvalidParameterValue.UnsupportedTransition"

	// 参数错误：无效的Url。
	INVALIDPARAMETERVALUE_URL = "InvalidParameterValue.Url"

	// 参数值错误：人脸用户自定义库过滤标签非法。
	INVALIDPARAMETERVALUE_USERDEFINELIBRARYLABELSET = "InvalidParameterValue.UserDefineLibraryLabelSet"

	// 参数错误：vcrf。
	INVALIDPARAMETERVALUE_VCRF = "InvalidParameterValue.Vcrf"

	// 参数错误：视频流码率。
	INVALIDPARAMETERVALUE_VIDEOBITRATE = "InvalidParameterValue.VideoBitrate"

	// 参数错误：视频流的编码格式。
	INVALIDPARAMETERVALUE_VIDEOCODEC = "InvalidParameterValue.VideoCodec"

	// 无效的 Vids 参数。
	INVALIDPARAMETERVALUE_VIDS = "InvalidParameterValue.Vids"

	// 参数值错误：点播会话。
	INVALIDPARAMETERVALUE_VODSESSIONKEY = "InvalidParameterValue.VodSessionKey"

	// 参数值错误：Watermarks 参数取值非法。
	INVALIDPARAMETERVALUE_WATERMARKS = "InvalidParameterValue.Watermarks"

	// 参数错误：宽度。
	INVALIDPARAMETERVALUE_WIDTH = "InvalidParameterValue.Width"

	// 水印原点距离视频图像坐标原点的水平位置。支持 %、px 两种格式。
	INVALIDPARAMETERVALUE_XPOS = "InvalidParameterValue.XPos"

	// 水印原点距离视频图像坐标原点的垂直位置。支持 %、px 两种格式。
	INVALIDPARAMETERVALUE_YPOS = "InvalidParameterValue.YPos"

	// 超过配额限制。
	LIMITEXCEEDED = "LimitExceeded"

	// 超过限制值：新旧打点信息个数和超过限制值。
	LIMITEXCEEDED_KEYFRAMEDESCCOUNTREACHMAX = "LimitExceeded.KeyFrameDescCountReachMax"

	// 超过限制值：新旧标签个数和超过限制值。
	LIMITEXCEEDED_TAGCOUNTREACHMAX = "LimitExceeded.TagCountReachMax"

	// 超过限制值：模板数超限。
	LIMITEXCEEDED_TOOMUCHTEMPLATE = "LimitExceeded.TooMuchTemplate"

	// 缺少参数错误。
	MISSINGPARAMETER = "MissingParameter"

	// 请求的次数超过了频率限制。
	REQUESTLIMITEXCEEDED = "RequestLimitExceeded"

	// 资源不存在。
	RESOURCENOTFOUND = "ResourceNotFound"

	// 资源不存在：封面不存在。
	RESOURCENOTFOUND_COVERURL = "ResourceNotFound.CoverUrl"

	// 资源不存在：文件不存在。
	RESOURCENOTFOUND_FILENOTEXIST = "ResourceNotFound.FileNotExist"

	// 资源不存在：人物。
	RESOURCENOTFOUND_PERSON = "ResourceNotFound.Person"

	// 没有开通服务。
	RESOURCENOTFOUND_SERVICENOTEXIST = "ResourceNotFound.ServiceNotExist"

	// 资源不存在：模板不存在。
	RESOURCENOTFOUND_TEMPLATENOTEXIST = "ResourceNotFound.TemplateNotExist"

	// 用户不存在。
	RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"

	// 资源不存在：关键词。
	RESOURCENOTFOUND_WORD = "ResourceNotFound.Word"

	// 资源不可用。
	RESOURCEUNAVAILABLE = "ResourceUnavailable"

	// 参数错误：不支持MasterPlaylist的M3u8。
	RESOURCEUNAVAILABLE_MASTERPLAYLIST = "ResourceUnavailable.MasterPlaylist"

	// 未授权操作。
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"

	// 未知参数错误。
	UNKNOWNPARAMETER = "UnknownParameter"

	// 操作不支持。
	UNSUPPORTEDOPERATION = "UnsupportedOperation"

	// 不支持删除非空分类。
	UNSUPPORTEDOPERATION_CLASSNOTEMPTY = "UnsupportedOperation.ClassNotEmpty"
)
