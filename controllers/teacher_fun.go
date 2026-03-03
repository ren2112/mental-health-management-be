package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mental-health-management-be/config"
	"mental-health-management-be/constants"
	"mental-health-management-be/converter"
	"mental-health-management-be/models"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
	"strconv"
)

func UpdateTeacherSelfInformation(c *gin.Context) {

	// ===== 身份校验 =====
	roleAny, _ := c.Get("role")
	role := roleAny.(int)

	if role != constants.RoleTeacher {
		response.CommonResp(c, 1, "身份异常", nil)
		return
	}

	// ===== JWT 获取用户 =====
	userIDAny, _ := c.Get("userID")
	teacherID := userIDAny.(int)

	var req struct {
		Name         string `json:"name"`
		Sex          *int8  `json:"sex"`
		Phone        string `json:"phone"`
		Email        string `json:"email" binding:"omitempty,email"`
		Introduction string `json:"introduction"`
		Workspace    string `json:"workspace"`
	}

	// ===== 参数绑定 =====
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 查询教师 =====
	var teacher models.Teacher
	if err := config.DB.First(&teacher, teacherID).Error; err != nil {
		response.CommonResp(c, 1, "教师不存在", nil)
		return
	}

	// ===== email 唯一校验 =====
	if req.Email != "" && req.Email != teacher.Email {

		var count int64
		config.DB.Model(&models.Teacher{}).
			Where("email = ? AND id <> ?", req.Email, teacherID).
			Count(&count)

		if count > 0 {
			response.CommonResp(c, 1, "邮箱已存在", nil)
			return
		}
	}

	// ===== phone 唯一校验 =====
	if req.Phone != "" && req.Phone != teacher.Phone {

		var count int64
		config.DB.Model(&models.Teacher{}).
			Where("phone = ? AND id <> ?", req.Phone, teacherID).
			Count(&count)

		if count > 0 {
			response.CommonResp(c, 1, "手机号已存在", nil)
			return
		}
	}

	// ===== 构造更新字段 =====
	updateData := utils.BuildUpdateMap(req)

	if sex, ok := updateData["sex"]; ok {
		if sex.(int8) != 0 && sex.(int8) != 1 {
			updateData["sex"] = 0
		}
	}

	if len(updateData) == 0 {
		response.CommonResp(c, 1, "没有可更新字段", nil)
		return
	}

	// ===== 更新 =====
	if err := config.DB.Model(&models.Teacher{}).
		Where("id = ?", teacherID).
		Updates(updateData).Error; err != nil {

		response.CommonResp(c, 1, "修改失败", nil)
		return
	}

	response.CommonResp(c, 0, "修改成功", nil)
}

func GetTeacherInformation(c *gin.Context) {

	// ===== Query 参数 =====
	var query struct {
		UserID int64 `form:"userID"`
	}

	_ = c.ShouldBindQuery(&query)

	// ===== 如果没传 userID -> 默认自己 =====
	if query.UserID == 0 {
		userIDAny, exists := c.Get("userID")
		if !exists {
			response.CommonResp(c, 1, "用户未登录", nil)
			return
		}
		query.UserID = userIDAny.(int64)
	}

	// ===== 查询教师 =====
	var teacher models.Teacher
	if err := config.DB.
		First(&teacher, query.UserID).Error; err != nil {

		response.CommonResp(c, 1, "教师不存在", nil)
		return
	}

	// ===== 转 VO =====
	teacherVO := converter.ToTeacherVO(teacher)

	// ===== 返回 =====
	response.CommonResp(c, 0, "获取成功", gin.H{
		"teacher": teacherVO,
	})
}

