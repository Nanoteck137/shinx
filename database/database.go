package database

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/shinx/utils"
)

type ToSQL interface {
	ToSQL() (sql string, params []interface{}, err error)
}

type Database struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Database {
	return &Database{
		conn: conn,
	}
}

var dialect = goqu.Dialect("postgres")

func (db *Database) Exec(ctx context.Context, query ToSQL) (pgconn.CommandTag, error) {
	sql, params, err := query.ToSQL()
	if err != nil {
		return pgconn.CommandTag{}, nil
	}

	return db.conn.Exec(ctx, sql, params...)
}

func (db *Database) CreateProject(ctx context.Context, id string) error {
	ds := dialect.Insert("projects").Rows(goqu.Record{
		"id": id,
	}).Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}

func (db *Database) CreateUser(ctx context.Context, username, password string) (string, error) {
	userId := utils.CreateId()
	ds := dialect.Insert("users").Rows(goqu.Record{
		"id": userId,
		"username": username,
		"password": password,
	}).Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return "", err
	}

	fmt.Printf("tag: %v\n", tag)

	return userId, nil
}

func (db *Database) CreateLink(ctx context.Context, projectId, userId string) (string, error) {
	linkId := utils.CreateId()
	ds := dialect.Insert("project_user_links").Rows(goqu.Record{
		"id": linkId,
		"user_id": userId,
		"project_id": projectId,
	}).Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return "", err
	}

	fmt.Printf("tag: %v\n", tag)

	return linkId, nil
}

func (db *Database) DeleteAll(ctx context.Context) error {
	tag, err := db.conn.Exec(ctx, `
		DELETE FROM project_user_links;
		DELETE FROM users;
		DELETE FROM projects;
	`)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}
