package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

// LoginHandler -- Login to each provider
func LoginHandler(c echo.Context) error {
	provider, err := gomniauth.Provider(c.Param("provider"))
	if err != nil {
		return err
	}

	authURL, err := provider.GetBeginAuthURL(nil, nil)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// CallbackHandler -- Provider called this handler after login
func CallbackHandler(c echo.Context) error {
	provider, err := gomniauth.Provider(c.Param("provider"))
	if err != nil {
		return err
	}

	omap, err := objx.FromURLQuery(c.QueryString())
	if err != nil {
		return err
	}

	creds, err := provider.CompleteAuth(omap)
	if err != nil {
		return err
	}

	user, err := provider.GetUser(creds)
	if err != nil {
		return err
	}

	authCookieValue := objx.New(map[string]interface{}{
		"name":      user.Name(),
		"email":     user.Email(),
		"avatarURL": user.AvatarURL(),
	}).MustBase64()

	cookie := &http.Cookie{
		Name:    "auth",
		Value:   authCookieValue,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)

	// return c.String(http.StatusOK, "Login Success!")
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
