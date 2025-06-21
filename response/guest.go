package response

type GuestResponse struct {
	CredentialType string `json:"credential_type" binding:"required"`
	IDNumber       string `json:"id_number" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	Phone          string `json:"phone_number" binding:"required"`
	Email          string `json:"email"`
}
