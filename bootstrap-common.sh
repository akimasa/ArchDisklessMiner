#!/bin/bash
mkdir -p /srv/arch/root/.ssh/
chmod 700 /srv/arch/root/.ssh/
cat sshkey.pub >> /srv/arch/root/.ssh/authorized_keys
chmod 600 /srv/arch/root/.ssh/authorized_keys
chroot /srv/arch/ systemctl enable sshd
chroot /srv/arch/ systemctl enable avahi-daemon
chroot /srv/arch/ systemctl enable systemd-resolved
chroot /srv/arch/ ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf
cp /srv/arch/usr/share/doc/avahi/ssh.service /srv/arch/etc/avahi/services/
echo "mininger" > /srv/arch/etc/mininger