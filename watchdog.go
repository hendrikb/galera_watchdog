package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	err error

	httpHost = kingpin.Flag("http-host", "Connect to host").Default("localhost").Short('H').String()
	httpPort = kingpin.Flag("http-port", "Connect to host").Default("9199").Short('P').String()

	mysqlHost = kingpin.Flag("sql-host", "Connect to host").Default("localhost").Short('h').String()
	mysqlUser = kingpin.Flag("sql-user", "User for login").Default("root").Short('u').String()
	mysqlPass = kingpin.Flag("sql-password", "Password to use when connecting to server").Short('p').String()

	donorOk = kingpin.Flag("donor-ok", "treat donor as regular working node").Short('d').Bool()
)

func main() {
	kingpin.Parse()
	DB, err = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", *mysqlUser, *mysqlPass, *mysqlHost),
	)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", statusHandler)
	http.ListenAndServe(*httpHost+":"+*httpPort, nil)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	err = DB.Ping()
	if err != nil {
		http.Error(w, "Galera Node is *down*. ("+err.Error()+")", 503)
		return
	}

	var key string
	var value int64
	err = DB.QueryRow("SHOW STATUS LIKE 'wsrep_local_state'").Scan(&key, &value)
	if err != nil {
		http.Error(w, "Galera Node is *down*. ("+err.Error()+")", 503)
		return
	}

	switch {
	case 4 == value:
		fmt.Fprintf(w, "Galera Node is running.")
		return
	case 2 == value && *donorOk == true:
		fmt.Fprintf(w, "Galera Node is running.")
		return
	default:
		http.Error(w, "Galera Node is *down*. (State Mismatch)", 503)
		return
	}
}
