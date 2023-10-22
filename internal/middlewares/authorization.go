package middlewares

import (
	"github.com/SevereCloud/vksdk/v2/vkapps"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/internal/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (ms *middlewareService) Authorization(ctx *gin.Context) {
	if strings.HasPrefix(ctx.Request.RequestURI, "/swagger/") {
		return
	}

	authString := ctx.GetHeader("Authorization")
	authString = strings.ReplaceAll(authString, "Bearer ", "")

	ok, err := vkapps.ParamsVerify(authString, config.Config.AppSecretToken)
	if err != nil {
		ms.logger.Error("Failed params verify", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, errs.NewInternalServer("An unknown error occurred while checking authorization"))
		ctx.Abort()

		return
	}

	if ok {
		u, _ := url.Parse(authString)
		params, _ := vkapps.NewParams(u)

		ts, err := strconv.Atoi(params.VkTs)
		if err != nil {
			ms.logger.Error("Failed convert VkTs (string) to int", zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, errs.NewInternalServer("An unknown error occurred while checking authorization"))
			ctx.Abort()

			return
		}

		timeCreateSign := time.Unix(int64(ts), 0)
		sinceCreateSign := time.Since(timeCreateSign)

		if sinceCreateSign.Hours() >= 1 && config.Config.ProdFlag {
			ms.logger.Debug("Authorization failed, signature expired", zap.Duration("SinceCreateSign", sinceCreateSign))

			ctx.JSON(http.StatusUnauthorized, errs.NotAuthorized("Authorization failed, signature expired"))
			ctx.Abort()

			return
		}

		ms.logger.Debug("Authorization success", zap.Bool("ParamsVerify", ok))

		ctx.Set("vk_params", params)
		ctx.Next()
	} else {
		ms.logger.Debug("Authorization failed", zap.Bool("ParamsVerify", ok))

		ctx.JSON(http.StatusUnauthorized, errs.NotAuthorized("Authorization failed"))
		ctx.Abort()

		return
	}
}
