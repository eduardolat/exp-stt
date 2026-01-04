package api

import (
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

func (h *handlers) registerProcVersionGet() {
	h.uforpcServer.Procs.VersionGet.Handle(func(c *uforpc.VersionGetHandlerContext[urpcProps]) (uforpc.VersionGetOutput, error) {
		return uforpc.VersionGetOutput{AppVersion: config.AppVersion}, nil
	})
}
