package fitness

func IMT(height int, weight float64) float64 {
	var meters = float64(height) / 100
	return weight / (meters * meters)
}
