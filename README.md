Jimiko
======

![](https://github.com/kimikimi714/jimiko/workflows/go-ci/badge.svg)

話しかけたらいい感じに家のことを回してくれるbotです。  
現時点ではGoogle Home(+Dialogflow)とSlackに対応しています。

## 機能

- 「今何がある?」と聞いたら在庫があるもの一覧を答えてくれる。
- 「○○ある?」と聞いたら「○○あるよ」もしくは「○○ないよ」と答えて在庫があるか、ないか答えてくれる。

## 必要要件

- Go 1.11+
- Cloud Functions
- Google Sheets
- 以下のいずれか
    - Google Home + Dialogflow
    - Slack

Google HomeとDialogflowを接続する方法などは[こちらの記事](https://kimikimi714.hatenablog.com/entry/2019/12/07/183000)をご覧ください。

## 使い方

1. Google Homeに「じみこにつないで」と話しかけて、じみこを呼び出す。
2. 聞きたいことをじみこに聞く。
    - 「今何がある?」で在庫があるもの一覧を答える。
    - 「今何がない?」で在庫がないもの一覧を答える。
    - 「○○ある?」で特定の商品の在庫があるか、ないか答える。
3. 「バイバイ」でじみことの会話を終了する。

# ユニットテスト

```
make test
```

# デプロイ

あらかじめ `.env.yaml` に必須パラメータを追加してディレクトリ直下においてください。  
必須パラメータは `.env.sample.yaml` を確認してください。

```
scripts/deploy.sh {dialogflow|slack}
```

## 作者

[kimikimi714](https://kimikimi714.hatenablog.com/)

## ライセンス

MIT