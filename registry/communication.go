package registry

// Struct used to send information to the services
type SendCommand struct {
	Command  string `json:"command"`
	User     string `json:"user"`
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Args     string `json:"args"`
}

type Response struct {
	Response string `json:"response"`
}
