# 🤖 AI駆動 GitHub Projects管理ガイド

このガイドでは、AIアシスタントと連携してGitHub Projectsを効率的に管理する方法を説明します。

## 🎯 概要

AI駆動プロジェクト管理システムは以下の機能を提供します：

- **AIアシスタント連携**: 自然言語でのプロジェクト操作
- **自動化ワークフロー**: GitHub Actionsによる自動プロジェクト管理
- **コマンドラインツール**: 効率的なプロジェクト操作
- **テンプレート**: 標準化されたプロジェクト構造

## 🚀 セットアップ

### 1. 初回セットアップ

```bash
# AI開発環境用プロジェクトを自動作成
./scripts/project-manager.sh ai-setup
```

このコマンドで以下が自動実行されます：
- 新しいプロジェクトの作成
- 標準フィールド（Priority, Sprint）の追加
- 開発テンプレートの適用

### 2. 権限の確認

GitHub CLIが適切な権限を持っていることを確認：

```bash
gh auth status
# project スコープが含まれていることを確認
```

## 🤖 AIアシスタントとの連携

### 基本的な指示パターン

#### プロジェクト作成
```
「新機能開発用のプロジェクトを作成してください」
「バックエンドAPI改善プロジェクトを立ち上げて」
```

#### タスク管理
```
「Issue #10をプロジェクト1に追加してください」
「Issue #5のステータスをIn Progressに変更して」
「緊急度高のバグ修正タスクを作成してプロジェクトに追加」
```

#### 進捗確認
```
「現在のプロジェクトの進捗を教えて」
「今週完了予定のタスクを確認して」
「プロジェクト2の未完了タスクをリストアップ」
```

### 高度な指示例

#### スプリント計画
```
「来週のスプリント計画を立てて。優先度高のタスクを中心に」
「現在のバックログから次の2週間のタスクを選択してスプリントに追加」
```

#### 自動化設定
```
「新しいPRが作成されたときに自動でプロジェクトに追加する設定を作って」
「バグラベルが付いたIssueを自動で緊急プロジェクトに追加したい」
```

## 💻 コマンドライン操作

### プロジェクト管理スクリプト

```bash
# 使用方法を表示
./scripts/project-manager.sh help

# プロジェクト一覧
./scripts/project-manager.sh list-projects

# 新しいプロジェクト作成
./scripts/project-manager.sh create-project "プロジェクト名" "説明"

# IssueをProjectに追加
./scripts/project-manager.sh add-issue 1 5

# 開発用プロジェクトの自動作成
./scripts/project-manager.sh create-dev-project
```

### GitHub CLI直接操作

```bash
# プロジェクト作成
gh project create --title "新プロジェクト" --body "説明"

# プロジェクト一覧
gh project list

# アイテム追加
gh project item-add 1 --url "https://github.com/owner/repo/issues/10"
```

## 🔄 GitHub Actions自動化

### 手動実行ワークフロー

**Actions** → **AI Project Management** → **Run workflow**

利用可能なアクション：
- `create-project`: 新しいプロジェクト作成
- `add-issue-to-project`: IssueをProjectに追加
- `update-status`: アイテムステータス更新
- `create-milestone`: マイルストーン作成

### 自動実行トリガー

以下のイベントで自動実行されます：
- **Issue作成時**: 自動でデフォルトプロジェクトに追加
- **PR作成時**: 自動でプロジェクトに追加し「In Progress」に設定
- **PRマージ時**: ステータスを「Done」に更新

## 📋 プロジェクト構造のベストプラクティス

### 推奨フィールド構成

1. **Status** (標準)
   - Todo
   - In Progress
   - In Review
   - Done

2. **Priority** (カスタム)
   - Low
   - Medium
   - High
   - Urgent

3. **Sprint** (カスタム)
   - Sprint 1
   - Sprint 2
   - Sprint 3
   - Backlog

4. **Assignee** (標準)
   - 担当者の設定

5. **Labels** (標準)
   - 分類用ラベル

### プロジェクトテンプレート

#### 機能開発プロジェクト
```
🚀 [機能名] 開発プロジェクト

## 🎯 目標
- 機能の設計・実装・テスト
- 期限: [日付]
- 担当チーム: [チーム名]

## 📋 マイルストーン
- [ ] 要件定義完了
- [ ] 設計レビュー完了
- [ ] 実装完了
- [ ] テスト完了
- [ ] リリース準備完了
```

