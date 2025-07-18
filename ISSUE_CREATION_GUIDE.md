# Issue作成ガイド

このプロジェクトでは、複数の方法でissueを作成できます。

## 🤖 AIアシスタントへの指示でIssue作成

AIアシスタント（このチャット）に以下のような指示をすることで、issueを作成するためのファイルやワークフローを準備してもらえます：

### 例：
```
「バックエンドのパフォーマンス改善について、緊急度高でissueを作成してください」
```

AIアシスタントが：
1. 適切なissue内容を生成
2. 必要に応じてワークフローを実行
3. 関連するドキュメントを作成

## 🖱️ GitHub UIでの手動作成

### 方法1: GitHub Actions ワークフローを使用
1. GitHubリポジトリの「Actions」タブに移動
2. 「Create Issue Manually」ワークフローを選択
3. 「Run workflow」をクリック
4. 以下の情報を入力：
   - **タイトル**: Issueのタイトル
   - **内容**: Issueの詳細な説明
   - **ラベル**: カンマ区切り（例：`bug,urgent`）
   - **担当者**: GitHubユーザー名（例：`username1,username2`）
   - **優先度**: low/medium/high/urgent から選択

### 方法2: Issue テンプレートを使用
1. GitHubリポジトリの「Issues」タブに移動
2. 「New issue」をクリック
3. 適切なテンプレートを選択：
   - 🐛 **バグレポート**: バグ報告用
   - ✨ **機能要求**: 新機能提案用

## 💻 コマンドラインでの作成

### 方法1: カスタムスクリプトを使用
```bash
./scripts/create-issue.sh "タイトル" "内容" "ラベル" "担当者"
```

#### 例：
```bash
./scripts/create-issue.sh \
  "APIのレスポンス時間改善" \
  "現在のAPIレスポンス時間が遅いため、最適化が必要です。目標は500ms以下にすることです。" \
  "performance,backend" \
  "developer1"
```

### 方法2: GitHub CLI 直接使用
```bash
gh issue create \
  --title "タイトル" \
  --body "内容" \
  --label "ラベル1,ラベル2" \
  --assignee "担当者"
```

## 🏷️ 推奨ラベル

プロジェクトで使用する標準的なラベル：

### 種類
- `bug`: バグ報告
- `enhancement`: 機能改善
- `feature`: 新機能
- `documentation`: ドキュメント関連

### 優先度
- `priority:low`: 低優先度
- `priority:medium`: 中優先度
- `priority:high`: 高優先度
- `priority:urgent`: 緊急

### 領域
- `frontend`: フロントエンド関連
- `backend`: バックエンド関連
- `database`: データベース関連
- `devops`: DevOps関連

## 🎯 Issue作成のベストプラクティス

1. **明確なタイトル**: 問題や要求を簡潔に表現
2. **詳細な説明**: 背景、期待する結果、受け入れ基準を含める
3. **適切なラベル**: 内容に応じたラベルを付与
4. **担当者の指定**: 可能であれば担当者を指定
5. **関連情報の添付**: スクリーンショット、ログ、参考資料など

## 🤝 AIアシスタントとの連携

AIアシスタントに以下のような指示ができます：

- 「新機能のissueを作成して」
- 「バグレポートのテンプレートでissueを作成」
- 「緊急度高でセキュリティ関連のissueを作成」
- 「週次レビュー用のissueを毎週月曜日に自動作成したい」

AIアシスタントが適切なワークフローファイルやスクリプトを作成・修正して、あなたの要求に応じたissue作成システムを構築します。 