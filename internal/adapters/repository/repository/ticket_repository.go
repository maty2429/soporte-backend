package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) ports.TicketRepository {
	if db == nil {
		return unavailableTicketRepository{}
	}
	return &TicketRepository{db: db}
}

func (r *TicketRepository) RunInTx(ctx context.Context, fn func(txRepo ports.TicketRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &TicketRepository{db: tx}
		return fn(txRepo)
	})
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	row := toTicketDB(*ticket)

	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("ticket", err)
	}

	ticket.ID = row.ID
	ticket.CreatedAt = row.CreatedAt
	ticket.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *TicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	updates := map[string]any{
		"cod_estado_ticket": ticket.CodEstadoTicket,
	}
	if ticket.IDTecnicoAsignado != nil {
		updates["id_tecnico_asignado"] = *ticket.IDTecnicoAsignado
	}
	if ticket.IDCatalogoFalla != nil {
		updates["id_catalogo_falla"] = *ticket.IDCatalogoFalla
	}
	if ticket.IDNivelPrioridad != nil {
		updates["id_nivel_prioridad"] = *ticket.IDNivelPrioridad
	}

	if err := r.db.WithContext(ctx).Model(&dbmodels.Ticket{}).Where("id = ?", ticket.ID).Updates(updates).Error; err != nil {
		return wrapUpdateError("ticket", err)
	}
	return nil
}

func (r *TicketRepository) UpdateFields(ctx context.Context, ticket *domain.Ticket, fields map[string]any) error {
	if err := r.db.WithContext(ctx).Model(&dbmodels.Ticket{}).Where("id = ?", ticket.ID).Updates(fields).Error; err != nil {
		return wrapUpdateError("ticket", err)
	}
	return nil
}

func (r *TicketRepository) ListTickets(ctx context.Context, f ports.ListTicketsFilters) ([]domain.Ticket, int64, error) {
	q := applyListTicketsFilters(r.db.WithContext(ctx).Model(&dbmodels.Ticket{}), f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, wrapListError("tickets", err)
	}

	var rows []dbmodels.Ticket
	dataQuery := preloadTicketRelations(applyListTicketsFilters(r.db.WithContext(ctx).Model(&dbmodels.Ticket{}), f))
	if err := dataQuery.Order("ticket.created_at DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("tickets", err)
	}

	items := make([]domain.Ticket, 0, len(rows))
	for _, row := range rows {
		items = append(items, toTicketDomain(row))
	}
	return items, total, nil
}

func (r *TicketRepository) GetByID(ctx context.Context, id int) (domain.Ticket, error) {
	var row dbmodels.Ticket
	if err := preloadTicketRelations(r.db.WithContext(ctx)).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.Ticket{}, wrapDBError("ticket", err)
	}
	return toTicketDomain(row), nil
}

func (r *TicketRepository) GetByNroTicket(ctx context.Context, nro string) (domain.Ticket, error) {
	var row dbmodels.Ticket
	if err := preloadTicketRelations(r.db.WithContext(ctx)).Where("nro_ticket = ?", nro).Take(&row).Error; err != nil {
		return domain.Ticket{}, wrapDBError("ticket", err)
	}
	return toTicketDomain(row), nil
}

func (r *TicketRepository) NroTicketExists(ctx context.Context, nro string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&dbmodels.Ticket{}).Where("nro_ticket = ?", nro).Count(&count).Error
	if err != nil {
		return false, wrapDBError("ticket", err)
	}
	return count > 0, nil
}

func (r *TicketRepository) GetEstadoTicketByCod(ctx context.Context, cod string) (domain.EstadoTicket, error) {
	var row dbmodels.EstadoTicket
	if err := r.db.WithContext(ctx).Where("cod_estado_ticket = ?", cod).Take(&row).Error; err != nil {
		return domain.EstadoTicket{}, wrapDBError("estado ticket", err)
	}
	return domain.EstadoTicket{
		ID:              row.ID,
		Descripcion:     row.Descripcion,
		CodEstadoTicket: row.CodEstadoTicket,
	}, nil
}

