# これは何？
マイニングツールです。

ただし、普通のWindows向けのマイニングツールとは違います。

Windowsが入っていない状態、HDDやSSDを接続していない状態のリグを、
LANケーブルだけ接続してネットワークにつなぎ、BIOS設定し、立ち上げるだけでマイニングできる状態になります。
# 使い方
1. [最新のRelease](https://github.com/akimasa/ArchDisklessMiner/releases/latest)をダウンロードして展開する。
2. tftp.exeを実行する。
3. リグでPXEブートする。
## PXEブート
ネットワークから起動する仕組みです。

機種によってはPXEではなく、IBAと書いてあったりします。Network Bootかもしれません。

BIOSでNetwork Bootを有効にし、**保存して再起動**しないと出てこない場合があります。
# よくある質問
Q: 既に入っているWindowsは消えるの？

A: 消えません。ただし、意図的に消すコマンドを入力した場合は消えます。
このマイニングツールは、すべてメモリ上に展開されて実行されます。
# 不明点があったら…
- わからないことがあったらIssueをたてるか、Twitterの[@akimasa2000](https://twitter.com/akimasa2000)にでも聞いてください。
- あなたにわからなかったことは、ほかの人にわからない事です。
- 言われるまで作者は気が付かないでしょう。
# 最新版
- [自動ビルドされた最新のPre-release](https://github.com/akimasa/ArchDisklessMiner/releases/)のarch.sfs, linux, initrdをダウンロードして置き換えることで最新版のArch Linuxを利用可能です。
- 自動ビルドは毎日実行されますが、何らかの原因でうまくアップロードされない場合があります。
- その場合は諦めて少し古いのを使ってください。