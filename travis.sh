#!/bin/bash
pacman -Syu --noconfirm
pacman --noconfirm -S squashfs-tools
mkdir /run/shm
cd /root
bash bootstrap-arch.sh
chroot /srv/arch/ pacman --noconfirm -R lvm2 man-db man-pages nano reiserfsprogs mdadm
bash setip.sh
#mksquashfs /srv/arch/ /srv/http/arch.sfs -comp xz -Xbcj x86 -b 1M -e /srv/arch/boot/ -e /srv/arch/usr/share/man/
mksquashfs /srv/arch/ /srv/http/arch.sfs -comp lzo -e /srv/arch/boot/ -e /srv/arch/usr/share/man/