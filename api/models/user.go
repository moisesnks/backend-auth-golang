package models

// User representa la estructura de datos de un usuario en Firebase
type User struct {
	UID          string `json:"uid"`
	DisplayName  string `json:"displayName"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	Email        string `json:"email"`
	Verified     bool   `json:"verified"`
}
