#!/usr/bin/python
# coding: utf-8

import sys, os, json

os.chdir("/www/server/panel")
sys.path.append("class/")
import public
from panelSite import panelSite
from files import files
from firewalls import firewalls

import bt_tools
print(f"bt_tools version: {bt_tools.version()}")

class oryx_main:
    # Normally the plugin is at:
    #       /www/server/panel/plugin/oryx
    # Other paths are:
    #       public.get_setup_path() is /www/server
    #       public.get_panel_path() is /www/server/panel
    #       public.get_site_path() is /www/wwwroot
    #       public.get_vhost_path() is /www/server/panel/vhost
    __plugin_path = "{}/panel/plugin/oryx".format(public.get_setup_path())
    __srs_service = "/usr/lib/systemd/system/oryx.service"
    __srs_home = '/usr/local/srs-stack'
    __r0_file = '/tmp/oryx_install.r0'
    __firewall = '/tmp/oryx_install.fw'
    __log_file = '/tmp/oryx_install.log'
    __ready_file = '{}/.bt_ready'.format(__plugin_path)
    __site = 'srs.stack.local'

    def __init__(self):
        pass

    def serviceStatus(self, args):
        status = {}
        status['tools_version'] = bt_tools.version()
        status['plugin_ready'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__ready_file))[0]
        status['srs_error'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        status['srs_ready'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        status['nginx'] = public.ExecShell('ls {}/nginx/sbin/nginx >/dev/null 2>&1 && echo -n ok'.format(public.get_setup_path()))[0]
        status['docker_manager'] = public.ExecShell('ls {}/panel/plugin/docker >/dev/null 2>&1 && echo -n ok'.format(public.get_setup_path()))[0]
        status['docker_installed'] = public.ExecShell('ls /usr/lib/systemd/system/docker.service >/dev/null 2>&1 && echo -n ok')[0]
        status['docker_running'] = public.ExecShell('systemctl status docker.service >/dev/null 2>&1 && echo -n ok')[0]
        status['firewall_ready'] = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__firewall))[0]
        status['site_created'] = public.ExecShell('ls {root}/{site} {www}/nginx/{site}.conf >/dev/null 2>&1 && echo -n ok'.format(
            root=public.get_site_path(), site=self.__site, www=public.get_vhost_path(),
        ))[0]

        # We use the default site.
        site = panelSite().get_site_info(self.__site)
        if site is not None and 'id' in site:
            domains = panelSite().GetSiteDomains(Params(id=site['id']))
            for item in domains['domains']:
                if item["name"] != self.__site:
                    status['site_domain'] = item["name"]

        # Whether site is setup ok.
        status['site_setup'] = public.ExecShell(
            f"grep -q 'SRS-PROXY-START' {public.get_vhost_path()}/nginx/{self.__site}.conf >/dev/null 2>&1 && echo -n ok"
        )[0]

        return public.returnMsg(True, json.dumps(status))

    def installTasks(self, args):
        echoMsg = files().GetTaskSpeed(Params())
        if 'status' in echoMsg:
            return public.returnMsg(True, json.dumps([]))

        tasks = echoMsg['task']
        return public.returnMsg(True, json.dumps(tasks))

    # If not set docker_installed.
    def installService(self, args):
        if args.service != 'docker':
            self.__trace(f"Error: Install invalid service {args.service}")
            return public.returnMsg(False, 'invalid service {}'.format(args.service))
        public.ExecShell('bash {}/do_docker.sh'.format(self.__plugin_path))
        return public.returnMsg(True, json.dumps('ok'))

    # If not set docker_running.
    def restartService(self, args):
        if args.service != 'docker':
            self.__trace(f"Error: Restart invalid service {args.service}")
            return public.returnMsg(False, 'invalid service {}'.format(args.service))
        public.ExecShell('systemctl restart {}'.format(args.service))
        return public.returnMsg(True, json.dumps('ok'))

    # If not srs_ready or srs_error.
    def installSrs(self, args):
        if 'start' not in args: args.start = 0
        if 'end' not in args: args.end = 0

        nginx = '{}/nginx/logs/nginx.pid'.format(public.get_setup_path())
        if not os.path.exists(nginx):
            self.__trace(f"Error: Install no nginx")
            return public.returnMsg(False, 'no nginx')
        if public.GetWebServer() != 'nginx':
            self.__trace(f"Error: Install not nginx, but {public.GetWebServer()}")
            return public.returnMsg(False, 'not nginx, but {}'.format(public.GetWebServer()))

        srs = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        running = public.ExecShell('ps aux |grep -v grep |grep oryx |grep setup >/dev/null 2>&1 && echo -n ok')[0]

        [tail, wc] = ['', 0]
        r0 = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        if running != 'ok' and srs != 'ok' and r0 != 'failed':
            public.ExecShell('nohup bash {plugin}/setup.sh --r0 "{r0}" --nginx "{nginx}" --www "{www}" --site "{site}" 1>{log} 2>&1 &'.format(
                plugin=self.__plugin_path, r0=self.__r0_file, nginx=nginx, www=public.get_site_path(),
                site=self.__site, log=self.__log_file,
            ))
        elif os.path.exists('/tmp/oryx_install.log'):
            tail = public.ExecShell('sed -n "{start},{end}p" {log}'.format(
                start=args.start, end=args.end, log=self.__log_file,
            ))[0]
            wc = public.ExecShell('wc -l {}'.format(self.__log_file))[0]

        return public.returnMsg(True, json.dumps({'srs': srs, 'running': running, 'r0': r0, 'wc': wc, 'tail': tail}))

    def querySrs(self, args):
        srs = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(self.__srs_service))[0]
        running = public.ExecShell('ps aux |grep -v grep |grep oryx |grep setup >/dev/null 2>&1 && echo -n ok')[0]

        r0 = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n failed'.format(self.__r0_file))[0]
        tail = public.ExecShell('tail {}'.format(self.__log_file))[0]

        return public.returnMsg(True, json.dumps({'srs': srs, 'running': running, 'r0': r0, 'tail': tail}))

    def cleanupIntall(self, args):
        public.ExecShell('rm -f {}'.format(self.__r0_file))
        return public.returnMsg(True, json.dumps('ok'))

    # If not set site_created.
    def createSrsSite(self, args):
        if 'domain' not in args or args.domain == '':
            self.__trace(f"Error: Empty Oryx domain.")
            return public.returnMsg(False, 'invalid domain')

        site = panelSite().AddSite(Params(
            webname = json.dumps({"domain": self.__site, "domainlist": [args.domain], "count": 0}),
            type = 'PHP',
            port = '80',
            ps = self.__site,
            path = os.path.join(public.get_site_path(), self.__site),
            type_id = 0,
            version = '00', # Static site.
            ftp = 'false',
            sql = 'false',
            codeing = 'utf8',
        ))
        if 'status' in site:
            return site

        return public.returnMsg(True, json.dumps(site))

    # If not set site_setup
    def setupSrsSite(self, args):
        # Setup the nginx config.
        confPath = f'{public.get_vhost_path()}/nginx/{self.__site}.conf'
        conf = files().GetFileBody(Params(path=confPath))
        confData = conf['data']
        confData = bt_tools.setup_site(confData)

        # Save the nginx config and reload it.
        r0 = files().SaveFileBody(Params(
            path=confPath, data=confData, encoding='utf-8',
        ))

        # We must rewrite the message if ok, for json not support chinese.
        if 'status' in r0 and r0['status'] == True:
            r0['msg'] = json.dumps({'path': confPath, 'conf': confData})

        if r0['status'] == False:
            self.__trace(f"Error: Failed to setup site, {r0['msg']}")
        return r0

    # If not set firewall_ready
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
        ok = public.ExecShell('systemctl status oryx.service >/dev/null 2>&1 && echo -n ok')[0]
        return public.returnMsg(True, json.dumps({'active': ok}))

    def __discover_path(self, general_path):
        real_path = None
        ok = public.ExecShell('ls {} >/dev/null 2>&1 && echo -n ok'.format(general_path))[0]
        if ok == 'ok':
            real_path = public.ExecShell('ls {}'.format(general_path))[0]
        if real_path is not None:
            real_path = real_path.strip()
        return [ok, real_path]

    def __trace(self, msg):
        with open(self.__log_file, 'a+') as file:
            file.write(f"{msg}\n")

class Params(object):
    def __init__(self, **kwargs):
        self.__dict__.update(kwargs)
    def __iter__(self):
        return self.__dict__.__iter__();
    def __getitem__(self, y):
        return self.__dict__.__getitem__(y)

