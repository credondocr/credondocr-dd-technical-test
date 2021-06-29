package dtos

type Params struct {
	From   string `form:"from" binding:"required"`
	Until  string `form:"until" binding:"required"`
	Artist string `form:"artist"`
}
