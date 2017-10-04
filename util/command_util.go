package util

// SimplePublicCommand is a very simple match sugar.
func SimplePublicCommand(desired string, supplied string, isUser bool) bool {
	if supplied == desired {
		return true
	}
	return false
}

// SimplePrivateCommand returns true if they match AND is the user.
func SimplePrivateCommand(desired string, supplied string, isUser bool) bool {
	if supplied == desired && isUser {
		return true
	}
	return false
}
