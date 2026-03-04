package utils

import "time"

func IsValidSlotStart(t time.Time) bool {
	if t.Second() != 0 || t.Nanosecond() != 0 {
		return false
	}

	if !(t.Minute() == 0 || t.Minute() == 30) {
		return false
	}

	h := t.Hour()
	return h >= 8 && h <= 17
}

func IsValidAppointmentDate(t time.Time) bool {

	now := time.Now().In(t.Location())

	// ===============================
	// 1️⃣ 不能早于当前时间
	// ===============================
	if !t.After(now) {
		return false
	}

	// ===============================
	// 2️⃣ 计算本周一
	// Go: Sunday = 0
	// ===============================
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	thisMonday := time.Date(
		now.Year(),
		now.Month(),
		now.Day()-weekday+1,
		0, 0, 0, 0,
		now.Location(),
	)

	// ===============================
	// 3️⃣ 下周日 23:59:59
	// ===============================
	nextSunday := thisMonday.
		AddDate(0, 0, 13).
		Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// ===============================
	// 4️⃣ 是否在允许范围
	// ===============================
	if t.Before(thisMonday) || t.After(nextSunday) {
		return false
	}

	// ===============================
	// 5️⃣ 必须工作日
	// ===============================
	w := t.Weekday()
	if w == time.Saturday || w == time.Sunday {
		return false
	}

	return true
}
func IsWithinWorkTime(t time.Time) bool {
	h := t.Hour()
	m := t.Minute()

	total := h*60 + m
	return total >= 8*60 && total < 17*60
}
