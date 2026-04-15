package domain

type Archivo struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Ruta   string `json:"ruta"`
}

type ArchivoRepository interface {
	Save(archivo *Archivo) error
	FindAll() ([]Archivo, error)
	FindByID(id int) (*Archivo, error)
	Update(archivo *Archivo) error
	Delete(id int) error
}