func PublishPost(c *gin.Context) {

	// ===== 1. 请求参数 =====
	var req struct {
		Title   string `json:"title" binding:"required"`
		Cover   string `json:"cover"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 2. 获取登录信息 =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "未登录", nil)
		return
	}

	roleAny, _ := c.Get("role")

	userID := userIDAny.(int)
	role := roleAny.(int)

	// ===== 3. 权限校验（仅老师允许发布）=====
	if role != 2 {
		response.CommonResp(c, 1, "只有教师可以发布帖子", nil)
		return
	}

	// ===== 4. 内容安全校验 =====
	if len(req.Title) > 200 {
		response.CommonResp(c, 1, "标题过长", nil)
		return
	}

	if len(req.Content) == 0 {
		response.CommonResp(c, 1, "内容不能为空", nil)
		return
	}

	// ===== 5. 构造帖子 =====
	post := models.Post{
		Title:    req.Title,
		Cover:    req.Cover,
		Content:  req.Content,
		AuthorID: userID,
	}

	// ===== 6. 入库 =====
	if err := config.DB.Create(&post).Error; err != nil {
		response.CommonResp(c, 1, "发布失败", nil)
		return
	}

	// ===== 7. 返回结果 =====
	response.CommonResp(c, 0, "发布成功", gin.H{
		"postID": post.ID,
	})
}

func GetPost(c *gin.Context) {

	// ===== 1. 获取参数 =====
	postIDStr := c.Query("postID")
	if postIDStr == "" {
		response.CommonResp(c, 1, "postID不能为空", nil)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		response.CommonResp(c, 1, "postID参数错误", nil)
		return
	}

	// ===== 2. 查询文章 + 作者 =====
	var post models.Post

	err = config.DB.
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&post, postID).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.CommonResp(c, 1, "文章不存在", nil)
			return
		}

		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ===== 3. 转换 VO =====
	postVO := converter.ToPostVO(post)

	// ===== 4. 返回 =====
	response.CommonResp(c, 0, "成功", gin.H{
		"post": postVO,
	})
}

func UpdatePost(c *gin.Context) {

	var req struct {
		PostID  int     `json:"postID" binding:"required"`
		Title   *string `json:"title"`
		Cover   *string `json:"cover"`
		Content *string `json:"content"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	updates := make(map[string]interface{})

	// ===== title =====
	if req.Title != nil {
		if *req.Title == "" {
			response.CommonResp(c, 1, "title不能为空", nil)
			return
		}
		updates["title"] = *req.Title
	}

	// ===== cover =====
	if req.Cover != nil {
		if *req.Cover == "" {
			response.CommonResp(c, 1, "cover不能为空", nil)
			return
		}
		updates["cover"] = *req.Cover
	}

	// ===== content =====
	if req.Content != nil {
		if *req.Content == "" {
			response.CommonResp(c, 1, "content不能为空", nil)
			return
		}
		updates["content"] = *req.Content
	}

	// 没有更新内容
	if len(updates) == 0 {
		response.CommonResp(c, 1, "没有需要更新的内容", nil)
		return
	}

	// ===== 登录信息（限制作者）=====
	userIDAny, _ := c.Get("userID")
	userID := userIDAny.(int)

	tx := config.DB.Model(&models.Post{}).
		Where("id = ? AND author_id = ?", req.PostID, userID).
		Updates(updates)

	if tx.Error != nil {
		response.CommonResp(c, 1, "修改失败", nil)
		return
	}

	if tx.RowsAffected == 0 {
		response.CommonResp(c, 1, "文章不存在或无权限", nil)
		return
	}

	response.CommonResp(c, 0, "修改成功", nil)
}

func DeletePost(c *gin.Context) {

	// ===== 1. 参数绑定 =====
	var req struct {
		PostID int `json:"postID" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 2. 获取登录信息 =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "未登录", nil)
		return
	}
	userID := userIDAny.(int)

	roleAny, _ := c.Get("role")
	role := roleAny.(int)

	// ===== 3. 构建删除条件 =====
	db := config.DB.Model(&models.Post{})

	// 管理员可以删除所有帖子
	if role == constants.RoleAdmin { // 假设常量 manager = 3
		db = db.Where("id = ?", req.PostID)
	} else {
		// 普通用户只能删除自己的帖子
		db = db.Where("id = ? AND author_id = ?", req.PostID, userID)
	}

	// ===== 4. 执行软删除 =====
	tx := db.Delete(&models.Post{})
	if tx.Error != nil {
		response.CommonResp(c, 1, "删除失败", nil)
		return
	}

	if tx.RowsAffected == 0 {
		response.CommonResp(c, 1, "文章不存在或无权限", nil)
		return
	}

	response.CommonResp(c, 0, "删除成功", nil)
}

func BrowsePosts(c *gin.Context) {
	var req struct {
		PageSize    int    `json:"pageSize" binding:"required,min=1"`
		PageNum     int    `json:"pageNum" binding:"required,min=1"`
		AuthorID    int    `json:"authorID"`
		SearchQuery string `json:"searchQuery"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	db := config.DB.Model(&models.Post{}).Preload("Author").Where("deleted_at IS NULL")

	// 按作者筛选
	if req.AuthorID > 0 {
		db = db.Where("author_id = ?", req.AuthorID)
	}

	// 搜索标题或内容
	if req.SearchQuery != "" {
		query := "%" + req.SearchQuery + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", query, query)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.CommonResp(c, 1, "查询总数失败", nil)
		return
	}

	// 分页
	offset := (req.PageNum - 1) * req.PageSize
	var posts []models.Post
	if err := db.Order("create_time DESC").
		Limit(req.PageSize).Offset(offset).
		Find(&posts).Error; err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// 转换 VO
	postVOs := converter.ToPostVOList(posts)

	// 返回结果
	response.CommonResp(c, 0, "查询成功", map[string]interface{}{
		"posts": postVOs,
		"total": total,
	})
}
