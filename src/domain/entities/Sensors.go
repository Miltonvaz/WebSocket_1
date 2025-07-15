package entities

type Message struct {
	IdUser        int     `json:"id_user"`
	Code          int     `json:"code"`
	PH            float64 `json:"ph"`
	Conductividad int     `json:"conductividad"`
	Turbuidez     int     `json:"turbuidez"`
	Temperatura   float64 `json:"temperatura"`
	Alcohol       int     `json:"alcohol"`
	Densidad      int     `json:"densidad"`
}
