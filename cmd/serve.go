package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
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
			return render(c, 200, view.Layout(view.Index()))
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
