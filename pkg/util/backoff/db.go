package backoff

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

func QueryContext(db *sql.DB, ctx context.Context, query string, i ...interface{}) (*sql.Rows, error) {
	entry := logrus.NewEntry(logrus.StandardLogger())

	var rows *sql.Rows = nil

	backoff := New("query-context", func() error {
		r, e := db.QueryContext(ctx, query, i...)
		if e != nil {
			return e
		} else {
			rows = r
			return nil
		}
	}, entry).WithMaxAttempt(25)

	execute := backoff.Execute()
	return rows, execute
}

func ExecContext(db *sql.DB, ctx context.Context, query string, i ...interface{}) error {
	entry := logrus.NewEntry(logrus.StandardLogger())

	backoff := New("query-context", func() error {
		_, e := db.QueryContext(ctx, query, i...)
		if e != nil {
			return e
		} else {
			return nil
		}
	}, entry).WithMaxAttempt(25)

	execute := backoff.Execute()
	return execute
}

func StmtQueryContext(stmt *sql.Stmt, ctx context.Context, i ...interface{}) (*sql.Rows, error) {
	entry := logrus.NewEntry(logrus.StandardLogger())

	var rows *sql.Rows = nil

	backoff := New("query-context", func() error {
		r, e := stmt.QueryContext(ctx, i...)
		if e != nil {
			return e
		} else {
			rows = r
			return nil
		}
	}, entry).WithMaxAttempt(25)

	execute := backoff.Execute()
	return rows, execute
}

func StmtExecContext(stmt *sql.Stmt, ctx context.Context, i ...interface{}) error {
	entry := logrus.NewEntry(logrus.StandardLogger())

	backoff := New("query-context", func() error {
		r, e := stmt.ExecContext(ctx, i...)
		fmt.Println(r)
		if e != nil {
			return e
		} else {
			return nil
		}
	}, entry).WithMaxAttempt(25)

	execute := backoff.Execute()
	return execute
}
