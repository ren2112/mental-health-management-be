package controllers

import (
	"github.com/gin-gonic/gin"
	"mental-health-management-be/config"
	"mental-health-management-be/constants"
	"mental-health-management-be/converter"
	"mental-health-management-be/models"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
)

func UpdateStudentSelfInformation(c *gin.Context) {

	// ===== 身份校验 =====
	roleAny, _ := c.Get("role")
	role := roleAny.(int)

	if role != constants.RoleStudent {
		response.CommonResp(c, 1, "身份异常", nil)
		return
	}

	// ===== 从 JWT 获取当前用户 =====
	userIDAny, _ := c.Get("userID")
	studentID := userIDAny.(int)

	var req struct {
		Name      string `json:"name"`
		Sex       *int8  `json:"sex"`
		StudentNo string `json:"studentNo"`
		Email     string `json:"email" binding:"omitempty,email"`
	}

	// ===== 参数绑定 =====
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 查询当前学生 =====
	var student models.Student
	if err := config.DB.First(&student, studentID).Error; err != nil {
		response.CommonResp(c, 1, "学生不存在", nil)
		return
	}

	// ===== email 唯一校验 =====
	if req.Email != "" && req.Email != student.Email {

		var count int64
		config.DB.Model(&models.Student{}).
			Where("email = ? AND id <> ?", req.Email, studentID).
			Count(&count)

		if count > 0 {
			response.CommonResp(c, 1, "邮箱已存在", nil)
			return
		}
	}

	// ===== 学号唯一校验 =====
	if req.StudentNo != "" && req.StudentNo != student.StudentNo {

		var count int64
		config.DB.Model(&models.Student{}).
			Where("student_no = ? AND id <> ?", req.StudentNo, studentID).
			Count(&count)

		if count > 0 {
			response.CommonResp(c, 1, "学号已存在", nil)
			return
		}
	}

	// ===== 构造更新字段 =====
	updateData := utils.BuildUpdateMap(req)

	// sex 合法化
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
	if err := config.DB.Model(&models.Student{}).
		Where("id = ?", studentID).
		Updates(updateData).Error; err != nil {

		response.CommonResp(c, 1, "修改失败", nil)
		return
	}

	response.CommonResp(c, 0, "修改成功", nil)
}

func GetStudentInformation(c *gin.Context) {

	// ===== Query 参数绑定 =====
	var query struct {
		UserID int64 `form:"userID" binding:"required"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		response.CommonResp(c, 1, "userID参数错误", nil)
		return
	}

	// ===== 查询数据库 =====
	var student models.Student
	if err := config.DB.
		First(&student, query.UserID).Error; err != nil {

		response.CommonResp(c, 1, "学生不存在", nil)
		return
	}

	// ===== 转换 VO =====
	studentVO := converter.ToStudentVO(student)

	// ===== 返回 =====
	response.CommonResp(c, 0, "获取成功", gin.H{
		"student": studentVO,
	})
}
