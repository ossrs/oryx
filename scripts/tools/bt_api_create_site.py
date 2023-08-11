#coding: utf-8
import os,time,hashlib,json,urllib.request
md5sum = lambda s: hashlib.md5(s.encode('utf-8')).hexdigest()

# See https://www.bt.cn/bbs/thread-20376-1-1.html
# See https://www.bt.cn/data/api-doc.pdf
BT_KEY = os.getenv('BT_KEY')
if not BT_KEY:
    print("BT_KEY is not set")
    exit(1)

domain = "bt.ossrs.net"
if os.getenv("DOMAIN") is not None:
    domain = os.getenv("DOMAIN")
print(f"Create site for domain={domain}")

now_time = int(time.time())
request_token = md5sum(str(now_time) + md5sum(BT_KEY))
req = urllib.request.Request(
    'http://localhost:7800/site?action=AddSite',
    urllib.parse.urlencode({
        'request_token': request_token,
        'request_time': now_time,
        'webname': json.dumps({"domain":"srs.stack.local", "domainlist":[domain],"count":0}),
        'path': '/www/wwwroot/srs.stack.local',
        'type_id': 0,
        'type': 'PHP',
        'version': '00',
        'port': 80,
        'ps': 'SRS Stack',
        'ftp': 'false',
        'sql': 'false',
        # For aaPanel.
        'set_ssl': 0,
    }).encode('utf-8')
)
res = urllib.request.urlopen(req)
result = res.read().decode('utf-8')
print(f"request_time={now_time}, request_token={request_token}, result={result})")

jr = json.loads(result)
if jr['siteStatus'] != True:
    print(f"failed to create site, result={result})")
    exit(1)
print(f"Create site success, id={jr['siteId']}")
