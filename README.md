# 動作環境
- Ubuntu 17.10
# 使い方
- bash bootstrap-ubuntu.sh
- sudo CHROOT="mount --bind /proc /srv/arch/proc; mount --bind /dev /srv/arch/dev; chroot" IP="10.0.2.12" HTTPROOT="/var/www/html/" NFSBOOT=1 bash setip.sh
- - CHROOTとHTTPROOTはUbuntuの場合の値にしてある。
- - *IPはサーバーのIPアドレスをip addrなどで調べて変更すること*

同じネットワークにつながっているリグを起動すれば、PXEブートが選べるはず。

初回起動後、必要なマイニングツールのインストールなどを実施

シャットダウン後、サーバーも念のため再起動する。
- cd /srv/arch/
- sudo mksquashfs ./ ../arch.sfs
- sudo mv ../arch.sfs ./

bootstrap-ubuntu.shのあるディレクトリに戻る

- sudo CHROOT="mount --bind /proc /srv/arch/proc; mount --bind /dev /srv/arch/dev; chroot" IP="10.0.2.12" HTTPROOT="/var/www/html/" bash setip.sh

先ほど実行したコマンドからNFSBOOT=1を抜いたものを実行する。
そして、PXEブートすれば、今度はsquashfsとoverlayfsで起動する。

# 使わせていただいたもの
- https://github.com/tokland/arch-bootstrap