package util

type Result struct {
	Error        interface{} `json:"error,omitempty"`
	Message      interface{} `json:"message,omitempty"`
	MessageTh    interface{} `json:"message_th,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Total        int         `json:"total,omitempty"`
	Count        int         `json:"count,omitempty"`
	Status       int         `json:"status,omitempty"`
	TotalHome    int         `json:"total_home,omitempty"`
	TotalUser    int         `json:"total_user,omitempty"`
	TotalUpLevel int         `json:"total_up_level,omitempty"`
}
