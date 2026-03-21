package models

type Solicitante struct {
	ID             int       `gorm:"column:id;primaryKey"`
	IDServicio     *int      `gorm:"column:id_servicio;index"`
	Servicio       *Servicio `gorm:"foreignKey:IDServicio;references:ID"`
	Correo         *string   `gorm:"column:correo;size:100;uniqueIndex"`
	Rut            string    `gorm:"column:rut;size:10;not null;uniqueIndex"`
	Dv             string    `gorm:"column:dv;size:1;not null"`
	NombreCompleto string    `gorm:"column:nombre_completo;not null"`
	Anexo          *int      `gorm:"column:anexo"`
	Estado         *bool     `gorm:"column:estado;default:true"`
}

func (Solicitante) TableName() string {
	return "solicitantes"
}
