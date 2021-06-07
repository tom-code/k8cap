package main

import (
  "log"
  "net/http"
)


type Handler struct {}
func (h* Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println(r.RequestURI)
  if r.RequestURI == "/reset" {
    pcap.reset()
    w.WriteHeader(200)
    return
  }
  out := pcap.get()
  w.Write(out)
}

func httpStart() {
  server := &http.Server{
    Addr:    ":80",
    Handler: &Handler{},
  }
  go server.ListenAndServe()
}