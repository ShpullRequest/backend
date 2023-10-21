package handlers

import (
	"github.com/SevereCloud/vksdk/v2/vkapps"
	"github.com/gin-gonic/gin"
)

func (hs *handlerService) GetVKParams(ctx *gin.Context) *vkapps.Params {
	vkParams, _ := ctx.Get("vk_params")

	return vkParams.(*vkapps.Params)
}
