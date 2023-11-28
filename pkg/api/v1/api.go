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

//===========================================================================
// Authentication Requests and Responses
//===========================================================================

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReply struct {
	ID    int64  `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ReauthenticateRequest struct {
	RefreshToken string `json:"refresh_token"`
}
