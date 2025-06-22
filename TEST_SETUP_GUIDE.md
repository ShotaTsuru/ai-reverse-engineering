# テスト環境設定・TDD実践ガイド

## 🧪 テスト環境の設定

### バックエンド（Go）テスト環境

#### 1. 必要なライブラリの追加

```bash
cd backend
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get github.com/testcontainers/testcontainers-go
```

#### 2. テスト設定ファイル作成

```go
// backend/internal/testutils/setup.go
package testutils

import (
    "database/sql"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/postgres"  
    "gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open mock sql db: %v", err)
    }
    
    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    
    if err != nil {
        t.Fatalf("Failed to open gorm db: %v", err)
    }
    
    return gormDB, mock
}
```

#### 3. Makefileにテストコマンド追加

```makefile
# 既存のMakefileに以下を追加
test-unit: ## Run unit tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests  
	cd backend && go test -v -tags=integration ./tests/integration/...

test-coverage: ## Show test coverage
	cd backend && go tool cover -html=coverage.out -o coverage.html
	open backend/coverage.html
```

### フロントエンド（Next.js）テスト環境

#### 1. E2Eテストライブラリの追加

```bash
cd frontend
npm install --save-dev @playwright/test
npx playwright install
```

#### 2. Playwrightの設定ファイル

```typescript
// frontend/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

#### 3. package.jsonのスクリプト更新

```json
{
  "scripts": {
    "test:e2e": "playwright test",
    "test:e2e:headed": "playwright test --headed",
    "test:e2e:ui": "playwright test --ui"
  }
}
```

---

## 🔄 TDD実践ワークフロー

### 1. 機能開発の流れ（例：プロジェクト作成機能）

#### Step 1: 受け入れ条件の定義
```gherkin
Feature: プロジェクト作成
  Scenario: 有効なデータでプロジェクトを作成する
    Given ユーザーがログインしている
    When 有効なプロジェクト名と説明を入力する
    And "作成"ボタンをクリックする
    Then プロジェクトが正常に作成される
    And 成功メッセージが表示される
    And プロジェクト一覧ページにリダイレクトされる
```

#### Step 2: E2Eテストの作成（Red）
```typescript
// frontend/tests/e2e/project-creation.spec.ts
import { test, expect } from '@playwright/test';

test('プロジェクト作成の成功フロー', async ({ page }) => {
  // Given: ログイン状態
  await page.goto('/login');
  await page.fill('[data-testid=email]', 'test@example.com');
  await page.fill('[data-testid=password]', 'password');
  await page.click('[data-testid=login-button]');
  
  // When: プロジェクト作成
  await page.goto('/projects/new');
  await page.fill('[data-testid=project-name]', 'テストプロジェクト');
  await page.fill('[data-testid=project-description]', 'テスト用のプロジェクトです');
  await page.click('[data-testid=create-button]');
  
  // Then: 成功の確認
  await expect(page.locator('[data-testid=success-message]')).toBeVisible();
  await expect(page).toHaveURL('/projects');
  await expect(page.locator('[data-testid=project-item]').first()).toContainText('テストプロジェクト');
});
```

#### Step 3: バックエンド単体テストの作成（Red）
```go
// backend/controllers/project_controller_test.go
package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreateProject(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    mockService := new(MockProjectService)
    controller := NewProjectController(nil, nil)
    controller.service = mockService
    
    projectData := map[string]interface{}{
        "name": "テストプロジェクト",
        "description": "テスト用のプロジェクトです",
    }
    
    jsonData, _ := json.Marshal(projectData)
    
    mockService.On("CreateProject", mock.AnythingOfType("*models.Project")).Return(nil)
    
    // Act
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonData))
    c.Request.Header.Set("Content-Type", "application/json")
    
    controller.CreateProject(c)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

#### Step 4: 最小限の実装（Green）
```go
// backend/controllers/project_controller.go
func (pc *ProjectController) CreateProject(c *gin.Context) {
    var project models.Project
    
    if err := c.ShouldBindJSON(&project); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := pc.service.CreateProject(&project); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Project created successfully",
        "project": project,
    })
}
```

#### Step 5: リファクタリング（Refactor）
- バリデーション追加
- エラーハンドリング改善
- ログ出力追加
- パフォーマンス最適化

---

## 🛠️ テストユーティリティ

### 共通テストヘルパー
```go
// backend/internal/testutils/helpers.go
package testutils

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    
    "github.com/gin-gonic/gin"
)

func MakeRequest(handler gin.HandlerFunc, method, url string, body interface{}) *httptest.ResponseRecorder {
    gin.SetMode(gin.TestMode)
    
    var buf bytes.Buffer
    if body != nil {
        json.NewEncoder(&buf).Encode(body)
    }
    
    req, _ := http.NewRequest(method, url, &buf)
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    handler(c)
    return w
}
```

### フィクスチャデータ
```go
// backend/internal/testutils/fixtures.go
package testutils

import "reverse-engineering-backend/models"

func CreateTestProject() *models.Project {
    return &models.Project{
        Name:        "テストプロジェクト",
        Description: "テスト用のプロジェクトです",
        Status:      "active",
    }
}

func CreateTestProjects(count int) []*models.Project {
    projects := make([]*models.Project, count)
    for i := 0; i < count; i++ {
        projects[i] = &models.Project{
            Name:        fmt.Sprintf("プロジェクト%d", i+1),
            Description: fmt.Sprintf("テスト用プロジェクト%d", i+1),
            Status:      "active",
        }
    }
    return projects
}
```

---

## 📊 テスト品質の管理

### カバレッジレポート
```bash
# バックエンドのカバレッジ計測
make test-coverage

# フロントエンドのカバレッジ計測
cd frontend && npm run test:coverage
```

### テスト品質チェックリスト
- [ ] **テストケースの充実性**
  - [ ] 正常系のテストケース
  - [ ] 異常系のテストケース  
  - [ ] 境界値のテストケース
  - [ ] エラーハンドリングのテストケース

- [ ] **テストの保守性**
  - [ ] 理解しやすいテスト名
  - [ ] DRY原則の適用
  - [ ] テストデータの再利用性
  - [ ] モックの適切な使用

- [ ] **パフォーマンステスト**
  - [ ] API応答時間の計測
  - [ ] データベース負荷テスト
  - [ ] フロントエンドの描画速度テスト

---

## 🚀 CI/CDパイプライン

### GitHub Actions設定例
```yaml
# .github/workflows/test.yml
name: Test Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.24
      - name: Run unit tests
        run: make test-unit
      - name: Run integration tests  
        run: make test-integration

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Run E2E tests
        run: cd frontend && npm run test:e2e
```

---

## 📚 参考資料・ベストプラクティス

### Go テスト
- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Test Patterns](https://github.com/golang/go/wiki/TestComments)

### Frontend E2E テスト
- [Playwright Documentation](https://playwright.dev/)
- [Testing Best Practices](https://playwright.dev/docs/best-practices)

### TDD
- [Test-Driven Development: By Example](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530)
- [Growing Object-Oriented Software, Guided by Tests](https://www.amazon.com/Growing-Object-Oriented-Software-Guided-Tests/dp/0321503627)
