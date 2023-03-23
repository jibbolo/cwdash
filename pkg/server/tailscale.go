package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"tailscale.com/tsnet"
)

func TailscaleServer(authkey, hostname string, h http.Handler) error {
	s := &tsnet.Server{
		Ephemeral: true,
		Hostname:  hostname,
		AuthKey:   authkey,
		Logf:      func(string, ...any) {},
	}
	defer s.Close()

	ln, err := s.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}
	ln = tls.NewListener(ln, &tls.Config{
		GetCertificate: lc.GetCertificate,
	})
	log.Print("Starting default server on https://" + hostname)
	return http.Serve(ln, h)
}
