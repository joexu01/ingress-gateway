package middleware

import (
	"github.com/gin-gonic/gin"
)

func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//serviceInterface, ok := c.Get("service")
		//if !ok {
		//	middleware.ResponseError(c, 2001, errors.New("service not found"))
		//	c.Abort()
		//	return
		//}
		//serviceDetail := serviceInterface.(*service.Detail)
		//
		//token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		//matched := false
		//if token != "" {
		//	claims, err := public.JwtDecode(token)
		//	if err != nil {
		//		middleware.ResponseError(c, 2002, err)
		//		c.Abort()
		//		return
		//	}
		//	appList := dao.AppManagerHandler.GetAppList()
		//	for _, appInfo := range appList {
		//		if appInfo.AppID == claims.Issuer {
		//			c.Set("app_detail", appInfo)
		//			matched = true
		//			break
		//		}
		//	}
		//}
		//if serviceDetail.AccessControl.OpenAuth == 1 && !matched {
		//	middleware.ResponseError(c, 2003, errors.New("failed to match app info"))
		//	c.Abort()
		//	return
		//}

		c.Next()
	}
}
