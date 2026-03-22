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

func registerTicketRoutes(rg *gin.RouterGroup, h handlers.TicketHandler) {
	group := rg.Group("/tickets")

	// --- ticket: CRUD y consultas ---
	group.GET("", h.ListTickets)
	group.GET("/:id", h.GetByID)
	group.GET("/nro/:nro", h.GetByNroTicket)
	group.POST("", h.Create)
	group.PATCH("/:id", h.UpdateTicket)

	// --- ticket: flujo de estados ---
	group.PATCH("/:id/asignar", h.Assign)
	group.PATCH("/:id/estado", h.ChangeEstado)
	group.PATCH("/:id/cerrar", h.Close)

	// --- ticket: bitácora ---
	group.GET("/:id/bitacora", h.ListBitacora)
	group.POST("/:id/bitacora", h.CreateBitacora)

	// --- ticket: pausas ---
	group.GET("/:id/pausas", h.ListPausas)
	group.POST("/:id/pausas", h.CreatePausa)
	group.PATCH("/:id/reanudar", h.ReanudarTicket)

	// --- ticket: traspasos ---
	group.GET("/:id/traspasos", h.ListTraspasos)
	group.POST("/:id/traspasos", h.CreateTraspaso)

	// --- pausas (por id de pausa, no de ticket) ---
	pausas := rg.Group("/pausas")
	pausas.PATCH("/:id/resolver", h.ResolverPausa)

	// --- traspasos (por id de traspaso, no de ticket) ---
	traspasos := rg.Group("/traspasos")
	traspasos.PATCH("/:id/resolver", h.ResolverTraspaso)
}

func registerTecnicoRoutes(rg *gin.RouterGroup, h handlers.TecnicoHandler) {
	// --- técnicos ---
	group := rg.Group("/tecnicos")
	group.GET("", h.List)
	group.GET("/rut/:rut", h.GetByRut)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PATCH("/:id", h.Update)

	// --- configuración horarios turno ---
	cht := rg.Group("/configuracion-horarios-turno")
	cht.GET("", h.ListHorariosTurno)
	cht.POST("", h.CreateHorarioTurno)
	cht.PATCH("/:id", h.UpdateHorarioTurno)
}

func registerServicioRoutes(rg *gin.RouterGroup, h handlers.ServicioHandler) {
	group := rg.Group("/servicios")
	group.GET("", h.List)
	group.POST("", h.Create)
	group.PATCH("/:id", h.Update)
}

func registerCatalogoFallaRoutes(rg *gin.RouterGroup, h handlers.CatalogoFallaHandler) {
	group := rg.Group("/catalogo-fallas")
	group.GET("", h.List)
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
