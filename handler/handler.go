package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/objx"
)

// MainPageHandler -- top page handler
func MainPageHandler(c echo.Context) error {
	auth, err := c.Cookie("auth")
	if err != nil {
		return c.Render(http.StatusOK, "welcome", map[string]interface{}{
			"title": "Welcome",
		})
	}

	userData := objx.MustFromBase64(auth.Value)

	return c.Render(http.StatusOK, "top", map[string]interface{}{
		"name":      userData["name"],
		"email":     userData["email"],
		"avatarURL": userData["avatarURL"],
		"title":     "TopPage",
	})
}