func (r *TicketRepository) CreateTrazabilidad(ctx context.Context, t *domain.TrazabilidadTicket) error {
	row := dbmodels.TrazabilidadTicket{
		IDTicket:        t.IDTicket,
		CodEstadoTicket: t.CodEstadoTicket,
		RutResponsable:  t.RutResponsable,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("trazabilidad ticket", err)
	}
	t.ID = row.ID
	t.FechaTrazabilidad = row.FechaTrazabilidad
	return nil
}

func (r *TicketRepository) ListTrazabilidad(ctx context.Context, idTicket int) ([]domain.TrazabilidadTicket, error) {
	type trazRow struct {
		ID                int
		IDTicket          int
		CodEstadoTicket   string
		DescripcionEstado string
		RutResponsable    string
		FechaTrazabilidad time.Time
	}

	var rows []trazRow
	err := r.db.WithContext(ctx).
		Table("trazabilidad_ticket t").
		Select("t.id, t.id_ticket, t.cod_estado_ticket, e.descripcion AS descripcion_estado, t.rut_responsable, t.fecha_trazabilidad").
		Joins("LEFT JOIN estado_ticket e ON e.cod_estado_ticket = t.cod_estado_ticket").
		Where("t.id_ticket = ?", idTicket).
		Order("t.fecha_trazabilidad ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, wrapListError("trazabilidad ticket", err)
	}

	items := make([]domain.TrazabilidadTicket, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.TrazabilidadTicket{
			ID:                row.ID,
			IDTicket:          row.IDTicket,
			CodEstadoTicket:   row.CodEstadoTicket,
			DescripcionEstado: row.DescripcionEstado,
			RutResponsable:    row.RutResponsable,
			FechaTrazabilidad: row.FechaTrazabilidad,
		})
	}
	return items, nil
}

func (r *TicketRepository) ListBitacora(ctx context.Context, idTicket int) ([]domain.BitacoraTicket, error) {
	var rows []dbmodels.BitacoraTicket
	if err := r.db.WithContext(ctx).Where("id_ticket = ?", idTicket).Order("fecha_registro ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("bitacora ticket", err)
	}
	items := make([]domain.BitacoraTicket, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.BitacoraTicket{
			ID:            row.ID,
			IDTicket:      row.IDTicket,
			RutAutor:      row.RutAutor,
			Comentario:    row.Comentario,
			FechaRegistro: row.FechaRegistro,
		})
	}
	return items, nil
}

func (r *TicketRepository) CreateBitacora(ctx context.Context, b *domain.BitacoraTicket) error {
	row := dbmodels.BitacoraTicket{
		IDTicket:   b.IDTicket,
		RutAutor:   b.RutAutor,
		Comentario: b.Comentario,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("bitacora ticket", err)
	}
	b.ID = row.ID
	b.FechaRegistro = row.FechaRegistro
	return nil
}

func (r *TicketRepository) CreateValorizacion(ctx context.Context, v *domain.Valorizacion) error {
	row := dbmodels.Valorizacion{
		IDTicket:      v.IDTicket,
		IDTecnico:     v.IDTecnico,
		IDSolicitante: v.IDSolicitante,
		Nota:          v.Nota,
		Comentarios:   v.Comentarios,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("valorizacion", err)
	}
	v.ID = row.ID
	v.CreatedAt = row.CreatedAt
	return nil
}

func (r *TicketRepository) CreatePausa(ctx context.Context, p *domain.TicketPausa) error {
	row := dbmodels.TicketPausa{
		IDTicket:       p.IDTicket,
		IDTecnicoPausa: p.IDTecnicoPausa,
		EstadoPausa:    p.EstadoPausa,
		IDMotivoPausa:  p.IDMotivoPausa,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("ticket pausa", err)
	}
	p.ID = row.ID
	p.FechaInicioPausa = row.FechaInicioPausa
	return nil
}

func (r *TicketRepository) GetPausaByID(ctx context.Context, id int) (domain.TicketPausa, error) {
	var row dbmodels.TicketPausa
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.TicketPausa{}, wrapDBError("ticket pausa", err)
	}
	return toPausaDomain(row), nil
}

