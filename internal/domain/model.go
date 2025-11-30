package domain

type Record struct {
	Links map[string]string `json:"links"`
	ID    int64             `json:"links_num"`
}

type TempRecord struct {
	Links []string `json:"links"`
	ID    int64    `json:"links_num"`
}
