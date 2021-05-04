package pluginproviderdemo_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/traefik/genconf/dynamic"
	"github.com/traefik/genconf/dynamic/tls"
	"github.com/traefik/pluginproviderdemo"
)

func TestNew(t *testing.T) {
	config := pluginproviderdemo.CreateConfig()
	config.PollInterval = "1s"

	provider, err := pluginproviderdemo.New(context.Background(), config, "test")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = provider.Stop()
		if err != nil {
			t.Fatal(err)
		}
	})

	err = provider.Init()
	if err != nil {
		t.Fatal(err)
	}

	cfgChan := make(chan json.Marshaler)

	err = provider.Provide(cfgChan)
	if err != nil {
		t.Fatal(err)
	}

	data := <-cfgChan

	expected := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"pp-route-01": {
					EntryPoints: []string{"web"},
					Service:     "pp-service-01",
					Rule:        "Host(`example.com`)",
				},
			},
			Services: map[string]*dynamic.Service{
				"pp-service-01": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{
								URL: "http://localhost:9090",
							},
						},
						PassHostHeader: boolPtr(true),
					},
				},
			},
			Middlewares:       make(map[string]*dynamic.Middleware),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
	}

	if time.Now().Minute()%2 == 0 {
		expected.HTTP.Routers["pp-route-02"] = &dynamic.Router{
			EntryPoints: []string{"web"},
			Service:     "pp-service-02",
			Rule:        "Host(`another.example.com`)",
		}

		expected.HTTP.Services["pp-service-02"] = &dynamic.Service{
			LoadBalancer: &dynamic.ServersLoadBalancer{
				Servers: []dynamic.Server{
					{
						URL: "http://localhost:9091",
					},
				},
				PassHostHeader: boolPtr(true),
			},
		}
	}

	expectedJSON, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expectedJSON, dataJSON) {
		t.Fatalf("got %s, want: %s", string(dataJSON), string(expectedJSON))
	}
}

func boolPtr(v bool) *bool {
	return &v
}
