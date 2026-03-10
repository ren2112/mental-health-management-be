package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"mental-health-management-be/config"
	"mental-health-management-be/constants"
	"mental-health-management-be/converter"
	"mental-health-management-be/models"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
	"strings"
	"time"
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

func CreateAppointment(c *gin.Context) {

	// ===== 1. 请求参数 =====
	var req struct {
		TeacherID int    `json:"teacherID" binding:"required"`
		Title     string `json:"title" binding:"required"`
		Detail    string `json:"detail" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 2. 登录信息 =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "未登录", nil)
		return
	}

	roleAny, _ := c.Get("role")

	userID := userIDAny.(int)
	role := roleAny.(int)

	// ===== 3. 权限校验 =====
	if role != constants.RoleStudent {
		response.CommonResp(c, 1, "只有学生可以发起预约", nil)
		return
	}

	// ===== 4. 时间转换（毫秒时间戳）=====
	start := time.UnixMilli(req.StartTime).Local()
	end := time.UnixMilli(req.EndTime).Local()

	// ===== 5. 基础校验 =====
	if !start.Before(end) {
		response.CommonResp(c, 1, "结束时间必须晚于开始时间", nil)
		return
	}

	// ===== 6. 时间合法性 =====
	// start/end 必须是 slot 起点
	if !utils.IsValidSlotStart(start) ||
		!utils.IsValidSlotStart(end) {
		response.CommonResp(c, 1, "时间必须为整点或半点", nil)
		return
	}

	// 覆盖范围必须在工作时间
	if !utils.IsWithinWorkTime(start) ||
		!utils.IsWithinWorkTime(end.Add(-time.Second)) {
		response.CommonResp(c, 1, "仅允许08:00-17:00预约", nil)
		return
	}

	if !utils.IsValidAppointmentDate(start) ||
		!utils.IsValidAppointmentDate(end.Add(-time.Second)) {
		response.CommonResp(c, 1, "只能预约当前周和下周工作日", nil)
		return
	}

	// ===== 7. 生成 slots =====
	slots := utils.GenerateSlots(start, end)

	var appointmentID int

	// ====================================
	// ===== 8. 事务（核心并发安全） =====
	// ====================================
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ 创建预约主记录
		appoint := models.Appointment{
			Title:     req.Title,
			Detail:    req.Detail,
			StudentID: userID,
			TeacherID: req.TeacherID,
			Status:    0,
			StartTime: start,
			EndTime:   end,
		}

		if err := tx.Create(&appoint).Error; err != nil {
			return err
		}

		// 2️⃣ 插入 slots（⭐ 冲突检测发生在这里）
		for _, s := range slots {

			slot := models.AppointmentSlot{
				TeacherID:     req.TeacherID,
				SlotTime:      s,
				AppointmentID: appoint.ID,
			}

			if err := tx.Create(&slot).Error; err != nil {

				// MySQL 唯一键冲突
				if strings.Contains(err.Error(), "Duplicate") {
					return errors.New("时间段已被预约")
				}

				return err
			}
		}

		appointmentID = appoint.ID
		return nil
	})

	// ===== 9. 事务结果 =====
	if err != nil {

		if err.Error() == "时间段已被预约" {
			response.CommonResp(c, 1, "该时间段已被预约", nil)
			return
		}

		response.CommonResp(c, 1, "预约失败", nil)
		return
	}

	// ===== 10. 返回 =====
	response.CommonResp(c, 0, "预约成功", gin.H{
		"appointmentID": appointmentID,
	})
}

func WithdrawAppointment(c *gin.Context) {

	// ===== 1. 请求参数 =====
	var req struct {
		AppointmentID int `json:"appointmentID" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 2. 登录信息 =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "未登录", nil)
		return
	}

	roleAny, _ := c.Get("role")

	userID := userIDAny.(int)
	role := roleAny.(int)

	if role != constants.RoleStudent {
		response.CommonResp(c, 1, "只有学生可以撤销预约", nil)
		return
	}

	// ==================================================
	// ⭐ 事务：保证 slot + appointment 一致删除
	// ==================================================
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		var appoint models.Appointment

		// ===== 3. 行锁（防并发撤销/审核）=====
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&appoint, req.AppointmentID).Error; err != nil {
			return errors.New("预约不存在")
		}

		// ===== 4. 权限校验 =====
		if appoint.StudentID != userID {
			return errors.New("无权限操作该预约")
		}

		// ===== 5. 已删除检查 =====
		if appoint.DeletedAt.Valid {
			return errors.New("预约已撤销")
		}

		// ===== 6. 删除 Slot（释放时间资源）=====
		if err := tx.
			Where("appointment_id = ?", appoint.ID).
			Delete(&models.AppointmentSlot{}).Error; err != nil {
			return err
		}

		// ===== 7. 删除 Appointment（软删除）=====
		if err := tx.Delete(&appoint).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		response.CommonResp(c, 1, err.Error(), nil)
		return
	}

	response.CommonResp(c, 0, "撤销成功", nil)
}

