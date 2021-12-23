package domain

type Language struct {
	Code       string `json:"code"`
	IsoName    string `json:"isoName"`
	NativeName string `json:"nativeName"`
}

type LangRepo interface {
	List() []*Language
}
