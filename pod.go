package main

import (
  "context"
  "fmt"
  "io/ioutil"
  "log"
  "strconv"
  "strings"

  corev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
)


func podPid(pod *corev1.Pod) int {
  cid := pod.Status.ContainerStatuses[0].ContainerID
  cidspl := strings.Split(cid, "/")
  cid = cidspl[len(cidspl)-1]
  log.Printf("container id %s", cid)
  procfiles, err := ioutil.ReadDir("/proc")
  if err != nil {
    log.Println(err)
    return -1
  }
  for _, fi := range procfiles {
    if !fi.IsDir() {
      continue
    }
    filename := fmt.Sprintf("/proc/%s/cgroup", fi.Name())
    cgrpcontent, err := ioutil.ReadFile(filename)
    if err != nil {
      log.Printf("can't open %s", filename)
      continue
    }
    if strings.Contains(string(cgrpcontent), cid) {
      pid, err := strconv.Atoi(fi.Name())
      if err == nil {
        return pid
      } else {
        log.Printf("can't convert to number %s", fi.Name())
      }
    }
  }
  return -1
}

func podCreate(cs *kubernetes.Clientset, origPod *corev1.Pod, image string) {
  podpid := podPid(origPod)
  log.Printf("pid of original pod is %d", podpid)
  if podpid < 1 {
    return
  }
  b := true
  pod := corev1.Pod {
    ObjectMeta: metav1.ObjectMeta {
      Name: "capture-"+origPod.Name,
      Namespace: "default",
    },
    Spec: corev1.PodSpec{
      NodeName: origPod.Spec.NodeName,
      Containers: []corev1.Container { {
          Name: "c1",
          Image: image,
          Command: []string {"/k8cap"},
          Args: []string {"capturesvc", origPod.Namespace, origPod.Name},
          SecurityContext: &corev1.SecurityContext{
            Privileged: &b,
          },
        },
      },
      HostPID: true,
      ServiceAccountName: "capture",
    },
  }
  pn, err := cs.CoreV1().Pods("default").Create(context.Background(), &pod, metav1.CreateOptions{})
  if err != nil {
    log.Println(err)
  }
  log.Println(pn)
}