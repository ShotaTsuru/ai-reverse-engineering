#!/bin/bash

# Issue作成スクリプト
# 使用方法: ./scripts/create-issue.sh "タイトル" "内容" "ラベル" "担当者"

set -e

TITLE="$1"
BODY="$2"
LABELS="$3"
ASSIGNEES="$4"

if [ -z "$TITLE" ]; then
    echo "❌ エラー: タイトルが必要です"
    echo "使用方法: ./scripts/create-issue.sh \"タイトル\" \"内容\" \"ラベル\" \"担当者\""
    exit 1
fi

if [ -z "$BODY" ]; then
    echo "❌ エラー: 内容が必要です"
    echo "使用方法: ./scripts/create-issue.sh \"タイトル\" \"内容\" \"ラベル\" \"担当者\""
    exit 1
fi

echo "🚀 Issue作成中..."

# コマンドの構築
CMD="gh issue create --title \"$TITLE\" --body \"$BODY\""

if [ -n "$LABELS" ]; then
    CMD="$CMD --label \"$LABELS\""
fi

if [ -n "$ASSIGNEES" ]; then
    CMD="$CMD --assignee \"$ASSIGNEES\""
fi

# Issue作成
eval $CMD

echo "✅ Issue「$TITLE」を作成しました" 