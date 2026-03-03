package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"mental-health-management-be/config"
	"mental-health-management-be/constants"
	"mental-health-management-be/converter"
	"mental-health-management-be/models"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
	"time"
)

func StudentList(c *gin.Context) {

	var req struct {
		StudentName string `json:"studentName" `
		StudentNo   string `json:"studentNo"`
		Email       string `json:"email"`
		PageSize    int    `json:"pageSize" binding:"required"`
		PageNum     int    `json:"pageNum"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 分页 =====
	pageSize := 10
	pageNum := 1

	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	if req.PageNum > 0 {
		pageNum = req.PageNum
	}

	offset := (pageNum - 1) * pageSize

	// ===== 查询 =====
	db := config.DB.Model(&models.Student{})

	// ⭐ 多条件模糊查询（不再需要 searchType）
	if req.StudentName != "" {
		db = db.Where("name LIKE ?", "%"+req.StudentName+"%")
	}
	if req.StudentNo != "" {
		db = db.Where("student_no LIKE ?", "%"+req.StudentNo+"%")
	}
	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}

	// ===== 总数 =====
	var total int64
	db.Count(&total)

	// ===== 数据 =====
	var students []models.Student
	err := db.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&students).Error

	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ⭐ 转 VO
	studentVOs := converter.ToStudentVOList(students)

	// ===== 返回 =====
	response.CommonResp(c, 0, "成功", gin.H{
		"students": studentVOs,
		"total":    total,
	})
}

func TeacherList(c *gin.Context) {

	var req struct {
		TeacherName string `json:"teacherName"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		PageSize    int    `json:"pageSize" binding:"required"`
		PageNum     int    `json:"pageNum"`
	}

	// ===== JSON 绑定 =====
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 分页保护 =====
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	offset := (req.PageNum - 1) * req.PageSize

	// ===== 构建查询 =====
	db := config.DB.Model(&models.Teacher{})

	// ⭐ 条件查询（非空才拼接）
	if req.TeacherName != "" {
		db = db.Where("name LIKE ?", "%"+req.TeacherName+"%")
	}

	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}

	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	// ===== 查询总数 =====
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ===== 查询列表 =====
	var teachers []models.Teacher

	err := db.
		Order("id DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&teachers).Error

	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ===== 转 VO =====
	teacherVOs := converter.ToTeacherVOList(teachers)

	// ===== 返回 =====
	response.CommonResp(c, 0, "成功", gin.H{
		"teachers": teacherVOs,
		"total":    total,
	})
}

