package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/joexu01/ingress-gateway/public"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	zhTranslations "gopkg.in/go-playground/validator.v9/translations/zh"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// TranslationMiddleware 负责国际化和参数验证
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			log.Println("English Mode")
			break
		default:
			_ = zhTranslations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go

			//验证username
			_ = val.RegisterValidation("validate_username",
				func(fl validator.FieldLevel) bool {
					return fl.Field().String() == "admin"
				},
			)

			//自定义验证器
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			_ = val.RegisterTranslation(
				"validate_username",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_username", "{0} 填写不正确", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_username", fe.Field())
					return t
				},
			)

			//验证service_name
			_ = val.RegisterValidation("validate_service_name",
				func(fl validator.FieldLevel) bool {
					ok, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
					return ok
				},
			)

			_ = val.RegisterTranslation(
				"validate_service_name",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_service_name", "{0} 不符合输入格式", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_service_name", fe.Field())
					return t
				},
			)

			//验证http转发rule
			_ = val.RegisterValidation("validate_rule",
				func(fl validator.FieldLevel) bool {
					ok, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
					return ok
				},
			)

			_ = val.RegisterTranslation(
				"validate_rule",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_rule", "{0} 必须是非空字符", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_rule", fe.Field())
					return t
				},
			)

			//验证url重写规则
			_ = val.RegisterValidation("validate_url_rewrite",
				func(fl validator.FieldLevel) bool {
					if fl.Field().String() == "" {
						return true
					}
					for _, rule := range strings.Split(fl.Field().String(), ",") {
						if len(strings.Split(rule, " ")) != 2 {
							return false
						}
					}
					return true
				},
			)

			_ = val.RegisterTranslation(
				"validate_url_rewrite",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_url_rewrite", "{0} 不符合输入格式", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_url_rewrite", fe.Field())
					return t
				},
			)

			//验证Header Transform规则
			_ = val.RegisterValidation("validate_header_transform",
				func(fl validator.FieldLevel) bool {
					if fl.Field().String() == "" {
						return true
					}
					for _, rule := range strings.Split(fl.Field().String(), ",") {
						if len(strings.Split(rule, " ")) != 3 {
							return false
						}
					}
					return true
				},
			)

			_ = val.RegisterTranslation(
				"validate_header_transform",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_header_transform", "{0} 不符合输入格式", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_header_transform", fe.Field())
					return t
				},
			)

			//验证集群IP列表
			_ = val.RegisterValidation("validate_ip_port_list",
				func(fl validator.FieldLevel) bool {
					for _, addr := range strings.Split(fl.Field().String(), ",") {
						if matched, _ := regexp.Match(`^\S+:\d+$`, []byte(addr)); !matched {
							return false
						}
					}
					return true
				},
			)

			_ = val.RegisterTranslation(
				"validate_ip_port_list",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_ip_port_list", "{0} 不符合输入格式", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_ip_port_list", fe.Field())
					return t
				},
			)

			//验证权重weight
			_ = val.RegisterValidation("validate_weight_list",
				func(fl validator.FieldLevel) bool {
					for _, weight := range strings.Split(fl.Field().String(), ",") {
						if matched, _ := regexp.Match(`^\d+$`, []byte(weight)); !matched {
							return false
						}
					}
					return true
				},
			)

			_ = val.RegisterTranslation(
				"validate_weight_list",
				trans,
				func(ut ut.Translator) error {
					return ut.Add("validate_weight_list", "{0} 不符合输入格式", true)
				},
				func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validate_weight_list", fe.Field())
					return t
				},
			)

			//
			_ = val.RegisterValidation("validate_ip_list", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, item := range strings.Split(fl.Field().String(), ",") {
					matched, _ := regexp.Match(`\S+`, []byte(item)) //ip_addr
					if !matched {
						return false
					}
				}
				return true
			})

			//验证service name
			_ = val.RegisterValidation("validate_service_name", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return matched
			})

			_ = val.RegisterTranslation("validate_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("validate_service_name", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validate_service_name", fe.Field())
				return t
			})
			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}
