package dblayer

import (
	"database/sql"
	"os"
	"testing"
)

func newTstDB(t *testing.T) *sql.DB {

	t.Helper()
	db, err := sql.Open("mysql", "tst_http:t5thttP@tcp(localhost:3306)/tst_gists?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/init.sql")

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/drop.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