func (r *TicketRepository) GetPausaActiva(ctx context.Context, idTicket int) (domain.TicketPausa, error) {
	var row dbmodels.TicketPausa
	if err := r.db.WithContext(ctx).
		Where("id_ticket = ? AND fecha_fin_pausa IS NULL AND estado_pausa = ?", idTicket, "APROBADA").
		Take(&row).Error; err != nil {
		return domain.TicketPausa{}, wrapDBError("ticket pausa", err)
	}
	return toPausaDomain(row), nil
}

func (r *TicketRepository) UpdatePausa(ctx context.Context, p *domain.TicketPausa) error {
	updates := map[string]any{
		"estado_pausa": p.EstadoPausa,
	}
	if p.IDTecnicoAutorizado != nil {
		updates["id_tecnico_autorizado"] = *p.IDTecnicoAutorizado
	}
	if p.FechaFinPausa != nil {
		updates["fecha_fin_pausa"] = *p.FechaFinPausa
	}
	if p.FechaResolucion != nil {
		updates["fecha_resolucion"] = *p.FechaResolucion
	}
	if err := r.db.WithContext(ctx).Model(&dbmodels.TicketPausa{}).Where("id = ?", p.ID).Updates(updates).Error; err != nil {
		return wrapUpdateError("ticket pausa", err)
	}
	return nil
}

func (r *TicketRepository) ListPausas(ctx context.Context, f ports.ListPausasFilters) ([]domain.TicketPausa, int64, error) {
	q := r.db.WithContext(ctx).Model(&dbmodels.TicketPausa{}).Where("id_ticket = ?", f.IDTicket)

	if f.Estado != "" {
		q = q.Where("estado_pausa = ?", f.Estado)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, wrapListError("ticket pausas", err)
	}

	var rows []dbmodels.TicketPausa
	if err := q.Order("fecha_inicio_pausa DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("ticket pausas", err)
	}

	items := make([]domain.TicketPausa, 0, len(rows))
	for _, row := range rows {
		items = append(items, toPausaDomain(row))
	}
	return items, total, nil
}

func toPausaDomain(row dbmodels.TicketPausa) domain.TicketPausa {
	return domain.TicketPausa{
		ID:                  row.ID,
		IDTicket:            row.IDTicket,
		FechaInicioPausa:    row.FechaInicioPausa,
		FechaFinPausa:       row.FechaFinPausa,
		IDTecnicoPausa:      row.IDTecnicoPausa,
		EstadoPausa:         row.EstadoPausa,
		IDMotivoPausa:       row.IDMotivoPausa,
		IDTecnicoAutorizado: row.IDTecnicoAutorizado,
		FechaResolucion:     row.FechaResolucion,
	}
}

// --- traspasos ---

func (r *TicketRepository) CreateTraspaso(ctx context.Context, t *domain.TicketTraspaso) error {
	row := dbmodels.TicketTraspaso{
		IDTicket:         t.IDTicket,
		IDTecnicoOrigen:  t.IDTecnicoOrigen,
		IDTecnicoDestino: t.IDTecnicoDestino,
		EstadoTraspaso:   t.EstadoTraspaso,
		Motivo:           t.Motivo,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("ticket traspaso", err)
	}
	t.ID = row.ID
	t.FechaSolicitud = row.FechaSolicitud
	return nil
}

func (r *TicketRepository) GetTraspasoByID(ctx context.Context, id int) (domain.TicketTraspaso, error) {
	var row dbmodels.TicketTraspaso
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.TicketTraspaso{}, wrapDBError("ticket traspaso", err)
	}
	return toTraspasoDomain(row), nil
}

func (r *TicketRepository) GetTraspasoPendiente(ctx context.Context, idTicket int) (domain.TicketTraspaso, error) {
	var row dbmodels.TicketTraspaso
	if err := r.db.WithContext(ctx).
		Where("id_ticket = ? AND cod_traslado = ?", idTicket, "PENDIENTE").
		Take(&row).Error; err != nil {
		return domain.TicketTraspaso{}, wrapDBError("ticket traspaso", err)
	}
	return toTraspasoDomain(row), nil
}

