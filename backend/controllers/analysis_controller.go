package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"reverse-engineering-backend/models"
	"reverse-engineering-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AnalysisController struct {
	db            *gorm.DB
	redis         *redis.Client
	aiService     *services.AIService
	analysisQueue string
}

func NewAnalysisController(db *gorm.DB, redis *redis.Client) *AnalysisController {
	return &AnalysisController{
		db:            db,
		redis:         redis,
		aiService:     services.NewAIService(),
		analysisQueue: "analysis:queue",
	}
}

func (ac *AnalysisController) StartAnalysis(c *gin.Context) {
	var request struct {
		ProjectID uint     `json:"project_id" binding:"required"`
		Types     []string `json:"types" binding:"required"` // code_analysis, dependency_map, documentation, pattern_detection
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// プロジェクトの存在確認
	var project models.Project
	if err := ac.db.Preload("Files").First(&project, request.ProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to verify project",
			})
		}
		return
	}

	if len(project.Files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No files found in project",
		})
		return
	}

	var createdAnalyses []models.Analysis

	// 解析タスクの作成
	for _, analysisType := range request.Types {
		analysis := models.Analysis{
			ProjectID: request.ProjectID,
			Type:      analysisType,
			Status:    "pending",
		}

		if err := ac.db.Create(&analysis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create analysis task",
			})
			return
		}

		// Redisキューに解析タスクを追加
		taskData := map[string]interface{}{
			"analysis_id": analysis.ID,
			"project_id":  request.ProjectID,
			"type":        analysisType,
		}

		taskJSON, _ := json.Marshal(taskData)
		if err := ac.redis.LPush(context.Background(), ac.analysisQueue, taskJSON).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to queue analysis task",
			})
			return
		}

		createdAnalyses = append(createdAnalyses, analysis)
	}

	// プロジェクトのステータスを「分析中」に更新
	if err := ac.db.Model(&project).Update("status", "analyzing").Error; err != nil {
		// エラーはログに記録するが、レスポンスは成功とする
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Analysis started successfully",
		"analyses": createdAnalyses,
	})
}

func (ac *AnalysisController) GetAnalysisByProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("project_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var analyses []models.Analysis
	if err := ac.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&analyses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch analyses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analyses": analyses,
	})
}

func (ac *AnalysisController) GetAnalysis(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid analysis ID",
		})
		return
	}

	var analysis models.Analysis
	if err := ac.db.Preload("Project").Preload("File").First(&analysis, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Analysis not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch analysis",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis": analysis,
	})
}

func (ac *AnalysisController) GetAnalysisStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid analysis ID",
		})
		return
	}

	var analysis models.Analysis
	if err := ac.db.Select("id, status, created_at, updated_at").First(&analysis, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Analysis not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch analysis status",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         analysis.ID,
		"status":     analysis.Status,
		"created_at": analysis.CreatedAt,
		"updated_at": analysis.UpdatedAt,
	})
}
