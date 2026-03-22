package types

type CreateCatalogoFallaCommand struct {
	CodigoFalla          string
	DescripcionFalla     string
	Complejidad          int
	RequiereVisitaFisica bool
	IDDepartamento       *int
	Categoria            string
	Subcategoria         string
}

type UpdateCatalogoFallaCommand struct {
	ID                   int
	CodigoFalla          *string
	DescripcionFalla     *string
	Complejidad          *int
	RequiereVisitaFisica *bool
	IDDepartamento       *int
	Categoria            *string
	Subcategoria         *string
}
