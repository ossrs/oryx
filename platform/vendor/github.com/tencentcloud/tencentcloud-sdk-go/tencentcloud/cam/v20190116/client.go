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

import (
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2019-01-16"

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


func NewAddUserRequest() (request *AddUserRequest) {
    request = &AddUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "AddUser")
    
    
    return
}

func NewAddUserResponse() (response *AddUserResponse) {
    response = &AddUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AddUser
// 创建子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_PASSWORDVIOLATEDRULES = "InvalidParameter.PasswordViolatedRules"
//  INVALIDPARAMETER_SUBUSERFULL = "InvalidParameter.SubUserFull"
//  INVALIDPARAMETER_SUBUSERNAMEINUSE = "InvalidParameter.SubUserNameInUse"
//  INVALIDPARAMETER_USERNAMEILLEGAL = "InvalidParameter.UserNameIllegal"
//  REQUESTLIMITEXCEEDED_CREATEUSER = "RequestLimitExceeded.CreateUser"
func (c *Client) AddUser(request *AddUserRequest) (response *AddUserResponse, err error) {
    return c.AddUserWithContext(context.Background(), request)
}

// AddUser
// 创建子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_PASSWORDVIOLATEDRULES = "InvalidParameter.PasswordViolatedRules"
//  INVALIDPARAMETER_SUBUSERFULL = "InvalidParameter.SubUserFull"
//  INVALIDPARAMETER_SUBUSERNAMEINUSE = "InvalidParameter.SubUserNameInUse"
//  INVALIDPARAMETER_USERNAMEILLEGAL = "InvalidParameter.UserNameIllegal"
//  REQUESTLIMITEXCEEDED_CREATEUSER = "RequestLimitExceeded.CreateUser"
func (c *Client) AddUserWithContext(ctx context.Context, request *AddUserRequest) (response *AddUserResponse, err error) {
    if request == nil {
        request = NewAddUserRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AddUser require credential")
    }

    request.SetContext(ctx)
    
    response = NewAddUserResponse()
    err = c.Send(request, response)
    return
}

func NewAddUserToGroupRequest() (request *AddUserToGroupRequest) {
    request = &AddUserToGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "AddUserToGroup")
    
    
    return
}

func NewAddUserToGroupResponse() (response *AddUserToGroupResponse) {
    response = &AddUserToGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AddUserToGroup
// 用户加入到用户组
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPNOTEXIST = "InvalidParameter.GroupNotExist"
//  INVALIDPARAMETER_GROUPUSERFULL = "InvalidParameter.GroupUserFull"
//  INVALIDPARAMETER_USERGROUPFULL = "InvalidParameter.UserGroupFull"
//  INVALIDPARAMETER_USERUINANDUINNOTALLNULL = "InvalidParameter.UserUinAndUinNotAllNull"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AddUserToGroup(request *AddUserToGroupRequest) (response *AddUserToGroupResponse, err error) {
    return c.AddUserToGroupWithContext(context.Background(), request)
}

// AddUserToGroup
// 用户加入到用户组
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPNOTEXIST = "InvalidParameter.GroupNotExist"
//  INVALIDPARAMETER_GROUPUSERFULL = "InvalidParameter.GroupUserFull"
//  INVALIDPARAMETER_USERGROUPFULL = "InvalidParameter.UserGroupFull"
//  INVALIDPARAMETER_USERUINANDUINNOTALLNULL = "InvalidParameter.UserUinAndUinNotAllNull"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AddUserToGroupWithContext(ctx context.Context, request *AddUserToGroupRequest) (response *AddUserToGroupResponse, err error) {
    if request == nil {
        request = NewAddUserToGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AddUserToGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewAddUserToGroupResponse()
    err = c.Send(request, response)
    return
}

func NewAttachGroupPolicyRequest() (request *AttachGroupPolicyRequest) {
    request = &AttachGroupPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "AttachGroupPolicy")
    
    
    return
}

func NewAttachGroupPolicyResponse() (response *AttachGroupPolicyResponse) {
    response = &AttachGroupPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AttachGroupPolicy
// 本接口（AttachGroupPolicy）可用于绑定策略到用户组。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AttachGroupPolicy(request *AttachGroupPolicyRequest) (response *AttachGroupPolicyResponse, err error) {
    return c.AttachGroupPolicyWithContext(context.Background(), request)
}

// AttachGroupPolicy
// 本接口（AttachGroupPolicy）可用于绑定策略到用户组。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AttachGroupPolicyWithContext(ctx context.Context, request *AttachGroupPolicyRequest) (response *AttachGroupPolicyResponse, err error) {
    if request == nil {
        request = NewAttachGroupPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AttachGroupPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewAttachGroupPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewAttachRolePolicyRequest() (request *AttachRolePolicyRequest) {
    request = &AttachRolePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "AttachRolePolicy")
    
    
    return
}

func NewAttachRolePolicyResponse() (response *AttachRolePolicyResponse) {
    response = &AttachRolePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AttachRolePolicy
// 本接口（AttachRolePolicy）用于绑定策略到角色。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) AttachRolePolicy(request *AttachRolePolicyRequest) (response *AttachRolePolicyResponse, err error) {
    return c.AttachRolePolicyWithContext(context.Background(), request)
}

// AttachRolePolicy
// 本接口（AttachRolePolicy）用于绑定策略到角色。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) AttachRolePolicyWithContext(ctx context.Context, request *AttachRolePolicyRequest) (response *AttachRolePolicyResponse, err error) {
    if request == nil {
        request = NewAttachRolePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AttachRolePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewAttachRolePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewAttachUserPolicyRequest() (request *AttachUserPolicyRequest) {
    request = &AttachUserPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "AttachUserPolicy")
    
    
    return
}

func NewAttachUserPolicyResponse() (response *AttachUserPolicyResponse) {
    response = &AttachUserPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// AttachUserPolicy
// 本接口（AttachUserPolicy）可用于绑定到用户的策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AttachUserPolicy(request *AttachUserPolicyRequest) (response *AttachUserPolicyResponse, err error) {
    return c.AttachUserPolicyWithContext(context.Background(), request)
}

