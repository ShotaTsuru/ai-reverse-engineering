# 汎用Webアプリケーション開発テンプレート

Go + Next.js + PostgreSQL + Redis による現代的なWebアプリケーション開発のためのテンプレートプロジェクトです。

## 🎯 プロジェクト概要

このテンプレートは、チーム開発に適した構成で、様々なWebアプリケーションプロジェクトに適用できる汎用的な基盤を提供します。

## 🏗️ アーキテクチャ

### フロントエンド (React + Next.js)
- 📁 `frontend/` - Next.js 14ベースのWebアプリケーション
- 🎨 TypeScript + Tailwind CSS による型安全でモダンなUI開発
- 📊 React Query によるサーバーステート管理

### バックエンド (Go)
- 📁 `backend/` - Goベースの高性能RESTful APIサーバー
- 🏛️ クリーンアーキテクチャパターンの実装
- 📝 ルーティング、コントローラー、サービス層の明確な分離

## 🚀 主要機能・特徴

- 🔧 **チーム開発対応**: Docker Compose による統一開発環境
- 📈 **スケーラブル**: PostgreSQL + Redis によるデータ層
- 🔒 **セキュリティ**: CORS設定、環境変数管理
- ⚡ **高速開発**: ホットリロード対応開発サーバー
- 🧪 **テスト準備**: フロントエンド・バックエンド両方のテスト環境
- 🤖 **AI Issue作成システム**: GitHub Actionsとテンプレートによる効率的なissue管理

## 🛠️ 技術スタック

### フロントエンド
- **Next.js 14** - React フレームワーク
- **TypeScript** - 型安全性
- **Tailwind CSS** - スタイリング
- **React Query** - データフェッチング・キャッシュ

### バックエンド
- **Go 1.21+** - 高性能バックエンド
- **Gin** - Webフレームワーク
- **GORM** - ORM
- **PostgreSQL** - メインデータベース
- **Redis** - キャッシュ・セッション管理

## 📦 プロジェクトセットアップ

### 前提条件
- Node.js 18+
- Go 1.21+
- Docker & Docker Compose

### 新規プロジェクト作成手順

1. **テンプレートをクローン**
```bash
git clone <repository-url> your-project-name
cd your-project-name
```

2. **プロジェクト固有の設定**
```bash
# プロジェクト名を変更
find . -name "*.go" -exec sed -i 's/reverse-engineering-backend/your-project-backend/g' {} \;
find . -name "*.json" -exec sed -i 's/"frontend"/"your-project-frontend"/g' {} \;

# Gitの初期化
rm -rf .git
git init
git add .
git commit -m "Initial commit from template"
```

3. **環境設定**
```bash
cp env.example .env
# .envファイルを編集してプロジェクト固有の設定を追加
```

4. **依存関係のインストール**
```bash
make install
```

5. **開発環境の起動**
```bash
make dev
```

## 🎯 開発ワークフロー

### 日常開発
```bash
# 開発環境全体を起動
make dev

# フロントエンドのみ
make dev-frontend

# バックエンドのみ  
make dev-backend
```

### テスト実行
```bash
# 全テスト実行
make test

# フロントエンドテスト
make test-frontend

# バックエンドテスト
make test-backend
```

## 🔧 カスタマイズガイド

### プロジェクト名の変更
1. `go.mod` のモジュール名を変更
2. import文の更新
3. `package.json` の name フィールド更新
4. Docker Compose のコンテナ名変更

### 新機能の追加
1. **バックエンド**: `controllers/`, `services/`, `models/` にファイル追加
2. **フロントエンド**: `src/components/`, `src/pages/` にコンポーネント追加
3. **API**: `routes/routes.go` にエンドポイント追加

### データベーススキーマ
1. `models/` でGORMモデル定義
2. マイグレーション実行

## 🚀 本番デプロイ

```bash
# 本番ビルド
make build-prod

# Docker本番イメージ作成
docker-compose -f docker-compose.prod.yml build
```

## 🤖 AI駆動プロジェクト管理システム

プロジェクトには効率的なissue管理とGitHub Projects連携のためのAI対応システムが組み込まれています。

### 🔧 AI Issue作成システム
- **GitHub Actions ワークフロー**: 手動でissueを作成
- **Issue テンプレート**: バグレポートと機能要求の標準化
- **コマンドラインスクリプト**: 素早いissue作成
- **AI アシスタント連携**: 自然言語での指示によるissue作成

### 📋 AI Projects管理システム
- **プロジェクト自動作成**: AI指示による効率的なプロジェクト立ち上げ
- **自動タスク管理**: Issue/PRの自動プロジェクト追加とステータス更新
- **スプリント管理**: 優先度とスプリント計画の自動化
- **進捗追跡**: リアルタイムでの開発進捗確認

### 🔗 GitHub MCP サーバー統合
- **Model Context Protocol 対応**: 自然言語でのGitHub操作
- **Visual Studio Code 統合**: エディタ内でのGitHub操作
- **高度なAI連携**: コンテキストを理解したGitHub操作

### 使用方法
- Issue作成: [`ISSUE_CREATION_GUIDE.md`](./ISSUE_CREATION_GUIDE.md)
- Projects管理: [`AI_PROJECT_MANAGEMENT_GUIDE.md`](./AI_PROJECT_MANAGEMENT_GUIDE.md)
- MCP設定: [`MCP_SETUP_GUIDE.md`](./MCP_SETUP_GUIDE.md)

### クイックスタート
```bash
# AI開発環境のセットアップ
./scripts/project-manager.sh ai-setup

# プロジェクト一覧確認
./scripts/project-manager.sh list-projects
```

## 📋 チーム開発のベストプラクティス

### 1. 環境統一
- Docker Compose を使用して全員同じ開発環境で作業
- `.env.example` を参考に個人の `.env` を設定

### 2. コード品質
- ESLint/TypeScript による静的解析
- Go の標準的なフォーマット・規約に従う
- プルリクエスト前にテスト実行

### 3. Git ワークフロー
- feature ブランチでの開発
- プルリクエストによるコードレビュー
- main ブランチの保護

### 4. Issue管理
- Issue作成システムを活用した効率的なタスク管理
- 適切なラベルと優先度の設定
- テンプレートを使用した情報の標準化

## 🤝 コントリビューション

このテンプレートの改善にご協力ください。プルリクエストやイシューの報告を歓迎いたします。

## 📄 ライセンス

MIT License 