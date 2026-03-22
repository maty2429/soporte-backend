package models

type Servicio struct {
	ID                      int    `gorm:"column:id;primaryKey"`
	Edificio                string `gorm:"column:edificio;default:SIN ESPECIFICAR"`
	Piso                    int    `gorm:"column:piso;default:0"`
	Servicios               string `gorm:"column:servicios;default:SIN ESPECIFICAR"`
	Ubicacion               string `gorm:"column:ubicacion;default:SIN ESPECIFICAR"`
	Unidades                string `gorm:"column:unidades;default:SIN ESPECIFICAR"`
	IDNivelPrioridadDefault *int   `gorm:"column:id_nivel_prioridad_default"`
}

func (Servicio) TableName() string {
	return "servicio"
}
