package core

type LoadBalancerRule struct {
	ID              string     `json:"id,omitempty"`
	Algorithm       string     `json:"algorithm,omitempty"` // TODO: replace with constrained type?
	DestinationPort int        `json:"destination_port,omitempty"`
	ListenPort      int        `json:"listen_port,omitempty"`
	Protocol        string     `json:"protocol,omitempty"`     // TODO: replace with type?
	Certificates    []struct{} `json:"certificates,omitempty"` // TODO: is this the same certificate type as certificate.go
	BackendSSL      bool       `json:"backend_ssl,omitempty"`
	PassthroughSSL  bool       `json:"passthrough_ssl,omitempty"`
	CheckEnabled    bool       `json:"check_enabled,omitempty"`
	CheckFall       int        `json:"check_fall,omitempty"`
	CheckInterval   int        `json:"check_interval,omitempty"`
	CheckPath       string     `json:"check_path,omitempty"`
	CheckProtocol   string     `json:"check_protocol,omitempty"` // TODO: replace with type?
	CheckRise       int        `json:"check_rise,omitempty"`
	CheckTimeout    int        `json:"check_timeout,omitempty"`
}

type LoadBalancerRuleCreateArguments struct {
}
