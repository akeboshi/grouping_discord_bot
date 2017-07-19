## 概要
トリガーワードを発言したユーザが参加しているボイスチャンネルのユーザを2チームに分けるDiscordのbot

## 使い方
config.ymlにdiscord botのTokenとトリガーワードを設定して、実行する


```sh
go run grouping

```

Discordで実行したいサーバに対象のbotを参加させ、任意のテキストチャンネルでトリガーに設定したワードを発言すると、チーム分けされたテキストが発言される。

例(user1が参加しているボイスチャンネルにuser2-4がいる場合):
```
[user1] shuffle
[bot] red_team: user1, user3
[bot] blue_team: user2, user4
```
