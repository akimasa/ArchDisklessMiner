#!/bin/bash
if [ -z "$IP" ]; then IP=$(hostname -i | cut -f1 -d' '); fi
if [ -z "$HTTPROOT" ]; then HTTPROOT="/srv/http"; fi
if [ -z "$CHROOT" ]; then CHROOT="arch-chroot"; fi
sed -e "s/MYIPADDR/${IP}/g" default.ipxe > ${HTTPROOT}/default.ipxe
sed -e "s/MYIPADDR/${IP}/g" dnsmasq.conf > /etc/dnsmasq.conf
sed -e "s/MYIPADDR/${IP}/g" net_nfs4 > /srv/arch/usr/lib/initcpio/hooks/net_nfs4
if [ "$NFSBOOT" == "1" ]; then cp nfsboot /srv/arch/usr/lib/initcpio/hooks/net_nfs4 ; fi
cp install-net_nfs4 /srv/arch/usr/lib/initcpio/install/net_nfs4
cp mkinitcpio.conf /srv/arch/etc/
# mkdir and copy ipxe to tftp root
mkdir /srv/tftp/
cp ipxe-default.pxe /srv/tftp/
chown -R dnsmasq:dnsmasq /srv/tftp/
# generate initramfs
if [ ! "$SKIPINITCPIO" == "1" ]; then eval $CHROOT /srv/arch/ mkinitcpio -p linux; fi
cp /srv/arch/boot/initramfs-linux-fallback.img ${HTTPROOT}/initrd
cp /srv/arch/boot/vmlinuz-linux ${HTTPROOT}/linux
systemctl restart dnsmasq
