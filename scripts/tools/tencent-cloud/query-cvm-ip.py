#coding: utf-8
import dotenv, os, sys, tools, argparse

parser = argparse.ArgumentParser(description="TencentCloud")
parser.add_argument("--instance", type=str, required=False, help="The CVM instance id")
parser.add_argument("--id", type=str, required=False, help="Write ID result to this file")

args = parser.parse_args()

if os.path.exists(f'{os.getenv("HOME")}/.lighthouse/.env'):
    dotenv.load_dotenv(dotenv.find_dotenv(filename=f'{os.getenv("HOME")}/.lighthouse/.env'))
else:
    dotenv.load_dotenv(dotenv.find_dotenv())

if os.getenv("VM_INSTANCE") is not None and args.instance is None:
    args.instance = os.getenv("VM_INSTANCE")
if args.instance == None:
    raise Exception("Please set --instance")

if os.getenv("SECRET_ID") == None:
    print("Please set SECRET_ID in .env or ~/.lighthouse/.env file")
    exit(1)
if os.getenv("SECRET_KEY") == None:
    print("Please set SECRET_KEY in .env or ~/.lighthouse/.env file")
    exit(1)

region = "ap-beijing"
instance_id = args.instance
print(f"Query CVM instance={instance_id}, region={region}, id={args.id}")

instance_details = tools.query_instance_detail(region, instance_id)['InstanceSet']
if len(instance_details) != 1:
    raise Exception(f"Instance {instance_id} detail not found")
instance_detail = instance_details[0]
public_ip = instance_detail['PublicIpAddresses'][0]
private_ip = instance_detail['PrivateIpAddresses'][0]
print(f"Instance {instance_id}, public ip={public_ip}, private ip={private_ip}")

if args.id != None:
    with open(args.id, 'w') as f:
        print(public_ip, file=f)
