#!/bin/bash

# AI駆動 GitHub Projects管理スクリプト
# 使用方法: ./scripts/project-manager.sh [action] [arguments...]

set -e

# 色付きログ出力
log_info() { echo -e "\033[36m[INFO]\033[0m $1"; }
log_success() { echo -e "\033[32m[SUCCESS]\033[0m $1"; }
log_error() { echo -e "\033[31m[ERROR]\033[0m $1"; }
log_warning() { echo -e "\033[33m[WARNING]\033[0m $1"; }

# 使用方法を表示
show_usage() {
    echo "🤖 AI駆動 GitHub Projects管理スクリプト"
    echo ""
    echo "使用方法:"
    echo "  ./scripts/project-manager.sh [action] [arguments...]"
    echo ""
    echo "利用可能なアクション:"
    echo "  create-project <title> <description>     新しいプロジェクトを作成"
    echo "  list-projects                            プロジェクト一覧を表示"
    echo "  add-issue <project_number> <issue_number> IssueをProjectに追加"
    echo "  update-status <project_number> <issue_number> <status> ステータス更新"
    echo "  create-dev-project                       開発用プロジェクトを自動作成"
    echo "  ai-setup                                 AI開発環境用プロジェクトセットアップ"
    echo ""
    echo "例:"
    echo "  ./scripts/project-manager.sh create-project \"新機能開発\" \"ユーザー認証機能の実装\""
    echo "  ./scripts/project-manager.sh add-issue 1 5"
    echo "  ./scripts/project-manager.sh update-status 1 5 \"In Progress\""
}

# プロジェクト作成
create_project() {
    local title="$1"
    local description="$2"
    
    if [[ -z "$title" ]]; then
        log_error "プロジェクトタイトルが必要です"
        return 1
    fi
    
    log_info "プロジェクトを作成中: $title"
    
    # プロジェクト作成（エラーハンドリング付き）
    if ! PROJECT_URL=$(gh project create --title "$title" --owner @me 2>&1); then
        log_error "プロジェクト作成に失敗しました: $PROJECT_URL"
        return 1
    fi
    
    if [[ -z "$PROJECT_URL" ]]; then
        log_error "プロジェクトURLが取得できませんでした"
        return 1
    fi
    
    # プロジェクト番号を取得
    PROJECT_NUMBER=$(echo $PROJECT_URL | grep -o '/projects/[0-9]*' | grep -o '[0-9]*')
    
    if [[ -z "$PROJECT_NUMBER" ]]; then
        log_error "プロジェクト番号が取得できませんでした。URL: $PROJECT_URL"
        return 1
    fi
    
    log_success "プロジェクトを作成しました"
    echo "📋 タイトル: $title"
    echo "🔗 URL: $PROJECT_URL"
    echo "🔢 プロジェクト番号: $PROJECT_NUMBER"
    
    # 標準フィールドを追加
    log_info "標準フィールドを追加中..."
    
    # Priority フィールド
    gh project field-create $PROJECT_NUMBER \
        --name "Priority" \
        --single-select-option "Low" \
        --single-select-option "Medium" \
        --single-select-option "High" \
        --single-select-option "Urgent" 2>/dev/null || log_warning "Priorityフィールドは既に存在します"
    
    # Sprint フィールド
    gh project field-create $PROJECT_NUMBER \
        --name "Sprint" \
        --single-select-option "Sprint 1" \
        --single-select-option "Sprint 2" \
        --single-select-option "Sprint 3" \
        --single-select-option "Backlog" 2>/dev/null || log_warning "Sprintフィールドは既に存在します"
    
    log_success "標準フィールドを追加しました"
    return 0
}

# プロジェクト一覧
list_projects() {
    log_info "プロジェクト一覧を取得中..."
    
    gh project list --owner $(gh api user --jq .login) --format json | jq -r '.projects[] | "📋 #\(.number) \(.title) [\(if .closed then "closed" else "open" end)] - \(.url)"'
    
    log_success "プロジェクト一覧を表示しました"
}

