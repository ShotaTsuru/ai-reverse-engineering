# GitHub MCP サーバー設定ガイド

## 概要

このプロジェクトでは、Model Context Protocol (MCP) を使用してGitHubをより柔軟に操作できます。従来のGitHub CLI (`gh`)の代わりに、自然言語でGitHubの様々な操作を実行できます。

## 前提条件

- GitHub Copilot へのアクセス
- Visual Studio Code（または対応するMCPエディタ）
- GitHub Personal Access Token

## 設定手順

### 1. MCP設定ファイルの作成

```bash
# 自動的に設定ファイルを作成
make mcp-setup

# または手動でサンプルファイルからコピー
cp .vscode/mcp.json.sample .vscode/mcp.json
```

### 2. GitHub Personal Access Token の作成

1. GitHubの「Settings」→「Developer settings」→「Personal access tokens」→「Tokens (classic)」
2. 「Generate new token」をクリック
3. 必要なスコープを選択：
   - `repo` - リポジトリアクセス
   - `read:packages` - パッケージ読み取り
   - `workflow` - GitHub Actions（必要に応じて）
   - `read:org` - 組織情報（必要に応じて）

### 3. Visual Studio Code での設定

1. プロジェクトをVisual Studio Codeで開く
2. `.vscode/mcp.json` ファイルが自動的に認識される
3. 「Start」ボタンをクリックしてMCPサーバーを開始
4. プロンプトが表示されたらPersonal Access Tokenを入力

### 4. 使用方法

1. Copilot Chat を開く（アイコンをクリック）
2. 「Agent」モードを選択
3. ツールアイコンから「MCP Server: GitHub」を確認
4. 自然言語でGitHub操作を指示

#### 使用例

```
「〇〇の機能に関するIssueを作成して」
「このリポジトリのPull Requestを一覧表示して」
「リポジトリの設定を確認して」
「最新のコミットを確認して」
```

### 5. 設定管理コマンド

```bash
# 設定の確認
make mcp-check

# 設定の初期化（サンプルファイルからコピー）
make mcp-setup

# 設定のリセット（既存の設定を削除してサンプルから再作成）
make mcp-reset

# 設定の削除
make mcp-clean
```

## 利用可能な機能

- **Issue管理**
  - Issue の作成・更新・削除
  - Issue の検索・フィルタリング
  - ラベル管理

- **Pull Request管理**
  - Pull Request の作成・更新
  - レビュー・マージ
  - コメント管理

- **リポジトリ操作**
  - リポジトリ情報の取得
  - ブランチ管理
  - コミット履歴の確認

- **コラボレーション**
  - コラボレーター管理
  - 通知管理
  - プロジェクト管理

## 従来のスクリプトとの比較

### 従来の方法 (`scripts/create-issue.sh`)

```bash
./scripts/create-issue.sh "タイトル" "内容" "ラベル" "担当者"
```

### MCP を使用した方法

```
「バックエンドのAPIエンドポイントに関するIssueを作成して。タイトルは『API改善』、内容は『パフォーマンス向上が必要』、ラベルは『enhancement』を付けて」
```

## トラブルシューティング

### 設定ファイルが見つからない

```bash
# 設定ファイルをサンプルから作成
make mcp-setup

# 設定ファイルをリセット
make mcp-reset
```

### 認証エラー

- Personal Access Token の有効期限を確認
- 必要なスコープが含まれているか確認
- トークンの権限が適切か確認

### MCP サーバーが起動しない

- Visual Studio Code を再起動
- `.vscode/mcp.json` の構文を確認
- GitHub Copilot 拡張機能が有効化されているか確認
- 設定ファイルをリセット: `make mcp-reset`

### ツールが表示されない

- Agent モードが選択されているか確認
- MCP サーバーが正常に起動しているか確認
- ネットワーク接続を確認
- 設定の確認: `make mcp-check`

## セキュリティ上の注意

- Personal Access Token は安全に保管し、他人と共有しない
- 最小限必要なスコープのみを付与する
- 定期的にトークンを更新する
- 使用しなくなったトークンは削除する

## 追加リソース

- [GitHub MCP サーバー公式ドキュメント](https://docs.github.com/en/copilot/customizing-copilot/using-model-context-protocol/using-the-github-mcp-server)
- [Model Context Protocol 公式サイト](https://modelcontextprotocol.io/)
- [GitHub Copilot ドキュメント](https://docs.github.com/en/copilot)
