#!/usr/bin/python
# coding: utf-8

import sys, os, json

os.chdir("/www/server/panel")
sys.path.append("class/")
import public
from panelSite import panelSite
from files import files
from firewalls import firewalls

class srs_cloud_main:
    __plugin_path = "{}/panel/plugin/srs_cloud".format(public.get_setup_path())
    __srs_service = "/usr/lib/systemd/system/srs-cloud.service"
    __deploy = '/usr/local/lighthouse/softwares'
    __srs_home = '{}/srs-cloud'.format(__deploy)
    __r0_file = '/tmp/srs_cloud_install.r0'
    __firewall = '/tmp/srs_cloud_install.fw'
    __log_file = '/tmp/srs_cloud_install.log'
    __site = 'srs.cloud.local'

    def __init__(self):
        pass

    def serviceStatus(self, args):
        status = {}
        status['r0'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        status['srs'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        status['nginx'] = public.ExecShell('ls {}/nginx/sbin/nginx >/dev/null 2>&1 && echo -n ok'.format(public.get_setup_path()))[0]
        status['docker'] = public.ExecShell('ls /usr/lib/systemd/system/docker.service >/dev/null 2>&1 && echo -n ok')[0]
        status['docker_running'] = public.ExecShell('systemctl status docker.service >/dev/null 2>&1 && echo -n ok')[0]
        status['node_manager'] = public.ExecShell('ls {}/panel/plugin/pm2 >/dev/null 2>&1 && echo -n ok'.format(public.get_setup_path()))[0]
        status['node'] = self.__node()[0]
        status['firewall'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__firewall))[0]
        status['site'] = public.ExecShell('ls {root}/{site} >/dev/null 2>&1 && echo -n ok'.format(
            root=public.get_site_path(), site=self.__site,
        ))[0]

        # We use the default site.
        site = panelSite().GetDefaultSite(None)
        status['default_site_available'] = site['defaultSite']
        if site['defaultSite'] == '' or site['defaultSite'] == False or site['defaultSite'] == self.__site:
            status['default_site_available'] = 'ok'

        # Whether site is setup ok.
        status['site_setup'] = public.ExecShell('grep -q \'{pattern}\' {www}/{site}.conf >/dev/null 2>&1 && echo -n ok'.format(
            pattern='{}/mgmt/containers/conf/default.d'.format(self.__srs_home),
            www='{}/nginx'.format(public.get_vhost_path()), site=self.__site,
        ))[0]

        return public.returnMsg(True, json.dumps(status))

    def installTasks(self, args):
        echoMsg = files().GetTaskSpeed(Params())
        if 'status' in echoMsg:
            return public.returnMsg(True, json.dumps([]))

        tasks = echoMsg['task']
        return public.returnMsg(True, json.dumps(tasks))

    def restartService(self, args):
        if args.service != 'docker':
            return public.returnMsg(False, 'invalid service {}'.format(args.service))
        public.ExecShell('systemctl restart {}'.format(args.service))
        return public.returnMsg(True, json.dumps('ok'))

    def installSrs(self, args):
        if 'start' not in args: args.start = 0
        if 'end' not in args: args.end = 0

        nodejs = self.__node()[1]
        if nodejs is None:
            return public.returnMsg(False, 'no nodejs')

        nginx = '{}/nginx/logs/nginx.pid'.format(public.get_setup_path())
        if not os.path.exists(nginx):
            return public.returnMsg(False, 'no nginx')
        if public.GetWebServer() != 'nginx':
            return public.returnMsg(False, 'not nginx, but {}'.format(public.GetWebServer()))

        srs = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        running = public.ExecShell('ps aux |grep -v grep |grep srs_cloud |grep setup >/dev/null 2>&1 && echo -n ok')[0]

        [tail, wc] = ['', 0]
        r0 = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        if running != 'ok' and srs != 'ok' and r0 != 'failed':
            public.ExecShell('nohup bash {plugin}/setup.sh {r0} {node} {nginx} {www} {site} 1>{log} 2>&1 &'.format(
                plugin=self.__plugin_path, r0=self.__r0_file, node=nodejs, nginx=nginx, www=public.get_site_path(),
                site=self.__site, log=self.__log_file,
            ))
        elif os.path.exists('/tmp/srs_cloud_install.log'):
            tail = public.ExecShell('sed -n "{start},{end}p" {log}'.format(
                start=args.start, end=args.end, log=self.__log_file,
            ))[0]
            wc = public.ExecShell('wc -l {}'.format(self.__log_file))[0]

        return public.returnMsg(True, json.dumps({'srs': srs, 'running': running, 'r0': r0, 'wc': wc, 'tail': tail}))

    def querySrs(self, args):
        srs = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        running = public.ExecShell('ps aux |grep -v grep |grep srs_cloud |grep setup >/dev/null 2>&1 && echo -n ok')[0]

        r0 = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        tail = public.ExecShell('tail {}'.format(self.__log_file))[0]

        return public.returnMsg(True, json.dumps({'srs': srs, 'running': running, 'r0': r0, 'tail': tail}))

    def cleanupIntall(self, args):
        public.ExecShell('rm -f {}'.format(self.__r0_file))
        return public.returnMsg(True, json.dumps('ok'))

    def createSrsSite(self, args):
        site = panelSite().AddSite(Params(
            webname = json.dumps({"domain": self.__site, "domainlist": [], "count": 0}),
            type = 'PHP',
            port = '80',
            ps = self.__site,
            path = os.path.join(public.get_site_path(), self.__site),
            type_id = 0,
            version = '00',
            ftp = 'false',
            sql = 'false',
            codeing = 'utf8',
        ))
        if 'status' in site:
            return site

        return public.returnMsg(True, json.dumps(site))

    def setupSrsSite(self, args):
        # Setup the default site to SRS.
        site = panelSite().GetDefaultSite(None)
        if site['defaultSite'] != '' and site['defaultSite'] != False and site['defaultSite'] != self.__site:
            return public.returnMsg(False, 'default site is {}'.format(site['defaultSite']))

        r0 = panelSite().SetDefaultSite(Params(name=self.__site))
        if 'status' in r0 and r0['status'] != True:
            return r0

        # Setup the nginx config.
        confPath = '{www}/{site}.conf'.format(
            www='{}/nginx'.format(public.get_vhost_path()), site=self.__site,
        )
        conf = files().GetFileBody(Params(path=confPath))
        confData = conf['data'];

        # Process the http section of nginx.
        if confData.find('#SRS-HTTP-START') == -1:
            srsHttpConf = 'include {}/mgmt/containers/conf/conf.d/nginx.http.conf;'.format(self.__srs_home)
            confData = '#SRS-HTTP-START\n{}\n#SRS-HTTP-END\n\n'.format(srsHttpConf) + confData

        # Process the server section of nginx.
        if confData.find('#SRS-SERVER-START') == -1:
            srsServerConf = 'include {}/mgmt/containers/conf/default.d/*.conf;'.format(self.__srs_home)
            confData = confData.replace(
                '#SSL-START',
                '#SRS-SERVER-START\n    {}\n    #SRS-SERVER-END\n\n    #SSL-START'.format(srsServerConf),
            )

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

        # Save the nginx config and reload it.
        r0 = files().SaveFileBody(Params(
            path=confPath, data=confData, encoding='utf-8',
        ))

        # We must rewrite the message if ok, for json not support chinese.
        if 'status' in r0 and r0['status'] == True:
            r0['msg'] = json.dumps({'path': confPath, 'conf': confData})
        return r0

    def setupFirewall(self, args):
        rtmp = firewalls().AddAcceptPortAll('1935', None)
        webrtc = firewalls().AddAcceptPortAll('8000', None)
        srt = firewalls().AddAcceptPortAll('10080', None)
        sip = firewalls().AddAcceptPortAll('5060', None)
        gb = firewalls().AddAcceptPortAll('9000', None)
        mgmt = firewalls().AddAcceptPortAll('2022', None)
        if rtmp is True and webrtc is True and srt is True and sip is True and gb is True and mgmt is True:
            public.ExecShell('touch {}'.format(self.__firewall))
        return public.returnMsg(True, json.dumps({
            'rtmp': rtmp, 'webrtc': webrtc, 'srt': srt, 'sip': sip, 'gb': gb,
        }))

    def querySrsService(self, args):
        ok = public.ExecShell('systemctl status srs-cloud.service >/dev/null 2>&1 && echo -n ok')[0]
        return public.returnMsg(True, json.dumps({'active': ok}))

    def __node(self):
        # Try to detect node by pm2 manager, if not ok, detect by nodejs manager.
        [ok, node] = self.__discover_path('{}/nvm/versions/node/v16*/bin/node'.format(public.get_setup_path()))
        if ok != 'ok':
            [ok, node] = self.__discover_path('{}/nodejs/v16*/bin/node'.format(public.get_setup_path()))
        if ok != 'ok':
            ok = public.ExecShell('which node >/dev/null 2>&1 && echo -n ok')[0]
            node = public.ExecShell('which node 2>/dev/null')[0]
        return [ok, node]

    def __discover_path(self, general_path):
        real_path = None
        ok = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(general_path))[0]
        if ok == 'ok':
            real_path = public.ExecShell('ls {}'.format(general_path))[0]
        if real_path is not None:
            real_path = real_path.strip()
        return [ok, real_path]

class Params(object):
    def __init__(self, **kwargs):
        self.__dict__.update(kwargs)
    def __iter__(self):
        return self.__dict__.__iter__();
    def __getitem__(self, y):
        return self.__dict__.__getitem__(y)

