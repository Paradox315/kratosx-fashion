package middleware

import (
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
	"sync"
)

var _ kmw.FiberMiddleware = (*Logger)(nil)

const (
	sysPrefix      = "/api/system/v1/"
	pubPrefix      = "pub/"
	userPrefix     = "user/"
	rolePrefix     = "role/"
	resourcePrefix = "resource/"
)

type Logger struct {
	once     sync.Once
	logRepo  biz.UserLogRepo
	eventMap map[string]string
}

func NewLoggerHook(logRepo biz.UserLogRepo) *Logger {
	l := &Logger{
		logRepo: logRepo,
		eventMap: map[string]string{
			sysPrefix + pubPrefix + "register":        "用户注册",
			sysPrefix + pubPrefix + "login":           "用户登录",
			sysPrefix + pubPrefix + "logout":          "用户退出",
			sysPrefix + pubPrefix + "refresh-token":   "刷新token",
			sysPrefix + userPrefix + "reset-password": "重置密码",
			sysPrefix + userPrefix + "password":       "修改密码",
		},
	}
	l.once.Do(func() {
		kmw.RegisterMiddleware(l)
	})
	return l
}

func (l *Logger) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			errHandler fiber.ErrorHandler
		)
		// Set error handler once
		l.once.Do(func() {
			// override error handler
			errHandler = c.App().ErrorHandler
		})
		chainErr := c.Next()
		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		userLog := &model.UserLog{
			Ip:        c.IP(),
			UID:       c.Locals("uid").(uint),
			Method:    c.Method(),
			Path:      c.Path(),
			Status:    uint32(c.Response().StatusCode()),
			UserAgent: c.Get(fiber.HeaderUserAgent),
		}
		loc, err := l.logRepo.SelectLocation(c.UserContext(), c.IP())
		if err != nil {
			return err
		}
		userLog.Country = loc.Country
		userLog.Region = loc.Region
		userLog.City = loc.City
		userLog.SetPosition(loc.Position)
		agent, err := l.logRepo.SelectAgent(c.UserContext(), c.Get(fiber.HeaderUserAgent))
		if err != nil {
			return err
		}
		userLog.Client = agent.Name
		userLog.OS = agent.OS
		userLog.Device = agent.Device
		userLog.DeviceType = agent.DeviceType
		userLog.Type = l.eventMap[c.Path()]
		if err = l.logRepo.Insert(c.UserContext(), userLog); err != nil {
			return err
		}
		return nil
	}
}

func (l *Logger) Name() string {
	return kmw.LoggerCfg
}
