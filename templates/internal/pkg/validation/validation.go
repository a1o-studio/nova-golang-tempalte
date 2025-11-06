package validation

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/a1ostudio/nova/internal/logger"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func NewValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 json tag
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 校验手机号
		if err := v.RegisterValidation("phone", validatePhone); err != nil {
			logger.L().Panic("failed to register phone validation", zap.Error(err))
			return
		}

		// 校验密码
		if err := v.RegisterValidation("password", validatePassword); err != nil {
			logger.L().Panic("failed to register password validation", zap.Error(err))
			return
		}
	}
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}

	// 正则匹配手机
	pattern := `^1(?:3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(phone)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	// 正则匹配密码
	pattern := `^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':",.<>/?\\|~]+$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(password)
}
