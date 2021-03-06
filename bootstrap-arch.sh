#!/bin/bash
mkdir /srv/arch
pacstrap -c /srv/arch base nvidia nvidia-settings opencl-nvidia ocl-icd libcurl-compat xorg-server xorg-xinit xorg-twm xterm mkinitcpio-nfs-utils nfs-utils openssh avahi curl libmicrohttpd hwloc binutils
./bootstrap-common.sh
echo finished bootstraping /srv/arch
echo please execute ./setip.sh

#sudo apt install -y nfs-kernel-server dnsmasq lighttpd squashfs-tools
#echo "/srv/arch *(rw,no_root_squash,no_subtree_check)" | sudo tee -a /etc/exports

#sudo systemctl restart lighttpd
#sudo systemctl restart nfs-server