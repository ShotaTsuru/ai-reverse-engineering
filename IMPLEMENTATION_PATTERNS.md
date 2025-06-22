# 実装パターン・ベストプラクティス集

このドキュメントは、汎用Webアプリケーションテンプレートで採用している実装パターンとベストプラクティスを記録したものです。新規プロジェクトでの実装時の参考資料として活用してください。

## 🏛️ アーキテクチャパターン

### 1. レイヤードアーキテクチャ（バックエンド）

```
presentation/   # controllers/ - HTTPハンドラー
    ↓
business/       # services/ - ビジネスロジック
    ↓
persistence/    # models/ - データアクセス
    ↓
database/       # PostgreSQL, Redis
```

**実装原則：**
- 上位レイヤーは下位レイヤーのみに依存
- ビジネスロジックはcontrollersに書かない
- データベース操作はmodels層で抽象化

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

## 🔧 Go バックエンド実装パターン

### HTTP ハンドラーパターン

```go
// 標準的なコントローラー構造
type ProjectController struct {
    db      *gorm.DB
    redis   *redis.Client
    service *services.ProjectService
}

// コンストラクタパターン
func NewProjectController(db *gorm.DB, redis *redis.Client) *ProjectController {
    return &ProjectController{
        db:      db,
        redis:   redis,
        service: services.NewProjectService(db),
    }
}

// ハンドラーメソッドパターン
func (pc *ProjectController) GetProjects(c *gin.Context) {
    // 1. リクエストバリデーション
    // 2. サービス層呼び出し
    // 3. レスポンス生成
    // 4. エラーハンドリング
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

### データベース接続パターン

```go
// シングルトンパターンでDB接続管理
var (
    db   *gorm.DB
    once sync.Once
)

func GetDB() *gorm.DB {
    once.Do(func() {
        var err error
        db, err = initDB()
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
    })
    return db
}

// 環境変数からの設定読み込み
func initDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(getLogLevel()),
    })
}
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

### データベースクエリ最適化

```go
// プリロードによるN+1問題解決
func (ps *ProjectService) GetProjectsWithFiles() ([]models.Project, error) {
    var projects []models.Project
    err := ps.db.
        Preload("Files").
        Where("status = ?", "active").
        Order("created_at DESC").
        Find(&projects).Error
    
    return projects, err
}

// ページネーション実装
func (ps *ProjectService) GetProjectsPaginated(page, size int) (*PaginatedResult, error) {
    var projects []models.Project
    var total int64
    
    offset := (page - 1) * size
    
    err := ps.db.Model(&models.Project{}).Count(&total).Error
    if err != nil {
        return nil, err
    }
    
    err = ps.db.
        Offset(offset).
        Limit(size).
        Find(&projects).Error
    
    return &PaginatedResult{
        Data:       projects,
        Total:      total,
        Page:       page,
        Size:       size,
        TotalPages: int(math.Ceil(float64(total) / float64(size))),
    }, err
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

### Go テストパターン

```go
// テーブル駆動テスト
func TestProjectService_CreateProject(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateProjectRequest
        want    *Project
        wantErr bool
    }{
        {
            name: "valid project creation",
            input: CreateProjectRequest{
                Name:        "Test Project",
                Description: "Test Description",
            },
            want: &Project{
                Name:        "Test Project",
                Description: "Test Description",
                Status:      "active",
            },
            wantErr: false,
        },
        // その他のテストケース...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}

// モックパターン
type MockProjectService struct {
    projects map[string]*Project
}

func (m *MockProjectService) GetProject(id string) (*Project, error) {
    if project, exists := m.projects[id]; exists {
        return project, nil
    }
    return nil, errors.New("project not found")
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

これらのパターンを参考に、プロジェクトの要件に応じてカスタマイズしてください。一貫性のあるコードベースと高い保守性を維持できます。 