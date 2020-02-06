package route

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Amakuchisan/OAuth-Login/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"
)

// TemplateRenderer -- custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render -- Render a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Echo -- instance for initialization
var Echo *echo.Echo

func init() {
	e := echo.New()

	err := setupOAuth()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Group("", authCheckMiddleware())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	e.GET("/", handler.MainPageHandler)
	e.GET("/auth/login/:provider", handler.LoginHandler)
	e.GET("/auth/callback/:provider", handler.CallbackHandler)

	Echo = e
}

func authCheckMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := c.Cookie("auth")
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}

			return next(c)
		}
	}
}

func setupOAuth() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	host := os.Getenv("SAMPLE_HOST")
	googleCallbackURL := fmt.Sprintf("http://%s/auth/callback/google", host)

	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			googleCallbackURL,
		),
	)
	return nil
}
