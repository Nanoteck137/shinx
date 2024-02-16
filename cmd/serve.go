package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/shinx/database"
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
		godotenv.Load()

		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL not set")
		}

		conn, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		db := database.New(conn)

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

		app.GET("/login", func(c echo.Context) error {
			return render(c, 200, view.Layout(view.AuthLogin("")))
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

			pretty.Println(body)

			url, err := url.Parse(c.Request().Header.Get("HX-Current-URL"))
			if err != nil {
				return err
			}

			user, err := db.GetUserByUsername(context.Background(), body.Username)
			if err != nil {
				if err == pgx.ErrNoRows {
					return render(c, 200, view.AuthLogin("Incorrect username or password"))
				}

				return err
			}

			if user.Password != body.Password {
				return render(c, 200, view.AuthLogin("Incorrect username or password"))
			}

			pretty.Println(user)

			projectId := url.Query().Get("projectId")
			fmt.Printf("projectId: %v\n", projectId)

			redirectUrl := url.Query().Get("redirectUrl")
			fmt.Printf("redirectUrl: %v\n", redirectUrl)

			c.Response().Header().Set("HX-Redirect", redirectUrl + "?code=123")
			c.Response().WriteHeader(204)
			return nil
		})

		err = app.Start(":3000")
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
