package goserver

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginswagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/marmotedu/goserver/internal/goserver/controller/v1/post"
	"github.com/marmotedu/goserver/internal/goserver/controller/v1/user"
	"github.com/marmotedu/goserver/internal/goserver/store"
	"github.com/marmotedu/goserver/internal/pkg/known"
	"github.com/marmotedu/goserver/internal/pkg/middleware"
	"github.com/marmotedu/goserver/pkg/core"
	"github.com/marmotedu/goserver/pkg/errno"
	"github.com/marmotedu/goserver/pkg/token"
)

// loadRouter loads the middlewares, routes, handlers.
func loadRouter(g *gin.Engine, mw ...gin.HandlerFunc) {
	installMiddleware(g, mw...)
	installController(g)
}

// installMiddleware install Middlewares.
func installMiddleware(g *gin.Engine, mw ...gin.HandlerFunc) {
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
}

// installController install controllers.
func installController(g *gin.Engine) {
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// swagger api docs
	g.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	pprof.Register(g)

	uc := user.New(store.S)
	pc := post.New(store.S)

	// api for authentication functionalities
	g.POST("/login", uc.Login)

	// The user handlers, requiring authentication
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.Use(authMiddleware())
			userv1.DELETE(":name", uc.Delete)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.PUT(":name", uc.Update)
			userv1.GET("", uc.List)
			userv1.GET(":name", uc.Get)
		}

		// post RESTful resource
		postv1 := v1.Group("/posts", authMiddleware())
		{
			postv1.POST("", pc.Create)
			postv1.DELETE("", pc.DeleteCollection)
			postv1.DELETE(":postID", pc.Delete)
			postv1.PUT(":postID", pc.Update)
			postv1.GET("", pc.List)
			postv1.GET(":postID", pc.Get)
		}
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		username, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()

			return
		}

		c.Set(known.XUsernameKey, username)
		c.Next()
	}
}
