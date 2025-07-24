package entities

type Message struct {
	IdUser        int     `json:"id_user"`
	Code          int     `json:"code"`
	PH            float64 `json:"ph"`
	Conductividad float64 `json:"conductividad"`
	Turbuidez     float64 `json:"turbuidez"`
	Temperatura   float64 `json:"temperatura"`
	Alcohol       float64 `json:"alcohol"`
	Densidad      float64 `json:"densidad"`
	Rpm           int     `json:"rpm"`
}
