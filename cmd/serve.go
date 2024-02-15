package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/shinx/public"
	"github.com/nanoteck137/shinx/view"
	"github.com/spf13/cobra"
)

func render(c echo.Context, status int, component templ.Component) error {
	c.Response().Writer.WriteHeader(status)

	err := component.Render(context.Background(), c.Response().Writer)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		// godotenv.Load()
		//
		// dbUrl := os.Getenv("DB_URL")
		// if dbUrl == "" {
		// 	log.Fatal("DB_URL not set")
		// }
		//
		// db, err := pgxpool.New(context.Background(), dbUrl)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//
		// _ = db

		app := echo.New()

		app.StaticFS("/public", public.Content)

		app.GET("/", func(c echo.Context) error {
			cookie, _ := c.Cookie("access-token")
			pretty.Println(cookie)
			if cookie == nil {
				return c.Redirect(http.StatusPermanentRedirect, "/login")
			}

			return render(c, 200, view.Layout(view.Index()))
		})

		app.GET("/setup", func(c echo.Context) error {
			return nil
		})

		app.GET("/login", func(c echo.Context) error {
			cookie, _ := c.Cookie("access-token")
			if cookie != nil {
				return c.Redirect(http.StatusPermanentRedirect, "/")
			}

			return render(c, 200, view.Layout(view.Login(view.LoginError{})))
		})

		app.POST("/login", func(c echo.Context) error {
			type Body struct {
				Username string `form:"username"`
				Password string `form:"password"`
			}

			var body Body
			err := c.Bind(&body)
			if err != nil {
				return err
			}

			if body.Username != "admin" && body.Password != "admin" {
				loginErr := view.LoginError{Username: "Incorrect user credentials"}
				return render(c, 200, view.Login(loginErr))
			}

			cookie := &http.Cookie{
				Name:    "access-token",
				Value:   "logged in",
				SameSite: http.SameSiteStrictMode,
				Expires: time.Now().Add(1 * time.Minute),
			}

			c.SetCookie(cookie)

			c.Response().Header().Set("HX-Redirect", "/")
			c.Response().WriteHeader(204)

			return nil
		})

		err := app.Start(":3000")
		if err != nil {
			log.Fatal(err)
		}

		// api := api.New(db)
		// if err := api.Start(":3000"); err != nil {
		// 	log.Fatal(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
