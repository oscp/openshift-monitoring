package models

type BaseModel struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}
