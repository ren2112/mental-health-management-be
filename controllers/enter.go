package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mental-health-management-be/config"
	"mental-health-management-be/constants"
	"mental-health-management-be/converter"
	"mental-health-management-be/models"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
	"time"
)

func GetValidateCode(c *gin.Context) {

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	ctx := config.Ctx

	// ===============================
	// 1️⃣ 防刷限制（60秒）
	// ===============================
	limitKey := "validate_limit:" + req.Email

	exist, err := config.RDB.Exists(ctx, limitKey).Result()
	if err != nil {
		response.CommonResp(c, 1, "系统错误", nil)
		return
	}

	if exist == 1 {
		response.CommonResp(c, 1, "请求过于频繁，请60秒后再试", nil)
		return
	}

	// ===============================
	// 2️⃣ 生成验证码
	// ===============================
	code := utils.GenerateCode()

	// ===============================
	// 3️⃣ 存 Redis（5分钟有效）
	// ===============================
	codeKey := "validate_code:" + req.Email

	err = config.RDB.Set(
		ctx,
		codeKey,
		code,
		5*time.Minute,
	).Err()

	if err != nil {
		response.CommonResp(c, 1, "验证码保存失败", nil)
		return
	}

	// ===============================
	// 4️⃣ 设置防刷锁（60秒）
	// ===============================
	config.RDB.Set(ctx, limitKey, 1, time.Minute)

	// ===============================
	// 5️⃣ 发送邮件
	// ===============================
	err = utils.SendEmail(req.Email, code)
	if err != nil {
		response.CommonResp(c, 1, "邮件发送失败", nil)
		return
	}

	// ===============================
	// 6️⃣ 返回
	// ===============================
	response.CommonResp(c, 0, "验证码发送成功", nil)
}

func Register(c *gin.Context) {

	var req struct {
		Email        string `json:"email" binding:"required,email"`
		Password     string `json:"password" binding:"required"`
		ValidateCode string `json:"validateCode" binding:"required"`
		Name         string `json:"name"  binding:"required"`
		Sex          int8   `json:"sex" binding:"required,oneof=0 1"`
		StudentNo    string `json:"studentNo"  binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	key := "validate_code:" + req.Email

	// 1. 取redis验证码
	code, err := config.RDB.Get(config.Ctx, key).Result()
	if err != nil {
		response.CommonResp(c, 1, "验证码不存在或过期", nil)
		return
	}

	// 2. 校验
	if code != req.ValidateCode {
		response.CommonResp(c, 1, "验证码错误", nil)
		return
	}

	// 3. 删除验证码（⭐必须）
	config.RDB.Del(config.Ctx, key)

	// 4. 密码加密
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	student := models.Student{
		Name:       req.Name,
		Sex:        req.Sex,
		StudentNo:  req.StudentNo,
		Email:      req.Email,
		Password:   string(hash),
		CreateTime: time.Now(),
	}

	// 5. 入库
	if err := config.DB.Create(&student).Error; err != nil {
		response.CommonResp(c, 1, "注册失败，邮箱可能已存在", nil)
		return
	}

	response.CommonResp(c, 0, "注册成功", nil)
}

func Login(c *gin.Context) {

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Type     int    `json:"type" binding:"required"` // 0管理员 1学生 2老师
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	var (
		userID  int
		hashPwd string
		userVO  interface{} // ⭐ 返回不同类型
		err     error
	)

	// ⭐ 根据身份查询不同表
	switch req.Type {

	case constants.RoleStudent:

		var student models.Student
		err = config.DB.Where("email = ?", req.Email).
			First(&student).Error

		if err == nil {
			userID = student.ID
			hashPwd = student.Password
			userVO = converter.ToStudentVO(student) // ⭐ 转换VO
		}

	case constants.RoleTeacher:

		var teacher models.Teacher
		err = config.DB.Where("email = ?", req.Email).
			First(&teacher).Error

		if err == nil {
			userID = teacher.ID
			hashPwd = teacher.Password
			userVO = converter.ToTeacherVO(teacher)
		}

	case constants.RoleAdmin:

		var manager models.Manager
		err = config.DB.Where("email = ?", req.Email).
			First(&manager).Error

		if err == nil {
			userID = manager.ID
			hashPwd = manager.Password
			userVO = converter.ToManagerVO(manager)
		}

	default:
		response.CommonResp(c, 1, "用户类型错误", nil)
		return
	}

	// 用户不存在
	if err != nil {
		response.CommonResp(c, 1, "用户不存在", nil)
		return
	}

	// ⭐ 校验密码
	if bcrypt.CompareHashAndPassword(
		[]byte(hashPwd),
		[]byte(req.Password),
	) != nil {
		response.CommonResp(c, 1, "密码错误", nil)
		return
	}

	// ⭐ 生成JWT
	token, err := utils.GenerateJWT(userID, req.Type)
	if err != nil {
		response.CommonResp(c, 1, "生成Token失败", nil)
		return
	}

	// ⭐ 返回 token + user
	response.CommonResp(c, 0, "登录成功", gin.H{
		"token": token,
		"user":  userVO,
		"type":  req.Type,
	})
}

func UpdatePassword(c *gin.Context) {

	// ===== 1. 请求参数 =====
	var req struct {
		NewPassword  string `json:"newPassword" binding:"required"`
		ValidateCode string `json:"validateCode" binding:"required"`
		Email        string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 2. 获取 token 信息 =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "未登录", nil)
		return
	}

	roleAny, _ := c.Get("role")

	userID := userIDAny.(int)
	role := roleAny.(int)

	// ===== 3. 校验 Redis 验证码 =====
	key := "validate_code:" + req.Email

	code, err := config.RDB.Get(config.Ctx, key).Result()
	if err != nil {
		response.CommonResp(c, 1, "验证码不存在或已过期", nil)
		return
	}

	if code != req.ValidateCode {
		response.CommonResp(c, 1, "验证码错误", nil)
		return
	}

	// ===== 4. 新密码加密 =====
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		response.CommonResp(c, 1, "密码加密失败", nil)
		return
	}

	// ===== 5. 根据角色更新密码（ID + Email 双条件）=====
	var result *gorm.DB

	switch role {

	case 1:
		result = config.DB.Model(&models.Student{}).
			Where("id = ? AND email = ?", userID, req.Email).
			Update("password", string(hash))

	case 2:
		result = config.DB.Model(&models.Teacher{}).
			Where("id = ? AND email = ?", userID, req.Email).
			Update("password", string(hash))

	case 3:
		result = config.DB.Model(&models.Manager{}).
			Where("id = ? AND email = ?", userID, req.Email).
			Update("password", string(hash))

	default:
		response.CommonResp(c, 1, "非法角色", nil)
		return
	}

	// ===== SQL执行错误 =====
	if result.Error != nil {
		response.CommonResp(c, 1, "密码修改失败", nil)
		return
	}

	// ===== 没匹配到用户（重点）=====
	if result.RowsAffected == 0 {
		response.CommonResp(c, 1, "用户不存在或邮箱不匹配", nil)
		return
	}

	// ===== 6. 删除验证码 =====
	config.RDB.Del(config.Ctx, key)

	response.CommonResp(c, 0, "密码修改成功", nil)
}
