package util

import (
	"time"
)

var jst, _ = time.LoadLocation("Asia/Tokyo")

// ConvertJSTStringTime : timeをJSTタイムゾーンに変換したのち、timeをyyyy-MM-dd HH:mm:ss形式の文字列に変換
func ConvertJSTStringTime(t time.Time) string {
	jstTime := JST(t)
	return ConvertStringTime(jstTime)
}

// ConvertStringTime : timeをyyyy-MM-dd HH:mm:ss形式の文字列に変換
func ConvertStringTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ConvertJSTStringDate : timeをJSTタイムゾーンに変換したのち、yyyy-MM-dd形式の文字列に変換
func ConvertJSTStringDate(t time.Time) string {
	jstTime := JST(t)
	return ConvertStringDate(jstTime)
}

// ConvertStringDate : timeをyyyy-MM-dd形式の文字列に変換
func ConvertStringDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func JST(t time.Time) time.Time {
	return t.In(jst)
}
