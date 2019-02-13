**このツールはもうメンテナンスされていません。内容が古くなっている可能性が高いです。**
# これは何？
マイニングツールです。

ただし、普通のWindows向けのマイニングツールとは違います。

Windowsが入っていない状態、HDDやSSDを接続していない状態のリグを、
LANケーブルだけ接続してネットワークにつなぎ、BIOS設定し、立ち上げるだけでマイニングできる状態になります。
# 使い方
1. [最新のRelease](https://github.com/akimasa/ArchDisklessMiner/releases/latest)をダウンロードして展開する。
2. tftp.exeを実行する。
3. リグでPXEブートする。
4. [管理ツール](https://github.com/akimasa/DisklessMinerAdmin/releases)を利用して設定する。
## PXEブート
ネットワークから起動する仕組みです。

機種によってはPXEではなく、IBAと書いてあったりします。Network Bootかもしれません。

BIOSでNetwork Bootを有効にし、**保存して再起動**しないと出てこない場合があります。
# 起動しない
Skypeがインストールされていると、TCP Port 80が既に使われていてリグが起動ができないことがあります。
その場合は、Skypeを終了するか、[Skype がポート 80 と 443 を使用しないようにする](http://blog.nnasaki.com/entry/2015/11/20/151532)
を参考にTCP Port 80を使われないようにしてください。

起動時にファイアーウォールの許可を行わなかった場合もリグが正常に起動できない原因のことがあります。
tftp.exeの起動時に許可すればよいので通常は必要ありませんが、間違って拒否した場合などは
[Windowsファイアウォールの例外にアプリケーションを追加する方法](http://faq.buffalo.jp/app/answers/detail/a_id/792)
を参考に許可を行うと使えるかもしれません。
ダメだった場合はIssueをたてるか、Twitterの[@akimasa2000](https://twitter.com/akimasa2000)にでも聞いてください。
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
