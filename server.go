package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"go-echo-tutorial/internal/routers"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type Routes []*echo.Route

func (p Routes) Len() int {
	return len(p)
}
func (p Routes) Less(i, j int) bool {
	diff := strings.Compare(p[i].Path, p[j].Path)
	if diff != 0 {
		return diff > 0
	}
	return strings.Compare(p[i].Method, p[j].Method) > 0
}
func (p Routes) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	jwtConfig := middleware.JWTConfig{
		SigningKey: []byte("secret"),
		Claims:     &routers.JWTCustomClaims{},
	}

	userGroup := e.Group("/users")
	userGroup.Use(middleware.JWTWithConfig(jwtConfig))
	routers.UserRouter(userGroup)

	routers.AuthRouter(e.Group("/auth"))

	fmt.Println("Routes ->")
	maxPathLength := 0
	for _, route := range e.Routes() {
		if maxPathLength < len(route.Path) {
			maxPathLength = len(route.Path)
		}
	}
	routes := Routes(e.Routes())
	sort.Sort(routes)
	for i, route := range routes {
		fmt.Printf("%2d %8s: %-*s %s\n", i, route.Method, maxPathLength+2, route.Path, route.Name)
	}

	e.Logger.Fatal(e.Start(":1323"))

}
