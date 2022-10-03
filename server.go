package main

import (
	"fmt"
	"net/http"

	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	type User struct {
		Id   string
		Name string
	}
	users := []*User{}
	userGroup := e.Group("/users")
	userGroup.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, users)
	})
	type CreateUserBody struct {
		Name string `json:"name"`
	}
	userGroup.POST("", func(c echo.Context) error {
		body := new(CreateUserBody)
		if err := c.Bind(body); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		canonicID, err := nanoid.Standard(21)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		user := User{
			Name: body.Name,
			Id:   canonicID(),
		}
		users = append(users, &user)
		return c.JSON(http.StatusCreated, user)
	})
	type UpdateUserBody struct {
		Name string `json:"name"`
	}
	userGroup.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		body := new(UpdateUserBody)
		if err := c.Bind(body); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		for _, user := range users {
			if user.Id == id {
				user.Name = body.Name
				return c.JSON(http.StatusAccepted, user)
			}
		}
		return c.JSON(http.StatusNotFound, "User not found.")
	})
	userGroup.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")

		for i, user := range users {
			if user.Id == id {
				users = append(users[:i], users[i+1:]...)
				return c.JSON(http.StatusAccepted, user)
			}
		}
		return c.JSON(http.StatusNotFound, "User not found.")
	})

	fmt.Println("Routes ->")
	maxPathLength := 0
	for _, route := range e.Routes() {
		if maxPathLength < len(route.Path) {
			maxPathLength = len(route.Path)
		}
	}
	for i, route := range e.Routes() {
		fmt.Printf("%2d %6s: %-*s %s\n", i, route.Method, maxPathLength+2, route.Path, route.Name)
	}

	e.Logger.Fatal(e.Start(":1323"))
}
