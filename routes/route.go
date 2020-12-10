package routes

import (
	"myapp/middleware"

	"myapp/controllers"

	"github.com/kataras/iris/v12"
)

func App(api *iris.Application) {
	api.UseRouter(middleware.CrsAuth())
	//api.UseRouter(middleware.CrsAuth())
	app := api.Party("/").AllowMethods(iris.MethodOptions)
	{
		// 二进制模式 ， 启用项目入口
		// if libs.Config.Bindata {
		// 	app.Get("/", func(ctx iris.Context) { // 首页模块
		// 		_ = ctx.View("index.html")
		// 	})
		// }

		v1 := app.Party("/v1")
		{
			v1.PartyFunc("/admin", func(admin iris.Party) {
				admin.PartyFunc("/users", func(users iris.Party) {
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"
					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"

				})
			})
		}
	}

}
