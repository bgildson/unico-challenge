package response

// Generic represents a generic response for errors
type Generic struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
