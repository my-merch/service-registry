// service-registry/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
)

func registerService(client *api.Client, serviceID, serviceName string, servicePort int) {
	registration := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		Port: servicePort,
		Check: &api.AgentServiceCheck{ // Health check
			Interval:                       string(15 * time.Second),
			HTTP:                           "http://localhost:" + strconv.Itoa(servicePort) + "/health",
			Timeout:                        string(3 * time.Second),
			DeregisterCriticalServiceAfter: string(30 * time.Second),
		},
	}

	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register %s service: %v", serviceName, err)
	}
}

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Register services
	registerService(client, "product-service", "product-service", 8081)
	registerService(client, "order-service", "order-service", 8082)
	registerService(client, "user-service", "user-service", 8083)

	// Handle interrupt signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// Deregister services on shutdown
	client.Agent().ServiceDeregister("product-service")
	client.Agent().ServiceDeregister("order-service")
	client.Agent().ServiceDeregister("user-service")
}
