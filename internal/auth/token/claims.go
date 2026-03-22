package token

type ClaimsPayload struct {
	UserID uint64
	Email  string
	Role   string
	JTI    string
	Exp    int64
}
