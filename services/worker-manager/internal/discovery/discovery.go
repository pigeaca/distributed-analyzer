package discovery

type Service struct {
	Name      string
	Port      string
	Addr      string
	Namespace string
}

type ServiceDiscovery interface {
	DiscoverServices(serviceName string) ([]*Service, error)
}
