package protocol

type Request struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type Response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}
