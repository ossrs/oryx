#coding: utf-8
import dotenv, os, tools, argparse

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

region = "ap-beijing"
image_id = args.image
print(f"Remove CVM image id={image_id}, region={region}")

shares = tools.get_image_share(region, image_id)['SharePermissionSet']
account_ids = [item['AccountId'] for item in shares]
if len(account_ids) > 0:
    tools.cancel_image_share(region, image_id, account_ids)
    print(f"Image={image_id}, Shares={len(shares)}, Accounts={account_ids}, Canceled")

tools.delete_image_and_snapshot(region, image_id)
print(f"Image {image_id} deleted")