func ManagerList(c *gin.Context) {

	// 请求参数结构体
	var req struct {
		ManagerName string `json:"managerName"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		PageSize    int    `json:"pageSize"`
		PageNum     int    `json:"pageNum"`
	}

	// 绑定请求数据
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// 分页保护
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	offset := (req.PageNum - 1) * req.PageSize

	// 构建查询
	db := config.DB.Model(&models.Manager{})

	// 条件查询
	if req.ManagerName != "" {
		db = db.Where("name LIKE ?", "%"+req.ManagerName+"%")
	}
	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	// 查询总数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// 查询数据
	var managers []models.Manager
	err = db.
		Order("id DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&managers).Error

	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// 转换为 VO
	managerVOs := converter.ToManagerVOList(managers)

	// 返回数据
	response.CommonResp(c, 0, "成功", gin.H{
		"managers": managerVOs,
		"total":    total,
	})
}

// 软删除学生
func DelStudent(c *gin.Context) {
	var req struct {
		StudentID int `json:"studentID" binding:"required"`
	}

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// 查找学生
	var student models.Student
	err := config.DB.Where("id = ? AND deleted_at IS NULL", req.StudentID).First(&student).Error
	if err != nil {
		response.CommonResp(c, 1, "学生不存在或已被删除", nil)
		return
	}

	// 软删除操作
	err = config.DB.Model(&student).UpdateColumn("deleted_at", time.Now()).Error
	if err != nil {
		response.CommonResp(c, 1, "删除失败", nil)
		return
	}

	response.CommonResp(c, 0, "删除成功", nil)
}

// 软删除教师
func DelTeacher(c *gin.Context) {
	var req struct {
		TeacherID int `json:"teacherID" binding:"required"`
	}

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// 查找教师
	var teacher models.Teacher
	err := config.DB.Where("id = ? AND deleted_at IS NULL", req.TeacherID).First(&teacher).Error
	if err != nil {
		response.CommonResp(c, 1, "教师不存在或已被删除", nil)
		return
	}

	// 软删除操作
	err = config.DB.Model(&teacher).UpdateColumn("deleted_at", time.Now()).Error
	if err != nil {
		response.CommonResp(c, 1, "删除失败", nil)
		return
	}

	response.CommonResp(c, 0, "删除成功", nil)
}

// 软删除管理员
func DelManager(c *gin.Context) {
	var req struct {
		ManagerID int `json:"managerID" binding:"required"`
	}

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// 获取当前登录的管理员ID（从JWT中获取）
	currentManagerIDAny, _ := c.Get("userID")
	currentManagerID, ok := currentManagerIDAny.(int)
	if !ok {
		response.CommonResp(c, 1, "权限异常", nil)
		return
	}

	// 如果尝试删除自己
	if req.ManagerID == currentManagerID {
		response.CommonResp(c, 1, "不能删除自己", nil)
		return
	}

	// 查找管理员
	var manager models.Manager
	err := config.DB.Where("id = ? AND deleted_at IS NULL", req.ManagerID).First(&manager).Error
	if err != nil {
		response.CommonResp(c, 1, "管理员不存在或已被删除", nil)
		return
	}

	// 软删除操作
	err = config.DB.Model(&manager).UpdateColumn("deleted_at", time.Now()).Error
	if err != nil {
		response.CommonResp(c, 1, "删除失败", nil)
		return
	}

	response.CommonResp(c, 0, "删除成功", nil)
}

func AddTeacher(c *gin.Context) {

	var req struct {
		Name         string `json:"name" binding:"required"`
		Sex          int8   `json:"sex"`
		Email        string `json:"email" binding:"required,email"`
		Introduction string `json:"introduction"`
		Phone        string `json:"phone"`
		Workspace    string `json:"workspace"`
	}

	// ===== 参数绑定 =====
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== sex 校验（可选字段）=====
	if req.Sex != 0 && req.Sex != 1 {
		req.Sex = 0
	}

	// ===== 邮箱是否已存在 =====
	var count int64
	config.DB.Model(&models.Teacher{}).
		Where("email = ?", req.Email).
		Count(&count)

	if count > 0 {
		response.CommonResp(c, 1, "邮箱已存在", nil)
		return
	}

	// ===== 默认密码加密 =====
	hashPwd, err := bcrypt.GenerateFromPassword(
		[]byte(constants.DefaultUserPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		response.CommonResp(c, 1, "密码生成失败", nil)
		return
	}

	// ===== 创建教师 =====
	teacher := models.Teacher{
		Name:         req.Name,
		Sex:          req.Sex,
		Email:        req.Email,
		Introduction: req.Introduction,
		Phone:        req.Phone,
		Workspace:    req.Workspace,
		Password:     string(hashPwd),
	}

	if err := config.DB.Create(&teacher).Error; err != nil {
		response.CommonResp(c, 1, "创建失败", nil)
		return
	}

	response.CommonResp(c, 0, "教师创建成功", nil)
}

func AddStudent(c *gin.Context) {

	var req struct {
		Name      string `json:"name" binding:"required"`
		Sex       int8   `json:"sex"`
		Email     string `json:"email" binding:"required,email"`
		StudentNo string `json:"studentNo"`
	}

	// ===== 参数绑定 =====
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== sex 限制 =====
	if req.Sex != 0 && req.Sex != 1 {
		req.Sex = 0
	}

	// ===== 邮箱唯一校验 =====
	var count int64
	config.DB.Model(&models.Student{}).
		Where("email = ?", req.Email).
		Count(&count)

	if count > 0 {
		response.CommonResp(c, 1, "邮箱已存在", nil)
		return
	}

	// ===== 学号唯一校验（如果传了）=====
	if req.StudentNo != "" {
		config.DB.Model(&models.Student{}).
			Where("student_no = ?", req.StudentNo).
			Count(&count)

		if count > 0 {
			response.CommonResp(c, 1, "学号已存在", nil)
			return
		}
	}

	// ===== 默认密码加密 =====
	hashPwd, err := bcrypt.GenerateFromPassword(
		[]byte(constants.DefaultUserPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		response.CommonResp(c, 1, "密码生成失败", nil)
		return
	}

	// ===== 创建学生 =====
	student := models.Student{
		Name:      req.Name,
		Sex:       req.Sex,
		Email:     req.Email,
		StudentNo: req.StudentNo,
		Password:  string(hashPwd),
	}

	if err := config.DB.Create(&student).Error; err != nil {
		response.CommonResp(c, 1, "创建失败", nil)
		return
	}

	response.CommonResp(c, 0, "学生创建成功", nil)
}

func AddManager(ctx *gin.Context) {

	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone"`
	}

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.CommonResp(ctx, 1, "参数错误", nil)
		return
	}

	// ===== 校验邮箱是否存在 =====
	var count int64
	if err := config.DB.Model(&models.Manager{}).
		Where("email = ?", req.Email).
		Count(&count).Error; err != nil {

		response.CommonResp(ctx, 1, "查询失败", nil)
		return
	}

	if count > 0 {
		response.CommonResp(ctx, 1, "邮箱已存在", nil)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(constants.DefaultUserPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		response.CommonResp(ctx, 1, "密码加密失败", nil)
		return
	}

	// ===== 创建管理员 =====
	manager := models.Manager{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hash),
	}

	if err := config.DB.Create(&manager).Error; err != nil {
		response.CommonResp(ctx, 1, "添加管理员失败", nil)
		return
	}

	response.CommonResp(ctx, 0, "添加成功", nil)
}

