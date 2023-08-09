#coding: utf-8
import dotenv, os, tools, argparse, time

parser = argparse.ArgumentParser(description="TencentCloud")
parser.add_argument("--image", type=str, required=False, help="The CVM image id")

args = parser.parse_args()

if os.path.exists(f'{os.getenv("HOME")}/.lighthouse/.env'):
    dotenv.load_dotenv(dotenv.find_dotenv(filename=f'{os.getenv("HOME")}/.lighthouse/.env'))
else:
    dotenv.load_dotenv(dotenv.find_dotenv())

if os.getenv("VM_IMAGE") is not None and args.image is None:
    args.image = os.getenv("VM_IMAGE")
if args.image == None:
    raise Exception("Please set --image")

if os.getenv("SECRET_ID") == None:
    print("Please set SECRET_ID in .env or ~/.lighthouse/.env file")
    exit(1)
if os.getenv("SECRET_KEY") == None:
    print("Please set SECRET_KEY in .env or ~/.lighthouse/.env file")
    exit(1)
if os.getenv("LH_ACCOUNT") == None:
    print("Please set LH_ACCOUNT in .env or ~/.lighthouse/.env file")
    exit(1)

region = "ap-beijing"
image_id = args.image
account_id = os.getenv("LH_ACCOUNT")
print(f"Share image id={image_id}, region={region} to account={account_id}")

while True:
    info = tools.query_image(region, image_id)['ImageSet']
    if len(info) != 1:
        raise Exception(f"Image {image_id} not found")

    if info[0]['ImageState'] == 'NORMAL':
        break

    print(f"Image {image_id} state is {info[0]['ImageState']}, wait 5 seconds")
    time.sleep(5)

tools.share_image(region, image_id, account_id)
print(f"Image {image_id} shared to account {account_id}")
