#!/bin/bash
pacman -Syu --noconfirm
pacman --noconfirm -S squashfs-tools
mkdir /run/shm
cd /root
echo travis_fold:start:bootstrap
bash bootstrap-arch.sh
echo travis_fold:end:bootstrap
echo travis_fold:start:mkinitcpio
bash setip.sh
echo travis_fold:end:mkinitcpio
echo travis_fold:start:mksquashfs
mksquashfs /srv/arch/ /srv/http/arch.sfs -comp xz -Xbcj x86 -b 1M -e /srv/arch/boot/
echo travis_fold:end:mksquashfs