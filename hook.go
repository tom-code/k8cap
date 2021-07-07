package main

import (
  "context"
  "encoding/base64"
  "log"
  "time"

  aregv1 "k8s.io/api/admissionregistration/v1"
  corev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/util/intstr"
  "k8s.io/client-go/kubernetes"
)

func hookDelete(cs *kubernetes.Clientset) {
  err := cs.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.Background(),
                  "hook.test.com", metav1.DeleteOptions{})
  if err != nil {
    log.Println(err)
  }

  err = cs.CoreV1().Services("default").Delete(context.Background(), "hook1", metav1.DeleteOptions{})
  if err != nil {
    log.Println(err)
  }

  err = cs.CoreV1().Pods("default").Delete(context.Background(), "hook1", metav1.DeleteOptions{})
  if err != nil {
    log.Println(err)
  }
}

func hookInstall(cs *kubernetes.Clientset) {
  var err error

  scope := aregv1.NamespacedScope
  sideeffect := aregv1.SideEffectClassNone
  timeout := int32(5)
  failurepolicy := aregv1.Ignore
  ca := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3akNDQWRhZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFRTVE0d0RBWURWUVFLRXdWdmQyNWwKY2pBZUZ3MHlNVEEzTURZeE9ERXlNREJhRncweU16QTJNRFl4T0RFeU1EQmFNQkF4RGpBTUJnTlZCQW9UQlc5MwpibVZ5TUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF0SnMzUkpqVVRKYzRRUGUrCmtHcHhrMGx4d1pyUWFiNlpwNTJVOXBsR0dTSVZZZ1BaQzhWTnlpYUEwbU0zUjV5aWpxaGtmN29IK0tEUSs0aloKTWNXTzh4T2JTZ2FnRjBTZlpOWERqL2lCWG5ocFJ6cWJ0VWhqWDRkTGRtczFxcHNjZjRDY2dzZ0tRN0o4c2RkTwpzaCtTTjl5a3U1V2p5ZTA2L1M3MkhGeVhqemtzWHcwTmJhblE0NllRSVRpSEFFbmwxd043bXliZUJQcVo3SFJjCmlFdm5zQUMxZjRDVFRDL1E5YnRTVkl5QkloR3FuR1V0RWtvNFNsVW8rWXc2RWYxbFZTeXEzaFovUkJaa3Naak4KcUE1SkdGa1BzRmlRRW1vSGYwbjFuUXVrYlk0R0NCYkZJQWR3aGorby9TZWdJWmVBSHlnM09mYVQ0L3hEZmo4dQpQazRUR1FJREFRQUJvMU13VVRBT0JnTlZIUThCQWY4RUJBTUNCYUF3RXdZRFZSMGxCQXd3Q2dZSUt3WUJCUVVICkF3RXdEQVlEVlIwVEFRSC9CQUl3QURBY0JnTlZIUkVFRlRBVGdoRm9iMjlyTVM1a1pXWmhkV3gwTG5OMll6QU4KQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBbkc1NXIrc1dYT0ZsSUttT2NxY3NsM1lRVC9DazArclh2QUExRXJDdAo4aUxtWlJjamhzMk5vSHhiUzByV3N6b2g2QjZ6Vkx3K2JsQjlNRHN2ZGVoTEV0eko0dlBJRkV1bEMxbDdhY09OCjVlSDhORXkySEJzTGZUcDNTMmtpVHJqZUMzVDRISitodWZIWTJuU1A3SnBCVDZEcXl5cEhiamRheWdvN1Y5V0MKN3hZSXlob2hBMXBqUU1ZQVZrR1E4ekxNZ284L2kxUGtERnZWQ3JxSFlNZHBJak5uNGdyVE9oaURUek90VHFBSAozenJNaHoxaGhPZm9mQXN2V0ZKTHM5c3krS1FDSHJhcXlRZXdqTERWeThpLzBvSEdXaGpieUFJRjZBMTVucVQyCjhVUkE2b1R0d2J5dXhvVnpIOWZ4N1NwMnA2OVRtWXdMRVFqMmRmbnFRRnVOZ1E9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
  cahex, _ := base64.RawStdEncoding.DecodeString(ca)
  areg := aregv1.MutatingWebhookConfiguration {
    ObjectMeta: metav1.ObjectMeta {
      Name: "hook.test.com",
      Namespace: "default",
    },
    Webhooks: []aregv1.MutatingWebhook {
      {
        Name: "hook.test.com",
        FailurePolicy: &failurepolicy,
        Rules: []aregv1.RuleWithOperations {
          {
            Operations: []aregv1.OperationType {aregv1.Create},
            Rule: aregv1.Rule{
              APIGroups: []string{""},
              APIVersions: []string{"v1"},
              Resources: []string{"pods"},
              Scope: &scope,
            },
          },
        },
        ClientConfig: aregv1.WebhookClientConfig{
          Service: &aregv1.ServiceReference{
            Name: "hook1",
            Namespace: "default",
          },
          CABundle: cahex,
        },
        AdmissionReviewVersions: []string{"v1", "v1beta1"},
        SideEffects: &sideeffect,
        TimeoutSeconds: &timeout,
      },
    }, 
  }



  svc := corev1.Service {
    ObjectMeta: metav1.ObjectMeta {
      Name: "hook1",
      Namespace: "default",
    },
    Spec: corev1.ServiceSpec{
      Selector: map[string]string{
        "app": "hook1",
      },
      Ports: []corev1.ServicePort {
        {
          Protocol: corev1.ProtocolTCP,
          Port: 443,
          TargetPort: intstr.FromInt(443),
        },
      },
    },
  }
  _, err = cs.CoreV1().Services("default").Create(context.Background(), &svc, metav1.CreateOptions{})
  if err != nil {
    log.Println(err)
  }

  pod := corev1.Pod {
    ObjectMeta: metav1.ObjectMeta {
      Name: "hook1",
      Namespace: "default",
      Labels: map[string]string {
        "app": "hook1",
      },
    },
    Spec: corev1.PodSpec{
      Containers: []corev1.Container { {
          Name: "c1",
          Image: "cap:1",
          Command: []string {"/k8cap"},
          Args: []string {"hookserver"},
        },
      },
      HostPID: true,
      ServiceAccountName: "capture",
    },
  }
  _, err = cs.CoreV1().Pods("default").Create(context.Background(), &pod, metav1.CreateOptions{})
  if err != nil {
    log.Println(err)
  }

  time.Sleep(3*time.Second)

  _, err = cs.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.Background(),
  &areg, metav1.CreateOptions{})
  if err != nil {
    log.Println(err)
  }
}

