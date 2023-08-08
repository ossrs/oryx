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

print(f"Run with regions: {regions}, SECRET={len(os.getenv('SECRET_KEY'))}B")

index=0
for region in regions:
    snapshorts = tools.get_snapshorts(region)
    for snapshort in snapshorts["SnapshotSet"]:
        snapshot_id = snapshort['SnapshotId']
        images = snapshort['Images']
        print(f"#{index} Region={region}, Snapshot={snapshot_id}, Images={len(images)}, DiskId={snapshort['DiskId']}")
        index+=1
