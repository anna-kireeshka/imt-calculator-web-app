package server

type BMRBody struct {
	Weight float64 `json:"weight"`
	Height int     `json:"height"`
	DOB    int     `json:"dob"`
	Gender bool    `json:"gender"`
}
