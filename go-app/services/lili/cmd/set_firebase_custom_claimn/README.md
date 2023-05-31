# firebase カスタムクレイム登録スクリプト

## 背景

- firebase emulator を使っている場合 firebase カスタムクレイムは自由に登録・編集できる
  - emulator UI 上から、もしくはシーダーとして登録・編集できる
- 一方 Cloud 上の firebase auth に接続している場合は、カスタムクレイムの登録・編集が簡単にできない
  - firebase admin sdk 経由でしか登録・編集できない

## スクリプト説明

firebase admin sdk を使い、firebase カスタムクレイムを登録・編集するためのスクリプト

### 使い方

```sh
$ go run . -u="{{Firebase UID}}" -s="{{firebase service account file path}}"
Input custom claim:
{
  "is_valid": true
}
```

s オプションのデフォルト値は`../../service-account-file.json`です
