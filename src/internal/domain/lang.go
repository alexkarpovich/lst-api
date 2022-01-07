package domain

type Language struct {
	Code       string `json:"code" db:"code"`
	IsoName    string `json:"isoName" db:"iso_name"`
	NativeName string `json:"nativeName" db:"native_name"`
}

type LangRepo interface {
	List() ([]*Language, error)
}
