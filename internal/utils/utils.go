package utils

import(
	"fmt"
	"context"
	"github.com/gin-gonic/gin"
	"graphql-server/internal/model"
)

var(
	STATUS_CODE_400 = "Bad Request"
	STATUS_CODE_403 = "Forbidden"
	STATUS_CODE_500 = "Internal Server Error"
)

func ContextValueChecksum(ctx context.Context, keys ...string) (claims map[string]string, errors []*model.Errors) {
	if len(keys) < 1 {
		return
	}

	var value string
	var ok bool
	claims = make(map[string]string)

	for _, key := range keys {
		if value, ok = ctx.Value(key).(string); !ok {
			errors = append(errors, MakeErrors(500, fmt.Sprintf("Faild get HTTP header value: %s", key)))
		}
		if value == "" {
			errors = append(errors, MakeErrors(400, fmt.Sprintf("HTTP header value is empty: %s", key)))
		}
		claims[key] = value
	}
	return
}

func MakeErrors(code int, msg string) (errors *model.Errors) {
	switch code {
	case 400:
		errors = &model.Errors{
			Code: 400,
			Message: STATUS_CODE_400,
			Description: msg,
		}
	case 403:
		errors = &model.Errors{
			Code: 403,
			Message: STATUS_CODE_403,
			Description: msg,
		}
	case 500:
		errors = &model.Errors{
			Code: 500,
			Message: STATUS_CODE_500,
			Description: msg,
		}
	}
	return
}

func CastStringPointer(str string) *string {
	return &str
}


func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}