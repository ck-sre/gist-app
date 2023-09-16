package main

import (
	"database/sql"
	"flag"
	"gistapp.ck89.net/internal/dblayer"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type mission struct {
	logger *slog.Logger
	gists  *dblayer.Gistdblayer
	//DB    *sql.DB
	//eLog *log.Logger
	//iLog *log.Logger
}

//
//func (m *mission) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	writer.Write([]byte("This is my home page"))
//}

func main() {
	//Initialize new servemux register landing as a handler

	type cfg struct {
		port   string
		dbconn string
	}
	var cf cfg

	flag.StringVar(&cf.port, "port", ":9100", "port to listen on")
	flag.StringVar(&cf.dbconn, "dbconn", "gistuser:pwd/gistapp?parseTime=true", "connection string for mysql")
	flag.Parse()

	//Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

	mysqlDB, err := createDB(cf.dbconn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer mysqlDB.Close()

	msn := &mission{
		logger: logger,
		gists:  &dblayer.Gistdblayer{DB: mysqlDB},
	}

	customSvr := &http.Server{
		Addr: cf.port,
		//ErrorLog: logErr,
		Handler: msn.paths(),
	}

	logger.Info("Listening on", slog.Any("port", cf.port))
	err = customSvr.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func createDB(dbconn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbconn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
