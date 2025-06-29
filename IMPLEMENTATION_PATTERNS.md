# 実装パターン・ベストプラクティス集

このドキュメントは、汎用Webアプリケーションテンプレートで採用している実装パターンとベストプラクティスを記録したものです。新規プロジェクトでの実装時の参考資料として活用してください。

## 🏛️ アーキテクチャパターン

### 1. クリーンアーキテクチャ（バックエンド）

```
外部レイヤー (Frameworks & Drivers)
    ↓
インターフェースアダプター (Interface Adapters)
    ↓
アプリケーションビジネスルール (Application Business Rules)
    ↓
エンタープライズビジネスルール (Enterprise Business Rules)
```

**レイヤー構成：**
```
controllers/     # インターフェースアダプター - HTTPハンドラー
    ↓ (依存関係逆転)
services/        # アプリケーションビジネスルール - ユースケース
    ↓ (依存関係逆転)
domain/          # エンタープライズビジネスルール - エンティティ
    ↑ (インターフェース実装)
repositories/    # データアクセス抽象化
    ↑ (実装)
infrastructure/  # 外部レイヤー - データベース、外部API
```

**実装原則：**
- 依存関係は内側に向かってのみ流れる
- 内側の層は外側の層を知らない
- インターフェースによる依存関係逆転
- ドメインロジックはドメイン層に集約

### 2. コンポーネント指向設計（フロントエンド）

```
pages/          # ページコンポーネント
    ↓
templates/      # レイアウトテンプレート
    ↓
organisms/      # 複合UIコンポーネント
    ↓
molecules/      # 基本UIコンポーネント
    ↓
atoms/          # プリミティブコンポーネント
```

## 🔧 Go クリーンアーキテクチャ実装パターン

### ドメイン層（エンティティ）

```go
// domain/entities/project.go
package entities

import "time"

type Project struct {
    ID          string
    Name        string
    Description string
    Status      ProjectStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ProjectStatus string

const (
    StatusPending    ProjectStatus = "pending"
    StatusActive     ProjectStatus = "active"
    StatusCompleted  ProjectStatus = "completed"
)

// ビジネスルール
func (p *Project) CanBeDeleted() bool {
    return p.Status == StatusPending || p.Status == StatusCompleted
}
```

### リポジトリインターフェース（ドメイン層）

```go
// domain/repositories/project_repository.go
package repositories

import "context"

type ProjectRepository interface {
    Save(ctx context.Context, project *entities.Project) error
    FindByID(ctx context.Context, id string) (*entities.Project, error)
    FindAll(ctx context.Context) ([]*entities.Project, error)
    Delete(ctx context.Context, id string) error
}
```

### ユースケース層（アプリケーションビジネスルール）

```go
// application/usecases/project_usecase.go
package usecases

type ProjectUsecase struct {
    projectRepo repositories.ProjectRepository
}

func NewProjectUsecase(projectRepo repositories.ProjectRepository) *ProjectUsecase {
    return &ProjectUsecase{
        projectRepo: projectRepo,
    }
}

func (pu *ProjectUsecase) CreateProject(ctx context.Context, name, description string) (*entities.Project, error) {
    // ビジネスルール適用
    project := &entities.Project{
        ID:          generateID(),
        Name:        name,
        Description: description,
        Status:      entities.StatusPending,
        CreatedAt:   time.Now(),
    }
    
    return project, pu.projectRepo.Save(ctx, project)
}
```

### インフラストラクチャ層（リポジトリ実装）

```go
// infrastructure/repositories/project_repository_impl.go
package repositories

type ProjectRepositoryImpl struct {
    db *gorm.DB
}

func NewProjectRepositoryImpl(db *gorm.DB) repositories.ProjectRepository {
    return &ProjectRepositoryImpl{db: db}
}

func (r *ProjectRepositoryImpl) Save(ctx context.Context, project *entities.Project) error {
    model := toGormModel(project)
    return r.db.WithContext(ctx).Save(model).Error
}
```

### インターフェースアダプター層（コントローラー）

```go
// interfaces/controllers/project_controller.go
package controllers

type ProjectController struct {
    projectUsecase *usecases.ProjectUsecase
}

func NewProjectController(projectUsecase *usecases.ProjectUsecase) *ProjectController {
    return &ProjectController{
        projectUsecase: projectUsecase,
    }
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
    var req CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    project, err := pc.projectUsecase.CreateProject(c.Request.Context(), req.Name, req.Description)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"project": toResponse(project)})
}
```

### エラーハンドリングパターン

