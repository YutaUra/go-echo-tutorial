package routers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(c echo.Context) error {
	body := new(LoginBody)
	if err := c.Bind(body); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if body.Username != "jon" || body.Password != "shhh!" {
		return echo.ErrUnauthorized
	}

	claims := &JWTCustomClaims{
		"Jon Snow",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func AuthRouter(g *echo.Group) {
	g.POST("/login", login)

}
