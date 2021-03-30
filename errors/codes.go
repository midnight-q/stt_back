package errors

const DefaultErrorLanguageId = 1

// Error codes
const (
	ErrorCodeUndefined               ErrorCode = 0
	ErrorCodeNotFound                ErrorCode = 100
	ErrorCodeInvalidAuthorize        ErrorCode = 200
	ErrorCodeUnsupportedFunctionType ErrorCode = 300
	ErrorCodeParseId                 ErrorCode = 400
	ErrorCodeInvalidPerPage          ErrorCode = 500
	ErrorCodeNotValid                ErrorCode = 700
	ErrorCodeSqlError                ErrorCode = 750
	ErrorCodeInvalidCurrentPage      ErrorCode = 800
	ErrorCodeNotEmpty                ErrorCode = 900
	ErrorCodeFieldLengthTooShort     ErrorCode = 1100
	ErrorCodeUnsupportedFileFormat   ErrorCode = 1200
)

// Error messages EN
const (
	ErrorEnMessageUndefined               = "Undefined error"
	ErrorEnMessageNotFound                = "Not found"
	ErrorEnMessageInvalidAuthorize        = "Invalid authorize"
	ErrorEnMessageUnsupportedFunctionType = "Unsupported function type"
	ErrorEnMessageParseId                 = "Error parse Id"
	ErrorEnMessageInvalidPerPage          = "Invalid PerPage"
	ErrorEnMessageNotValid                = "Field not valid"
	ErrorEnMessageSqlError                = "Sql error"
	ErrorEnMessageInvalidCurrentPage      = "Invalid CurrentPage"
	ErrorEnMessageNotEmpty                = "Field in not empty"
	ErrorEnMessageFieldLengthTooShort     = "Field value has too short length"
	ErrorEnMessageUnsupportedFileFormat   = "Unsupported file format"
)
