package domain

type TipoTicket struct {
	ID             int
	CodTipoTicket  string
	Descripcion    string
}

type NivelPrioridad struct {
	ID          int
	Descripcion string
}

type TipoTecnico struct {
	ID          int
	Descripcion string
}

type DepartamentoSoporte struct {
	ID              int
	CodDepartamento string
	Descripcion     string
}

type MotivoPausa struct {
	ID                    int
	MotivoPausa           string
	RequiereAutorizacion  bool
}

type CatalogoFalla struct {
	ID                       int
	CodigoFalla              string
	DescripcionFalla         string
	Complejidad              int
	TiempoResolucionEstimado string
	RequiereVisitaFisica     bool
	IDDepartamento           *int
	Categoria                string
	Subcategoria             string
}

type TipoTurno struct {
	ID          int
	Nombre      string
	Descripcion string
	Estado      bool
}