```go
// 統一エラーレスポンス形式
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}

// エラーハンドリングのベストプラクティス
func (pc *ProjectController) HandleError(c *gin.Context, err error) {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusNotFound, ErrorResponse{
            Error: "Resource not found",
        })
        return
    }
    
    // その他のエラー処理...
    c.JSON(http.StatusInternalServerError, ErrorResponse{
        Error: "Internal server error",
    })
}
```

### 依存性注入パターン（クリーンアーキテクチャ）

```go
// infrastructure/di/container.go
package di

type Container struct {
    projectRepo    repositories.ProjectRepository
    projectUsecase *usecases.ProjectUsecase
    projectController *controllers.ProjectController
}

func NewContainer(db *gorm.DB) *Container {
    // インフラストラクチャ層の実装を注入
    projectRepo := repositories.NewProjectRepositoryImpl(db)
    
    // ユースケース層にリポジトリを注入
    projectUsecase := usecases.NewProjectUsecase(projectRepo)
    
    // コントローラー層にユースケースを注入
    projectController := controllers.NewProjectController(projectUsecase)
    
    return &Container{
        projectRepo:       projectRepo,
        projectUsecase:    projectUsecase,
        projectController: projectController,
    }
}

func (c *Container) GetProjectController() *controllers.ProjectController {
    return c.projectController
}
```

### データベース接続の抽象化

```go
// infrastructure/database/connection.go
package database

type DatabaseConnection interface {
    GetDB() *gorm.DB
    Close() error
    Migrate(models ...interface{}) error
}

type PostgreSQLConnection struct {
    db *gorm.DB
}

func NewPostgreSQLConnection(config *Config) (DatabaseConnection, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        config.Host,
        config.User,
        config.Password,
        config.Name,
        config.Port,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(getLogLevel()),
    })
    
    if err != nil {
        return nil, err
    }
    
    return &PostgreSQLConnection{db: db}, nil
}

func (p *PostgreSQLConnection) GetDB() *gorm.DB {
    return p.db
}

func (p *PostgreSQLConnection) Close() error {
    sqlDB, err := p.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
```
```

## ⚛️ React/Next.js 実装パターン

### カスタムHooksパターン

```typescript
// データフェッチング用カスタムHook
export function useProjects() {
    return useQuery({
        queryKey: ['projects'],
        queryFn: async () => {
            const response = await fetch('/api/projects');
            if (!response.ok) {
                throw new Error('Failed to fetch projects');
            }
            return response.json();
        },
        staleTime: 5 * 60 * 1000, // 5分
    });
}

// フォーム管理用カスタムHook
export function useProjectForm(initialData?: Project) {
    const [formData, setFormData] = useState(initialData || {});
    const [errors, setErrors] = useState<Record<string, string>>({});
    
    const validate = useCallback(() => {
        // バリデーションロジック
    }, [formData]);
    
    return { formData, setFormData, errors, validate };
}
```

### コンポーネント設計パターン

```typescript
// Props型定義
interface ButtonProps {
    variant?: 'primary' | 'secondary' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
    children: React.ReactNode;
    onClick?: () => void;
}

// Compound Component パターン
export const Card = {
    Root: ({ children, className, ...props }: CardRootProps) => (
        <div className={cn("rounded-lg border", className)} {...props}>
            {children}
        </div>
    ),
    Header: ({ children, className, ...props }: CardHeaderProps) => (
        <div className={cn("p-4 border-b", className)} {...props}>
            {children}
        </div>
    ),
    Content: ({ children, className, ...props }: CardContentProps) => (
        <div className={cn("p-4", className)} {...props}>
            {children}
        </div>
    ),
};
```

### 状態管理パターン

```typescript
// Zustand を使用した状態管理例
interface AppState {
    user: User | null;
    projects: Project[];
    setUser: (user: User | null) => void;
    addProject: (project: Project) => void;
    removeProject: (id: string) => void;
}

export const useAppStore = create<AppState>((set) => ({
    user: null,
    projects: [],
    setUser: (user) => set({ user }),
    addProject: (project) => 
        set((state) => ({ projects: [...state.projects, project] })),
    removeProject: (id) => 
        set((state) => ({ 
            projects: state.projects.filter(p => p.id !== id) 
        })),
}));
```

## 🐳 Docker ベストプラクティス

### マルチステージビルドパターン

```dockerfile
# 本番用Dockerfile例
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Docker Compose開発最適化

