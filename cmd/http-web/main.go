package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"gistapp.ck89.net/internal/dblayer"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mission struct {
	debug     bool
	logger    *slog.Logger
	gists     dblayer.GistModelIface
	tmplCache map[string]*template.Template
	formDcdr  *form.Decoder
	snMgr     *scs.SessionManager
	usrs      dblayer.UserModelIface
}

func main() {
	//Initialize new servemux register landing as a handler

	type cfg struct {
		port   string
		dbconn string
		debug  bool
	}
	var cf cfg

	flag.StringVar(&cf.port, "port", ":9100", "port to listen on")
	flag.StringVar(&cf.dbconn, "dbconn", "gistuser:pwd@(localhost:3306)/gistapp?parseTime=true", "connection string for mysql")
	flag.BoolVar(&cf.debug, "debug", false, "Enable debug flag")
	flag.Parse()

	//Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

	mysqlDB, err := createDB(cf.dbconn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer mysqlDB.Close()

	tmplCache, err := newTmplCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDcdr := form.NewDecoder()

	snMgr := scs.New()
	snMgr.Store = mysqlstore.New(mysqlDB)
	snMgr.Lifetime = 12 * time.Hour
	snMgr.Cookie.Secure = true

	msn := &mission{
		debug:     cf.debug,
		logger:    logger,
		gists:     &dblayer.Gistdblayer{DB: mysqlDB},
		tmplCache: tmplCache,
		formDcdr:  formDcdr,
		snMgr:     snMgr,
		usrs:      &dblayer.UserLayer{MysqlDB: mysqlDB},
	}

	tlsCfg := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, //nolint:staticcheck
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   //nolint:staticcheck
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,  //nolint:staticcheck
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,    //nolint:staticcheck
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, //nolint:staticcheck
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,   //nolint:staticcheck

		},
	}

	customSvr := &http.Server{
		Addr:           cf.port,
		ErrorLog:       slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Handler:        msn.paths(),
		TLSConfig:      tlsCfg,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	logger.Info("Listening on", slog.Any("port", cf.port))
	err = customSvr.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
