#coding: utf-8
import json, dotenv, os, time, sys

from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.cvm.v20170312 import cvm_client, models as cvm_models
from tencentcloud.cbs.v20170312 import cbs_client, models as cbs_models

def get_zones(region):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeZonesRequest()
    req.from_json_string(json.dumps({}))

    resp = client.DescribeZones(req)
    return json.loads(resp.to_json_string())

def get_instance_type(region):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeInstanceTypeConfigsRequest()
    req.from_json_string(json.dumps({}))

    resp = client.DescribeInstanceTypeConfigs(req)
    return json.loads(resp.to_json_string())

def get_images(region, filter_image_name):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeImagesRequest()
    params = {
        "Filters": [
            {
                "Name": "image-type",
                "Values": [ "PUBLIC_IMAGE" ]
            },
            {
                "Name": "image-name",
                "Values": [ filter_image_name ]
            }
        ],
        "Limit": 100
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeImages(req)
    return json.loads(resp.to_json_string())

def get_zone_instance(region):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeZoneInstanceConfigInfosRequest()
    params = {
        "Filters": [
            {
                "Name": "instance-charge-type",
                "Values": [ "POSTPAID_BY_HOUR" ]
            },
            {
                "Name": "sort-keys",
                "Values": [ "cpu:asc" ]
            }
        ]
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeZoneInstanceConfigInfos(req)
    return json.loads(resp.to_json_string())

def create_instance(region, zone, instance_type, image_id, disk, network, password):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.RunInstancesRequest()
    params = {
        "InstanceChargeType": "POSTPAID_BY_HOUR",
        "Placement": {
            "Zone": zone
        },
        "InstanceType": instance_type,
        "ImageId": image_id,
        "SystemDisk": {
            "DiskType": "CLOUD_BSSD",
            "DiskSize": disk
        },
        "InternetAccessible": {
            "InternetChargeType": "TRAFFIC_POSTPAID_BY_HOUR",
            "InternetMaxBandwidthOut": network,
            "PublicIpAssigned": True
        },
        "LoginSettings": {
            "Password": password
        }
    }
    req.from_json_string(json.dumps(params))

    resp = client.RunInstances(req)
    return json.loads(resp.to_json_string())

def delete_instance(region, instance_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    # 实例化一个请求对象,每个接口都会对应一个request对象
    req = cvm_models.TerminateInstancesRequest()
    params = {
        "InstanceIds": [ instance_id ],
        "ReleasePrepaidDataDisks": True
    }
    req.from_json_string(json.dumps(params))

    resp = client.TerminateInstances(req)
    return json.loads(resp.to_json_string())

def query_instance_status(region, instance_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    # 实例化一个client选项，可选的，没有特殊需求可以跳过
    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    # 实例化要请求产品的client对象,clientProfile是可选的
    client = cvm_client.CvmClient(cred, region, clientProfile)

    # 实例化一个请求对象,每个接口都会对应一个request对象
    req = cvm_models.DescribeInstancesStatusRequest()
    params = {
        "InstanceIds": [ instance_id ]
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeInstancesStatus(req)
    return json.loads(resp.to_json_string())

def query_instance_detail(region, instance_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeInstancesRequest()
    params = {
        "InstanceIds": [ instance_id ]
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeInstances(req)
    return json.loads(resp.to_json_string())

def delete_instance(region, instance_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    # 实例化一个请求对象,每个接口都会对应一个request对象
    req = cvm_models.TerminateInstancesRequest()
    params = {
        "InstanceIds": [ instance_id ],
        "ReleasePrepaidDataDisks": True
    }
    req.from_json_string(json.dumps(params))

    resp = client.TerminateInstances(req)
    return json.loads(resp.to_json_string())

def query_instances(region):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeInstancesRequest()
    params = {
        "Limit": 100
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeInstances(req)
    return json.loads(resp.to_json_string())

def get_regions():
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, "", clientProfile)

    req = cvm_models.DescribeRegionsRequest()
    req.from_json_string(json.dumps({}))

    resp = client.DescribeRegions(req)
    return json.loads(resp.to_json_string())

def get_snapshorts(region):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cbs.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cbs_client.CbsClient(cred, region, clientProfile)

    req = cbs_models.DescribeSnapshotsRequest()
    params = {
        "Limit": 100
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeSnapshots(req)
    snapshorts = json.loads(resp.to_json_string())
    return snapshorts

def get_image_share(region, image_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeImageSharePermissionRequest()
    params = {
        "ImageId": image_id
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeImageSharePermission(req)
    imageShare = json.loads(resp.to_json_string())
    return imageShare

def cancel_image_share(region, image_id, account_ids):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.ModifyImageSharePermissionRequest()
    params = {
        "ImageId": image_id,
        "AccountIds": account_ids,
        "Permission": "CANCEL"
    }
    req.from_json_string(json.dumps(params))

    resp = client.ModifyImageSharePermission(req)
    return json.loads(resp.to_json_string())

def delete_image_and_snapshot(region, image_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DeleteImagesRequest()
    params = {
        "ImageIds": [ image_id ],
        "DeleteBindedSnap": True
    }
    req.from_json_string(json.dumps(params))

    resp = client.DeleteImages(req)
    return json.loads(resp.to_json_string())

def delete_snapshot_and_image(region, snapshot_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cbs.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cbs_client.CbsClient(cred, region, clientProfile)

    req = cbs_models.DeleteSnapshotsRequest()
    params = {
        "SnapshotIds": [ snapshot_id ],
        "DeleteBindImages": True
    }
    req.from_json_string(json.dumps(params))

    resp = client.DeleteSnapshots(req)
    return json.loads(resp.to_json_string())

def create_image(region, instance_id, image_name, image_desc):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.CreateImageRequest()
    params = {
        "InstanceId": instance_id,
        "ImageName": image_name,
        "ImageDescription": image_desc,
        "ForcePoweroff": "TRUE"
    }
    req.from_json_string(json.dumps(params))

    resp = client.CreateImage(req)
    return json.loads(resp.to_json_string())

def share_image(region, image_id, account_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.ModifyImageSharePermissionRequest()
    params = {
        "ImageId": image_id,
        "AccountIds": [ account_id ],
        "Permission": "SHARE"
    }
    req.from_json_string(json.dumps(params))

    resp = client.ModifyImageSharePermission(req)
    return json.loads(resp.to_json_string())

def query_image(region, image_id):
    cred = credential.Credential(os.getenv("SECRET_ID"), os.getenv("SECRET_KEY"))
    httpProfile = HttpProfile()
    httpProfile.endpoint = "cvm.tencentcloudapi.com"

    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    client = cvm_client.CvmClient(cred, region, clientProfile)

    req = cvm_models.DescribeImagesRequest()
    params = {
        "ImageIds": [image_id]
    }
    req.from_json_string(json.dumps(params))

    resp = client.DescribeImages(req)
    return json.loads(resp.to_json_string())
