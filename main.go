package main

import (
	ServerConfig "Reveroxy/internal/config"
	Backend "Reveroxy/pkg/Backend"
	"net"
	"net/http"
)

func main() {

	err := Backend.AddApp(&ServerConfig.ServPydio)
	if err != nil {
		panic(err)
	}
	err = Backend.AddApp(&ServerConfig.ServMinio)
	if err != nil {
		panic(err)
	}
	// err = Backend.AddApp(&ServerConfig.ServDB)
	// if err != nil {
	// 	panic(err)
	// }
	// err = Backend.AddApp(&ServerConfig.ServNote)
	// if err != nil {
	// 	panic(err)
	// }
	err = Backend.AddApp(&ServerConfig.ServJupyter)
	if err != nil {
		panic(err)
	}
	// err = Backend.AddService(&ServerConfig.ServSSH)
	// if err != nil {
	// 	panic(err)
	// }

	// go http.ListenAndServeTLS(":443", "reveroxy.cert", "reveroxy.key", http.HandlerFunc(Backend.HandleTLS))
	// http.ListenAndServe(":80", http.HandlerFunc(Backend.RedirectToTLS))
	lnHTTP, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	defer lnHTTP.Close()

	lnHTTPs, err := net.Listen("tcp", ":443")
	if err != nil {
		panic(err)
	}
	defer lnHTTPs.Close()

	go Backend.Serve(lnHTTP, http.HandlerFunc(Backend.RedirectToTLS))
	Backend.ServeTLS(lnHTTPs, http.HandlerFunc(Backend.HandleTLS), "assets/secure/Reveroxy.cert.pem", "assets/secure/Reveroxy.key")

}
