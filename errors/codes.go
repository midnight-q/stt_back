
package errors

const DefaultErrorLanguageId = 1

// Error codes
const (
	ErrorCodeUndefined                          ErrorCode = 0
	ErrorCodeNotFound                           ErrorCode = 100
	ErrorCodeInvalidAuthorize                   ErrorCode = 200
	ErrorCodeUnsupportedFunctionType            ErrorCode = 300
	ErrorCodeParseId                            ErrorCode = 400
	ErrorCodeInvalidPerPage                     ErrorCode = 500
	ErrorCodeRabbitQueueNameNotSet              ErrorCode = 600
	ErrorCodeNotValid                           ErrorCode = 700
	ErrorCodeSqlError                           ErrorCode = 750
	ErrorCodeInvalidCurrentPage                 ErrorCode = 800
	ErrorCodeNotEmpty                           ErrorCode = 900
	ErrorCodeAgeToSmall                         ErrorCode = 1000
	ErrorCodeFieldLengthTooShort                ErrorCode = 1100
)

// Error messages EN
const (
	ErrorEnMessageUndefined                          = "Undefined error"
	ErrorEnMessageNotFound                           = "Not found"
	ErrorEnMessageInvalidAuthorize                   = "Invalid authorize"
	ErrorEnMessageUnsupportedFunctionType            = "Unsupported function type"
	ErrorEnMessageParseId                            = "Error parse Id"
	ErrorEnMessageInvalidPerPage                     = "Invalid PerPage"
	ErrorEnMessageRabbitQueueNameNotSet              = "RabbitQueueName not set for application"
	ErrorEnMessageNotValid                           = "Field not valid"
	ErrorEnMessageSqlError                           = "Sql error"
	ErrorEnMessageInvalidCurrentPage                 = "Invalid CurrentPage"
	ErrorEnMessageNotEmpty                           = "Field in not empty"
	ErrorEnMessageAgeToSmall                         = "You must be over 14 years old"
	ErrorEnMessageFieldLengthTooShort                = "Field value has too short length"
)
