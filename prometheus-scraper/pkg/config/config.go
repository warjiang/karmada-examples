package config

type ScrapeTarget struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Port      string `json:"port"`
}

type Config struct {
	IntervalSeconds int            `json:"intervalSeconds"`
	Targets         []ScrapeTarget `json:"targets"`
}
