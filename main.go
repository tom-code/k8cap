package main

import (
  "context"
  "crypto/tls"
  "fmt"
  "log"
  "net/http"

  "github.com/spf13/cobra"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func start(cmd *cobra.Command, args []string) {
  clientset := kubeConnect()
  namespace := args[0]
  podname := args[1]
  image, _ := cmd.Flags().GetString("image")
  origpod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podname, metav1.GetOptions{})
  if err != nil {
    log.Println(err)
  }
  podCreate(clientset, origpod, image)
}

/*
func capturesvc(cmd *cobra.Command, args []string) {
  pid, err := strconv.Atoi(args[0])
  if err != nil {
    log.Println(err)
    return
  }
  capture(pid)
}
*/
func capturesvc(cmd *cobra.Command, args []string) {
  namespace := args[0]
  podname := args[1]

  clientset := kubeConnect()
  origpod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podname, metav1.GetOptions{})
  if err != nil {
    log.Println(err)
    return
  }
  pid := podPid(origpod)
  if pid < 1 {
    log.Println("can't get podpid")
    return
  }
  httpStart()

  capture(pid)
}

func hookserver(cmd *cobra.Command, args []string) {
  hooksvc()
}
func hookinstall(cmd *cobra.Command, args []string) {
  clientset := kubeConnect()
  hookInstall(clientset)
}
func hookdele(cmd *cobra.Command, args []string) {
  clientset := kubeConnect()
  hookDelete(clientset)
}
func initcfunc(cmd *cobra.Command, args []string) {
  http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
  url := fmt.Sprintf("https://hook1.default.svc/start?namespace=%s&pod=%s",
                      args[0], args[1])
  resp, err := http.Get(url)
  if err != nil {
    log.Println(err)
  } else {
    log.Println(resp.Status)
  }
}
func main() {

  root := cobra.Command {
    Use: "k8cap",
    Short: "k8 capturer",
  }

  start := &cobra.Command {
    Use: "start namespace pod",
    Short: "start capture",
    Args: cobra.MinimumNArgs(2),
    Run: start,
  }
  start.Flags().StringP("image", "i", "cap:1", "capturer image")
  root.AddCommand(start)

  /*capturesvc := &cobra.Command {
    Use: "capturesvc pid",
    Short: "start capture service",
    Run: capturesvc,
  }
  root.AddCommand(capturesvc)*/

  capturesvc := &cobra.Command {
    Use: "capturesvc namespace pod",
    Short: "start capture service",
    Args: cobra.MinimumNArgs(2),
    Run: capturesvc,
  }
  root.AddCommand(capturesvc)

  hooksrv := &cobra.Command {
    Use: "hookserver",
    Short: "start hook server",
    //Args: cobra.MinimumNArgs(2),
    Run: hookserver,
  }
  root.AddCommand(hooksrv)

  hookinst := &cobra.Command {
    Use: "hookinstall",
    Short: "install hook server",
    //Args: cobra.MinimumNArgs(2),
    Run: hookinstall,
  }
  root.AddCommand(hookinst)

  hookdel := &cobra.Command {
    Use: "hookdel",
    Short: "delete hook server",
    //Args: cobra.MinimumNArgs(2),
    Run: hookdele,
  }
  root.AddCommand(hookdel)

  initc := &cobra.Command {
    Use: "init",
    Short: "init container service",
    Args: cobra.MinimumNArgs(2),
    Run: initcfunc,
  }
  root.AddCommand(initc)


  if err := root.Execute(); err != nil {
    log.Println(err)
  }
}