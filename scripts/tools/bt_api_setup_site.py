#coding: utf-8
import os, sys, bt_tools

confPath = '/www/server/panel/vhost/nginx/srs.stack.local.conf'
print(f"Nginx confPath: {confPath}")

with open(confPath, 'r') as f:
    confData = f.read()
print(f"Nginx original confData: {confData}")

confData = bt_tools.setup_site(confData)
print(f"Nginx new confData: {confData}")

with open(confPath, 'w') as f:
    f.write(confData)
print(f"Nginx save config success.")

r0 = os.system('/etc/init.d/nginx reload')
print(f"Nginx reload result: {r0}")
if r0 != 0:
    sys.exit(r0)

print("OK")
