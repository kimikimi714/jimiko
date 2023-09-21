Jimiko
======

![](https://github.com/kimikimi714/jimiko/workflows/CI/badge.svg)

話しかけたらいい感じに家のことを回してくれるbotです。  
現時点ではSlackに対応しています。

## 機能

- 「今何がある?」と聞いたら在庫があるもの一覧を答えてくれる。
- 「○○ある?」と聞いたら「○○あるよ」もしくは「○○ないよ」と答えて在庫があるか、ないか答えてくれる。

## 必要要件

- Go 1.20+
- Cloud Functions
- Google Sheets
- Slack

## 使い方

今在庫があるものだけ返す。

```
@<bot-name> 何がある?
```

今在庫がないものだけ返す。

```
@<bot-name> 何がない?
```

# ユニットテスト

```
make test
```

# デプロイ

あらかじめ `.env.yaml` に必須パラメータを追加してディレクトリ直下においてください。  
必須パラメータは `.env.sample.yaml` を確認してください。

```
make deploy
```

## 作者

[kimikimi714](https://kimikimi714.hatenablog.com/)

## ライセンス

MIT