#### バグ修正プロジェクト
```
🐛 [バグ修正] プロジェクト

## 🎯 目標
- 緊急度に応じたバグ修正
- 根本原因の特定と対策

## 📋 優先度
- Urgent: 即日対応
- High: 3日以内
- Medium: 1週間以内
- Low: 次回リリース時
```

## 🔧 カスタマイズ

### 独自フィールドの追加

```bash
# Estimate フィールド（工数見積もり）
gh project field-create PROJECT_NUMBER \
  --name "Estimate" \
  --single-select-option "1 day" \
  --single-select-option "3 days" \
  --single-select-option "1 week" \
  --single-select-option "2 weeks"

# Component フィールド（コンポーネント分類）
gh project field-create PROJECT_NUMBER \
  --name "Component" \
  --single-select-option "Frontend" \
  --single-select-option "Backend" \
  --single-select-option "Database" \
  --single-select-option "DevOps"
```

### ワークフローのカスタマイズ

`.github/workflows/project-management.yml`を編集して：
- 自動実行タイミングの調整
- 独自のステータス更新ロジック
- Slack通知の追加
- 他サービスとの連携

## 🚨 トラブルシューティング

### よくある問題

#### 1. 権限エラー
```bash
error: your authentication token is missing required scopes
```

**解決方法:**
```bash
gh auth refresh -s read:project,project
```

#### 2. プロジェクトが見つからない
```bash
could not resolve to a ProjectV2 with the number 999
```

**解決方法:**
```bash
# プロジェクト番号を確認
gh project list
```

#### 3. GraphQL API エラー
```bash
GraphQL error: Field 'fieldId' doesn't exist on type 'UpdateProjectV2ItemFieldValueInput'
```

**解決方法:**
- GitHub Actions経由でのステータス更新を使用
- APIの最新仕様を確認

### デバッグ方法

#### 1. 詳細ログの有効化
```bash
export GH_DEBUG=1
gh project list
```

#### 2. ワークフロー実行ログの確認
GitHub Actions → 該当ワークフロー → ログを確認

#### 3. APIレスポンスの確認
```bash
gh api graphql -f query='query { viewer { login } }'
```

## 🎯 開発フローの例

### 典型的なAI駆動開発フロー

1. **プロジェクト開始**
   ```
   AIアシスタント: 「新機能開発プロジェクトを作成して」
   ```

2. **タスク分解**
   ```
   AIアシスタント: 「ユーザー認証機能のタスクを分解してissueを作成し、プロジェクトに追加して」
   ```

3. **開発進行**
   ```
   AIアシスタント: 「Issue #15を開始したので、ステータスをIn Progressに変更して」
   ```

4. **レビュー段階**
   ```
   AIアシスタント: 「PR #20を作成したので、関連するIssueのステータスをIn Reviewに変更して」
   ```

5. **完了**
   ```
   AIアシスタント: 「Issue #15が完了したので、ステータスをDoneに変更して次のタスクを確認して」
   ```

## 📈 メトリクスと分析

### プロジェクト進捗の確認

GitHubのInsights機能を活用：
- バーンダウンチャート
- 速度グラフ
- 完了率の推移

### AIアシスタントとの分析

```
「今月のプロジェクト完了率を分析して」
「チームの生産性向上のための提案をして」
「ボトルネックになっているタスクを特定して」
```

## 🔮 今後の拡張計画

### 予定されている機能

1. **より高度なAI連携**
   - 自然言語処理による自動タスク分類
   - 工数見積もりの自動化
   - リスク予測とアラート

2. **外部サービス連携**
   - Slack通知
   - Jira連携
   - Time tracking

3. **分析機能**
   - 生産性ダッシュボード
   - 予測分析
   - レポート自動生成

### カスタマイズ要望

プロジェクトに特化した機能が必要な場合は、AIアシスタントに相談してください：

```
「プロジェクトに[具体的な機能]を追加したい」
「[業界/分野]に特化したプロジェクト管理機能を作りたい」
```

---

## 🤝 サポート

質問やカスタマイズ要望があれば、以下の方法でサポートを受けられます：

1. **AIアシスタントに直接相談**
2. **GitHub Issueでの報告**
3. **プロジェクトチームとの相談**

AIアシスタントが最適な解決策を提案し、必要に応じてワークフローやスクリプトを調整します。 