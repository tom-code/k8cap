package main

import (
  "context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "sync"

  admission "k8s.io/api/admission/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HookHandler struct {
  namespaces map[string]bool
  mutex sync.Mutex
}
func (h *HookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Printf("new request %s\n", r.RequestURI)
  if r.URL.Path == "/start" {
    name := r.URL.Query().Get("pod")
    namespace := r.URL.Query().Get("namespace")
    if (len(name) == 0) || (len(namespace) == 0) {
      w.WriteHeader(400)
      return
    }
    cs := kubeConnect()
    pod, err := cs.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
      log.Println(err)
      w.WriteHeader(500)
      return
    }
    podCreate(cs, pod, "cap:1")
    return
  }
  if r.URL.Path == "/arm" {
    namespace := r.URL.Query().Get("namespace")
    log.Println(namespace)
    if len(namespace) == 0 {
      log.Println("namespace not present")
      w.WriteHeader(400)
      return
    }
    h.mutex.Lock()
    defer h.mutex.Unlock()
    h.namespaces[namespace] = true
    log.Println(h.namespaces)
    return
  }
  if r.URL.Path == "/disarm" {
    namespace := r.URL.Query().Get("namespace")
    log.Println(namespace)
    if len(namespace) == 0 {
      log.Println("namespace not present")
      w.WriteHeader(400)
      return
    }
    h.mutex.Lock()
    defer h.mutex.Unlock()
    delete(h.namespaces, namespace)
    log.Println(h.namespaces)
    return
  }

  data, err := ioutil.ReadAll(r.Body)
  r.Body.Close()
  var review admission.AdmissionReview
  err = json.Unmarshal(data, &review)
  if err != nil {
    log.Println(err)
    return
  }

  armed := false
  h.mutex.Lock()
  ok1, ok2 := h.namespaces[review.Request.Namespace]
  if ok1 && ok2 {
    armed = true
  }
  h.mutex.Unlock()

  patchType := admission.PatchTypeJSONPatch
  patch := []byte("[]")
  if armed {
    pf := `
    [{
      "op": "add",
      "path": "/spec/initContainers",
      "value": [{
        "name": "capinit",
        "image": "cap:1",
        "command": ["/k8cap"],
        "args": ["init", "%s", "%s"]
      }]
    }]
    `
    p := fmt.Sprintf(pf, review.Request.Namespace, review.Request.Name)
    patch = []byte(p)
  }
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
    Handler: &HookHandler{
      namespaces: map[string]bool{},
    },
  }
  err := server.ListenAndServeTLS("/cert.pem", "/private.pem")
  if err != nil {
    panic(err)
  }
}