// AttachUserPolicy
// 本接口（AttachUserPolicy）可用于绑定到用户的策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) AttachUserPolicyWithContext(ctx context.Context, request *AttachUserPolicyRequest) (response *AttachUserPolicyResponse, err error) {
    if request == nil {
        request = NewAttachUserPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AttachUserPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewAttachUserPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewConsumeCustomMFATokenRequest() (request *ConsumeCustomMFATokenRequest) {
    request = &ConsumeCustomMFATokenRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ConsumeCustomMFAToken")
    
    
    return
}

func NewConsumeCustomMFATokenResponse() (response *ConsumeCustomMFATokenResponse) {
    response = &ConsumeCustomMFATokenResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ConsumeCustomMFAToken
// 验证自定义多因子Token
//
// 可能返回的错误码:
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
func (c *Client) ConsumeCustomMFAToken(request *ConsumeCustomMFATokenRequest) (response *ConsumeCustomMFATokenResponse, err error) {
    return c.ConsumeCustomMFATokenWithContext(context.Background(), request)
}

// ConsumeCustomMFAToken
// 验证自定义多因子Token
//
// 可能返回的错误码:
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
func (c *Client) ConsumeCustomMFATokenWithContext(ctx context.Context, request *ConsumeCustomMFATokenRequest) (response *ConsumeCustomMFATokenResponse, err error) {
    if request == nil {
        request = NewConsumeCustomMFATokenRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ConsumeCustomMFAToken require credential")
    }

    request.SetContext(ctx)
    
    response = NewConsumeCustomMFATokenResponse()
    err = c.Send(request, response)
    return
}

func NewCreateAccessKeyRequest() (request *CreateAccessKeyRequest) {
    request = &CreateAccessKeyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateAccessKey")
    
    
    return
}

func NewCreateAccessKeyResponse() (response *CreateAccessKeyResponse) {
    response = &CreateAccessKeyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateAccessKey
// 为CAM用户创建访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
func (c *Client) CreateAccessKey(request *CreateAccessKeyRequest) (response *CreateAccessKeyResponse, err error) {
    return c.CreateAccessKeyWithContext(context.Background(), request)
}

// CreateAccessKey
// 为CAM用户创建访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
func (c *Client) CreateAccessKeyWithContext(ctx context.Context, request *CreateAccessKeyRequest) (response *CreateAccessKeyResponse, err error) {
    if request == nil {
        request = NewCreateAccessKeyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateAccessKey require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateAccessKeyResponse()
    err = c.Send(request, response)
    return
}

func NewCreateGroupRequest() (request *CreateGroupRequest) {
    request = &CreateGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateGroup")
    
    
    return
}

func NewCreateGroupResponse() (response *CreateGroupResponse) {
    response = &CreateGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateGroup
// 创建用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPFULL = "InvalidParameter.GroupFull"
//  INVALIDPARAMETER_GROUPNAMEINUSE = "InvalidParameter.GroupNameInUse"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) CreateGroup(request *CreateGroupRequest) (response *CreateGroupResponse, err error) {
    return c.CreateGroupWithContext(context.Background(), request)
}

// CreateGroup
// 创建用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPFULL = "InvalidParameter.GroupFull"
//  INVALIDPARAMETER_GROUPNAMEINUSE = "InvalidParameter.GroupNameInUse"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) CreateGroupWithContext(ctx context.Context, request *CreateGroupRequest) (response *CreateGroupResponse, err error) {
    if request == nil {
        request = NewCreateGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateGroupResponse()
    err = c.Send(request, response)
    return
}

func NewCreateOIDCConfigRequest() (request *CreateOIDCConfigRequest) {
    request = &CreateOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateOIDCConfig")
    
    
    return
}

func NewCreateOIDCConfigResponse() (response *CreateOIDCConfigResponse) {
    response = &CreateOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateOIDCConfig
// 创建角色OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateOIDCConfig(request *CreateOIDCConfigRequest) (response *CreateOIDCConfigResponse, err error) {
    return c.CreateOIDCConfigWithContext(context.Background(), request)
}

// CreateOIDCConfig
// 创建角色OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateOIDCConfigWithContext(ctx context.Context, request *CreateOIDCConfigRequest) (response *CreateOIDCConfigResponse, err error) {
    if request == nil {
        request = NewCreateOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewCreatePolicyRequest() (request *CreatePolicyRequest) {
    request = &CreatePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreatePolicy")
    
    
    return
}

func NewCreatePolicyResponse() (response *CreatePolicyResponse) {
    response = &CreatePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreatePolicy
// 本接口（CreatePolicy）可用于创建策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) CreatePolicy(request *CreatePolicyRequest) (response *CreatePolicyResponse, err error) {
    return c.CreatePolicyWithContext(context.Background(), request)
}

// CreatePolicy
// 本接口（CreatePolicy）可用于创建策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) CreatePolicyWithContext(ctx context.Context, request *CreatePolicyRequest) (response *CreatePolicyResponse, err error) {
    if request == nil {
        request = NewCreatePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreatePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreatePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewCreatePolicyVersionRequest() (request *CreatePolicyVersionRequest) {
    request = &CreatePolicyVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreatePolicyVersion")
    
    
    return
}

func NewCreatePolicyVersionResponse() (response *CreatePolicyVersionResponse) {
    response = &CreatePolicyVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreatePolicyVersion
// 该接口（CreatePolicyVersion）用于新增策略版本，用户创建了一个策略版本之后可以方便的通过变更策略版本的方式来变更策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONFULL = "FailedOperation.PolicyVersionFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) CreatePolicyVersion(request *CreatePolicyVersionRequest) (response *CreatePolicyVersionResponse, err error) {
    return c.CreatePolicyVersionWithContext(context.Background(), request)
}

// CreatePolicyVersion
// 该接口（CreatePolicyVersion）用于新增策略版本，用户创建了一个策略版本之后可以方便的通过变更策略版本的方式来变更策略。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_POLICYFULL = "FailedOperation.PolicyFull"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONFULL = "FailedOperation.PolicyVersionFull"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) CreatePolicyVersionWithContext(ctx context.Context, request *CreatePolicyVersionRequest) (response *CreatePolicyVersionResponse, err error) {
    if request == nil {
        request = NewCreatePolicyVersionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreatePolicyVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreatePolicyVersionResponse()
    err = c.Send(request, response)
    return
}

func NewCreateRoleRequest() (request *CreateRoleRequest) {
    request = &CreateRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateRole")
    
    
    return
}

func NewCreateRoleResponse() (response *CreateRoleResponse) {
    response = &CreateRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateRole
// 本接口（CreateRole）用于创建角色。
//
// 可能返回的错误码:
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSCROSSERROR = "InvalidParameter.PrincipalQcsCrossError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLEFULL = "InvalidParameter.RoleFull"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENAMEINUSE = "InvalidParameter.RoleNameInUse"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) CreateRole(request *CreateRoleRequest) (response *CreateRoleResponse, err error) {
    return c.CreateRoleWithContext(context.Background(), request)
}

// CreateRole
// 本接口（CreateRole）用于创建角色。
//
// 可能返回的错误码:
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSCROSSERROR = "InvalidParameter.PrincipalQcsCrossError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLEFULL = "InvalidParameter.RoleFull"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENAMEINUSE = "InvalidParameter.RoleNameInUse"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) CreateRoleWithContext(ctx context.Context, request *CreateRoleRequest) (response *CreateRoleResponse, err error) {
    if request == nil {
        request = NewCreateRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateRoleResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSAMLProviderRequest() (request *CreateSAMLProviderRequest) {
    request = &CreateSAMLProviderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateSAMLProvider")
    
    
    return
}

func NewCreateSAMLProviderResponse() (response *CreateSAMLProviderResponse) {
    response = &CreateSAMLProviderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateSAMLProvider
// 创建SAML身份提供商
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateSAMLProvider(request *CreateSAMLProviderRequest) (response *CreateSAMLProviderResponse, err error) {
    return c.CreateSAMLProviderWithContext(context.Background(), request)
}

// CreateSAMLProvider
// 创建SAML身份提供商
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateSAMLProviderWithContext(ctx context.Context, request *CreateSAMLProviderRequest) (response *CreateSAMLProviderResponse, err error) {
    if request == nil {
        request = NewCreateSAMLProviderRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateSAMLProvider require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateSAMLProviderResponse()
    err = c.Send(request, response)
    return
}

func NewCreateServiceLinkedRoleRequest() (request *CreateServiceLinkedRoleRequest) {
    request = &CreateServiceLinkedRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateServiceLinkedRole")
    
    
    return
}

func NewCreateServiceLinkedRoleResponse() (response *CreateServiceLinkedRoleResponse) {
    response = &CreateServiceLinkedRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateServiceLinkedRole
// 创建服务相关角色
//
// 可能返回的错误码:
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENAMEINUSE = "InvalidParameter.RoleNameInUse"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) CreateServiceLinkedRole(request *CreateServiceLinkedRoleRequest) (response *CreateServiceLinkedRoleResponse, err error) {
    return c.CreateServiceLinkedRoleWithContext(context.Background(), request)
}

// CreateServiceLinkedRole
// 创建服务相关角色
//
// 可能返回的错误码:
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENAMEINUSE = "InvalidParameter.RoleNameInUse"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) CreateServiceLinkedRoleWithContext(ctx context.Context, request *CreateServiceLinkedRoleRequest) (response *CreateServiceLinkedRoleResponse, err error) {
    if request == nil {
        request = NewCreateServiceLinkedRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateServiceLinkedRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateServiceLinkedRoleResponse()
    err = c.Send(request, response)
    return
}

func NewCreateUserOIDCConfigRequest() (request *CreateUserOIDCConfigRequest) {
    request = &CreateUserOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateUserOIDCConfig")
    
    
    return
}

func NewCreateUserOIDCConfigResponse() (response *CreateUserOIDCConfigResponse) {
    response = &CreateUserOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateUserOIDCConfig
// 创建用户OIDC配置。只能创建一个用户OIDC身份提供商，并且创建用户OIDC配置之后会自动关闭用户SAML SSO身份提供商。
//
// 可能返回的错误码:
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYKEYERROR = "InvalidParameterValue.IdentityKeyError"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateUserOIDCConfig(request *CreateUserOIDCConfigRequest) (response *CreateUserOIDCConfigResponse, err error) {
    return c.CreateUserOIDCConfigWithContext(context.Background(), request)
}

// CreateUserOIDCConfig
// 创建用户OIDC配置。只能创建一个用户OIDC身份提供商，并且创建用户OIDC配置之后会自动关闭用户SAML SSO身份提供商。
//
// 可能返回的错误码:
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYKEYERROR = "InvalidParameterValue.IdentityKeyError"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
func (c *Client) CreateUserOIDCConfigWithContext(ctx context.Context, request *CreateUserOIDCConfigRequest) (response *CreateUserOIDCConfigResponse, err error) {
    if request == nil {
        request = NewCreateUserOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateUserOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateUserOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewCreateUserSAMLConfigRequest() (request *CreateUserSAMLConfigRequest) {
    request = &CreateUserSAMLConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "CreateUserSAMLConfig")
    
    
    return
}

func NewCreateUserSAMLConfigResponse() (response *CreateUserSAMLConfigResponse) {
    response = &CreateUserSAMLConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateUserSAMLConfig
// 创建用户SAML配置
//
// 可能返回的错误码:
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
func (c *Client) CreateUserSAMLConfig(request *CreateUserSAMLConfigRequest) (response *CreateUserSAMLConfigResponse, err error) {
    return c.CreateUserSAMLConfigWithContext(context.Background(), request)
}

// CreateUserSAMLConfig
// 创建用户SAML配置
//
// 可能返回的错误码:
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
func (c *Client) CreateUserSAMLConfigWithContext(ctx context.Context, request *CreateUserSAMLConfigRequest) (response *CreateUserSAMLConfigResponse, err error) {
    if request == nil {
        request = NewCreateUserSAMLConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateUserSAMLConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateUserSAMLConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAccessKeyRequest() (request *DeleteAccessKeyRequest) {
    request = &DeleteAccessKeyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteAccessKey")
    
    
    return
}

func NewDeleteAccessKeyResponse() (response *DeleteAccessKeyResponse) {
    response = &DeleteAccessKeyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteAccessKey
// 为CAM用户删除访问密钥。
//
// 此接口属于高风险操作，删除密钥后不可恢复，腾讯云将永久拒绝此密钥的所有请求，请谨慎使用。
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) DeleteAccessKey(request *DeleteAccessKeyRequest) (response *DeleteAccessKeyResponse, err error) {
    return c.DeleteAccessKeyWithContext(context.Background(), request)
}

// DeleteAccessKey
// 为CAM用户删除访问密钥。
//
// 此接口属于高风险操作，删除密钥后不可恢复，腾讯云将永久拒绝此密钥的所有请求，请谨慎使用。
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) DeleteAccessKeyWithContext(ctx context.Context, request *DeleteAccessKeyRequest) (response *DeleteAccessKeyResponse, err error) {
    if request == nil {
        request = NewDeleteAccessKeyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteAccessKey require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteAccessKeyResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteGroupRequest() (request *DeleteGroupRequest) {
    request = &DeleteGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteGroup")
    
    
    return
}

func NewDeleteGroupResponse() (response *DeleteGroupResponse) {
    response = &DeleteGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteGroup
// 删除用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DeleteGroup(request *DeleteGroupRequest) (response *DeleteGroupResponse, err error) {
    return c.DeleteGroupWithContext(context.Background(), request)
}

// DeleteGroup
// 删除用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DeleteGroupWithContext(ctx context.Context, request *DeleteGroupRequest) (response *DeleteGroupResponse, err error) {
    if request == nil {
        request = NewDeleteGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteOIDCConfigRequest() (request *DeleteOIDCConfigRequest) {
    request = &DeleteOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteOIDCConfig")
    
    
    return
}

func NewDeleteOIDCConfigResponse() (response *DeleteOIDCConfigResponse) {
    response = &DeleteOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteOIDCConfig
// 删除OIDC身份提供商
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) DeleteOIDCConfig(request *DeleteOIDCConfigRequest) (response *DeleteOIDCConfigResponse, err error) {
    return c.DeleteOIDCConfigWithContext(context.Background(), request)
}

// DeleteOIDCConfig
// 删除OIDC身份提供商
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) DeleteOIDCConfigWithContext(ctx context.Context, request *DeleteOIDCConfigRequest) (response *DeleteOIDCConfigResponse, err error) {
    if request == nil {
        request = NewDeleteOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDeletePolicyRequest() (request *DeletePolicyRequest) {
    request = &DeletePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeletePolicy")
    
    
    return
}

func NewDeletePolicyResponse() (response *DeletePolicyResponse) {
    response = &DeletePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeletePolicy
// 本接口（DeletePolicy）可用于删除策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) DeletePolicy(request *DeletePolicyRequest) (response *DeletePolicyResponse, err error) {
    return c.DeletePolicyWithContext(context.Background(), request)
}

// DeletePolicy
// 本接口（DeletePolicy）可用于删除策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) DeletePolicyWithContext(ctx context.Context, request *DeletePolicyRequest) (response *DeletePolicyResponse, err error) {
    if request == nil {
        request = NewDeletePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeletePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeletePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewDeletePolicyVersionRequest() (request *DeletePolicyVersionRequest) {
    request = &DeletePolicyVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeletePolicyVersion")
    
    
    return
}

func NewDeletePolicyVersionResponse() (response *DeletePolicyVersionResponse) {
    response = &DeletePolicyVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeletePolicyVersion
// 本接口（DeletePolicyVersion）可用于删除一个策略的策略版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONALREADYDEFAULT = "FailedOperation.PolicyVersionAlreadyDefault"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DeletePolicyVersion(request *DeletePolicyVersionRequest) (response *DeletePolicyVersionResponse, err error) {
    return c.DeletePolicyVersionWithContext(context.Background(), request)
}

// DeletePolicyVersion
// 本接口（DeletePolicyVersion）可用于删除一个策略的策略版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONALREADYDEFAULT = "FailedOperation.PolicyVersionAlreadyDefault"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DeletePolicyVersionWithContext(ctx context.Context, request *DeletePolicyVersionRequest) (response *DeletePolicyVersionResponse, err error) {
    if request == nil {
        request = NewDeletePolicyVersionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeletePolicyVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeletePolicyVersionResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRoleRequest() (request *DeleteRoleRequest) {
    request = &DeleteRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteRole")
    
    
    return
}

func NewDeleteRoleResponse() (response *DeleteRoleResponse) {
    response = &DeleteRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteRole
// 本接口（DeleteRole）用于删除指定角色。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ORGANIZATIONROLEOPERATEERROR = "InvalidParameter.OrganizationRoleOperateError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) DeleteRole(request *DeleteRoleRequest) (response *DeleteRoleResponse, err error) {
    return c.DeleteRoleWithContext(context.Background(), request)
}

// DeleteRole
// 本接口（DeleteRole）用于删除指定角色。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ORGANIZATIONROLEOPERATEERROR = "InvalidParameter.OrganizationRoleOperateError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) DeleteRoleWithContext(ctx context.Context, request *DeleteRoleRequest) (response *DeleteRoleResponse, err error) {
    if request == nil {
        request = NewDeleteRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteRoleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRolePermissionsBoundaryRequest() (request *DeleteRolePermissionsBoundaryRequest) {
    request = &DeleteRolePermissionsBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteRolePermissionsBoundary")
    
    
    return
}

func NewDeleteRolePermissionsBoundaryResponse() (response *DeleteRolePermissionsBoundaryResponse) {
    response = &DeleteRolePermissionsBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteRolePermissionsBoundary
// 删除角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDROLECANTUSEPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedRoleCantUsePermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteRolePermissionsBoundary(request *DeleteRolePermissionsBoundaryRequest) (response *DeleteRolePermissionsBoundaryResponse, err error) {
    return c.DeleteRolePermissionsBoundaryWithContext(context.Background(), request)
}

// DeleteRolePermissionsBoundary
// 删除角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDROLECANTUSEPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedRoleCantUsePermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteRolePermissionsBoundaryWithContext(ctx context.Context, request *DeleteRolePermissionsBoundaryRequest) (response *DeleteRolePermissionsBoundaryResponse, err error) {
    if request == nil {
        request = NewDeleteRolePermissionsBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteRolePermissionsBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteRolePermissionsBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSAMLProviderRequest() (request *DeleteSAMLProviderRequest) {
    request = &DeleteSAMLProviderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteSAMLProvider")
    
    
    return
}

func NewDeleteSAMLProviderResponse() (response *DeleteSAMLProviderResponse) {
    response = &DeleteSAMLProviderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteSAMLProvider
// 删除SAML身份提供商
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) DeleteSAMLProvider(request *DeleteSAMLProviderRequest) (response *DeleteSAMLProviderResponse, err error) {
    return c.DeleteSAMLProviderWithContext(context.Background(), request)
}

// DeleteSAMLProvider
// 删除SAML身份提供商
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) DeleteSAMLProviderWithContext(ctx context.Context, request *DeleteSAMLProviderRequest) (response *DeleteSAMLProviderResponse, err error) {
    if request == nil {
        request = NewDeleteSAMLProviderRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteSAMLProvider require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteSAMLProviderResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteServiceLinkedRoleRequest() (request *DeleteServiceLinkedRoleRequest) {
    request = &DeleteServiceLinkedRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteServiceLinkedRole")
    
    
    return
}

func NewDeleteServiceLinkedRoleResponse() (response *DeleteServiceLinkedRoleResponse) {
    response = &DeleteServiceLinkedRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteServiceLinkedRole
// 删除服务相关角色
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) DeleteServiceLinkedRole(request *DeleteServiceLinkedRoleRequest) (response *DeleteServiceLinkedRoleResponse, err error) {
    return c.DeleteServiceLinkedRoleWithContext(context.Background(), request)
}

// DeleteServiceLinkedRole
// 删除服务相关角色
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) DeleteServiceLinkedRoleWithContext(ctx context.Context, request *DeleteServiceLinkedRoleRequest) (response *DeleteServiceLinkedRoleResponse, err error) {
    if request == nil {
        request = NewDeleteServiceLinkedRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteServiceLinkedRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteServiceLinkedRoleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteUserRequest() (request *DeleteUserRequest) {
    request = &DeleteUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteUser")
    
    
    return
}

func NewDeleteUserResponse() (response *DeleteUserResponse) {
    response = &DeleteUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteUser
// 删除子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  OPERATIONDENIED_HAVEKEYS = "OperationDenied.HaveKeys"
//  REQUESTLIMITEXCEEDED_UINLIMITEXCEEDED = "RequestLimitExceeded.UinLimitExceeded"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION_DELETEAPIKEY = "UnauthorizedOperation.DeleteApiKey"
func (c *Client) DeleteUser(request *DeleteUserRequest) (response *DeleteUserResponse, err error) {
    return c.DeleteUserWithContext(context.Background(), request)
}

// DeleteUser
// 删除子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  OPERATIONDENIED_HAVEKEYS = "OperationDenied.HaveKeys"
//  REQUESTLIMITEXCEEDED_UINLIMITEXCEEDED = "RequestLimitExceeded.UinLimitExceeded"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION_DELETEAPIKEY = "UnauthorizedOperation.DeleteApiKey"
func (c *Client) DeleteUserWithContext(ctx context.Context, request *DeleteUserRequest) (response *DeleteUserResponse, err error) {
    if request == nil {
        request = NewDeleteUserRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteUser require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteUserResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteUserPermissionsBoundaryRequest() (request *DeleteUserPermissionsBoundaryRequest) {
    request = &DeleteUserPermissionsBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DeleteUserPermissionsBoundary")
    
    
    return
}

func NewDeleteUserPermissionsBoundaryResponse() (response *DeleteUserPermissionsBoundaryResponse) {
    response = &DeleteUserPermissionsBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteUserPermissionsBoundary
// 删除用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteUserPermissionsBoundary(request *DeleteUserPermissionsBoundaryRequest) (response *DeleteUserPermissionsBoundaryResponse, err error) {
    return c.DeleteUserPermissionsBoundaryWithContext(context.Background(), request)
}

// DeleteUserPermissionsBoundary
// 删除用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteUserPermissionsBoundaryWithContext(ctx context.Context, request *DeleteUserPermissionsBoundaryRequest) (response *DeleteUserPermissionsBoundaryResponse, err error) {
    if request == nil {
        request = NewDeleteUserPermissionsBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteUserPermissionsBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteUserPermissionsBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeOIDCConfigRequest() (request *DescribeOIDCConfigRequest) {
    request = &DescribeOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeOIDCConfig")
    
    
    return
}

func NewDescribeOIDCConfigResponse() (response *DescribeOIDCConfigResponse) {
    response = &DescribeOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeOIDCConfig
// 查询角色OIDC配置
//
// 可能返回的错误码:
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeOIDCConfig(request *DescribeOIDCConfigRequest) (response *DescribeOIDCConfigResponse, err error) {
    return c.DescribeOIDCConfigWithContext(context.Background(), request)
}

// DescribeOIDCConfig
// 查询角色OIDC配置
//
// 可能返回的错误码:
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeOIDCConfigWithContext(ctx context.Context, request *DescribeOIDCConfigRequest) (response *DescribeOIDCConfigResponse, err error) {
    if request == nil {
        request = NewDescribeOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRoleListRequest() (request *DescribeRoleListRequest) {
    request = &DescribeRoleListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeRoleList")
    
    
    return
}

func NewDescribeRoleListResponse() (response *DescribeRoleListResponse) {
    response = &DescribeRoleListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeRoleList
// 本接口（DescribeRoleList）用于获取账号下的角色列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) DescribeRoleList(request *DescribeRoleListRequest) (response *DescribeRoleListResponse, err error) {
    return c.DescribeRoleListWithContext(context.Background(), request)
}

// DescribeRoleList
// 本接口（DescribeRoleList）用于获取账号下的角色列表。
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) DescribeRoleListWithContext(ctx context.Context, request *DescribeRoleListRequest) (response *DescribeRoleListResponse, err error) {
    if request == nil {
        request = NewDescribeRoleListRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRoleList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRoleListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSafeAuthFlagRequest() (request *DescribeSafeAuthFlagRequest) {
    request = &DescribeSafeAuthFlagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeSafeAuthFlag")
    
    
    return
}

func NewDescribeSafeAuthFlagResponse() (response *DescribeSafeAuthFlagResponse) {
    response = &DescribeSafeAuthFlagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSafeAuthFlag
// 查询用户安全设置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlag(request *DescribeSafeAuthFlagRequest) (response *DescribeSafeAuthFlagResponse, err error) {
    return c.DescribeSafeAuthFlagWithContext(context.Background(), request)
}

// DescribeSafeAuthFlag
// 查询用户安全设置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlagWithContext(ctx context.Context, request *DescribeSafeAuthFlagRequest) (response *DescribeSafeAuthFlagResponse, err error) {
    if request == nil {
        request = NewDescribeSafeAuthFlagRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSafeAuthFlag require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSafeAuthFlagResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSafeAuthFlagCollRequest() (request *DescribeSafeAuthFlagCollRequest) {
    request = &DescribeSafeAuthFlagCollRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeSafeAuthFlagColl")
    
    
    return
}

func NewDescribeSafeAuthFlagCollResponse() (response *DescribeSafeAuthFlagCollResponse) {
    response = &DescribeSafeAuthFlagCollResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSafeAuthFlagColl
// 获取子账号安全设置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlagColl(request *DescribeSafeAuthFlagCollRequest) (response *DescribeSafeAuthFlagCollResponse, err error) {
    return c.DescribeSafeAuthFlagCollWithContext(context.Background(), request)
}

// DescribeSafeAuthFlagColl
// 获取子账号安全设置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlagCollWithContext(ctx context.Context, request *DescribeSafeAuthFlagCollRequest) (response *DescribeSafeAuthFlagCollResponse, err error) {
    if request == nil {
        request = NewDescribeSafeAuthFlagCollRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSafeAuthFlagColl require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSafeAuthFlagCollResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSafeAuthFlagIntlRequest() (request *DescribeSafeAuthFlagIntlRequest) {
    request = &DescribeSafeAuthFlagIntlRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeSafeAuthFlagIntl")
    
    
    return
}

func NewDescribeSafeAuthFlagIntlResponse() (response *DescribeSafeAuthFlagIntlResponse) {
    response = &DescribeSafeAuthFlagIntlResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSafeAuthFlagIntl
// 查询安全设置(国际站)
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlagIntl(request *DescribeSafeAuthFlagIntlRequest) (response *DescribeSafeAuthFlagIntlResponse, err error) {
    return c.DescribeSafeAuthFlagIntlWithContext(context.Background(), request)
}

// DescribeSafeAuthFlagIntl
// 查询安全设置(国际站)
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DescribeSafeAuthFlagIntlWithContext(ctx context.Context, request *DescribeSafeAuthFlagIntlRequest) (response *DescribeSafeAuthFlagIntlResponse, err error) {
    if request == nil {
        request = NewDescribeSafeAuthFlagIntlRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSafeAuthFlagIntl require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSafeAuthFlagIntlResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSubAccountsRequest() (request *DescribeSubAccountsRequest) {
    request = &DescribeSubAccountsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeSubAccounts")
    
    
    return
}

func NewDescribeSubAccountsResponse() (response *DescribeSubAccountsResponse) {
    response = &DescribeSubAccountsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeSubAccounts
// 通过子用户UIN列表查询子用户
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DescribeSubAccounts(request *DescribeSubAccountsRequest) (response *DescribeSubAccountsResponse, err error) {
    return c.DescribeSubAccountsWithContext(context.Background(), request)
}

// DescribeSubAccounts
// 通过子用户UIN列表查询子用户
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) DescribeSubAccountsWithContext(ctx context.Context, request *DescribeSubAccountsRequest) (response *DescribeSubAccountsResponse, err error) {
    if request == nil {
        request = NewDescribeSubAccountsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSubAccounts require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSubAccountsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeUserOIDCConfigRequest() (request *DescribeUserOIDCConfigRequest) {
    request = &DescribeUserOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeUserOIDCConfig")
    
    
    return
}

func NewDescribeUserOIDCConfigResponse() (response *DescribeUserOIDCConfigResponse) {
    response = &DescribeUserOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeUserOIDCConfig
// 查询用户OIDC配置
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeUserOIDCConfig(request *DescribeUserOIDCConfigRequest) (response *DescribeUserOIDCConfigResponse, err error) {
    return c.DescribeUserOIDCConfigWithContext(context.Background(), request)
}

// DescribeUserOIDCConfig
// 查询用户OIDC配置
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeUserOIDCConfigWithContext(ctx context.Context, request *DescribeUserOIDCConfigRequest) (response *DescribeUserOIDCConfigResponse, err error) {
    if request == nil {
        request = NewDescribeUserOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeUserOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeUserOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeUserSAMLConfigRequest() (request *DescribeUserSAMLConfigRequest) {
    request = &DescribeUserSAMLConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DescribeUserSAMLConfig")
    
    
    return
}

func NewDescribeUserSAMLConfigResponse() (response *DescribeUserSAMLConfigResponse) {
    response = &DescribeUserSAMLConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeUserSAMLConfig
// 查询用户SAML配置
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeUserSAMLConfig(request *DescribeUserSAMLConfigRequest) (response *DescribeUserSAMLConfigResponse, err error) {
    return c.DescribeUserSAMLConfigWithContext(context.Background(), request)
}

// DescribeUserSAMLConfig
// 查询用户SAML配置
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) DescribeUserSAMLConfigWithContext(ctx context.Context, request *DescribeUserSAMLConfigRequest) (response *DescribeUserSAMLConfigResponse, err error) {
    if request == nil {
        request = NewDescribeUserSAMLConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeUserSAMLConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeUserSAMLConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDetachGroupPolicyRequest() (request *DetachGroupPolicyRequest) {
    request = &DetachGroupPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DetachGroupPolicy")
    
    
    return
}

func NewDetachGroupPolicyResponse() (response *DetachGroupPolicyResponse) {
    response = &DetachGroupPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DetachGroupPolicy
// 本接口（DetachGroupPolicy）可用于解除绑定到用户组的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DetachGroupPolicy(request *DetachGroupPolicyRequest) (response *DetachGroupPolicyResponse, err error) {
    return c.DetachGroupPolicyWithContext(context.Background(), request)
}

// DetachGroupPolicy
// 本接口（DetachGroupPolicy）可用于解除绑定到用户组的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DetachGroupPolicyWithContext(ctx context.Context, request *DetachGroupPolicyRequest) (response *DetachGroupPolicyResponse, err error) {
    if request == nil {
        request = NewDetachGroupPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DetachGroupPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewDetachGroupPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewDetachRolePolicyRequest() (request *DetachRolePolicyRequest) {
    request = &DetachRolePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DetachRolePolicy")
    
    
    return
}

func NewDetachRolePolicyResponse() (response *DetachRolePolicyResponse) {
    response = &DetachRolePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DetachRolePolicy
// 本接口（DetachRolePolicy）用于解除绑定角色的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) DetachRolePolicy(request *DetachRolePolicyRequest) (response *DetachRolePolicyResponse, err error) {
    return c.DetachRolePolicyWithContext(context.Background(), request)
}

// DetachRolePolicy
// 本接口（DetachRolePolicy）用于解除绑定角色的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) DetachRolePolicyWithContext(ctx context.Context, request *DetachRolePolicyRequest) (response *DetachRolePolicyResponse, err error) {
    if request == nil {
        request = NewDetachRolePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DetachRolePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewDetachRolePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewDetachUserPolicyRequest() (request *DetachUserPolicyRequest) {
    request = &DetachUserPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DetachUserPolicy")
    
    
    return
}

func NewDetachUserPolicyResponse() (response *DetachUserPolicyResponse) {
    response = &DetachUserPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DetachUserPolicy
// 本接口（DetachUserPolicy）可用于解除绑定到用户的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DetachUserPolicy(request *DetachUserPolicyRequest) (response *DetachUserPolicyResponse, err error) {
    return c.DetachUserPolicyWithContext(context.Background(), request)
}

// DetachUserPolicy
// 本接口（DetachUserPolicy）可用于解除绑定到用户的策略。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) DetachUserPolicyWithContext(ctx context.Context, request *DetachUserPolicyRequest) (response *DetachUserPolicyResponse, err error) {
    if request == nil {
        request = NewDetachUserPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DetachUserPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewDetachUserPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewDisableUserSSORequest() (request *DisableUserSSORequest) {
    request = &DisableUserSSORequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "DisableUserSSO")
    
    
    return
}

func NewDisableUserSSOResponse() (response *DisableUserSSOResponse) {
    response = &DisableUserSSOResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DisableUserSSO
// 禁用用户SSO
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
func (c *Client) DisableUserSSO(request *DisableUserSSORequest) (response *DisableUserSSOResponse, err error) {
    return c.DisableUserSSOWithContext(context.Background(), request)
}

// DisableUserSSO
// 禁用用户SSO
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
func (c *Client) DisableUserSSOWithContext(ctx context.Context, request *DisableUserSSORequest) (response *DisableUserSSOResponse, err error) {
    if request == nil {
        request = NewDisableUserSSORequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DisableUserSSO require credential")
    }

    request.SetContext(ctx)
    
    response = NewDisableUserSSOResponse()
    err = c.Send(request, response)
    return
}

func NewGetAccountSummaryRequest() (request *GetAccountSummaryRequest) {
    request = &GetAccountSummaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetAccountSummary")
    
    
    return
}

func NewGetAccountSummaryResponse() (response *GetAccountSummaryResponse) {
    response = &GetAccountSummaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetAccountSummary
// 查询账户摘要
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
func (c *Client) GetAccountSummary(request *GetAccountSummaryRequest) (response *GetAccountSummaryResponse, err error) {
    return c.GetAccountSummaryWithContext(context.Background(), request)
}

// GetAccountSummary
// 查询账户摘要
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
func (c *Client) GetAccountSummaryWithContext(ctx context.Context, request *GetAccountSummaryRequest) (response *GetAccountSummaryResponse, err error) {
    if request == nil {
        request = NewGetAccountSummaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetAccountSummary require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetAccountSummaryResponse()
    err = c.Send(request, response)
    return
}

func NewGetCustomMFATokenInfoRequest() (request *GetCustomMFATokenInfoRequest) {
    request = &GetCustomMFATokenInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetCustomMFATokenInfo")
    
    
    return
}

func NewGetCustomMFATokenInfoResponse() (response *GetCustomMFATokenInfoResponse) {
    response = &GetCustomMFATokenInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetCustomMFATokenInfo
// 获取自定义多因子Token关联信息
//
// 可能返回的错误码:
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
func (c *Client) GetCustomMFATokenInfo(request *GetCustomMFATokenInfoRequest) (response *GetCustomMFATokenInfoResponse, err error) {
    return c.GetCustomMFATokenInfoWithContext(context.Background(), request)
}

// GetCustomMFATokenInfo
// 获取自定义多因子Token关联信息
//
// 可能返回的错误码:
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
func (c *Client) GetCustomMFATokenInfoWithContext(ctx context.Context, request *GetCustomMFATokenInfoRequest) (response *GetCustomMFATokenInfoResponse, err error) {
    if request == nil {
        request = NewGetCustomMFATokenInfoRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetCustomMFATokenInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetCustomMFATokenInfoResponse()
    err = c.Send(request, response)
    return
}

func NewGetGroupRequest() (request *GetGroupRequest) {
    request = &GetGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetGroup")
    
    
    return
}

func NewGetGroupResponse() (response *GetGroupResponse) {
    response = &GetGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetGroup
// 查询用户组详情
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) GetGroup(request *GetGroupRequest) (response *GetGroupResponse, err error) {
    return c.GetGroupWithContext(context.Background(), request)
}

// GetGroup
// 查询用户组详情
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) GetGroupWithContext(ctx context.Context, request *GetGroupRequest) (response *GetGroupResponse, err error) {
    if request == nil {
        request = NewGetGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetGroupResponse()
    err = c.Send(request, response)
    return
}

func NewGetPolicyRequest() (request *GetPolicyRequest) {
    request = &GetPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetPolicy")
    
    
    return
}

func NewGetPolicyResponse() (response *GetPolicyResponse) {
    response = &GetPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetPolicy
// 本接口（GetPolicy）可用于查询查看策略详情。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetPolicy(request *GetPolicyRequest) (response *GetPolicyResponse, err error) {
    return c.GetPolicyWithContext(context.Background(), request)
}

// GetPolicy
// 本接口（GetPolicy）可用于查询查看策略详情。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetPolicyWithContext(ctx context.Context, request *GetPolicyRequest) (response *GetPolicyResponse, err error) {
    if request == nil {
        request = NewGetPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewGetPolicyVersionRequest() (request *GetPolicyVersionRequest) {
    request = &GetPolicyVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetPolicyVersion")
    
    
    return
}

func NewGetPolicyVersionResponse() (response *GetPolicyVersionResponse) {
    response = &GetPolicyVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetPolicyVersion
// 该接口（GetPolicyVersion）用于查询策略版本详情
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetPolicyVersion(request *GetPolicyVersionRequest) (response *GetPolicyVersionResponse, err error) {
    return c.GetPolicyVersionWithContext(context.Background(), request)
}

// GetPolicyVersion
// 该接口（GetPolicyVersion）用于查询策略版本详情
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetPolicyVersionWithContext(ctx context.Context, request *GetPolicyVersionRequest) (response *GetPolicyVersionResponse, err error) {
    if request == nil {
        request = NewGetPolicyVersionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetPolicyVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetPolicyVersionResponse()
    err = c.Send(request, response)
    return
}

func NewGetRoleRequest() (request *GetRoleRequest) {
    request = &GetRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetRole")
    
    
    return
}

func NewGetRoleResponse() (response *GetRoleResponse) {
    response = &GetRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetRole
// 本接口（GetRole）用于获取指定角色的详细信息。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) GetRole(request *GetRoleRequest) (response *GetRoleResponse, err error) {
    return c.GetRoleWithContext(context.Background(), request)
}

// GetRole
// 本接口（GetRole）用于获取指定角色的详细信息。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) GetRoleWithContext(ctx context.Context, request *GetRoleRequest) (response *GetRoleResponse, err error) {
    if request == nil {
        request = NewGetRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetRoleResponse()
    err = c.Send(request, response)
    return
}

func NewGetRolePermissionBoundaryRequest() (request *GetRolePermissionBoundaryRequest) {
    request = &GetRolePermissionBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetRolePermissionBoundary")
    
    
    return
}

func NewGetRolePermissionBoundaryResponse() (response *GetRolePermissionBoundaryResponse) {
    response = &GetRolePermissionBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetRolePermissionBoundary
// 获取角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetRolePermissionBoundary(request *GetRolePermissionBoundaryRequest) (response *GetRolePermissionBoundaryResponse, err error) {
    return c.GetRolePermissionBoundaryWithContext(context.Background(), request)
}

// GetRolePermissionBoundary
// 获取角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetRolePermissionBoundaryWithContext(ctx context.Context, request *GetRolePermissionBoundaryRequest) (response *GetRolePermissionBoundaryResponse, err error) {
    if request == nil {
        request = NewGetRolePermissionBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetRolePermissionBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetRolePermissionBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewGetSAMLProviderRequest() (request *GetSAMLProviderRequest) {
    request = &GetSAMLProviderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetSAMLProvider")
    
    
    return
}

func NewGetSAMLProviderResponse() (response *GetSAMLProviderResponse) {
    response = &GetSAMLProviderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetSAMLProvider
// 查询SAML身份提供商详情
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) GetSAMLProvider(request *GetSAMLProviderRequest) (response *GetSAMLProviderResponse, err error) {
    return c.GetSAMLProviderWithContext(context.Background(), request)
}

// GetSAMLProvider
// 查询SAML身份提供商详情
//
// 可能返回的错误码:
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
func (c *Client) GetSAMLProviderWithContext(ctx context.Context, request *GetSAMLProviderRequest) (response *GetSAMLProviderResponse, err error) {
    if request == nil {
        request = NewGetSAMLProviderRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetSAMLProvider require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetSAMLProviderResponse()
    err = c.Send(request, response)
    return
}

func NewGetSecurityLastUsedRequest() (request *GetSecurityLastUsedRequest) {
    request = &GetSecurityLastUsedRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetSecurityLastUsed")
    
    
    return
}

func NewGetSecurityLastUsedResponse() (response *GetSecurityLastUsedResponse) {
    response = &GetSecurityLastUsedResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetSecurityLastUsed
// 获取密钥最近使用情况
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) GetSecurityLastUsed(request *GetSecurityLastUsedRequest) (response *GetSecurityLastUsedResponse, err error) {
    return c.GetSecurityLastUsedWithContext(context.Background(), request)
}

// GetSecurityLastUsed
// 获取密钥最近使用情况
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) GetSecurityLastUsedWithContext(ctx context.Context, request *GetSecurityLastUsedRequest) (response *GetSecurityLastUsedResponse, err error) {
    if request == nil {
        request = NewGetSecurityLastUsedRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetSecurityLastUsed require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetSecurityLastUsedResponse()
    err = c.Send(request, response)
    return
}

func NewGetServiceLinkedRoleDeletionStatusRequest() (request *GetServiceLinkedRoleDeletionStatusRequest) {
    request = &GetServiceLinkedRoleDeletionStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetServiceLinkedRoleDeletionStatus")
    
    
    return
}

func NewGetServiceLinkedRoleDeletionStatusResponse() (response *GetServiceLinkedRoleDeletionStatusResponse) {
    response = &GetServiceLinkedRoleDeletionStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetServiceLinkedRoleDeletionStatus
// 根据删除TaskId获取服务相关角色删除状态
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_DELETIONTASKNOTEXIST = "InvalidParameter.DeletionTaskNotExist"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) GetServiceLinkedRoleDeletionStatus(request *GetServiceLinkedRoleDeletionStatusRequest) (response *GetServiceLinkedRoleDeletionStatusResponse, err error) {
    return c.GetServiceLinkedRoleDeletionStatusWithContext(context.Background(), request)
}

// GetServiceLinkedRoleDeletionStatus
// 根据删除TaskId获取服务相关角色删除状态
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_DELETIONTASKNOTEXIST = "InvalidParameter.DeletionTaskNotExist"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) GetServiceLinkedRoleDeletionStatusWithContext(ctx context.Context, request *GetServiceLinkedRoleDeletionStatusRequest) (response *GetServiceLinkedRoleDeletionStatusResponse, err error) {
    if request == nil {
        request = NewGetServiceLinkedRoleDeletionStatusRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetServiceLinkedRoleDeletionStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetServiceLinkedRoleDeletionStatusResponse()
    err = c.Send(request, response)
    return
}

func NewGetUserRequest() (request *GetUserRequest) {
    request = &GetUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetUser")
    
    
    return
}

func NewGetUserResponse() (response *GetUserResponse) {
    response = &GetUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetUser
// 查询子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) GetUser(request *GetUserRequest) (response *GetUserResponse, err error) {
    return c.GetUserWithContext(context.Background(), request)
}

// GetUser
// 查询子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) GetUserWithContext(ctx context.Context, request *GetUserRequest) (response *GetUserResponse, err error) {
    if request == nil {
        request = NewGetUserRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetUser require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetUserResponse()
    err = c.Send(request, response)
    return
}

func NewGetUserAppIdRequest() (request *GetUserAppIdRequest) {
    request = &GetUserAppIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetUserAppId")
    
    
    return
}

func NewGetUserAppIdResponse() (response *GetUserAppIdResponse) {
    response = &GetUserAppIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetUserAppId
// 获取用户AppId
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
func (c *Client) GetUserAppId(request *GetUserAppIdRequest) (response *GetUserAppIdResponse, err error) {
    return c.GetUserAppIdWithContext(context.Background(), request)
}

// GetUserAppId
// 获取用户AppId
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
func (c *Client) GetUserAppIdWithContext(ctx context.Context, request *GetUserAppIdRequest) (response *GetUserAppIdResponse, err error) {
    if request == nil {
        request = NewGetUserAppIdRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetUserAppId require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetUserAppIdResponse()
    err = c.Send(request, response)
    return
}

func NewGetUserPermissionBoundaryRequest() (request *GetUserPermissionBoundaryRequest) {
    request = &GetUserPermissionBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "GetUserPermissionBoundary")
    
    
    return
}

func NewGetUserPermissionBoundaryResponse() (response *GetUserPermissionBoundaryResponse) {
    response = &GetUserPermissionBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// GetUserPermissionBoundary
// 获取用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetUserPermissionBoundary(request *GetUserPermissionBoundaryRequest) (response *GetUserPermissionBoundaryResponse, err error) {
    return c.GetUserPermissionBoundaryWithContext(context.Background(), request)
}

// GetUserPermissionBoundary
// 获取用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) GetUserPermissionBoundaryWithContext(ctx context.Context, request *GetUserPermissionBoundaryRequest) (response *GetUserPermissionBoundaryResponse, err error) {
    if request == nil {
        request = NewGetUserPermissionBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("GetUserPermissionBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewGetUserPermissionBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewListAccessKeysRequest() (request *ListAccessKeysRequest) {
    request = &ListAccessKeysRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListAccessKeys")
    
    
    return
}

func NewListAccessKeysResponse() (response *ListAccessKeysResponse) {
    response = &ListAccessKeysResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListAccessKeys
// 列出指定CAM用户的访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INTERNALERROR_UNKNOWNERROR = "InternalError.UnknownError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) ListAccessKeys(request *ListAccessKeysRequest) (response *ListAccessKeysResponse, err error) {
    return c.ListAccessKeysWithContext(context.Background(), request)
}

// ListAccessKeys
// 列出指定CAM用户的访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INTERNALERROR_UNKNOWNERROR = "InternalError.UnknownError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) ListAccessKeysWithContext(ctx context.Context, request *ListAccessKeysRequest) (response *ListAccessKeysResponse, err error) {
    if request == nil {
        request = NewListAccessKeysRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListAccessKeys require credential")
    }

    request.SetContext(ctx)
    
    response = NewListAccessKeysResponse()
    err = c.Send(request, response)
    return
}

func NewListAttachedGroupPoliciesRequest() (request *ListAttachedGroupPoliciesRequest) {
    request = &ListAttachedGroupPoliciesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListAttachedGroupPolicies")
    
    
    return
}

func NewListAttachedGroupPoliciesResponse() (response *ListAttachedGroupPoliciesResponse) {
    response = &ListAttachedGroupPoliciesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListAttachedGroupPolicies
// 本接口（ListAttachedGroupPolicies）可用于查询用户组关联的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedGroupPolicies(request *ListAttachedGroupPoliciesRequest) (response *ListAttachedGroupPoliciesResponse, err error) {
    return c.ListAttachedGroupPoliciesWithContext(context.Background(), request)
}

// ListAttachedGroupPolicies
// 本接口（ListAttachedGroupPolicies）可用于查询用户组关联的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedGroupPoliciesWithContext(ctx context.Context, request *ListAttachedGroupPoliciesRequest) (response *ListAttachedGroupPoliciesResponse, err error) {
    if request == nil {
        request = NewListAttachedGroupPoliciesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListAttachedGroupPolicies require credential")
    }

    request.SetContext(ctx)
    
    response = NewListAttachedGroupPoliciesResponse()
    err = c.Send(request, response)
    return
}

func NewListAttachedRolePoliciesRequest() (request *ListAttachedRolePoliciesRequest) {
    request = &ListAttachedRolePoliciesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListAttachedRolePolicies")
    
    
    return
}

func NewListAttachedRolePoliciesResponse() (response *ListAttachedRolePoliciesResponse) {
    response = &ListAttachedRolePoliciesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListAttachedRolePolicies
// 本接口（ListAttachedRolePolicies）用于获取角色绑定的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedRolePolicies(request *ListAttachedRolePoliciesRequest) (response *ListAttachedRolePoliciesResponse, err error) {
    return c.ListAttachedRolePoliciesWithContext(context.Background(), request)
}

// ListAttachedRolePolicies
// 本接口（ListAttachedRolePolicies）用于获取角色绑定的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedRolePoliciesWithContext(ctx context.Context, request *ListAttachedRolePoliciesRequest) (response *ListAttachedRolePoliciesResponse, err error) {
    if request == nil {
        request = NewListAttachedRolePoliciesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListAttachedRolePolicies require credential")
    }

    request.SetContext(ctx)
    
    response = NewListAttachedRolePoliciesResponse()
    err = c.Send(request, response)
    return
}

func NewListAttachedUserAllPoliciesRequest() (request *ListAttachedUserAllPoliciesRequest) {
    request = &ListAttachedUserAllPoliciesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListAttachedUserAllPolicies")
    
    
    return
}

func NewListAttachedUserAllPoliciesResponse() (response *ListAttachedUserAllPoliciesResponse) {
    response = &ListAttachedUserAllPoliciesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListAttachedUserAllPolicies
// 列出用户关联的策略（包括随组关联）
//
// 可能返回的错误码:
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ListAttachedUserAllPolicies(request *ListAttachedUserAllPoliciesRequest) (response *ListAttachedUserAllPoliciesResponse, err error) {
    return c.ListAttachedUserAllPoliciesWithContext(context.Background(), request)
}

// ListAttachedUserAllPolicies
// 列出用户关联的策略（包括随组关联）
//
// 可能返回的错误码:
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ListAttachedUserAllPoliciesWithContext(ctx context.Context, request *ListAttachedUserAllPoliciesRequest) (response *ListAttachedUserAllPoliciesResponse, err error) {
    if request == nil {
        request = NewListAttachedUserAllPoliciesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListAttachedUserAllPolicies require credential")
    }

    request.SetContext(ctx)
    
    response = NewListAttachedUserAllPoliciesResponse()
    err = c.Send(request, response)
    return
}

func NewListAttachedUserPoliciesRequest() (request *ListAttachedUserPoliciesRequest) {
    request = &ListAttachedUserPoliciesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListAttachedUserPolicies")
    
    
    return
}

func NewListAttachedUserPoliciesResponse() (response *ListAttachedUserPoliciesResponse) {
    response = &ListAttachedUserPoliciesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListAttachedUserPolicies
// 本接口（ListAttachedUserPolicies）可用于查询子账号关联的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedUserPolicies(request *ListAttachedUserPoliciesRequest) (response *ListAttachedUserPoliciesResponse, err error) {
    return c.ListAttachedUserPoliciesWithContext(context.Background(), request)
}

// ListAttachedUserPolicies
// 本接口（ListAttachedUserPolicies）可用于查询子账号关联的策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListAttachedUserPoliciesWithContext(ctx context.Context, request *ListAttachedUserPoliciesRequest) (response *ListAttachedUserPoliciesResponse, err error) {
    if request == nil {
        request = NewListAttachedUserPoliciesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListAttachedUserPolicies require credential")
    }

    request.SetContext(ctx)
    
    response = NewListAttachedUserPoliciesResponse()
    err = c.Send(request, response)
    return
}

func NewListCollaboratorsRequest() (request *ListCollaboratorsRequest) {
    request = &ListCollaboratorsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListCollaborators")
    
    
    return
}

func NewListCollaboratorsResponse() (response *ListCollaboratorsResponse) {
    response = &ListCollaboratorsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListCollaborators
// 获取协作者列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListCollaborators(request *ListCollaboratorsRequest) (response *ListCollaboratorsResponse, err error) {
    return c.ListCollaboratorsWithContext(context.Background(), request)
}

// ListCollaborators
// 获取协作者列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListCollaboratorsWithContext(ctx context.Context, request *ListCollaboratorsRequest) (response *ListCollaboratorsResponse, err error) {
    if request == nil {
        request = NewListCollaboratorsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListCollaborators require credential")
    }

    request.SetContext(ctx)
    
    response = NewListCollaboratorsResponse()
    err = c.Send(request, response)
    return
}

func NewListEntitiesForPolicyRequest() (request *ListEntitiesForPolicyRequest) {
    request = &ListEntitiesForPolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListEntitiesForPolicy")
    
    
    return
}

func NewListEntitiesForPolicyResponse() (response *ListEntitiesForPolicyResponse) {
    response = &ListEntitiesForPolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListEntitiesForPolicy
// 本接口（ListEntitiesForPolicy）可用于查询策略关联的实体列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ENTITYFILTERERROR = "InvalidParameter.EntityFilterError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
func (c *Client) ListEntitiesForPolicy(request *ListEntitiesForPolicyRequest) (response *ListEntitiesForPolicyResponse, err error) {
    return c.ListEntitiesForPolicyWithContext(context.Background(), request)
}

// ListEntitiesForPolicy
// 本接口（ListEntitiesForPolicy）可用于查询策略关联的实体列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ENTITYFILTERERROR = "InvalidParameter.EntityFilterError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
func (c *Client) ListEntitiesForPolicyWithContext(ctx context.Context, request *ListEntitiesForPolicyRequest) (response *ListEntitiesForPolicyResponse, err error) {
    if request == nil {
        request = NewListEntitiesForPolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListEntitiesForPolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewListEntitiesForPolicyResponse()
    err = c.Send(request, response)
    return
}

func NewListGroupsRequest() (request *ListGroupsRequest) {
    request = &ListGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListGroups")
    
    
    return
}

func NewListGroupsResponse() (response *ListGroupsResponse) {
    response = &ListGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListGroups
// 查询用户组列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListGroups(request *ListGroupsRequest) (response *ListGroupsResponse, err error) {
    return c.ListGroupsWithContext(context.Background(), request)
}

// ListGroups
// 查询用户组列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListGroupsWithContext(ctx context.Context, request *ListGroupsRequest) (response *ListGroupsResponse, err error) {
    if request == nil {
        request = NewListGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewListGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewListGroupsForUserRequest() (request *ListGroupsForUserRequest) {
    request = &ListGroupsForUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListGroupsForUser")
    
    
    return
}

func NewListGroupsForUserResponse() (response *ListGroupsForUserResponse) {
    response = &ListGroupsForUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListGroupsForUser
// 列出用户关联的用户组
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) ListGroupsForUser(request *ListGroupsForUserRequest) (response *ListGroupsForUserResponse, err error) {
    return c.ListGroupsForUserWithContext(context.Background(), request)
}

// ListGroupsForUser
// 列出用户关联的用户组
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) ListGroupsForUserWithContext(ctx context.Context, request *ListGroupsForUserRequest) (response *ListGroupsForUserResponse, err error) {
    if request == nil {
        request = NewListGroupsForUserRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListGroupsForUser require credential")
    }

    request.SetContext(ctx)
    
    response = NewListGroupsForUserResponse()
    err = c.Send(request, response)
    return
}

func NewListPoliciesRequest() (request *ListPoliciesRequest) {
    request = &ListPoliciesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListPolicies")
    
    
    return
}

func NewListPoliciesResponse() (response *ListPoliciesResponse) {
    response = &ListPoliciesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListPolicies
// 本接口（ListPolicies）可用于查询策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_GROUPIDERROR = "InvalidParameter.GroupIdError"
//  INVALIDPARAMETER_KEYWORDERROR = "InvalidParameter.KeywordError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_SCOPEERROR = "InvalidParameter.ScopeError"
//  INVALIDPARAMETER_SEARCHKEYWORDLENGTHERROR = "InvalidParameter.SearchKeywordLengthError"
//  INVALIDPARAMETER_SERVICETYPEERROR = "InvalidParameter.ServiceTypeError"
//  INVALIDPARAMETER_UINERROR = "InvalidParameter.UinError"
func (c *Client) ListPolicies(request *ListPoliciesRequest) (response *ListPoliciesResponse, err error) {
    return c.ListPoliciesWithContext(context.Background(), request)
}

// ListPolicies
// 本接口（ListPolicies）可用于查询策略列表。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_GROUPIDERROR = "InvalidParameter.GroupIdError"
//  INVALIDPARAMETER_KEYWORDERROR = "InvalidParameter.KeywordError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_SCOPEERROR = "InvalidParameter.ScopeError"
//  INVALIDPARAMETER_SEARCHKEYWORDLENGTHERROR = "InvalidParameter.SearchKeywordLengthError"
//  INVALIDPARAMETER_SERVICETYPEERROR = "InvalidParameter.ServiceTypeError"
//  INVALIDPARAMETER_UINERROR = "InvalidParameter.UinError"
func (c *Client) ListPoliciesWithContext(ctx context.Context, request *ListPoliciesRequest) (response *ListPoliciesResponse, err error) {
    if request == nil {
        request = NewListPoliciesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListPolicies require credential")
    }

    request.SetContext(ctx)
    
    response = NewListPoliciesResponse()
    err = c.Send(request, response)
    return
}

func NewListPoliciesGrantingServiceAccessRequest() (request *ListPoliciesGrantingServiceAccessRequest) {
    request = &ListPoliciesGrantingServiceAccessRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListPoliciesGrantingServiceAccess")
    
    
    return
}

func NewListPoliciesGrantingServiceAccessResponse() (response *ListPoliciesGrantingServiceAccessResponse) {
    response = &ListPoliciesGrantingServiceAccessResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListPoliciesGrantingServiceAccess
// 获取所有已授权服务
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListPoliciesGrantingServiceAccess(request *ListPoliciesGrantingServiceAccessRequest) (response *ListPoliciesGrantingServiceAccessResponse, err error) {
    return c.ListPoliciesGrantingServiceAccessWithContext(context.Background(), request)
}

// ListPoliciesGrantingServiceAccess
// 获取所有已授权服务
//
// 可能返回的错误码:
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
func (c *Client) ListPoliciesGrantingServiceAccessWithContext(ctx context.Context, request *ListPoliciesGrantingServiceAccessRequest) (response *ListPoliciesGrantingServiceAccessResponse, err error) {
    if request == nil {
        request = NewListPoliciesGrantingServiceAccessRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListPoliciesGrantingServiceAccess require credential")
    }

    request.SetContext(ctx)
    
    response = NewListPoliciesGrantingServiceAccessResponse()
    err = c.Send(request, response)
    return
}

func NewListPolicyVersionsRequest() (request *ListPolicyVersionsRequest) {
    request = &ListPolicyVersionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListPolicyVersions")
    
    
    return
}

func NewListPolicyVersionsResponse() (response *ListPolicyVersionsResponse) {
    response = &ListPolicyVersionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListPolicyVersions
// 该接口（ListPolicyVersions）用于获取策略版本列表
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) ListPolicyVersions(request *ListPolicyVersionsRequest) (response *ListPolicyVersionsResponse, err error) {
    return c.ListPolicyVersionsWithContext(context.Background(), request)
}

// ListPolicyVersions
// 该接口（ListPolicyVersions）用于获取策略版本列表
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) ListPolicyVersionsWithContext(ctx context.Context, request *ListPolicyVersionsRequest) (response *ListPolicyVersionsResponse, err error) {
    if request == nil {
        request = NewListPolicyVersionsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListPolicyVersions require credential")
    }

    request.SetContext(ctx)
    
    response = NewListPolicyVersionsResponse()
    err = c.Send(request, response)
    return
}

func NewListSAMLProvidersRequest() (request *ListSAMLProvidersRequest) {
    request = &ListSAMLProvidersRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListSAMLProviders")
    
    
    return
}

func NewListSAMLProvidersResponse() (response *ListSAMLProvidersResponse) {
    response = &ListSAMLProvidersResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListSAMLProviders
// 查询SAML身份提供商列表
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) ListSAMLProviders(request *ListSAMLProvidersRequest) (response *ListSAMLProvidersResponse, err error) {
    return c.ListSAMLProvidersWithContext(context.Background(), request)
}

// ListSAMLProviders
// 查询SAML身份提供商列表
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
func (c *Client) ListSAMLProvidersWithContext(ctx context.Context, request *ListSAMLProvidersRequest) (response *ListSAMLProvidersResponse, err error) {
    if request == nil {
        request = NewListSAMLProvidersRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListSAMLProviders require credential")
    }

    request.SetContext(ctx)
    
    response = NewListSAMLProvidersResponse()
    err = c.Send(request, response)
    return
}

func NewListUsersRequest() (request *ListUsersRequest) {
    request = &ListUsersRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListUsers")
    
    
    return
}

func NewListUsersResponse() (response *ListUsersResponse) {
    response = &ListUsersResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListUsers
// 拉取子用户
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListUsers(request *ListUsersRequest) (response *ListUsersResponse, err error) {
    return c.ListUsersWithContext(context.Background(), request)
}

// ListUsers
// 拉取子用户
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
func (c *Client) ListUsersWithContext(ctx context.Context, request *ListUsersRequest) (response *ListUsersResponse, err error) {
    if request == nil {
        request = NewListUsersRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListUsers require credential")
    }

    request.SetContext(ctx)
    
    response = NewListUsersResponse()
    err = c.Send(request, response)
    return
}

func NewListUsersForGroupRequest() (request *ListUsersForGroupRequest) {
    request = &ListUsersForGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListUsersForGroup")
    
    
    return
}

func NewListUsersForGroupResponse() (response *ListUsersForGroupResponse) {
    response = &ListUsersForGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListUsersForGroup
// 查询用户组关联的用户列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) ListUsersForGroup(request *ListUsersForGroupRequest) (response *ListUsersForGroupResponse, err error) {
    return c.ListUsersForGroupWithContext(context.Background(), request)
}

// ListUsersForGroup
// 查询用户组关联的用户列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) ListUsersForGroupWithContext(ctx context.Context, request *ListUsersForGroupRequest) (response *ListUsersForGroupResponse, err error) {
    if request == nil {
        request = NewListUsersForGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListUsersForGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewListUsersForGroupResponse()
    err = c.Send(request, response)
    return
}

func NewListWeChatWorkSubAccountsRequest() (request *ListWeChatWorkSubAccountsRequest) {
    request = &ListWeChatWorkSubAccountsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "ListWeChatWorkSubAccounts")
    
    
    return
}

func NewListWeChatWorkSubAccountsResponse() (response *ListWeChatWorkSubAccountsResponse) {
    response = &ListWeChatWorkSubAccountsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListWeChatWorkSubAccounts
// 获取企业微信子用户列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) ListWeChatWorkSubAccounts(request *ListWeChatWorkSubAccountsRequest) (response *ListWeChatWorkSubAccountsResponse, err error) {
    return c.ListWeChatWorkSubAccountsWithContext(context.Background(), request)
}

// ListWeChatWorkSubAccounts
// 获取企业微信子用户列表
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) ListWeChatWorkSubAccountsWithContext(ctx context.Context, request *ListWeChatWorkSubAccountsRequest) (response *ListWeChatWorkSubAccountsResponse, err error) {
    if request == nil {
        request = NewListWeChatWorkSubAccountsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ListWeChatWorkSubAccounts require credential")
    }

    request.SetContext(ctx)
    
    response = NewListWeChatWorkSubAccountsResponse()
    err = c.Send(request, response)
    return
}

func NewPutRolePermissionsBoundaryRequest() (request *PutRolePermissionsBoundaryRequest) {
    request = &PutRolePermissionsBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "PutRolePermissionsBoundary")
    
    
    return
}

func NewPutRolePermissionsBoundaryResponse() (response *PutRolePermissionsBoundaryResponse) {
    response = &PutRolePermissionsBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// PutRolePermissionsBoundary
// 设置角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDPOLICYCANTINPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedPolicyCantInPermissionBoundary"
//  INVALIDPARAMETER_SERVICELINKEDROLECANTUSEPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedRoleCantUsePermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PutRolePermissionsBoundary(request *PutRolePermissionsBoundaryRequest) (response *PutRolePermissionsBoundaryResponse, err error) {
    return c.PutRolePermissionsBoundaryWithContext(context.Background(), request)
}

// PutRolePermissionsBoundary
// 设置角色权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDPOLICYCANTINPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedPolicyCantInPermissionBoundary"
//  INVALIDPARAMETER_SERVICELINKEDROLECANTUSEPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedRoleCantUsePermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PutRolePermissionsBoundaryWithContext(ctx context.Context, request *PutRolePermissionsBoundaryRequest) (response *PutRolePermissionsBoundaryResponse, err error) {
    if request == nil {
        request = NewPutRolePermissionsBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("PutRolePermissionsBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewPutRolePermissionsBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewPutUserPermissionsBoundaryRequest() (request *PutUserPermissionsBoundaryRequest) {
    request = &PutUserPermissionsBoundaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "PutUserPermissionsBoundary")
    
    
    return
}

func NewPutUserPermissionsBoundaryResponse() (response *PutUserPermissionsBoundaryResponse) {
    response = &PutUserPermissionsBoundaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// PutUserPermissionsBoundary
// 设置用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDPOLICYCANTINPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedPolicyCantInPermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PutUserPermissionsBoundary(request *PutUserPermissionsBoundaryRequest) (response *PutUserPermissionsBoundaryResponse, err error) {
    return c.PutUserPermissionsBoundaryWithContext(context.Background(), request)
}

// PutUserPermissionsBoundary
// 设置用户权限边界
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_OPERATEENTITIESOVERLIMIT = "InvalidParameter.OperateEntitiesOverLimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_SERVICELINKEDPOLICYCANTINPERMISSIONBOUNDARY = "InvalidParameter.ServiceLinkedPolicyCantInPermissionBoundary"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) PutUserPermissionsBoundaryWithContext(ctx context.Context, request *PutUserPermissionsBoundaryRequest) (response *PutUserPermissionsBoundaryResponse, err error) {
    if request == nil {
        request = NewPutUserPermissionsBoundaryRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("PutUserPermissionsBoundary require credential")
    }

    request.SetContext(ctx)
    
    response = NewPutUserPermissionsBoundaryResponse()
    err = c.Send(request, response)
    return
}

func NewRemoveUserFromGroupRequest() (request *RemoveUserFromGroupRequest) {
    request = &RemoveUserFromGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "RemoveUserFromGroup")
    
    
    return
}

func NewRemoveUserFromGroupResponse() (response *RemoveUserFromGroupResponse) {
    response = &RemoveUserFromGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// RemoveUserFromGroup
// 从用户组删除用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_USERUINANDUINNOTALLNULL = "InvalidParameter.UserUinAndUinNotAllNull"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) RemoveUserFromGroup(request *RemoveUserFromGroupRequest) (response *RemoveUserFromGroupResponse, err error) {
    return c.RemoveUserFromGroupWithContext(context.Background(), request)
}

// RemoveUserFromGroup
// 从用户组删除用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_USERUINANDUINNOTALLNULL = "InvalidParameter.UserUinAndUinNotAllNull"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) RemoveUserFromGroupWithContext(ctx context.Context, request *RemoveUserFromGroupRequest) (response *RemoveUserFromGroupResponse, err error) {
    if request == nil {
        request = NewRemoveUserFromGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("RemoveUserFromGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewRemoveUserFromGroupResponse()
    err = c.Send(request, response)
    return
}

func NewSetDefaultPolicyVersionRequest() (request *SetDefaultPolicyVersionRequest) {
    request = &SetDefaultPolicyVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "SetDefaultPolicyVersion")
    
    
    return
}

func NewSetDefaultPolicyVersionResponse() (response *SetDefaultPolicyVersionResponse) {
    response = &SetDefaultPolicyVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SetDefaultPolicyVersion
// 本接口（SetDefaultPolicyVersion）可用于设置生效的策略版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONALREADYDEFAULT = "FailedOperation.PolicyVersionAlreadyDefault"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) SetDefaultPolicyVersion(request *SetDefaultPolicyVersionRequest) (response *SetDefaultPolicyVersionResponse, err error) {
    return c.SetDefaultPolicyVersionWithContext(context.Background(), request)
}

// SetDefaultPolicyVersion
// 本接口（SetDefaultPolicyVersion）可用于设置生效的策略版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  FAILEDOPERATION_POLICYVERSIONALREADYDEFAULT = "FailedOperation.PolicyVersionAlreadyDefault"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_POLICYVERSIONNOTEXISTS = "InvalidParameter.PolicyVersionNotExists"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) SetDefaultPolicyVersionWithContext(ctx context.Context, request *SetDefaultPolicyVersionRequest) (response *SetDefaultPolicyVersionResponse, err error) {
    if request == nil {
        request = NewSetDefaultPolicyVersionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SetDefaultPolicyVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewSetDefaultPolicyVersionResponse()
    err = c.Send(request, response)
    return
}

func NewSetMfaFlagRequest() (request *SetMfaFlagRequest) {
    request = &SetMfaFlagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "SetMfaFlag")
    
    
    return
}

func NewSetMfaFlagResponse() (response *SetMfaFlagResponse) {
    response = &SetMfaFlagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// SetMfaFlag
// 设置子用户的登录保护和敏感操作校验方式
//
// 可能返回的错误码:
//  FAILEDOPERATION_USERNOTBINDPHONE = "FailedOperation.UserNotBindPhone"
//  FAILEDOPERATION_USERNOTBINDWECHAT = "FailedOperation.UserNotBindWechat"
//  FAILEDOPERATION_USERUNBINDNOPERMISSION = "FailedOperation.UserUnbindNoPermission"
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) SetMfaFlag(request *SetMfaFlagRequest) (response *SetMfaFlagResponse, err error) {
    return c.SetMfaFlagWithContext(context.Background(), request)
}

// SetMfaFlag
// 设置子用户的登录保护和敏感操作校验方式
//
// 可能返回的错误码:
//  FAILEDOPERATION_USERNOTBINDPHONE = "FailedOperation.UserNotBindPhone"
//  FAILEDOPERATION_USERNOTBINDWECHAT = "FailedOperation.UserNotBindWechat"
//  FAILEDOPERATION_USERUNBINDNOPERMISSION = "FailedOperation.UserUnbindNoPermission"
//  INVALIDPARAMETER_MFATOKENERROR = "InvalidParameter.MFATokenError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) SetMfaFlagWithContext(ctx context.Context, request *SetMfaFlagRequest) (response *SetMfaFlagResponse, err error) {
    if request == nil {
        request = NewSetMfaFlagRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("SetMfaFlag require credential")
    }

    request.SetContext(ctx)
    
    response = NewSetMfaFlagResponse()
    err = c.Send(request, response)
    return
}

func NewTagRoleRequest() (request *TagRoleRequest) {
    request = &TagRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "TagRole")
    
    
    return
}

func NewTagRoleResponse() (response *TagRoleResponse) {
    response = &TagRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// TagRole
// 角色绑定标签
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) TagRole(request *TagRoleRequest) (response *TagRoleResponse, err error) {
    return c.TagRoleWithContext(context.Background(), request)
}

// TagRole
// 角色绑定标签
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_TAGRESOURCEFAILED = "FailedOperation.TagResourceFailed"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGLIMITEXCEEDED = "InvalidParameter.TagLimitExceeded"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) TagRoleWithContext(ctx context.Context, request *TagRoleRequest) (response *TagRoleResponse, err error) {
    if request == nil {
        request = NewTagRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("TagRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewTagRoleResponse()
    err = c.Send(request, response)
    return
}

func NewUntagRoleRequest() (request *UntagRoleRequest) {
    request = &UntagRoleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UntagRole")
    
    
    return
}

func NewUntagRoleResponse() (response *UntagRoleResponse) {
    response = &UntagRoleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UntagRole
// 角色解绑标签。
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_UNTAGRESOURCEFAILED = "FailedOperation.UnTagResourceFailed"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) UntagRole(request *UntagRoleRequest) (response *UntagRoleResponse, err error) {
    return c.UntagRoleWithContext(context.Background(), request)
}

// UntagRole
// 角色解绑标签。
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_UNTAGRESOURCEFAILED = "FailedOperation.UnTagResourceFailed"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENAMEERROR = "InvalidParameter.RoleNameError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_TAGPARAMERROR = "InvalidParameter.TagParamError"
func (c *Client) UntagRoleWithContext(ctx context.Context, request *UntagRoleRequest) (response *UntagRoleResponse, err error) {
    if request == nil {
        request = NewUntagRoleRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UntagRole require credential")
    }

    request.SetContext(ctx)
    
    response = NewUntagRoleResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateAccessKeyRequest() (request *UpdateAccessKeyRequest) {
    request = &UpdateAccessKeyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateAccessKey")
    
    
    return
}

func NewUpdateAccessKeyResponse() (response *UpdateAccessKeyResponse) {
    response = &UpdateAccessKeyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateAccessKey
// 为CAM用户更新访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) UpdateAccessKey(request *UpdateAccessKeyRequest) (response *UpdateAccessKeyResponse, err error) {
    return c.UpdateAccessKeyWithContext(context.Background(), request)
}

// UpdateAccessKey
// 为CAM用户更新访问密钥
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION_ACCESSKEY = "FailedOperation.Accesskey"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  OPERATIONDENIED_ACCESSKEYOVERLIMIT = "OperationDenied.AccessKeyOverLimit"
//  OPERATIONDENIED_SUBUIN = "OperationDenied.SubUin"
//  OPERATIONDENIED_UINNOTMATCH = "OperationDenied.UinNotMatch"
//  RESOURCENOTFOUND_SECRETNOTEXIST = "ResourceNotFound.SecretNotExist"
func (c *Client) UpdateAccessKeyWithContext(ctx context.Context, request *UpdateAccessKeyRequest) (response *UpdateAccessKeyResponse, err error) {
    if request == nil {
        request = NewUpdateAccessKeyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateAccessKey require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateAccessKeyResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateAssumeRolePolicyRequest() (request *UpdateAssumeRolePolicyRequest) {
    request = &UpdateAssumeRolePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateAssumeRolePolicy")
    
    
    return
}

func NewUpdateAssumeRolePolicyResponse() (response *UpdateAssumeRolePolicyResponse) {
    response = &UpdateAssumeRolePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateAssumeRolePolicy
// 本接口（UpdateAssumeRolePolicy）用于修改角色信任策略的策略文档。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSCROSSERROR = "InvalidParameter.PrincipalQcsCrossError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) UpdateAssumeRolePolicy(request *UpdateAssumeRolePolicyRequest) (response *UpdateAssumeRolePolicyResponse, err error) {
    return c.UpdateAssumeRolePolicyWithContext(context.Background(), request)
}

// UpdateAssumeRolePolicy
// 本接口（UpdateAssumeRolePolicy）用于修改角色信任策略的策略文档。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSCROSSERROR = "InvalidParameter.PrincipalQcsCrossError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
func (c *Client) UpdateAssumeRolePolicyWithContext(ctx context.Context, request *UpdateAssumeRolePolicyRequest) (response *UpdateAssumeRolePolicyResponse, err error) {
    if request == nil {
        request = NewUpdateAssumeRolePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateAssumeRolePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateAssumeRolePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateGroupRequest() (request *UpdateGroupRequest) {
    request = &UpdateGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateGroup")
    
    
    return
}

func NewUpdateGroupResponse() (response *UpdateGroupResponse) {
    response = &UpdateGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateGroup
// 更新用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPNAMEINUSE = "InvalidParameter.GroupNameInUse"
//  INVALIDPARAMETER_GROUPNOTEXIST = "InvalidParameter.GroupNotExist"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) UpdateGroup(request *UpdateGroupRequest) (response *UpdateGroupResponse, err error) {
    return c.UpdateGroupWithContext(context.Background(), request)
}

// UpdateGroup
// 更新用户组
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_GROUPNAMEINUSE = "InvalidParameter.GroupNameInUse"
//  INVALIDPARAMETER_GROUPNOTEXIST = "InvalidParameter.GroupNotExist"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
func (c *Client) UpdateGroupWithContext(ctx context.Context, request *UpdateGroupRequest) (response *UpdateGroupResponse, err error) {
    if request == nil {
        request = NewUpdateGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateGroupResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateOIDCConfigRequest() (request *UpdateOIDCConfigRequest) {
    request = &UpdateOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateOIDCConfig")
    
    
    return
}

func NewUpdateOIDCConfigResponse() (response *UpdateOIDCConfigResponse) {
    response = &UpdateOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateOIDCConfig
// 修改角色OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateOIDCConfig(request *UpdateOIDCConfigRequest) (response *UpdateOIDCConfigResponse, err error) {
    return c.UpdateOIDCConfigWithContext(context.Background(), request)
}

// UpdateOIDCConfig
// 修改角色OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  INVALIDPARAMETERVALUE_NAMEERROR = "InvalidParameterValue.NameError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateOIDCConfigWithContext(ctx context.Context, request *UpdateOIDCConfigRequest) (response *UpdateOIDCConfigResponse, err error) {
    if request == nil {
        request = NewUpdateOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewUpdatePolicyRequest() (request *UpdatePolicyRequest) {
    request = &UpdatePolicyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdatePolicy")
    
    
    return
}

func NewUpdatePolicyResponse() (response *UpdatePolicyResponse) {
    response = &UpdatePolicyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdatePolicy
// 本接口（UpdatePolicy ）可用于更新策略。
//
// 如果已存在策略版本，本接口会直接更新策略的默认版本，不会创建新版本，如果不存在任何策略版本，则直接创建一个默认版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) UpdatePolicy(request *UpdatePolicyRequest) (response *UpdatePolicyResponse, err error) {
    return c.UpdatePolicyWithContext(context.Background(), request)
}

// UpdatePolicy
// 本接口（UpdatePolicy ）可用于更新策略。
//
// 如果已存在策略版本，本接口会直接更新策略的默认版本，不会创建新版本，如果不存在任何策略版本，则直接创建一个默认版本。
//
// 可能返回的错误码:
//  FAILEDOPERATION_POLICYNAMEINUSE = "FailedOperation.PolicyNameInUse"
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_ACTIONERROR = "InvalidParameter.ActionError"
//  INVALIDPARAMETER_ACTIONMISS = "InvalidParameter.ActionMiss"
//  INVALIDPARAMETER_ACTIONNOTEXIST = "InvalidParameter.ActionNotExist"
//  INVALIDPARAMETER_ACTIONSERVICENOTEXIST = "InvalidParameter.ActionServiceNotExist"
//  INVALIDPARAMETER_ATTACHMENTFULL = "InvalidParameter.AttachmentFull"
//  INVALIDPARAMETER_CONDITIONCONTENTERROR = "InvalidParameter.ConditionContentError"
//  INVALIDPARAMETER_CONDITIONERROR = "InvalidParameter.ConditionError"
//  INVALIDPARAMETER_CONDITIONTYPEERROR = "InvalidParameter.ConditionTypeError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_EFFECTERROR = "InvalidParameter.EffectError"
//  INVALIDPARAMETER_NOTSUPPORTPRODUCT = "InvalidParameter.NotSupportProduct"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_POLICYDOCUMENTERROR = "InvalidParameter.PolicyDocumentError"
//  INVALIDPARAMETER_POLICYDOCUMENTLENGTHOVERLIMIT = "InvalidParameter.PolicyDocumentLengthOverLimit"
//  INVALIDPARAMETER_POLICYIDERROR = "InvalidParameter.PolicyIdError"
//  INVALIDPARAMETER_POLICYIDNOTEXIST = "InvalidParameter.PolicyIdNotExist"
//  INVALIDPARAMETER_POLICYNAMEERROR = "InvalidParameter.PolicyNameError"
//  INVALIDPARAMETER_PRINCIPALERROR = "InvalidParameter.PrincipalError"
//  INVALIDPARAMETER_PRINCIPALQCSERROR = "InvalidParameter.PrincipalQcsError"
//  INVALIDPARAMETER_PRINCIPALQCSNOTEXIST = "InvalidParameter.PrincipalQcsNotExist"
//  INVALIDPARAMETER_PRINCIPALSERVICENOTEXIST = "InvalidParameter.PrincipalServiceNotExist"
//  INVALIDPARAMETER_RESOURCECONTENTERROR = "InvalidParameter.ResourceContentError"
//  INVALIDPARAMETER_RESOURCEERROR = "InvalidParameter.ResourceError"
//  INVALIDPARAMETER_RESOURCEPROJECTERROR = "InvalidParameter.ResourceProjectError"
//  INVALIDPARAMETER_RESOURCEQCSERROR = "InvalidParameter.ResourceQcsError"
//  INVALIDPARAMETER_RESOURCEREGIONERROR = "InvalidParameter.ResourceRegionError"
//  INVALIDPARAMETER_RESOURCESERVICENOTEXIST = "InvalidParameter.ResourceServiceNotExist"
//  INVALIDPARAMETER_RESOURCEUINERROR = "InvalidParameter.ResourceUinError"
//  INVALIDPARAMETER_STATEMENTERROR = "InvalidParameter.StatementError"
//  INVALIDPARAMETER_USERNOTEXIST = "InvalidParameter.UserNotExist"
//  INVALIDPARAMETER_VERSIONERROR = "InvalidParameter.VersionError"
//  RESOURCENOTFOUND_GROUPNOTEXIST = "ResourceNotFound.GroupNotExist"
//  RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"
//  RESOURCENOTFOUND_POLICYIDNOTFOUND = "ResourceNotFound.PolicyIdNotFound"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) UpdatePolicyWithContext(ctx context.Context, request *UpdatePolicyRequest) (response *UpdatePolicyResponse, err error) {
    if request == nil {
        request = NewUpdatePolicyRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdatePolicy require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdatePolicyResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateRoleConsoleLoginRequest() (request *UpdateRoleConsoleLoginRequest) {
    request = &UpdateRoleConsoleLoginRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateRoleConsoleLogin")
    
    
    return
}

func NewUpdateRoleConsoleLoginResponse() (response *UpdateRoleConsoleLoginResponse) {
    response = &UpdateRoleConsoleLoginResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateRoleConsoleLogin
// 本接口（UpdateRoleConsoleLogin）用于修改角色是否可登录。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) UpdateRoleConsoleLogin(request *UpdateRoleConsoleLoginRequest) (response *UpdateRoleConsoleLoginResponse, err error) {
    return c.UpdateRoleConsoleLoginWithContext(context.Background(), request)
}

// UpdateRoleConsoleLogin
// 本接口（UpdateRoleConsoleLogin）用于修改角色是否可登录。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) UpdateRoleConsoleLoginWithContext(ctx context.Context, request *UpdateRoleConsoleLoginRequest) (response *UpdateRoleConsoleLoginResponse, err error) {
    if request == nil {
        request = NewUpdateRoleConsoleLoginRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateRoleConsoleLogin require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateRoleConsoleLoginResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateRoleDescriptionRequest() (request *UpdateRoleDescriptionRequest) {
    request = &UpdateRoleDescriptionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateRoleDescription")
    
    
    return
}

func NewUpdateRoleDescriptionResponse() (response *UpdateRoleDescriptionResponse) {
    response = &UpdateRoleDescriptionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateRoleDescription
// 本接口（UpdateRoleDescription）用于修改角色的描述信息。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) UpdateRoleDescription(request *UpdateRoleDescriptionRequest) (response *UpdateRoleDescriptionResponse, err error) {
    return c.UpdateRoleDescriptionWithContext(context.Background(), request)
}

// UpdateRoleDescription
// 本接口（UpdateRoleDescription）用于修改角色的描述信息。
//
// 可能返回的错误码:
//  INTERNALERROR_SYSTEMERROR = "InternalError.SystemError"
//  INVALIDPARAMETER_DESCRIPTIONLENGTHOVERLIMIT = "InvalidParameter.DescriptionLengthOverlimit"
//  INVALIDPARAMETER_PARAMERROR = "InvalidParameter.ParamError"
//  INVALIDPARAMETER_ROLENOTEXIST = "InvalidParameter.RoleNotExist"
func (c *Client) UpdateRoleDescriptionWithContext(ctx context.Context, request *UpdateRoleDescriptionRequest) (response *UpdateRoleDescriptionResponse, err error) {
    if request == nil {
        request = NewUpdateRoleDescriptionRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateRoleDescription require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateRoleDescriptionResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateSAMLProviderRequest() (request *UpdateSAMLProviderRequest) {
    request = &UpdateSAMLProviderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateSAMLProvider")
    
    
    return
}

func NewUpdateSAMLProviderResponse() (response *UpdateSAMLProviderResponse) {
    response = &UpdateSAMLProviderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateSAMLProvider
// 更新SAML身份提供商信息
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateSAMLProvider(request *UpdateSAMLProviderRequest) (response *UpdateSAMLProviderResponse, err error) {
    return c.UpdateSAMLProviderWithContext(context.Background(), request)
}

// UpdateSAMLProvider
// 更新SAML身份提供商信息
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateSAMLProviderWithContext(ctx context.Context, request *UpdateSAMLProviderRequest) (response *UpdateSAMLProviderResponse, err error) {
    if request == nil {
        request = NewUpdateSAMLProviderRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateSAMLProvider require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateSAMLProviderResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateUserRequest() (request *UpdateUserRequest) {
    request = &UpdateUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateUser")
    
    
    return
}

func NewUpdateUserResponse() (response *UpdateUserResponse) {
    response = &UpdateUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateUser
// 更新子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PASSWORDVIOLATEDRULES = "InvalidParameter.PasswordViolatedRules"
//  INVALIDPARAMETER_USERNAMEILLEGAL = "InvalidParameter.UserNameIllegal"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) UpdateUser(request *UpdateUserRequest) (response *UpdateUserResponse, err error) {
    return c.UpdateUserWithContext(context.Background(), request)
}

// UpdateUser
// 更新子用户
//
// 可能返回的错误码:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER_PASSWORDVIOLATEDRULES = "InvalidParameter.PasswordViolatedRules"
//  INVALIDPARAMETER_USERNAMEILLEGAL = "InvalidParameter.UserNameIllegal"
//  RESOURCENOTFOUND_USERNOTEXIST = "ResourceNotFound.UserNotExist"
func (c *Client) UpdateUserWithContext(ctx context.Context, request *UpdateUserRequest) (response *UpdateUserResponse, err error) {
    if request == nil {
        request = NewUpdateUserRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateUser require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateUserResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateUserOIDCConfigRequest() (request *UpdateUserOIDCConfigRequest) {
    request = &UpdateUserOIDCConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateUserOIDCConfig")
    
    
    return
}

func NewUpdateUserOIDCConfigResponse() (response *UpdateUserOIDCConfigResponse) {
    response = &UpdateUserOIDCConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateUserOIDCConfig
// 修改用户OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYKEYERROR = "InvalidParameterValue.IdentityKeyError"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateUserOIDCConfig(request *UpdateUserOIDCConfigRequest) (response *UpdateUserOIDCConfigResponse, err error) {
    return c.UpdateUserOIDCConfigWithContext(context.Background(), request)
}

// UpdateUserOIDCConfig
// 修改用户OIDC配置
//
// 可能返回的错误码:
//  FAILEDOPERATION = "FailedOperation"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_IDENTITYNAMEINUSE = "InvalidParameter.IdentityNameInUse"
//  INVALIDPARAMETERVALUE_IDENTITYKEYERROR = "InvalidParameterValue.IdentityKeyError"
//  INVALIDPARAMETERVALUE_IDENTITYURLERROR = "InvalidParameterValue.IdentityUrlError"
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  LIMITEXCEEDED_IDENTITYFULL = "LimitExceeded.IdentityFull"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateUserOIDCConfigWithContext(ctx context.Context, request *UpdateUserOIDCConfigRequest) (response *UpdateUserOIDCConfigResponse, err error) {
    if request == nil {
        request = NewUpdateUserOIDCConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateUserOIDCConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateUserOIDCConfigResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateUserSAMLConfigRequest() (request *UpdateUserSAMLConfigRequest) {
    request = &UpdateUserSAMLConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cam", APIVersion, "UpdateUserSAMLConfig")
    
    
    return
}

func NewUpdateUserSAMLConfigResponse() (response *UpdateUserSAMLConfigResponse) {
    response = &UpdateUserSAMLConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// UpdateUserSAMLConfig
// 修改用户SAML配置
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateUserSAMLConfig(request *UpdateUserSAMLConfigRequest) (response *UpdateUserSAMLConfigResponse, err error) {
    return c.UpdateUserSAMLConfigWithContext(context.Background(), request)
}

// UpdateUserSAMLConfig
// 修改用户SAML配置
//
// 可能返回的错误码:
//  INVALIDPARAMETERVALUE_METADATAERROR = "InvalidParameterValue.MetadataError"
//  RESOURCENOTFOUND_IDENTITYNOTEXIST = "ResourceNotFound.IdentityNotExist"
func (c *Client) UpdateUserSAMLConfigWithContext(ctx context.Context, request *UpdateUserSAMLConfigRequest) (response *UpdateUserSAMLConfigResponse, err error) {
    if request == nil {
        request = NewUpdateUserSAMLConfigRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("UpdateUserSAMLConfig require credential")
    }

    request.SetContext(ctx)
    
    response = NewUpdateUserSAMLConfigResponse()
    err = c.Send(request, response)
    return
}
