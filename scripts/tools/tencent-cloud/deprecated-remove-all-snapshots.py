#coding: utf-8
import dotenv, os, time, tools

if os.path.exists(f'{os.getenv("HOME")}/.lighthouse/.env'):
    dotenv.load_dotenv(dotenv.find_dotenv(filename=f'{os.getenv("HOME")}/.lighthouse/.env'))
else:
    dotenv.load_dotenv(dotenv.find_dotenv())
if os.getenv("SECRET_ID") == None:
    print("Please set SECRET_ID in .env or ~/.lighthouse/.env file")
    exit(1)
if os.getenv("SECRET_KEY") == None:
    print("Please set SECRET_KEY in .env or ~/.lighthouse/.env file")
    exit(1)

regions = tools.get_regions()
regions = [region['Region'] for region in regions['RegionSet']]

print(f"Warning!!! This script will delete all images and snapshots in {len(regions)} regions")
print(f"Warning!!! This script will delete all images and snapshots in {len(regions)} regions")
print(f"Warning!!! This script will delete all images and snapshots in {len(regions)} regions")

print("You have 5 seconds to cancel this script")
for i in range(5, 0, -1):
    print(i)
    time.sleep(1)
print(f"Run with regions: {regions}, SECRET={len(os.getenv('SECRET_KEY'))}B")

for region in regions:
    snapshorts = tools.get_snapshorts(region)
    for snapshort in snapshorts["SnapshotSet"]:
        snapshot_id = snapshort['SnapshotId']
        images = snapshort['Images']
        print(f"Snapshot {snapshot_id}, Images={len(images)}, DiskId={snapshort['DiskId']}")

        if len(images) == 0:
            tools.delete_snapshot_and_image(region, snapshot_id)
            print(f"Snapshot {snapshot_id}, Images={len(images)}, Deleted")
            time.sleep(1)
            continue

        for image_index, image in enumerate(images):
            image_id = image['ImageId']
            print(f"Snapshot {snapshot_id}, Images=#{image_index}, Image={image_id}, ImageName={image['ImageName']}")
            imageShare = tools.get_image_share(region, image_id)
            print(f"Snapshot {snapshot_id}, Images=#{image_index}, Image={image_id}, Shares={len(imageShare['SharePermissionSet'])}")
            account_ids = [item['AccountId'] for item in imageShare['SharePermissionSet']]
            if len(account_ids) > 0:
                tools.cancel_image_share(region, image_id, account_ids)
                print(f"Snapshot {snapshot_id}, Images=#{image_index}, Image={image_id}, Shares={len(imageShare['SharePermissionSet'])}, Canceled")
            tools.delete_image_and_snapshot(region, image_id)
            print(f"Snapshot {snapshot_id}, Images=#{image_index}, Image={image_id}, Deleted")
            time.sleep(1)


