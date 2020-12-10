package web_server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"myapp/libs"

	"github.com/kataras/iris/v12"

	"myapp/models"
	"myapp/routes"

	"github.com/kataras/iris/v12/context"
)

type Server struct {
	App       *iris.Application
	AssetFile http.FileSystem
}

func NewServer(assetFile http.FileSystem) *Server {
	app := iris.Default()
	return &Server{
		App:       app,
		AssetFile: assetFile,
	}
}

func (s *Server) Serve() error {
	if libs.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", libs.Config.Host, 443)
		if err := s.App.Run(iris.TLS(host, libs.Config.Certpath, libs.Config.Certkey)); err != nil {
			return err
		}
	} else {
		if err := s.App.Run(
			iris.Addr(fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) NewApp() {
	s.App.Logger().SetLevel(libs.Config.LogLevel)

	//if libs.Config.Bindata {
	//	s.App.RegisterView(iris.HTML(s.AssetFile, ".html"))
	//	s.App.HandleDir("/", s.AssetFile)
	//}

	libs.InitDb()
	models.Migrate()

	//iris.RegisterOnInterrupt(func() {
	//	_ = libs.Db
	//})

	routes.App(s.App) //注册 app 路由
}

type PathName struct {
	Name   string
	Path   string
	Method string
}

func getPathNames(routeReadOnly []context.RouteReadOnly) []*PathName {
	var pns []*PathName
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("routeReadOnly：%v", routeReadOnly))
	}
	for _, s := range routeReadOnly {
		pn := &PathName{
			Name:   s.Name(),
			Path:   s.Path(),
			Method: s.Method(),
		}
		pns = append(pns, pn)
	}

	return pns
}

// 过滤非必要权限
func isPermRoute(name string) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH", "payload"}
	for _, er := range exceptRouteName {
		if strings.Contains(name, er) {
			return true
		}
	}
	return false
}
