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

package v20190116

const (
	// 此产品的特有错误码

	// 没有权限。
	AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"

	// 操作失败。
	FAILEDOPERATION = "FailedOperation"

	// 操作访问密钥错误。
	FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"

	// 用户策略数超过上限。
	FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"

	// PolicyName字段指定的策略名已存在。
	FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"

	// 该策略版本已经是默认策略版本。
	FAILEDOPERATION_POLICYVERSIONALREADYDEFAULT = "FailedOperation.PolicyVersionAlreadyDefault"

	// 策略版本数已经达到上限。
	FAILEDOPERATION_POLICYVERSIONFULL = "FailedOperation.PolicyVersionFull"

	// 绑定标签失败。
	FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"

	// 解绑标签失败。
	FAILEDOPERATION_UNTAGRESOURCEFAILED = "FailedOperation.UnTagResourceFailed"

	// 用户未绑定手机。
	FAILEDOPERATION_USERNOTBINDPHONE = "FailedOperation.UserNotBindPhone"

	// 用户未绑定微信。
	FAILEDOPERATION_USERNOTBINDWECHAT = "FailedOperation.UserNotBindWechat"

	// 设置的用户没有权限。
	FAILEDOPERATION_USERUNBINDNOPERMISSION = "FailedOperation.UserUnbindNoPermission"

	// 内部错误。
	INTERNALERROR = "InternalError"

	// 内部错误。
	INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"

	// 参数错误。
	INVALIDPARAMETER = "InvalidParameter"

	// 策略文档的Action字段不合法。
	INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"

	// 策略语法中操作不存在。
	INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"

	// 策略语法中操作不存在。
	INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"

	// 策略语法中操作服务不存在。
	INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"

	// principal字段的授权对象关联策略数已达到上限。
	INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"

	// 策略语法中条件内容不正确。
	INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"

	// 策略文档的condition字段不合法。
	INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"

	// 策略语法中条件操作符不正确。
	INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"

	// 任务不存在。
	INVALIDPARAMETER_DELETIONTASKNOTEXIST = "InvalidParameter.DeletionTaskNotExist"

	// Description入参长度不能大于300字节。
	INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"

	// 策略文档的Effect字段不合法。
	INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"

	// EntityFilter字段不合法。
	INVALIDPARAMETER_ENTITYFILTERERROR = "InvalidParameter.EntityFilterError"

	// 用户组数量达到上限。
	INVALIDPARAMETER_GROUPFULL = "InvalidParameter.GroupFull"

	// GroupId字段不合法。
	INVALIDPARAMETER_GROUPIDERROR = "InvalidParameter.GroupIdError"

	// 用户组名称重复。
	INVALIDPARAMETER_GROUPNAMEINUSE = "InvalidParameter.GroupNameInUse"

	// 用户组不存在。
	INVALIDPARAMETER_GROUPNOTEXIST = "InvalidParameter.GroupNotExist"

	// 用户组中的子用户数量达到上限。
	INVALIDPARAMETER_GROUPUSERFULL = "InvalidParameter.GroupUserFull"

	// 身份提供商名称已经使用。
	INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"

	// Keyword字段不合法。
	INVALIDPARAMETER_KEYWORDERROR = "InvalidParameter.KeywordError"

	// 多因子Token验证失败。
	INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"

	// CAM不支持策略文档中所指定的资源类型。
	INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"

	// 一次操作实体数过多。
	INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"

	// 当前角色仅支持企业管理员操作，如需修改，请联系企业管理员。
	INVALIDPARAMETER_ORGANIZATIONROLEOPERATEERROR = "InvalidParameter.OrganizationRoleOperateError"

	// 非法入参。
	INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"

	// 密码不符合用户安全设置。
	INVALIDPARAMETER_PASSWORDVIOLATEDRULES = "InvalidParameter.PasswordViolatedRules"

	// PolicyDocument字段不合法。
	INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"

	// PolicyDocument字段超过长度限制。
	INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"

	// 输入参数PolicyId不合法。
	INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"

	// 策略ID不存在。
	INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"

	// PolicyName字段不合法。
	INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"

	// 策略版本不存在。
	INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"

	// 策略文档的principal字段不合法。
	INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"

	// 角色载体不允许跨站。
	INVALIDPARAMETER_PRINCIPALQCSCROSSERROR = "InvalidParameter.PrincipalQcsCrossError"

	// PrincipalQcsc错误。
	INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"

	// PrincipalQcs不存在。
	INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"

	// PrincipalService不存在。
	INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"

	// 策略语法中资源内容不正确。
	INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"

	// 策略文档的Resource字段不合法。
	INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"

	// 策略语法中资源项目不正确。
	INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"

	// 资源QCS错误。
	INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"

	// 策略语法中资源地域不正确。
	INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"

	// 策略语法中资源服务不存在。
	INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"

	// 策略语法中资源所属主账号不正确。
	INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"

	// 角色数量达到上限。
	INVALIDPARAMETER_ROLEFULL = "InvalidParameter.RoleFull"

	// 角色名不合法。
	INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"

	// 相同名称的角色已存在。
	INVALIDPARAMETER_ROLENAMEINUSE = "InvalidParameter.RoleNameInUse"

	// 角色不存在。
	INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"

	// Scope字段不合法。
	INVALIDPARAMETER_SCOPEERROR = "InvalidParameter.ScopeError"

	// 查询关键字长度错误。
	INVALIDPARAMETER_SEARCHKEYWORDLENGTHERROR = "InvalidParameter.SearchKeywordLengthError"

	// 权限边界不能用服务相关策略。
	INVALIDPARAMETER_SERVICELINKEDPOLICYCANTINPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedPolicyCantInPermissionBoundary"

	// 服务相关角色不能使用权限边界。
	INVALIDPARAMETER_SERVICELINKEDROLECANTUSEPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedRoleCantUsePermissionBoundary"

	// ServiceType字段不合法。
	INVALIDPARAMETER_SERVICETYPEERROR = "InvalidParameter.ServiceTypeError"

	// 策略文档的Statement字段不合法。
	INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"

	// 子帐号数量达到上限。
	INVALIDPARAMETER_SUBUSERFULL = "InvalidParameter.SubUserFull"

	// 子用户名称重复。
	INVALIDPARAMETER_SUBUSERNAMEINUSE = "InvalidParameter.SubUserNameInUse"

	// 超出标签限制。
	INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"

	// 标签参数错误。
	INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"

	// Uin字段不合法。
	INVALIDPARAMETER_UINERROR = "InvalidParameter.UinError"

	// 子用户加入的用户组数量达到上限。
	INVALIDPARAMETER_USERGROUPFULL = "InvalidParameter.UserGroupFull"

	// 用户名不合法。
	INVALIDPARAMETER_USERNAMEILLEGAL = "InvalidParameter.UserNameIllegal"

	// 用户对象不存在。
	INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"

	// 用户的Uin和Uid不能都为空。
	INVALIDPARAMETER_USERUINANDUINNOTALLNULL = "InvalidParameter.UserUinAndUinNotAllNull"

	// 策略文档的Version字段不合法。
	INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"

	// OIDC签名公钥错误。
	INVALIDPARAMETERVALUE_IDENTITYKEYERROR = "InvalidParameterValue.IdentityKeyError"

	// 身份提供商URL错误。
	INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"

	// 身份提供商元数据文档错误。
	INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"

	// 身份提供商名称错误。
	INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"

	// 身份提供商已达到上限。
	LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"

	// 每个账号最多支持两个AccessKey。
	OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"

	// 存在未删除的API密钥。
	OPERATIONDENIED_HAVEKEYS = "OperationDenied.HaveKeys"

	// 子用户不允许操作主账号密钥。
	OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"

	// 被操作密钥与账号不匹配。
	OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"

	// 创建子用户频率超过限制。
	REQUESTLIMITEXCEEDED_CREATEUSER = "RequestLimitExceeded.CreateUser"

	// 用户账号超出了限制。
	REQUESTLIMITEXCEEDED_UINLIMITEXCEEDED = "RequestLimitExceeded.UinLimitExceeded"

	// 用户组不存在。
	RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"

	// 身份提供商不存在。
	RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"

	// 资源不存在。
	RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"

	// PolicyId指定的资源不存在。
	RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"

	// 密钥不存在。
	RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"

	// 用户不存在。
	RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"

	// 未授权操作。
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"

	// 没有删除API密钥权限。
	UNAUTHORIZEDOPERATION_DELETEAPIKEY = "UnauthorizedOperation.DeleteApiKey"
)
