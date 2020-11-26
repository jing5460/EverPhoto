package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/winterssy/EverPhotoCheckin/internal/client"
	"github.com/winterssy/EverPhotoCheckin/internal/push"
)

const (
	EnvMobile   = "EverPhotoMobile"
	EnvPassword = "EverPhotoPassword"
	EnvToken    = "EverPhotoToken"
	EnvSCKey    = "SCKey"
)

var (
	_mobile   string
	_password string
	_token    string
)

func init() {
	flag.StringVar(&_mobile, "mobile", "", "your mobile phone number")
	flag.StringVar(&_password, "password", "", "your password")
	flag.StringVar(&_token, "token", "", "your token")
}

func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func createBot() (bot *client.Bot, err error) {
	_token = valueOrDefault(_token, os.Getenv(EnvToken))
	if _token != "" {
		bot = client.NewWithToken(_token)
		return
	}

	_mobile = valueOrDefault(_mobile, os.Getenv(EnvMobile))
	_password = valueOrDefault(_password, os.Getenv(EnvPassword))
	bot, err = client.New(_mobile, _password)
	return
}

var _scKey = os.Getenv(EnvSCKey)

func pushMessage(ok bool, desp string) error {
	if _scKey == "" {
		return nil
	}

	const (
		_textSuccess = "【时光相册】签到成功通知"
		_textFailure = "【时光相册】签到失败通知"
	)
	if ok {
		return push.Push(_scKey, _textSuccess, desp)
	} else {
		return push.Push(_scKey, _textFailure, "错误详情 >> "+desp)
	}
}

func errDesp(msg string, err error) string {
	return msg + "：" + err.Error()
}

func main() {
	flag.Parse()

	var desp string
	bot, err := createBot()
	if err != nil {
		desp = errDesp("登录失败", err)
		_ = pushMessage(false, desp)
		log.Fatal("【时光相册】" + desp)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cr, err := bot.Checkin(ctx)
	if err != nil {
		desp = errDesp("签到失败", err)
		_ = pushMessage(false, desp)
		log.Fatal("【时光相册】" + desp)
	}

	desp = fmt.Sprintf("你已连续签到%d天，累计获得空间%s，明天可白嫖%s，请继续保持(￣▽￣)",
		cr.Continuity, cr.TotalReward, cr.TomorrowReward)
	_ = pushMessage(true, desp)
	log.Print("【时光相册】" + desp)
}
