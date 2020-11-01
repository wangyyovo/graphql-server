package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/integration"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"graphql-server/internal/model"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	// generate package
	"graphql-server/internal/bff"
	"graphql-server/internal/resolver"
	"graphql-server/internal/server/controller"

	// gql-gen package
	"github.com/99designs/gqlgen/graphql/handler"
	gdebug "github.com/99designs/gqlgen/graphql/handler/debug"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/playground"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "9020"
	}
}

func initializeController() controller.Controller {
	var mb int64 = 1 << 20

	srv := handler.New(bff.NewExecutableSchema(resolver.New()))
	// 订阅
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: 4 * time.Second,
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				fmt.Println("websocket CheckOrigin")
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			fmt.Println("websocket InitFunc")
			return ctx, nil
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	// 上传文件配置
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 * mb,
		MaxUploadSize: 50 * mb,
	})

	//srv.SetQueryCache(lru.New(1000))
	srv.SetQueryCache(graphql.NoCache{})

	// 打开自省
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	// 限制查询复杂度
	//srv.Use(extension.FixedComplexityLimit(100))
	// 异常捕获
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		log.Print(err)
		debug.PrintStack()
		return errors.New("user message on panic")
	})


	srv.Use(&gdebug.Tracer{})
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		rc := graphql.GetFieldContext(ctx)
		fmt.Println("Entered", rc.Object, rc.Field.Name)
		res, err = next(ctx)
		fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
		return res, err
	})

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		fmt.Println(rc.OperationName)
		res := next(ctx)
		return res
	})

	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		graphql.RegisterExtension(ctx, "example", "value")

		resp := next(ctx)
		errors1 := graphql.GetErrors(ctx)
		resp.Errors = append(resp.Errors,errors1...)
		return resp
	})

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		var ie *integration.CustomError
		if errors.As(err, &ie) {
			return &gqlerror.Error{
				Message: ie.UserMessage,
				Path:    graphql.GetPath(ctx),
			}
		}
		return graphql.DefaultErrorPresenter(ctx, err)
	})

	return controller.Controller{
		PlaygroundServer: playground.Handler("GraphQL playground", "/query"),
		QueryServer: controller.QueryMiddleware(srv),
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("auth")
		c.Next()
	}
}

func setupRouter(ctrl controller.Controller) *gin.Engine {
	router := gin.Default()
	router.Use(AuthMiddleware())
	router.Use(GinContextToContextMiddleware())
	router.Use(model.LoaderMiddleware())
	router.GET("/", ctrl.PlaygroundHandler())
	router.POST("/query", ctrl.QueryHandler())
	router.GET("/query", ctrl.QueryHandler())
	return router
}

func main() {
	ctrl := initializeController()
	router := setupRouter(ctrl)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("failed launch router. err=%s", err)
		os.Exit(1)
	}
}