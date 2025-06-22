package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"reverse-engineering-backend/models"
	"reverse-engineering-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileController struct {
	db *gorm.DB
}

func NewFileController(db *gorm.DB) *FileController {
	return &FileController{
		db: db,
	}
}

func (fc *FileController) UploadFiles(c *gin.Context) {
	projectIDStr := c.PostForm("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required",
		})
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id",
		})
		return
	}

	// プロジェクトの存在確認
	var project models.Project
	if err := fc.db.First(&project, projectID).Error; err != nil {
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

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No files provided",
		})
		return
	}

	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	// アップロードディレクトリの作成
	projectUploadPath := filepath.Join(uploadPath, strconv.FormatUint(projectID, 10))
	if err := os.MkdirAll(projectUploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	var uploadedFiles []models.File

	for _, file := range files {
		// ファイル名の安全性チェック
		filename := filepath.Base(file.Filename)
		if filename == "." || filename == ".." || strings.Contains(filename, "/") {
			continue
		}

		// ファイルの保存
		savePath := filepath.Join(projectUploadPath, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file: " + filename,
			})
			return
		}

		// ファイル内容の読み取り
		content, err := fc.readFileContent(savePath)
		if err != nil {
			content = "" // バイナリファイルなどは内容を保存しない
		}

		// データベースへの保存
		fileModel := models.File{
			ProjectID: uint(projectID),
			Name:      filename,
			Path:      savePath,
			Size:      file.Size,
			MimeType:  file.Header.Get("Content-Type"),
			Content:   content,
			Language:  utils.DetectLanguage(filename),
		}

		if err := fc.db.Create(&fileModel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file metadata: " + filename,
			})
			return
		}

		uploadedFiles = append(uploadedFiles, fileModel)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"files":   uploadedFiles,
	})
}

func (fc *FileController) GetFilesByProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("project_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var files []models.File
	if err := fc.db.Where("project_id = ?", projectID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (fc *FileController) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}

	var file models.File
	if err := fc.db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch file",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": file,
	})
}

func (fc *FileController) DeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}

	var file models.File
	if err := fc.db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch file",
			})
		}
		return
	}

	// ファイルシステムからファイルを削除
	if err := os.Remove(file.Path); err != nil {
		// ファイルが既に存在しない場合は無視
	}

	// データベースからファイルを削除
	if err := fc.db.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

func (fc *FileController) readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// テキストファイルかどうかをチェック（簡単な実装）
	if utils.IsTextFile(content) {
		return string(content), nil
	}

	return "", nil // バイナリファイルは内容を保存しない
}
