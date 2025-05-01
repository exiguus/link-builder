package types

type LinkPreviewOutput struct {
	ID      int         `json:"id"`
	Date    string      `json:"date"`
	URL     string      `json:"url"`
	Preview interface{} `json:"preview"`
}