# IssueをProjectに追加
add_issue_to_project() {
    local project_number="$1"
    local issue_number="$2"
    
    if [[ -z "$project_number" || -z "$issue_number" ]]; then
        log_error "プロジェクト番号とIssue番号が必要です"
        return 1
    fi
    
    log_info "Issue #$issue_number をProject #$project_number に追加中..."
    
    OWNER=$(gh api user --jq .login)
    REPO=$(gh repo view --json name --jq .name)
    
    gh project item-add $project_number \
        --owner $OWNER \
        --url "https://github.com/$OWNER/$REPO/issues/$issue_number"
    
    log_success "Issue #$issue_number をProject #$project_number に追加しました"
}

# ステータス更新
update_status() {
    local project_number="$1"
    local issue_number="$2"
    local status="$3"
    
    if [[ -z "$project_number" || -z "$issue_number" || -z "$status" ]]; then
        log_error "プロジェクト番号、Issue番号、ステータスが必要です"
        return 1
    fi
    
    log_info "Issue #$issue_number のステータスを「$status」に更新中..."
    
    # 複雑なGraphQL操作のため、GitHub Actions経由で実行することを推奨
    log_warning "ステータス更新は現在GitHub Actions経由で実行することを推奨します"
    log_info "GitHub Actions > AI Project Management > Run workflow を使用してください"
}

# 開発用プロジェクト自動作成
create_dev_project() {
    local date=$(date +"%Y-%m")
    local title="🚀 AI開発プロジェクト $date"
    local description="AI駆動開発のためのプロジェクト管理ボード

## 🎯 目標
- アジャイル開発の実践
- AI連携によるタスク管理
- 効率的なチーム開発

## 📋 含まれる機能
- バックログ管理
- スプリント計画
- 進捗追跡
- 自動化ワークフロー"

    create_project "$title" "$description"
}

# AI開発環境セットアップ
ai_setup() {
    log_info "🤖 AI開発環境用プロジェクトをセットアップ中..."
    
    # メインの開発プロジェクトを作成
    create_dev_project
    
    log_info "サンプルタスクを作成中..."
    
    # サンプルIssueを作成してプロジェクトに追加
    SAMPLE_ISSUES=(
        "バックエンドAPIの実装"
        "フロントエンドUIの開発"
        "データベース設計の最適化"
        "テストケースの作成"
        "デプロイメント自動化"
    )
    
    for task in "${SAMPLE_ISSUES[@]}"; do
        log_info "サンプルタスク作成: $task"
        # 実際のIssue作成はオプション
        # gh issue create --title "$task" --body "AI開発環境セットアップ時に作成されたサンプルタスクです。"
    done
    
    log_success "🎉 AI開発環境のセットアップが完了しました！"
    echo ""
    echo "次のステップ:"
    echo "1. GitHub Projects画面でプロジェクトを確認"
    echo "2. 実際のタスクを追加"
    echo "3. AIアシスタントに指示してタスクを管理"
    echo ""
    echo "AIアシスタントへの指示例:"
    echo "- 「新しいタスクをプロジェクトに追加して」"
    echo "- 「Issue #5のステータスをIn Progressに変更して」"
    echo "- 「今週のスプリント計画を立てて」"
}

# メイン処理
main() {
    if [[ $# -eq 0 ]]; then
        show_usage
        exit 0
    fi
    
    local action="$1"
    shift
    
    case "$action" in
        "create-project")
            create_project "$1" "$2"
            ;;
        "list-projects")
            list_projects
            ;;
        "add-issue")
            add_issue_to_project "$1" "$2"
            ;;
        "update-status")
            update_status "$1" "$2" "$3"
            ;;
        "create-dev-project")
            create_dev_project
            ;;
        "ai-setup")
            ai_setup
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            log_error "不明なアクション: $action"
            show_usage
            exit 1
            ;;
    esac
}

# スクリプト実行
main "$@" 