func GetSelfAppointments(c *gin.Context) {

	// ===== 请求参数 =====
	var req struct {
		PageSize      int  `json:"pageSize"`
		PageNum       int  `json:"pageNum"`
		ApproveStatus *int `json:"approveStatus"`
	}

	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	// ===== 获取当前学生ID =====
	userIDAny, exists := c.Get("userID")
	if !exists {
		response.CommonResp(c, 1, "用户未登录", nil)
		return
	}
	studentID := userIDAny.(int)

	db := config.DB

	query := db.Model(&models.Appointment{}).
		Where("student_id = ?", studentID)

	// ===== 状态筛选（可选）=====
	if req.ApproveStatus != nil {
		query = query.Where("status = ?", *req.ApproveStatus)
	}

	// ===== 总数 =====
	var total int64
	query.Count(&total)

	// ===== 分页 =====
	offset := (req.PageNum - 1) * req.PageSize

	var appointments []models.Appointment
	err := query.
		Preload("Teacher").
		Order("start_time desc").
		Limit(req.PageSize).
		Offset(offset).
		Find(&appointments).Error

	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ===== VO转换 =====
	voList := converter.ToAppointmentVOList(appointments)

	response.CommonResp(c, 0, "success", gin.H{
		"appointments": voList,
		"total":        total,
	})
}

func GetTeacherAppointment(c *gin.Context) {

	// ===== Query 参数 =====
	var query struct {
		TeacherID int64 `form:"teacherID"`
	}

	_ = c.ShouldBindQuery(&query)

	// ===== 如果没传 teacherID -> 默认自己 =====
	if query.TeacherID == 0 {
		userIDAny, exists := c.Get("userID")
		if !exists {
			response.CommonResp(c, 1, "用户未登录", nil)
			return
		}
		query.TeacherID = userIDAny.(int64)
	}

	// =============================
	// 计算 本周 + 下周 时间范围
	// =============================
	now := time.Now()

	// Go中 Sunday=0
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// 本周一
	startOfWeek := time.Date(
		now.Year(),
		now.Month(),
		now.Day()-(weekday-1),
		0, 0, 0, 0,
		now.Location(),
	)

	// 下下周一（结束时间）
	endTime := startOfWeek.AddDate(0, 0, 14)

	db := config.DB

	var appointments []models.Appointment

	err := db.Model(&models.Appointment{}).
		Where("teacher_id = ?", query.TeacherID).
		Where("start_time >= ? AND start_time < ?", startOfWeek, endTime).
		Where("status IN ?", []int{0, 1}).
		Order("start_time asc").
		Find(&appointments).Error

	if err != nil {
		response.CommonResp(c, 1, "查询失败", nil)
		return
	}

	// ===== VO转换 =====
	voList := converter.ToAppointmentVOList(appointments)

	response.CommonResp(c, 0, "success", gin.H{
		"appointments": voList,
	})
}

func AIBot(c *gin.Context) {

	// ===== 1. 获取参数 =====
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.CommonResp(c, 1, "参数错误", nil)
		return
	}

	content := req.Content
	if content == "" {
		response.CommonResp(c, 1, "content不能为空", nil)
		return
	}

	// ===== 2. 提示词拼接 =====
	prompt := fmt.Sprintf(`
你是一名专业的大学生心理健康助手，请用温和、共情、积极的语气回复学生的问题。

要求：
1. 给予情绪理解
2. 提供积极建议
3. 回答简洁清晰
4. 如果问题涉及严重心理问题，建议寻求专业心理咨询师

学生的问题：
%s
`, content)

	// ===== 3. 创建AI客户端 =====
	configAI := openai.DefaultConfig("ms-bd647c5f-813c-4eec-abf1-02bf94c30ca7")
	configAI.BaseURL = "https://api-inference.modelscope.cn/v1"

	client := openai.NewClientWithConfig(configAI)

	// ===== 4. 创建流式请求 =====
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "Qwen/Qwen3.5-27B",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		response.CommonResp(c, 1, "AI服务调用失败", nil)
		return
	}

	defer stream.Close()

	// ===== 5. 设置 SSE 流式返回 =====
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Writer.Flush()

	// ===== 6. 循环读取AI返回 =====
	for {

		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		content := resp.Choices[0].Delta.Content
		if content == "" {
			continue
		}

		// 构造返回
		respData := gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"response": content,
			},
		}

		jsonData, _ := json.Marshal(respData)

		fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)

		c.Writer.Flush()
	}
}
