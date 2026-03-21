package domain

type Solicitante struct {
	ID             int
	IDServicio     *int
	Servicio       *Servicio
	Correo         string
	Rut            string
	Dv             string
	NombreCompleto string
	Anexo          *int
	Estado         bool
}
