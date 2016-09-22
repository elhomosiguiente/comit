package main

import (
	"flag"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/tmsp/server"
	. "github.com/zballs/3ii/actions"
	. "github.com/zballs/3ii/app"
	. "github.com/zballs/3ii/types"
	"net/http"
)

func main() {

	addrPtr := flag.String("addr", "tcp://0.0.0.0:46658", "Listen address")
	tmspPtr := flag.String("tmsp", "socket", "socket | grpc")
	flag.Parse()

	// Start the listener
	app := NewApplication()
	_, err := server.NewServer(*addrPtr, *tmspPtr, app)
	if err != nil {
		Exit(err.Error())
	}

	RegisterTemplates(
		"create_account.html",
		"create_admin.html",
		"remove_account.html",
		"remove_admin.html",
		"submit_form.html",
		"find_form.html",
		"resolve_form.html",
		"search_forms.html",
		"feed.html",
		"admin.html",
		"metrics.html",
	)

	CreatePages(
		"create_account",
		"create_admin",
		"remove_account",
		"remove_admin",
		"submit_form",
		"find_form",
		"resolve_form",
		"search_forms",
		"feed",
		"admin",
		"metrics",
	)

	action_listener, err := StartActionListener()
	if err != nil {
		Exit(err.Error())
	}

	go action_listener.FeedUpdates()

	action_listener.Run(app)

	js := JustFiles{http.Dir("static/")}
	http.Handle("/", action_listener)
	http.HandleFunc("/create_account", CustomHandler("create_account.html"))
	http.HandleFunc("/create_admin", CustomHandler("create_admin.html"))
	http.HandleFunc("/remove_account", CustomHandler("remove_account.html"))
	http.HandleFunc("/remove_admin", CustomHandler("remove_admin.html"))
	http.HandleFunc("/submit_form", CustomHandler("submit_form.html"))
	http.HandleFunc("/find_form", CustomHandler("find_form.html"))
	http.HandleFunc("/resolve_form", CustomHandler("resolve_form.html"))
	http.HandleFunc("/search_forms", CustomHandler("search_forms.html"))
	http.HandleFunc("/feed", CustomHandler("feed.html"))
	http.HandleFunc("/admin", CustomHandler("admin.html"))
	http.HandleFunc("/metrics", CustomHandler("metrics.html"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(js)))
	http.ListenAndServe(":8888", nil)

	// Wait forever
	TrapSignal(func() {
		// Cleanup
	})

}
