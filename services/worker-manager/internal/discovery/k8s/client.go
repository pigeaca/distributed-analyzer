package k8s

import (
	"context"
	"distributed-analyzer/services/worker-manager/internal/discovery"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ServiceDiscoveryImpl struct {
	clientset *kubernetes.Clientset
}

func NewServiceDiscovery() discovery.ServiceDiscovery {
	clientset, err := initKubeClient()
	if err != nil {
		log.Fatalf("failed to init kube client: %v", err)
	}
	return &ServiceDiscoveryImpl{clientset: clientset}
}

func initKubeClient() (*kubernetes.Clientset, error) {
	var cfg *rest.Config
	var err error

	if _, err = os.Stat("/var/run/secrets/kubernetes.io/worker/token"); err == nil {
		cfg, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func (d *ServiceDiscoveryImpl) DiscoverServices(serviceName string) ([]*discovery.Service, error) {
	svcs, err := d.clientset.CoreV1().Services(serviceName).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("failed to list services: %v", err)
	}

	services := make([]*discovery.Service, 0)

	for _, svc := range svcs.Items {
		for _, port := range svc.Spec.Ports {
			address := fmt.Sprintf("%s.%s.svc:%d", svc.Name, svc.Namespace, port.Port)
			services = append(services, &discovery.Service{
				Name:      svc.Name,
				Port:      port.Name,
				Addr:      address,
				Namespace: svc.Namespace,
			})
		}
	}

	return services, nil
}
