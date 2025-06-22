package routes

import (
	"reverse-engineering-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, redis *redis.Client) {
	// コントローラーの初期化
	projectController := controllers.NewProjectController(db, redis)
	fileController := controllers.NewFileController(db)
	analysisController := controllers.NewAnalysisController(db, redis)

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AI Reverse Engineering API is running",
		})
	})

	// API v1 グループ
	v1 := r.Group("/api/v1")
	{
		// プロジェクト管理
		projects := v1.Group("/projects")
		{
			projects.GET("/", projectController.GetProjects)
			projects.POST("/", projectController.CreateProject)
			projects.GET("/:id", projectController.GetProject)
			projects.PUT("/:id", projectController.UpdateProject)
			projects.DELETE("/:id", projectController.DeleteProject)
		}

		// ファイル管理
		files := v1.Group("/files")
		{
			files.POST("/upload", fileController.UploadFiles)
			files.GET("/project/:project_id", fileController.GetFilesByProject)
			files.GET("/:id", fileController.GetFile)
			files.DELETE("/:id", fileController.DeleteFile)
		}

		// AI解析
		analysis := v1.Group("/analysis")
		{
			analysis.POST("/start", analysisController.StartAnalysis)
			analysis.GET("/project/:project_id", analysisController.GetAnalysisByProject)
			analysis.GET("/:id", analysisController.GetAnalysis)
			analysis.GET("/:id/status", analysisController.GetAnalysisStatus)
		}
	}
}
