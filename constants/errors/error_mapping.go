package errors

func MappingError(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, ServiceErrorMessage...)

	// check errors
	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
