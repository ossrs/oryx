#!/bin/sh

# For SSH.
echo "ufw allow ssh"
ufw limit ssh

# For Nginx HTTP/HTTPS, also proxy for srs-cloud and SRS streaming.
echo "ufw allow HTTP/HTTPS"
ufw allow 80
ufw allow 443

# For RTMP streaming.
echo "ufw allow RTMP/WebRTC/SRT streaming"
ufw allow 1935/tcp
# For WebRTC streaming.
ufw allow 8000/udp
# For SRT streaming.
ufw allow 10080/udp

# For GB28181 with SIP.
echo "ufw allow GB/SIP streaming"
ufw allow 9000/tcp
ufw allow 9000/udp
ufw allow 5060/tcp
ufw allow 5060/udp

# Note that we must expose these port listen at PrivateIP or lo, for container to access it.
echo "ufw allow mgmt"
# For srs-cloud API.
ufw allow 2022/tcp

# Note that we must expose the exporter for prometheus to access it.
echo "ufw allow exporter for prometheus"
# For node_exporter
ufw allow 9100/tcp

# Apply force firewall.
echo "ufw apply force"
ufw --force enable

