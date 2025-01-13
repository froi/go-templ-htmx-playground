package types

type SignupInputFormParams struct {
	Email                       string
	Password1                   string
	Password2                   string
	ShowInvalidPasswordFlag     bool
	ShowNonMatchingPasswordFlag bool
	ShowInvalidEmailFlag        bool
	ShowTakenEmailFlag          bool
	SubmitButtonDisabled        bool
}

func (obj *SignupInputFormParams) FormAppearsValid() bool {
	allValid := (!obj.ShowInvalidPasswordFlag &&
		!obj.ShowInvalidEmailFlag &&
		!obj.ShowTakenEmailFlag &&
		!obj.SubmitButtonDisabled &&
		!obj.ShowNonMatchingPasswordFlag)
	return allValid
}
