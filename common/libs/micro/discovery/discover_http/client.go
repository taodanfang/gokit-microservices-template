package discover_http

import (
	"cctable/common/tools"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"strconv"
	"sync"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_discovery_http_client interface {
	Register(service_name, instance_id string, host string, port string) bool
	Deregister(instance_id string) bool
	Discover_services(service_name string) []interface{}
}

type Kit_discovery_http_client struct {
	Host          string
	Port          string
	client        consul.Client
	config        *api.Config
	mutex         sync.Mutex
	instances_map sync.Map
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_kit_discovery_http_client(host string, port string) (I_discovery_http_client, error) {

	consul_config := api.DefaultConfig()
	consul_config.Address = host + ":" + port
	api_client, err := api.NewClient(consul_config)
	if err != nil {
		return nil, err
	}

	client := consul.NewClient(api_client)

	return &Kit_discovery_http_client{
		Host:   host,
		Port:   port,
		config: consul_config,
		client: client,
	}, nil
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (c *Kit_discovery_http_client) Register(service_name, instance_id string, host string, port string) bool {

	int_port, _ := strconv.Atoi(port)
	service_registration := &api.AgentServiceRegistration{
		ID:      instance_id,
		Name:    service_name,
		Address: host,
		Port:    int_port,
		Tags:    []string{"http"},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + host + ":" + port + "/health",
			Interval:                       "15s",
		},
	}

	err := c.client.Register(service_registration)
	if err != nil {
		tools.Log("Register service " + service_name + " failed!")
		return false
	}

	tools.Log("Register service " + service_name + " success!")
	return true
}

func (c *Kit_discovery_http_client) Deregister(instance_id string) bool {

	service_registration := &api.AgentServiceRegistration{
		ID: instance_id,
	}

	err := c.client.Deregister(service_registration)
	if err != nil {
		tools.Log("Deregister service failed!")
		return false
	}

	tools.Log("Deregister service success!")
	return true
}

func (c *Kit_discovery_http_client) Discover_services(service_name string) []interface{} {

	instance_list, ok := c.instances_map.Load(service_name)
	if ok {
		return instance_list.([]interface{})
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	instance_list, ok = c.instances_map.Load(service_name)
	if ok {
		return instance_list.([]interface{})
	}

	go func() {
		params := make(map[string]interface{})
		params["type"] = "service"
		params["service"] = service_name
		plan, _ := watch.Parse(params)
		plan.Handler = func(u uint64, i interface{}) {
			if i == nil {
				return
			}

			v, ok := i.([]*api.ServiceEntry)
			if !ok {
				return
			}

			if len(v) == 0 {
				c.instances_map.Store(service_name, []interface{}{})
			}

			var health_services []interface{}
			for _, service := range v {
				if service.Checks.AggregatedStatus() == api.HealthPassing {
					health_services = append(health_services, service.Service)
				}
			}
			c.instances_map.Store(service_name, health_services)
		}

		defer plan.Stop()
		_ = plan.Run(c.config.Address)
	}()

	entries, _, err := c.client.Service(service_name, "http", false, nil)
	if err != nil {
		c.instances_map.Store(service_name, []interface{}{})
		tools.Log("Discover service error!")
		return nil
	}

	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}

	c.instances_map.Store(service_name, instances)
	return instances
}
