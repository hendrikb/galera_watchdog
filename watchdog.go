package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vharitonsky/iniflags"
)

var (
	DB  *sql.DB
	err error

	httpHost = flag.String("HTTP_HOST", "localhost", "Connect to host")
	httpPort = flag.String("HTTP_PORT", "9199", "Connect to host")

	sqlHost = flag.String("MYSQL_HOST", "localhost", "Connect to host")
	sqlPort = flag.String("MYSQL_PORT", "3306", "Connect to host")
	sqlUser = flag.String("MYSQL_USER", "root", "User for login to MySQL")
	sqlPass = flag.String("MYSQL_PASS", "", "Password for login to MySQL")

	donorOk = flag.Bool("DONOR_OK", false, "treat donor as regular working node")
)

func main() {
	iniflags.Parse()
	DB, err = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql", *sqlUser, *sqlPass, *sqlHost, *sqlPort),
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
	case 2 == value && *donorOk:
		fmt.Fprintf(w, "Galera Node is running.")
		return
	default:
		http.Error(w, "Galera Node is *down*. (State Mismatch)", 503)
		return
	}
}
