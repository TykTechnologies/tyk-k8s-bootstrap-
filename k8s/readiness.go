package k8s

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

func (c *Client) CheckIfRequiredDeploymentsAreReady() error {
	time.Sleep(5 * time.Second)

	var attemptCount int
	for {
		attemptCount++
		if attemptCount > 180 {
			return errors.New("attempted readiness check too many times")
		}
		pods, err := c.clientSet.
			CoreV1().
			Pods(c.AppArgs.TykPodNamespace).
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		
		fmt.Printf("There are %d other pods in the cluster\n", len(pods.Items)-1)

		var requiredPods []v1.Pod
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "dashboard") ||
				strings.Contains(pod.Name, "redis") {
				requiredPods = append(requiredPods, pod)
			}
		}

		notReadyPods := make(map[string]struct{})
		for _, pod := range requiredPods {
			podStatus := pod.Status
			for container := range pod.Spec.Containers {
				if !podStatus.ContainerStatuses[container].Ready {
					notReadyPods[pod.Name] = struct{}{}
				}
			}
		}
		fmt.Printf("The following pods have containers that are not ready: ")
		for pod, _ := range notReadyPods {
			fmt.Println(pod)
		}

		if len(notReadyPods) != 0 {
			fmt.Println("TYK PRO IS NOT READY")
		} else {
			fmt.Println("TYK PRO IS READY TO BE BOOTSTRAPPED")
			return nil
		}
		time.Sleep(2 * time.Second)
	}
}