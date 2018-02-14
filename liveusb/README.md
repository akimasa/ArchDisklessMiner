重要
====
変更があったので、今はそのままでは使えません！

Live USB作成向けのファイル
========================
1. default.ipxe, dnsmasq.conf, setip.shを/root/に設置
2. setip.serviceを/etc/systemd/system/に設置
3. dnsmasq, nfs-server, setipを有効化

既知の問題
========
- ノートPCなどEthernet以外のネットワークボード(Wi-Fiなど)を積んでいるPCだと、systemd-networkd-wait-online.targetが2分経ってタイムアウトするまで待たないとsetip.serviceが実行されない