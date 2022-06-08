package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/calebpalmer/simpleftsservice/internal/config"
	"github.com/calebpalmer/simpleftsservice/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config, err := config.New("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	if err := handlers.RegisterHandlers(router, config); err != nil {
		log.Fatal(err)
	}

	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", config.HttpConfig.Port)

	fmt.Printf("Starting server on %s\n", addr)
	httpServer := http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
