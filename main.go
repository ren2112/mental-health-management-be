package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mental-health-management-be/config"
	"mental-health-management-be/controllers"
	"mental-health-management-be/middleware"
)

const gport = 8080

func main() {
	fmt.Println("start main...")
	config.InitDB()
	config.InitRedis()

	r := gin.Default()
	r.Static("/public", "./public")

	noAuth := r.Group("/api/noauth")
	{
		noAuth.POST("/validate-code", controllers.GetValidateCode)
		noAuth.POST("/regist", controllers.Register)
		noAuth.POST("/login", controllers.Login)
	}
	auth := r.Group("/api/auth")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/upd-stu-self-information", controllers.UpdateStudentSelfInformation)
		auth.POST("/upd-teach-self-information", controllers.UpdateTeacherSelfInformation)
		auth.POST("/upd-manager-self-information", controllers.UpdateManagerSelfInformation)

		auth.GET("/get-stu-information", controllers.GetStudentInformation)
		auth.GET("/get-teach-information", controllers.GetTeacherInformation)
		auth.GET("/manager-information", controllers.GetManagerInformation)
		auth.POST("/update-password", controllers.UpdatePassword)

		auth.POST("/upload-cover", controllers.UploadCover)
		auth.POST("/publish-post", controllers.PublishPost)
		auth.GET("/post", controllers.GetPost)
		auth.POST("/update-post", controllers.UpdatePost)
		auth.POST("/delete-post", controllers.DeletePost)
		auth.POST("/browse", controllers.BrowsePosts)

		auth.POST("/appoint", controllers.CreateAppointment)
		auth.POST("/withdraw-appointment", controllers.WithdrawAppointment)
		auth.POST("/approve-appointment", controllers.ApproveAppointment)
		auth.POST("/self-appointments", controllers.GetSelfAppointments)
		auth.GET("/get-teacher-appointment", controllers.GetTeacherAppointment)
	}

	manage := r.Group("/api/auth/manage")
	manage.Use(middleware.JWTAuthMiddleware(), middleware.AdminAuth())
	{
		manage.POST("/student-list", controllers.StudentList)
		manage.POST("/teacher-list", controllers.TeacherList)
		manage.POST("/manager-list", controllers.ManagerList)
		manage.POST("/del-student", controllers.DelStudent) // 删除学生
		manage.POST("/del-teacher", controllers.DelTeacher) // 删除教师
		manage.POST("/del-manager", controllers.DelManager) // 删除管理员
		manage.POST("/add-teacher", controllers.AddTeacher)
		manage.POST("/add-student", controllers.AddStudent)
		manage.POST("/add-manager", controllers.AddManager)
	}

	addr := fmt.Sprintf("localhost:%d", gport)
	r.Run(addr)
}
