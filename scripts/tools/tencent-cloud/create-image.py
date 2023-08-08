#coding: utf-8
import dotenv, os, tools, argparse, sys

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
image_name = "srs"
instance_id = args.instance
print(f"Create CVM instance={instance_id}, region={region}, image={image_name}")

r0 = tools.create_image(region, instance_id, image_name)
image_id = r0['ImageId']
print(f"Image {image_name} created id={image_id}")

# print the instance public ip to stderr.
print(image_id, file=sys.stderr)
