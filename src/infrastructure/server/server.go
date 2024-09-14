package server

import (
	"context"
	"time"

	container "github.com/Go_CleanArch/infrastructure/container"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/Go_CleanArch/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type contextKey string

// Init is initialize server
func Init() {
	r := router()
	r.Run()
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), contextKey("context"), c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CustomLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 日本時間のタイムゾーンを取得
		loc, _ := time.LoadLocation("Asia/Tokyo")

		// リクエスト開始時の日本時間
		startTime := time.Now().In(loc)

		// リクエスト開始前のログ
		log.WithFields(log.Fields{
			"メッセージ":  "リクエスト開始",
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Received request")

		// リクエスト処理
		c.Next()

		// リクエスト終了時の日本時間
		endTime := time.Now().In(loc)

		// ステータスコードに基づいてログレベルを決定
		statusCode := c.Writer.Status()
		fields := log.Fields{
			"メッセージ":        "リクエスト終了",
			"ステータス":        statusCode,
			"メソッド":         c.Request.Method,
			"エンドポイント":      c.Request.URL.Path,
			"クライアントIPアドレス": c.ClientIP(),
			"実行開始時刻":       startTime.Format("2006-01-02 15:04:05"),
			"実行終了時刻":       endTime.Format("2006-01-02 15:04:05"),
			"処理時間":         endTime.Sub(startTime).String(), // 処理時間
		}
		if statusCode >= 400 && statusCode < 500 {
			log.WithFields(fields).Warn("Client error")
		} else if statusCode >= 500 {
			log.WithFields(fields).Error("Server error")
		} else {
			log.WithFields(fields).Info("Completed request")
		}
	}
}

func router() *gin.Engine {
	route := gin.Default()

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	route.Use(GinContextToContextMiddleware())
	route.Use(CustomLoggingMiddleware())

	route.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"*",
		},
		// アクセスを許可したいHTTPメソッド
		AllowMethods: []string{
			"POST",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	// Initialize dependencies
	ctx := context.Background()
	cont, err := container.NewContainer(ctx)
	if err != nil {
		return nil
	}

	// ヘルスチェックエンドポイント
	healthRoute := route.Group("/")
	{
		healthRoute.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		})
	}

	userRoute := route.Group("/api/users")
	{
		ctrl := cont.UserContainer.UserController
		userRoute.POST("/login", ctrl.LoginControler)
		userRoute.POST("", ctrl.UserController)
	}

	return route
}
