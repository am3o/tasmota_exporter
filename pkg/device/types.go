package device

type Device struct {
	Status struct {
		Name string `json:"DeviceName"`
		On   int    `json:"Power"`
	}
}

type Version struct {
	Status struct {
		SDK     string `json:"SDK"`
		Version string `json:"Version"`
	} `json:"StatusFWR"`
}

type Network struct {
	Status struct {
		Hostname string `json:"Hostname"`
		Address  string `json:"IPAddress"`
	} `json:"StatusNET"`
}

type PowerStatus struct {
	Status struct {
		Energy struct {
			Total         float64
			Yesterday     float64
			Today         float64
			Power         int
			ApparentPower int
			ReactivePower int
			Factor        float64
			Voltage       int
			Current       float64
		} `json:"ENERGY"`
	} `json:"StatusSNS"`
}
