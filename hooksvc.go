package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"

  admission "k8s.io/api/admission/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HookHandler struct {}
func (h *HookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Printf("new request %s\n", r.RequestURI)

  data, err := ioutil.ReadAll(r.Body)
  r.Body.Close()
  var review admission.AdmissionReview
  err = json.Unmarshal(data, &review)
  if err != nil {
    log.Println(err)
    return
  }

  patchType := admission.PatchTypeJSONPatch
  patch := []byte("[]")
  resp := admission.AdmissionReview{
    TypeMeta: metav1.TypeMeta {
      Kind: "AdmissionReview",
      APIVersion: "admission.k8s.io/v1" },
    Response: &admission.AdmissionResponse {
      UID: review.Request.UID,
      Allowed: true,
      PatchType: &patchType,
      Patch: patch,
    },
  }
  respb, err := json.Marshal(&resp)
  if err != nil {
    log.Println(err)
  }
  log.Println(string(respb))
  w.Write(respb)
}

func hooksvc() {
  server := &http.Server{
    Addr:    ":443",
    Handler: &HookHandler{},
  }
  err := server.ListenAndServeTLS("/cert.pem", "/private.pem")
  if err != nil {
    panic(err)
  }
}