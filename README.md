# kotob
GeminiAPIを通して任意の言語・形式に変換するCLI翻訳ツール。  
CLI translation tool that converts to any language and format through GeminiAPI.

---
## 目次 - Table of Contents
* [インストール - Installation](#インストール---installation)
* [準備 - Preparation](#準備---preparation)
* [使い方 - Usage](#使い方---usage)
* [フラグ - Flags](#フラグ---flags)
* [ライセンス - License](#ライセンス---license)
* [Special Thanks](#special-thanks)

# インストール - Installation
## Goを使用する場合 (For Go Users)
```golang
go install github.com/kotob-project/kotob@latest
```
kotob コマンドが認識されない場合は、以下のディレクトリを環境変数 PATH に追加する必要があります。  

・Windows: `%USERPROFILE%\go\bin`  
・macOS / Linux: `~/go/bin`

## バイナリをダウンロードする場合 (Direct Download)
[Releases](https://github.com/kotob-project/kotob/releases) ページから、ご自身の環境に合った最新の実行ファイルをダウンロードしてください。

1. 実行ファイルの配置:
ダウンロードしたバイナリ（Windowsの場合は `.exe`）を任意の場所に配置してください。
例：  
・Windows: `%USERPROFILE%\kotob\`  
・macOS / Linux: `~\kotob\`

>環境変数 PATH に追加されたディレクトリである必要があります。  
>PATHの通し方が分からない場合は各自で調べるようお願いします。

2. リネーム
ファイル名を kotob（Windowsなら kotob.exe）にリネームしてください。

3. 実行権限の付与 (macOS / Linux):
ダウンロードしたファイルに実行権限がない場合は、以下のコマンドを実行してください。
```bash
chmod +x "kotob"
```

# 準備 - Preparation

kotob を使用するには、Gemini APIキーが必要です。 Google AI Studio でキーを取得し、環境変数 `KOTOB_API_KEY` に設定してください。

To use kotob, you need a Gemini API key. Obtain your key from Google AI Studio and set it as the environment variable `KOTOB_API_KEY`.

モデルの選択などの詳細な設定は以下を参照してください。

For detailed settings such as model selection, please refer to the following.

・モデルについて

モデルは環境変数 `KOTOB_MODEL`に利用可能なGeminiのモデル名を指定することで変更できます。

デフォルトは `gemini-2.5-flash-lite` です。

・設定ファイルについて

kotobが参照するのは実行ディレクトリの `kotob.json` と `~/.config/kotob/kotob.json` です。

以下の形式で記述することで、フラグのデフォルト値や、設定ファイル専用の項目を指定できます。

```json
{
    "api-key": "YOUR_API_KEY",
    "to": "Japanese",
    "from": "auto",
    "model": "gemini-2.5-flash-lite",
    "system": "",
    "json": false,
    "no-stream": false,
    "explain": false,
    "explain-lang": "Japanese",
    "explain-model": "gemini-2.5-flash"
}
```

> **Note:** `explain-lang` (解説出力時の言語) および `explain-model` (解説モード有効時の専用AIモデル) は、コマンドラインフラグからは指定できない設定ファイル専用のカスタマイズ項目です。

値の適用は以下の順に優先されます。

1. コマンドライン引数 (Flags)
2. 環境変数 (Environment Variables)
3. 実行ディレクトリの設定ファイル
4. `~/.config/kotob/` の設定ファイル

# 使い方 - Usage
基本
```bash
kotob -t ja "Hello, world!"
# こんにちは、世界！
```
ファイルの読み込み
```bash
kotob -F ./test/message.txt
```
システムプロンプトの設定
```bash
kotob -t ja -s "カジュアルに翻訳" "Hello! How are you?"
# やあ！元気？
```
json出力
```bash
kotob -t ja --json "Hello! How are you?"

# {
#  "source": "auto",
#   "target": "ja",
#   "input": "Hello! How are you?",
#   "translated": "こんにちは！お元気ですか？",
#   "model": "gemini-2.5-flash-lite"
# }

```
解説付き出力
```bash
kotob -e -t ja "Hello! How are you?"

# こんにちは！お元気ですか？
# ---
# * 「Hello!」は一般的な挨拶なので「こんにちは！」と訳しました。
# * 「How are you?」は相手の安否を尋ねる表現で、日本語では「お元気ですか？」が自然です。
```

その他の機能は [フラグ - Flags](#フラグ---flags) を参照してください。

# フラグ - Flags

kotobの動作を制御するためのフラグです。

| 短縮 | フルパス | 説明 | デフォルト値 |
| --- | --- | --- | --- |
| `-k` | `--api-key` | Gemini API key | - |
| `-t` | `--to` | 翻訳先の言語 (Target language) | `Japanese` |
| `-f` | `--from` | 翻訳元の言語 (Source language) | `auto` |
| `-F` | `--file` | ファイルの読み込み | - |
| `-s` | `--system` | AIへの指示/制約 (System Prompt) | - |
| `-j` | `--json` | 出力結果を構造化データ(JSON)で取得 | `false` |
| `-e` | `--explain` | 翻訳の補足説明を追加 | `false` |
| `-m` | `--model` | 使用するAIモデルの指定 | `gemini-2.5-flash-lite` |
| `-S` | `--no-stream` | ストリーミングを無効化し、一括出力する | `false` |
| `-h` | `--help` | ヘルプを表示 | - |

> 注意: `-e/--explain` を使う場合、`-m/--model` が明示指定されていればその値を優先します。そうでなければ、設定ファイル内の `explain-model` があればそれを使用し、なければ解説モード用のデフォルトモデル `gemini-2.5-flash` を使用します。

**優先順位:** コマンド実行時に指定したフラグは、設定ファイル (`kotob.json`) や環境変数よりも優先して適用されます。

# ライセンス - License

**Apache License, Version 2.0** の下でライセンスされています。
全文については [LICENSE](./LICENSE) を参照してください。

Licensed under the **Apache License, Version 2.0**.
See [LICENSE](./LICENSE) for the full license text.

# Special Thanks

**tanahiro2010** : 仕様、リリースに関する相談、開発協力
