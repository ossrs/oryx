# NGINX-HLS-CDN

Construct a small CDN for HLS streaming using NGINX and SRS Stack.

> Note: If you want to use SRS, please refer to [Nginx for HLS](https://ossrs.io/lts/en-us/docs/v5/doc/nginx-for-hls) for more information.

## Step 1: Create an SRS Stack

Create a Digital Ocean droplet, or use another VPS:

```bash
doctl compute droplet create srs-stack --image ubuntu-20-04-x64 \
    --region sgp1 --size s-2vcpu-2gb --wait
```

> Note: You can also access [DO Droplets](https://cloud.digitalocean.com/droplets) to create a new droplet.

You will receive an internet IP address, for example, `47.92.161.70`. Log in and run:

```bash
apt-get update -y && apt-get install -y docker.io curl nginx net-tools ffmpeg pcp &&
curl -L https://github.com/ossrs/srs-stack/releases/latest/download/linux-srs_stack-en.tar.gz |tar -xz &&
bash srs_stack/scripts/setup-ubuntu/install.sh
```

After installation, you can access SRS Stack at [http://47.92.161.70:2022](http://47.92.161.70:2022).
Please configure the SRS Stack using the web console.

Open `System > HLS > Delivery HLS in High Performance mode` and click the `Submit` button. This will enable
high-performance HLS and allow NGINX to cache the m3u8 and ts files.

Next, set up NGINX to proxy port 80 to the SRS Stack. Edit `/etc/nginx/sites-enabled/default` and modify
the `location /` to proxy to the SRS Stack:

```nginx
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    location / {
        proxy_pass http://127.0.0.1:2022;
        proxy_set_header Host $host;
    }
}
```

Reload NGINX with:

```bash
/etc/init.d/nginx reload
```

Now, you can access the SRS Stack at [http://47.92.161.70](http://47.92.161.70) via the default HTTP port.
Follow the instructions from [Scenarios / Streaming / OBS or vMix](http://47.92.161.70/mgmt/en/routers-scenario?tab=live)
to publish a live stream using OBS, and play the stream via HLS.

A HLS stream should be available at [http://47.92.161.70/live/livestream.m3u8](http://47.92.161.70/live/livestream.m3u8),
and in the following steps, we will use NGINX to deliver it.

## Step 2: Create an NGINX server

Create a Digital Ocean droplet, or use another VPS:

```bash
doctl compute droplet create srs-stack-nginx01 --image ubuntu-20-04-x64 \
    --region sgp1 --size s-4vcpu-8gb --wait
```

> Note: You can also access [DO Droplets](https://cloud.digitalocean.com/droplets) to create a new droplet.

> Note: We create a VPS with `4CPU 8GB` and `160GB` disk, because the load is larger than SRS Stack.

You will receive an internet IP address, for example, `39.100.82.53`. Log in and run:

```bash
apt-get update -y && apt-get install -y docker.io curl nginx net-tools ffmpeg pcp
```

Create a directory for NGINX to cache the HLS stream:

```bash
mkdir -p /data/nginx-cache
```

Next, set up NGINX to proxy port 80 to the SRS Stack. Edit `/etc/nginx/sites-enabled/default` and modify
the `location /` to proxy to the SRS Stack:

```nginx
proxy_cache_path  /data/nginx-cache levels=1:2 keys_zone=srs_cache:8m max_size=1000m inactive=600m;
proxy_temp_path /data/nginx-cache/tmp;

server {
    listen 80 default_server;
    listen [::]:80 default_server;
    
    proxy_cache_valid  404      10s;
    proxy_cache_lock on;
    proxy_cache_lock_age 300s;
    proxy_cache_lock_timeout 300s;
    proxy_cache_min_uses 1;

    location ~ /.+/.*\.(m3u8)$ {
        proxy_pass http://47.92.161.70$request_uri;
        
        proxy_cache srs_cache;
        proxy_cache_key $scheme$proxy_host$uri$args;
        proxy_cache_valid  200 302  10s;
    }
    location ~ /.+/.*\.(ts)$ {
        proxy_pass http://47.92.161.70$request_uri;
        
        proxy_cache srs_cache;
        proxy_cache_key $scheme$proxy_host$uri;
        proxy_cache_valid  200 302  60m;
    }
}
```

> Note: Please replace the IP with yours.

Reload NGINX with:

```bash
/etc/init.d/nginx reload
```

A HLS stream should be available at [http://39.100.82.53/live/livestream.m3u8](http://39.100.82.53/live/livestream.m3u8),
and in the following steps, you can create more NGINX servers to deliver HLS stream.

## Step 3: Test the NGINX server

You can use srs-bench to test the NGINX server, for example, to simulate 500 clients to play HLS stream:

```bash
docker run --rm -d ossrs/srs:sb ./objs/sb_hls_load -c 500 -r http://39.100.82.53/live/livestream.m3u8
```

Check the bandwidth by dstat:

```bash
dstat -Nlo,eth0,eth1
```

You will find all the load is on the nginx01 server, and the SRS Stack server should be idle.

