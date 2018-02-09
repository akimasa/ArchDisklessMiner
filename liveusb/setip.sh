#!/bin/bash
cd /root/
if [ -z "$IP" ]; then IP=$(hostname -i | cut -f1 -d' '); fi
sed -e "s/MYIPADDR/${IP}/g" default.ipxe > /srv/http/default.ipxe
sed -e "s/MYIPADDR/${IP}/g" dnsmasq.conf > /etc/dnsmasq.conf
chown -R dnsmasq:dnsmasq /srv/tftp/
# generate initramfs
systemctl restart dnsmasq
