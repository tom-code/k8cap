package main

import (
  "log"
  "os"
  "time"

  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
)



func kubeCheckConnection(cs *kubernetes.Clientset) bool {
  ver, err := cs.Discovery().ServerVersion()
  if err == nil {
    log.Printf("k8s server version %s.%s\n", ver.Major, ver.Minor)
    return true
  }
  log.Println(err.Error())
  time.Sleep(1*time.Second)
  log.Println("can't contact k8s")
  return false
}

func kubeConnect() *kubernetes.Clientset {

  file := os.Getenv("KUBECONFIG")
  if len(file) == 0 {
    file = "admin.conf"
  }
  config, err := clientcmd.BuildConfigFromFlags("", file)
  if err != nil {
    log.Println("Error reading admin.conf. Trying in-cluster config")
    config, err = rest.InClusterConfig()
    if err != nil {
      panic(err)
    }
  }
  
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }
  
  log.Println("clientSet created")
  
  for {
    ok := kubeCheckConnection(clientset)
    if ok {
      break
    }
    time.Sleep(1*time.Second)
  }
  
  return clientset
}