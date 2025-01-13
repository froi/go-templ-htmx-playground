package types

type LoginInputFormParams struct {
	Email                string
	Password             string
	ShowFailedLoginFlag  bool
	SubmitButtonDisabled bool
}

func (obj *LoginInputFormParams) IsValid() bool {
	allValid := (!obj.ShowFailedLoginFlag)
	return allValid
}
