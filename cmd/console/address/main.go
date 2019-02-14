package main

import (
	"os"
	"time"

	"github.com/apex/gateway"
	"github.com/go-chi/chi"
	"github.com/linkai-io/am/am"
	"github.com/linkai-io/am/pkg/lb/consul"
	"github.com/linkai-io/frontend/api/console/address"
	"github.com/linkai-io/frontend/pkg/initializers"
	"github.com/linkai-io/frontend/pkg/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var addrClient am.AddressService
var scanGroupClient am.ScanGroupService

func init() {
	zerolog.TimeFieldFormat = ""
	log.Logger = log.With().Str("lambda", "Address").Logger()

	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	consul.RegisterDefault(time.Second*5, consulAddr) // Address comes from CONSUL_HTTP_ADDR or from aws metadata

	scanGroupClient = initializers.ScanGroupClient()
	addrClient = initializers.AddressClient()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.UserCtx)
	addrHandlers := address.New(addrClient, scanGroupClient)

	r.Route("/address", func(r chi.Router) {
		r.Get("/group/{id}", addrHandlers.GetAddresses)
		r.Get("/group/{id}/hosts", addrHandlers.GetHostList)
		r.Put("/group/{id}/initial", addrHandlers.PutInitialAddresses)
		r.Get("/group/{id}/count", addrHandlers.GetGroupCount)
		r.Post("/group/{id}/download", addrHandlers.ExportAddresses)
		r.Patch("/group/{id}/delete", addrHandlers.DeleteAddresses)
		r.Patch("/group/{id}/ignore", addrHandlers.IgnoreAddresses)
	})

	err := gateway.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
