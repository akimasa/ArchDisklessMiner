#!/bin/bash
if [ -z "$CHROOT" ]; then CHROOT="arch-chroot"; fi
mkdir -p /srv/arch/root/.ssh/
chmod 700 /srv/arch/root/.ssh/
cat sshkey.pub >> /srv/arch/root/.ssh/authorized_keys
chmod 600 /srv/arch/root/.ssh/authorized_keys
$chroot /srv/arch/ systemctl enable sshd
$chroot /srv/arch/ systemctl enable avahi-daemon
cp /srv/arch/usr/share/doc/avahi/ssh.service /srv/arch/etc/avahi/services/
echo "mininger" > /srv/arch/etc/mininger