package routers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id   string
	Name string
}

func UserRouter(g *echo.Group) {
	users := []*User{}

	g.GET("/", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JWTCustomClaims)
		name := claims.Name
		fmt.Printf("Access From %s\n", name)
		return c.JSON(http.StatusOK, users)
	})
	type CreateUserBody struct {
		Name string `json:"name" validate:"required,min=1"`
	}
	g.POST("/", func(c echo.Context) error {
		body := new(CreateUserBody)
		if err := c.Bind(body); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		if err := c.Validate(body); err != nil {
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
	g.PUT("/:id", func(c echo.Context) error {
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
	g.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")

		for i, user := range users {
			if user.Id == id {
				users = append(users[:i], users[i+1:]...)
				return c.JSON(http.StatusAccepted, user)
			}
		}
		return c.JSON(http.StatusNotFound, "User not found.")
	})

}
