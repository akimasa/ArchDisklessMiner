#!/bin/bash
pacman -Syu --noconfirm
pacman --noconfirm -S squashfs-tools
mkdir /run/shm
cd /root
echo -e "travis_fold:start:bootstrap\r"
bash bootstrap-arch.sh
echo -e "travis_fold:end:bootstrap\r"
echo -e "travis_fold:start:mkinitcpio\r"
bash setip.sh
echo -e "travis_fold:end:mkinitcpio\r"
echo -e "travis_fold:start:mksquashfs\r"
mksquashfs /srv/arch/ /srv/http/arch.sfs -comp xz -Xbcj x86 -b 1M -e /srv/arch/boot/
echo -e "travis_fold:end:mksquashfs\r"