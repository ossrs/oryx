#coding: utf-8
import dotenv, os, tools, argparse, sys, time, datetime

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
image_name = "srs"
instance_id = args.instance
image_desc = f"{image_name} from {instance_id} at {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
print(f"Create CVM instance={instance_id}, region={region}, image={image_name}, desc={image_desc}, id={args.id}")

r0 = tools.create_image(region, instance_id, image_name, image_desc)
image_id = r0['ImageId']
print(f"Image {image_name} created id={image_id}")

while True:
    info = tools.query_image(region, image_id)['ImageSet']
    if len(info) != 1:
        raise Exception(f"Image {image_id} not found")

    if info[0]['ImageState'] == 'NORMAL':
        break

    print(f"Image {image_id} state is {info[0]['ImageState']}, wait 5 seconds")
    time.sleep(5)

if args.id != None:
    with open(args.id, 'w') as f:
        print(image_id, file=f)