```yaml
# 開発効率最大化の設定例
services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    volumes:
      - ./backend:/app:cached  # macOS最適化
      - /app/vendor           # ベンダーキャッシュ
    environment:
      - CGO_ENABLED=0
    depends_on:
      postgres:
        condition: service_healthy
    
  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
```

## 🔒 セキュリティパターン

### 環境変数管理

```go
// 機密情報の適切な管理
type Config struct {
    JWTSecret    string `env:"JWT_SECRET,required"`
    DatabaseURL  string `env:"DATABASE_URL,required"`
    RedisURL     string `env:"REDIS_URL,required"`
    Environment  string `env:"GO_ENV" envDefault:"development"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    
    // 開発環境でのデフォルト値チェック
    if cfg.Environment == "development" && cfg.JWTSecret == "dev-secret" {
        log.Warn("Using default JWT secret in development")
    }
    
    return cfg, nil
}
```

### CORS設定パターン

```go
// 環境に応じたCORS設定
func setupCORS(r *gin.Engine) {
    config := cors.DefaultConfig()
    
    if os.Getenv("GO_ENV") == "production" {
        config.AllowOrigins = strings.Split(os.Getenv("CORS_ORIGINS"), ",")
    } else {
        config.AllowOrigins = []string{
            "http://localhost:3000",
            "http://127.0.0.1:3000",
        }
    }
    
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
    config.AllowCredentials = true
    
    r.Use(cors.New(config))
}
```

## 📊 パフォーマンス最適化パターン

### データベースクエリ最適化（クリーンアーキテクチャ）

```go
// リポジトリでのプリロード実装
func (r *ProjectRepositoryImpl) FindAllWithFiles(ctx context.Context, status entities.ProjectStatus) ([]*entities.Project, error) {
    var models []ProjectModel
    err := r.db.WithContext(ctx).
        Preload("Files").
        Where("status = ?", string(status)).
        Order("created_at DESC").
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    return toEntities(models), nil
}

// ページネーション専用リポジトリメソッド
func (r *ProjectRepositoryImpl) FindPaginated(ctx context.Context, page, size int) (*repositories.PaginatedResult[*entities.Project], error) {
    var models []ProjectModel
    var total int64
    
    offset := (page - 1) * size
    
    // 総件数取得
    err := r.db.WithContext(ctx).Model(&ProjectModel{}).Count(&total).Error
    if err != nil {
        return nil, err
    }
    
    // データ取得
    err = r.db.WithContext(ctx).
        Offset(offset).
        Limit(size).
        Order("created_at DESC").
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    return &repositories.PaginatedResult[*entities.Project]{
        Data:       toEntities(models),
        Total:      total,
        Page:       page,
        Size:       size,
        TotalPages: int(math.Ceil(float64(total) / float64(size))),
    }, nil
}

// ユースケース層でのクエリ最適化
func (pu *ProjectUsecase) GetActiveProjectsWithFiles(ctx context.Context) ([]*entities.Project, error) {
    return pu.projectRepo.FindAllWithFiles(ctx, entities.StatusActive)
}

func (pu *ProjectUsecase) GetProjectsPaginated(ctx context.Context, page, size int) (*repositories.PaginatedResult[*entities.Project], error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 20 // デフォルトサイズ
    }
    
    return pu.projectRepo.FindPaginated(ctx, page, size)
}
```

### 仕様パターン（複雑なクエリ条件）

```go
// domain/specifications/project_specification.go
package specifications

type ProjectSpecification interface {
    IsSatisfiedBy(project *entities.Project) bool
    ToSQL() (string, []interface{})
}

type ActiveProjectSpec struct{}

func (s *ActiveProjectSpec) IsSatisfiedBy(project *entities.Project) bool {
    return project.Status == entities.StatusActive
}

func (s *ActiveProjectSpec) ToSQL() (string, []interface{}) {
    return "status = ?", []interface{}{string(entities.StatusActive)}
}

// リポジトリでの仕様パターン使用
func (r *ProjectRepositoryImpl) FindBySpecification(ctx context.Context, spec specifications.ProjectSpecification) ([]*entities.Project, error) {
    query, args := spec.ToSQL()
    
    var models []ProjectModel
    err := r.db.WithContext(ctx).
        Where(query, args...).
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    return toEntities(models), nil
}
```

### React パフォーマンス最適化

```typescript
// メモ化パターン
const ProjectCard = memo(({ project }: { project: Project }) => {
    const handleClick = useCallback(() => {
        // クリック処理
    }, [project.id]);
    
    return (
        <Card onClick={handleClick}>
            <Card.Header>
                <h3>{project.name}</h3>
            </Card.Header>
            <Card.Content>
                <p>{project.description}</p>
            </Card.Content>
        </Card>
    );
});

