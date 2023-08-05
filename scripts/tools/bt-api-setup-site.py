#coding: utf-8
import os, sys, srsTools

confPath = '/www/server/panel/vhost/nginx/srs.cloud.local.conf'
print(f"Nginx confPath: {confPath}")

with open(confPath, 'r') as f:
    confData = f.read()
print(f"Nginx original confData: {confData}")

confData = srsTools.setup_site(confData)
print(f"Nginx new confData: {confData}")

with open(confPath, 'w') as f:
    f.write(confData)
print(f"Nginx save config success.")

r0 = os.system('/etc/init.d/nginx reload')
print(f"Nginx reload result: {r0}")
if r0 != 0:
    sys.exit(r0)

print("OK")
