package structs

type ValidateLogin struct {
	Phone string `json:"phone" form:"phone" binding:"required,min=4"`
}