func (r *TicketRepository) UpdateTraspaso(ctx context.Context, t *domain.TicketTraspaso) error {
	updates := map[string]any{
		"cod_traslado": t.EstadoTraspaso,
	}
	if t.ComentarioResolucion != "" {
		updates["motivo_respuesta"] = t.ComentarioResolucion
	}
	if t.FechaResolucion != nil {
		updates["fecha_respuesta"] = *t.FechaResolucion
	}
	if err := r.db.WithContext(ctx).Model(&dbmodels.TicketTraspaso{}).Where("id = ?", t.ID).Updates(updates).Error; err != nil {
		return wrapUpdateError("ticket traspaso", err)
	}
	return nil
}

func (r *TicketRepository) ListTraspasos(ctx context.Context, f ports.ListTraspasosFilters) ([]domain.TicketTraspaso, int64, error) {
	q := r.db.WithContext(ctx).Model(&dbmodels.TicketTraspaso{}).Where("id_tecnico_destino = ?", f.IDTecnicoDestino)

	if f.Estado != "" {
		q = q.Where("cod_traslado = ?", f.Estado)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, wrapListError("ticket traspasos", err)
	}

	var rows []dbmodels.TicketTraspaso
	if err := q.Order("fecha_solicitud DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("ticket traspasos", err)
	}

	items := make([]domain.TicketTraspaso, 0, len(rows))
	for _, row := range rows {
		items = append(items, toTraspasoDomain(row))
	}
	return items, total, nil
}

func toTraspasoDomain(row dbmodels.TicketTraspaso) domain.TicketTraspaso {
	return domain.TicketTraspaso{
		ID:                   row.ID,
		IDTicket:             row.IDTicket,
		IDTecnicoOrigen:      row.IDTecnicoOrigen,
		IDTecnicoDestino:     row.IDTecnicoDestino,
		EstadoTraspaso:       row.EstadoTraspaso,
		Motivo:               row.Motivo,
		ComentarioResolucion: row.ComentarioResolucion,
		FechaSolicitud:       row.FechaSolicitud,
		FechaResolucion:      row.FechaResolucion,
	}
}

// --- conversiones ---

func toTicketDB(t domain.Ticket) dbmodels.Ticket {
	return dbmodels.Ticket{
		NroTicket:             t.NroTicket,
		IDSolicitante:         t.IDSolicitante,
		IDTecnicoAsignado:     t.IDTecnicoAsignado,
		IDServicio:            t.IDServicio,
		IDTipoTicket:          t.IDTipoTicket,
		CodEstadoTicket:       t.CodEstadoTicket,
		IDNivelPrioridad:      t.IDNivelPrioridad,
		IDCatalogoFalla:       t.IDCatalogoFalla,
		IDDepartamentoSoporte: t.IDDepartamentoSoporte,
		Critico:               t.Critico,
		DetalleFallaReportada: t.DetalleFallaReportada,
		UbicacionObs:          t.UbicacionObs,
	}
}

func toTicketDomain(row dbmodels.Ticket) domain.Ticket {
	return domain.Ticket{
		ID:                    row.ID,
		NroTicket:             row.NroTicket,
		IDSolicitante:         row.IDSolicitante,
		IDTecnicoAsignado:     row.IDTecnicoAsignado,
		IDServicio:            row.IDServicio,
		IDTipoTicket:          row.IDTipoTicket,
		CodEstadoTicket:       row.CodEstadoTicket,
		IDNivelPrioridad:      row.IDNivelPrioridad,
		IDCatalogoFalla:       row.IDCatalogoFalla,
		IDDepartamentoSoporte: row.IDDepartamentoSoporte,
		Critico:               row.Critico,
		DetalleFallaReportada: row.DetalleFallaReportada,
		UbicacionObs:          row.UbicacionObs,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
		FechaInicioTrabajo:    row.FechaInicioTrabajo,
		FechaFinTrabajo:       row.FechaFinTrabajo,
		Solicitante:           toSolicitantePtr(row.Solicitante),
		TecnicoAsignado:       toTecnicoPtr(row.TecnicoAsignado),
		Servicio:              toServicioPtr(row.Servicio),
		TipoTicket:            toTipoTicketPtr(row.TipoTicket),
		EstadoTicket:          toEstadoTicketPtr(row.EstadoTicket),
		NivelPrioridad:        toNivelPrioridadPtr(row.NivelPrioridad),
		CatalogoFalla:         toCatalogoFallaPtr(row.CatalogoFalla),
		DepartamentoSoporte:   toDepartamentoSoportePtr(row.DepartamentoSoporte),
	}
}

