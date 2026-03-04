package utils

import "gopkg.in/gomail.v2"

// 163 邮箱账户信息
var hostEmail = "laibao2112@163.com"
var hostPassword = "FXvuLVHvyZGng543"

// SendEmail 发送验证码邮件（只负责发送）
func SendEmail(toEmail string, code string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", hostEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "心理健康管理平台：你的注册验证码到了！")

	m.SetBody("text/plain", "你的验证码是："+code+"，5分钟内有效。")

	d := gomail.NewDialer("smtp.163.com", 465, hostEmail, hostPassword)
	d.SSL = true

	return d.DialAndSend(m)
}
