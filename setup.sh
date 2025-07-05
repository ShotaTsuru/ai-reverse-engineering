#!/bin/bash

# 汎用Webアプリケーションテンプレート セットアップスクリプト
# 使用方法: ./setup.sh <project-name> <project-description>

set -e

PROJECT_NAME=${1:-"my-web-app"}
PROJECT_DESCRIPTION=${2:-"A modern web application"}

echo "🚀 新規プロジェクト「$PROJECT_NAME」のセットアップを開始します..."

# プロジェクト名のバリデーション
if [[ ! $PROJECT_NAME =~ ^[a-z0-9-]+$ ]]; then
    echo "❌ プロジェクト名は小文字、数字、ハイフンのみ使用可能です"
    exit 1
fi

# Go モジュール名を更新
echo "📝 Go モジュール名を更新中..."
find . -name "*.go" -type f -exec sed -i.bak "s/reverse-engineering-backend/${PROJECT_NAME}-backend/g" {} \;
find . -name "*.bak" -delete

# package.json を更新
echo "📝 フロントエンドパッケージ名を更新中..."
sed -i.bak "s/\"frontend\"/\"${PROJECT_NAME}-frontend\"/g" frontend/package.json
rm -f frontend/package.json.bak

# go.mod を更新
echo "📝 go.mod を更新中..."
sed -i.bak "s/module reverse-engineering-backend/module ${PROJECT_NAME}-backend/g" backend/go.mod
rm -f backend/go.mod.bak

# Docker Compose のコンテナ名を更新
echo "📝 Docker Compose 設定を更新中..."
sed -i.bak "s/reverse-eng-/${PROJECT_NAME}-/g" docker-compose.yml
rm -f docker-compose.yml.bak

# データベース名を更新
sed -i.bak "s/reverse_engineering_db/${PROJECT_NAME//-/_}_db/g" docker-compose.yml
rm -f docker-compose.yml.bak

# 環境変数ファイルを作成
echo "📝 環境変数ファイルを作成中..."
cp env.example .env
sed -i.bak "s/your_database_name/${PROJECT_NAME//-/_}_db/g" .env
sed -i.bak "s/Your App Name/${PROJECT_DESCRIPTION}/g" .env
rm -f .env.bak

# README.mdのプレースホルダーを更新
echo "📝 README.md を更新中..."
sed -i.bak "s/your-project-name/${PROJECT_NAME}/g" README.md
sed -i.bak "s/your-project-backend/${PROJECT_NAME}-backend/g" README.md
sed -i.bak "s/your-project-frontend/${PROJECT_NAME}-frontend/g" README.md
rm -f README.md.bak

# Git の初期化（既存の.gitがある場合）
if [ -d ".git" ]; then
    echo "🔄 Gitリポジトリを再初期化中..."
    rm -rf .git
    git init
    git add .
    git commit -m "Initial commit: $PROJECT_NAME project from template"
fi

# 依存関係のインストール
echo "📦 依存関係をインストール中..."
cd backend && go mod tidy && cd ..
cd frontend && npm install && cd ..

echo "✅ セットアップが完了しました！"
echo ""
echo "🎯 次のステップ:"
echo "   1. .env ファイルを確認・編集してください"
echo "   2. 開発環境を起動: make dev"
echo "   3. http://localhost:3000 でアプリケーションにアクセス"
echo ""
echo "🤖 AI機能の設定（オプション）:"
echo "   GitHub MCP サーバーの設定: make mcp-setup"
echo "   詳細は MCP_SETUP_GUIDE.md をご覧ください"
echo ""
echo "📚 詳細な使用方法は README.md をご覧ください" 