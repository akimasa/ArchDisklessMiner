#!/bin/bash
sudo bash ./arch-bootstrap.sh -r http://ftp.jaist.ac.jp/pub/Linux/ArchLinux/ /srv/arch
sudo chroot /srv/arch/ pacman -S --noconfirm base nvidia opencl-nvidia ocl-icd libcurl-compat xorg-server xorg-xinit xorg-twm xterm mkinitcpio-nfs-utils nfs-utils

sudo apt update
sudo apt install -y nfs-kernel-server dnsmasq lighttpd squashfs-tools
echo "/srv/arch *(rw,no_root_squash,no_subtree_check)" | sudo tee -a /etc/exports

sudo systemctl restart lighttpd
sudo systemctl restart nfs-server