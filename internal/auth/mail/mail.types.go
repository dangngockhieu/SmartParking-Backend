package mail

type VerificationEmailData struct {
	FirstName string
	VerifyURL string
}

type ResetPasswordEmailData struct {
	FirstName string
	CodeID    string
}

type VerifiedPageData struct {
	Year int
}
