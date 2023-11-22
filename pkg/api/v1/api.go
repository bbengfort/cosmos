package api

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are embedded in most API responses
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusReply is returned on status requests and heartbeats
type StatusReply struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
}
