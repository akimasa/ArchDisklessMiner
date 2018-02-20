#!/bin/bash
pacman -Syu --noconfirm
pacman --noconfirm -S squashfs-tools
mkdir /run/shm
cd /root
bash bootstrap-arch.sh
bash setip.sh
mksquashfs /srv/arch/ /root/arch.sfs -comp xz -Xbcj x86 -b 1M -e /srv/arch/boot/
ls -lh /root/arch.sfs