// 仮想化リスト（大量データ対応）
import { FixedSizeList as List } from 'react-window';

const ProjectList = ({ projects }: { projects: Project[] }) => (
    <List
        height={600}
        itemCount={projects.length}
        itemSize={120}
        itemData={projects}
    >
        {({ index, style, data }) => (
            <div style={style}>
                <ProjectCard project={data[index]} />
            </div>
        )}
    </List>
);
```

## 🧪 テスト戦略パターン

### Go クリーンアーキテクチャテストパターン

```go
// ドメイン層テスト（エンティティ）
func TestProject_CanBeDeleted(t *testing.T) {
    tests := []struct {
        name   string
        status entities.ProjectStatus
        want   bool
    }{
        {"pending project can be deleted", entities.StatusPending, true},
        {"completed project can be deleted", entities.StatusCompleted, true},
        {"active project cannot be deleted", entities.StatusActive, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            project := &entities.Project{Status: tt.status}
            got := project.CanBeDeleted()
            assert.Equal(t, tt.want, got)
        })
    }
}

// ユースケース層テスト（モックリポジトリ使用）
func TestProjectUsecase_CreateProject(t *testing.T) {
    mockRepo := &MockProjectRepository{}
    usecase := usecases.NewProjectUsecase(mockRepo)
    
    project, err := usecase.CreateProject(context.Background(), "Test Project", "Description")
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Project", project.Name)
    assert.Equal(t, entities.StatusPending, project.Status)
    assert.True(t, mockRepo.SaveCalled)
}

// リポジトリ層テスト（インメモリDB使用）
func TestProjectRepositoryImpl_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := repositories.NewProjectRepositoryImpl(db)
    
    project := &entities.Project{
        ID:   "test-id",
        Name: "Test Project",
    }
    
    err := repo.Save(context.Background(), project)
    assert.NoError(t, err)
    
    // 保存されたデータの検証
    saved, err := repo.FindByID(context.Background(), "test-id")
    assert.NoError(t, err)
    assert.Equal(t, project.Name, saved.Name)
}
```

### モックとテストダブルパターン

```go
// リポジトリモック
type MockProjectRepository struct {
    projects   map[string]*entities.Project
    SaveCalled bool
}

func NewMockProjectRepository() *MockProjectRepository {
    return &MockProjectRepository{
        projects: make(map[string]*entities.Project),
    }
}

func (m *MockProjectRepository) Save(ctx context.Context, project *entities.Project) error {
    m.SaveCalled = true
    m.projects[project.ID] = project
    return nil
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id string) (*entities.Project, error) {
    if project, exists := m.projects[id]; exists {
        return project, nil
    }
    return nil, domain.ErrProjectNotFound
}

// コントローラー統合テスト
func TestProjectController_CreateProject(t *testing.T) {
    // テスト用の依存関係を構築
    mockRepo := NewMockProjectRepository()
    usecase := usecases.NewProjectUsecase(mockRepo)
    controller := controllers.NewProjectController(usecase)
    
    // テスト用Ginエンジン
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/projects", controller.CreateProject)
    
    // テストリクエスト
    requestBody := `{"name":"Test Project","description":"Test Description"}`
    req := httptest.NewRequest("POST", "/projects", strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // レスポンス検証
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "Test Project")
}
```

### React テストパターン

```typescript
// コンポーネントテスト
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const renderWithProviders = (component: React.ReactElement) => {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    });
    
    return render(
        <QueryClientProvider client={queryClient}>
            {component}
        </QueryClientProvider>
    );
};

test('ProjectCard displays project information', () => {
    const mockProject = {
        id: '1',
        name: 'Test Project',
        description: 'Test Description',
    };
    
    renderWithProviders(<ProjectCard project={mockProject} />);
    
    expect(screen.getByText('Test Project')).toBeInTheDocument();
    expect(screen.getByText('Test Description')).toBeInTheDocument();
});
```

## 🚀 デプロイメントパターン

### CI/CD パイプライン

```yaml
# GitHub Actions 例
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21
          
      - name: Run backend tests
        run: |
          cd backend
          go test ./...
          
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18
          
      - name: Run frontend tests
        run: |
          cd frontend
          npm ci
          npm test
```

これらのクリーンアーキテクチャパターンを参考に、プロジェクトの要件に応じてカスタマイズしてください。依存関係逆転の原則により、テスタブルで保守性の高いコードベースを維持できます。 