package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/shinx/database"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
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
		ctx := context.Background()

		err = db.DeleteAll(ctx)
		if err != nil {
			log.Fatal(err)
		}

		projectId := "sewaddle"
		err = db.CreateProject(ctx, projectId)
		if err != nil {
			log.Fatal(err)
		}

		userId, err := db.CreateUser(ctx, "test", "test")
		if err != nil {
			log.Fatal(err)
		}

		linkId, err := db.CreateLink(ctx, projectId, userId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("linkId: %v\n", linkId)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
