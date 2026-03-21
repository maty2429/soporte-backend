package routes

import (
	"github.com/gin-gonic/gin"
	"soporte/internal/delivery/http/handlers"
)

func registerSolicitanteRoutes(rg *gin.RouterGroup, h handlers.SolicitanteHandler) {
	group := rg.Group("/solicitantes")
	group.GET("", h.List)
	group.GET("/rut/:rut", h.GetByRut)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PATCH("/:id", h.Update)
}

func registerCatalogoRoutes(rg *gin.RouterGroup, h handlers.CatalogoHandler) {
	// tipos_ticket
	tt := rg.Group("/tipos-ticket")
	tt.GET("", h.ListTiposTicket)
	tt.POST("", h.CreateTipoTicket)
	tt.PATCH("/:id", h.UpdateTipoTicket)

	// niveles_prioridad
	np := rg.Group("/niveles-prioridad")
	np.GET("", h.ListNivelesPrioridad)
	np.POST("", h.CreateNivelPrioridad)
	np.PATCH("/:id", h.UpdateNivelPrioridad)

	// tipo_tecnico
	tc := rg.Group("/tipos-tecnico")
	tc.GET("", h.ListTiposTecnico)
	tc.POST("", h.CreateTipoTecnico)
	tc.PATCH("/:id", h.UpdateTipoTecnico)

	// departamentos_soporte
	ds := rg.Group("/departamentos-soporte")
	ds.GET("", h.ListDepartamentosSoporte)
	ds.POST("", h.CreateDepartamentoSoporte)
	ds.PATCH("/:id", h.UpdateDepartamentoSoporte)

	// motivos_pausa
	mp := rg.Group("/motivos-pausa")
	mp.GET("", h.ListMotivosPausa)
	mp.POST("", h.CreateMotivoPausa)
	mp.PATCH("/:id", h.UpdateMotivoPausa)

	// tipos_turno
	tu := rg.Group("/tipos-turno")
	tu.GET("", h.ListTiposTurno)
	tu.POST("", h.CreateTipoTurno)
	tu.PATCH("/:id", h.UpdateTipoTurno)
}
