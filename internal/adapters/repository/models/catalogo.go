package models

type TipoTicket struct {
	ID            int    `gorm:"column:id;primaryKey"`
	CodTipoTicket string `gorm:"column:cod_tipo_ticket;not null;uniqueIndex"`
	Descripcion   string `gorm:"column:descripcion;not null;uniqueIndex"`
}

func (TipoTicket) TableName() string { return "tipo_ticket" }

type NivelPrioridad struct {
	ID          int    `gorm:"column:id;primaryKey"`
	Descripcion string `gorm:"column:descripcion;not null;uniqueIndex"`
}

func (NivelPrioridad) TableName() string { return "niveles_prioridad" }

type TipoTecnico struct {
	ID          int    `gorm:"column:id;primaryKey"`
	Descripcion string `gorm:"column:descripcion;not null;uniqueIndex"`
}

func (TipoTecnico) TableName() string { return "tipo_tecnico" }

type DepartamentoSoporte struct {
	ID              int    `gorm:"column:id;primaryKey"`
	CodDepartamento string `gorm:"column:cod_departamento;not null;uniqueIndex"`
	Descripcion     string `gorm:"column:descripcion;not null;uniqueIndex"`
}

func (DepartamentoSoporte) TableName() string { return "departamentos_soporte" }

type MotivoPausa struct {
	ID                    int    `gorm:"column:id;primaryKey"`
	MotivoPausa           string `gorm:"column:motivo_pausa;size:255;not null"`
	RequiereAutorizacion *bool   `gorm:"column:requiere_autorizacion;default:true"`
}

func (MotivoPausa) TableName() string { return "motivos_pausa" }

type CatalogoFalla struct {
	ID                       int    `gorm:"column:id;primaryKey"`
	CodigoFalla              string `gorm:"column:codigo_falla;not null;uniqueIndex"`
	DescripcionFalla         string `gorm:"column:descripcion_falla;not null"`
	Complejidad              int    `gorm:"column:complejidad;not null;default:1"`
	TiempoResolucionEstimado string `gorm:"column:tiempo_resolucion_estimado"`
	RequiereVisitaFisica     *bool  `gorm:"column:requiere_visita_fisica;default:true"`
	IDDepartamento           *int   `gorm:"column:id_departamento"`
	Categoria                string `gorm:"column:categoria"`
	Subcategoria             string `gorm:"column:subcategoria"`
}

func (CatalogoFalla) TableName() string { return "catalogo_fallas" }

type TipoTurno struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement"`
	Nombre      string `gorm:"column:nombre;size:50;not null;uniqueIndex"`
	Descripcion string `gorm:"column:descripcion;size:255"`
	Estado      *bool  `gorm:"column:estado;default:true"`
}

func (TipoTurno) TableName() string { return "tipos_turno" }