func applyListTicketsFilters(q *gorm.DB, f ports.ListTicketsFilters) *gorm.DB {
	if f.CodEstadoTicket != "" {
		q = q.Where("cod_estado_ticket = ?", f.CodEstadoTicket)
	}
	if f.IDTecnicoAsignado > 0 {
		q = q.Where("id_tecnico_asignado = ?", f.IDTecnicoAsignado)
	}
	if f.RutTecnico != "" && f.DVTecnico != "" {
		q = q.Joins("JOIN tecnicos ON tecnicos.id = id_tecnico_asignado").
			Where("tecnicos.rut = ? AND tecnicos.dv = ?", f.RutTecnico, f.DVTecnico)
	}
	if f.IDSolicitante > 0 {
		q = q.Where("id_solicitante = ?", f.IDSolicitante)
	}
	if f.IDDepartamentoSoporte > 0 {
		q = q.Where("id_departamento_soporte = ?", f.IDDepartamentoSoporte)
	}
	if f.Critico != nil {
		q = q.Where("critico = ?", *f.Critico)
	}
	return q
}

func preloadTicketRelations(q *gorm.DB) *gorm.DB {
	return q.
		Preload("Solicitante").
		Preload("TecnicoAsignado").
		Preload("Servicio").
		Preload("TipoTicket").
		Preload("EstadoTicket").
		Preload("NivelPrioridad").
		Preload("CatalogoFalla").
		Preload("DepartamentoSoporte")
}

func toSolicitantePtr(row *dbmodels.Solicitante) *domain.Solicitante {
	if row == nil {
		return nil
	}
	item := toSolicitanteDomain(*row)
	return &item
}

func toTecnicoPtr(row *dbmodels.Tecnico) *domain.Tecnico {
	if row == nil {
		return nil
	}
	item := toTecnicoDomain(*row)
	return &item
}

func toTipoTicketPtr(row *dbmodels.TipoTicket) *domain.TipoTicket {
	if row == nil {
		return nil
	}
	return &domain.TipoTicket{
		ID:            row.ID,
		CodTipoTicket: row.CodTipoTicket,
		Descripcion:   row.Descripcion,
	}
}

func toEstadoTicketPtr(row *dbmodels.EstadoTicket) *domain.EstadoTicket {
	if row == nil {
		return nil
	}
	return &domain.EstadoTicket{
		ID:              row.ID,
		Descripcion:     row.Descripcion,
		CodEstadoTicket: row.CodEstadoTicket,
	}
}

func toNivelPrioridadPtr(row *dbmodels.NivelPrioridad) *domain.NivelPrioridad {
	if row == nil {
		return nil
	}
	return &domain.NivelPrioridad{
		ID:          row.ID,
		Descripcion: row.Descripcion,
	}
}

func toCatalogoFallaPtr(row *dbmodels.CatalogoFalla) *domain.CatalogoFalla {
	if row == nil {
		return nil
	}
	item := toCatalogoFallaDomain(*row)
	return &item
}

func toDepartamentoSoportePtr(row *dbmodels.DepartamentoSoporte) *domain.DepartamentoSoporte {
	if row == nil {
		return nil
	}
	return &domain.DepartamentoSoporte{
		ID:              row.ID,
		CodDepartamento: row.CodDepartamento,
		Descripcion:     row.Descripcion,
	}
}
