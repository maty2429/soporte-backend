package models

type Tecnico struct {
	ID                    int                  `gorm:"column:id;primaryKey"`
	Rut                   string               `gorm:"column:rut;size:10;not null;uniqueIndex"`
	Dv                    string               `gorm:"column:dv;size:1;not null"`
	NombreCompleto        string               `gorm:"column:nombre_completo;not null"`
	IDTipoTecnico         *int                 `gorm:"column:id_tipo_tecnico"`
	IDDepartamentoSoporte *int                 `gorm:"column:id_departamento"`
	Activo                *bool                `gorm:"column:activo;default:true"`
	DepartamentoSoporte   *DepartamentoSoporte `gorm:"-:migration;foreignKey:IDDepartamentoSoporte;references:ID"`
}

func (Tecnico) TableName() string { return "tecnicos" }

type ConfiguracionHorarioTurno struct {
	ID          int    `gorm:"column:id;primaryKey"`
	IDTipoTurno int    `gorm:"column:id_tipo_turno;not null"`
	DiaSemana   int    `gorm:"column:dia_semana;not null"`
	HoraInicio  string `gorm:"column:hora_inicio;not null"`
	HoraFin     string `gorm:"column:hora_fin;not null"`
}

func (ConfiguracionHorarioTurno) TableName() string { return "configuracion_horarios_turno" }
