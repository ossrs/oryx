# coding: utf-8
import os

def version():
    return "1.0.0"

def setup_site(confData):
    # Include the nginx.http.conf for http(global) level.
    if os.path.exists('/data/config/nginx.http.conf') and confData.find('#SRS-HTTP-START') == -1:
        srsConf = [
            '#SRS-HTTP-START\n',
            'include /data/config/nginx.http.conf;\n',
            '#SRS-HTTP-END\n',
        ]
        confData = f"{''.join(srsConf)}\n{confData}"

    # Include the nginx.server.conf for server(vhost) level.
    if os.path.exists('/data/config/nginx.server.conf') and confData.find('#SRS-SERVER-START') == -1:
        srsConf = [
            '#SRS-SERVER-START\n',
            'include /data/config/nginx.server.conf;\n',
            '#SRS-SERVER-END\n',
        ]
        confData = confData.replace('#SSL-START', f"{'    '.join(srsConf)}\n    #SSL-START")

    # Proxy all to SRS Stack.
    if confData.find('#SRS-PROXY-START') == -1:
        srsConf = [
            '#SRS-PROXY-START\n',
            'location / {\n',
            '    proxy_pass http://127.0.0.1:2022;\n',
            '    proxy_set_header Host $host;\n',
            '}\n',
            '#SRS-PROXY-END\n',
        ]
        confData = confData.replace('#SSL-START', f"{'    '.join(srsConf)}\n    #SSL-START")

    # Disable the location section of nginx, we will handle it.
    if confData.find('location ~ /disabled.by.srs/.*\.(js|css)?$\n') == -1:
        confData = confData.replace(
            'location ~ .*\.(js|css)?$\n',
            'location ~ /disabled.by.srs/.*\.(js|css)?$\n',
        )
    if confData.find('#location ~ /disabled.by.srs/.*\.(gif|jpg|jpeg|png|bmp|swf)$\n') == -1:
        confData = confData.replace(
            'location ~ .*\.(gif|jpg|jpeg|png|bmp|swf)$\n',
            'location ~ /disabled.by.srs/.*\.(gif|jpg|jpeg|png|bmp|swf)$\n',
        );

    return confData
