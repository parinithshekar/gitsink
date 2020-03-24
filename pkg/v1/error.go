package v1

const (
	// Counter errors.
	ErrCounterCreateFailed  = Error("Failed to create counter.")
	ErrCounterDestroyFailed = Error("Failed to remove counter.")

	// CounterVec errors.
	ErrCounterVecCreateFailed = Error("Failed to create counter vector.")

	// Errors used for testing.
	ErrNoMatch = Error("No alert matched in alert config.")

	// User Registry errors.
	ErrUserNotFoundByID       = Error("User not found by ID.")
	ErrUserNotFoundByUsername = Error("User not found by Username.")

	// Authorization errors.
	ErrUnauthorizedAPI = Error("Unauthorized API access.")

	// Generic gRPC errors.
	ErrRecordAlreadyExists        = Error("Record already exists.")
	ErrRecordConversionToDBFailed = Error("Record conversion to DB format failed.")
	ErrRecordConversionToPBFailed = Error("Record conversion to PB format failed.")
	ErrRecordDeleteFailed         = Error("Record delete failed.")
	ErrRecordDoesNotExist         = Error("Record does not exist.")
	ErrRecordInsertFailed         = Error("Record insert failed.")
	ErrRecordUpdateFailed         = Error("Record update failed.")

	// Message validation errors.
	ErrClientMsgValidationFailed = Error("Client message validation failed.")
	ErrServerMsgValidationFailed = Error("Server message validation failed.")
)

// Error represents an error.
type Error string

// Error returns the error as a string.
func (e Error) Error() string { return string(e) }