func UpdateManagerSelfInformation(c *gin.Context) {

	// ===== 身份校验 =====
	roleAny, _ := c.Get("role")
	role := roleAny.(int)

	if role != constants.RoleAdmin {
		response.CommonResp(c, 1, "身份异常", nil)
		return
	}

	// ===== 当前登录用户ID =====
	userIDAny, _ := c.Get("userID")
	managerID := userIDAny.(int)

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"omitempty,email"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 自动生成更新字段 =====
	updateData := utils.BuildUpdateMap(req)

	if len(updateData) == 0 {
		response.CommonResp(c, 1, "没有可更新字段", nil)
		return
	}

	// ===== 更新 =====
	if err := config.DB.Model(&models.Manager{}).
		Where("id = ?", managerID).
		Updates(updateData).Error; err != nil {

		response.CommonResp(c, 1, "修改失败", nil)
		return
	}

	response.CommonResp(c, 0, "修改成功", nil)
}

func GetManagerInformation(c *gin.Context) {

	// ===== Query 参数 =====
	var query struct {
		UserID int `form:"userID"`
	}

	_ = c.ShouldBindQuery(&query)

	// ===== 如果没传 userID -> 默认自己 =====
	if query.UserID == 0 {
		userIDAny, exists := c.Get("userID")
		if !exists {
			response.CommonResp(c, 1, "用户未登录", nil)
			return
		}
		query.UserID = userIDAny.(int)
	}

	// ===== 查询管理员 =====
	var manager models.Manager
	if err := config.DB.
		First(&manager, query.UserID).Error; err != nil {

		response.CommonResp(c, 1, "管理员不存在", nil)
		return
	}

	// ===== 转 VO =====
	managerVO := converter.ToManagerVO(manager)

	// ===== 返回 =====
	response.CommonResp(c, 0, "获取成功", gin.H{
		"manager": managerVO,
	})
}
