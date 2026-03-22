package domain

type Servicio struct {
	ID                      int
	Edificio                string
	Piso                    int
	Servicios               string
	Ubicacion               string
	Unidades                string
	IDNivelPrioridadDefault *int
}
