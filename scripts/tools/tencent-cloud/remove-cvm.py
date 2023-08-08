#coding: utf-8
import dotenv, os, tools, argparse

parser = argparse.ArgumentParser(description="TencentCloud")
parser.add_argument("--instance", type=str, required=False, help="The CVM instance id")

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
print(f"Remove CVM instance={instance_id}, region={region}")

tools.delete_instance(region, instance_id)
print(f"Instance {instance_id} deleted")
