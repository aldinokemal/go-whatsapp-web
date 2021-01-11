package structs

type ValidateLogin struct {
	AppID string `json:"app_id" form:"app_id" binding:"required,min=4"`
}
