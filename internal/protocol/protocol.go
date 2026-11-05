package protocol

type OpenRequest struct {
	User string `json:"user"`
	Host string `json:"host"`
	Path string `json:"path"`
}
