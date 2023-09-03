package utils

const (
	BadRequest             = "bad request"
	InternalError          = "internal error"
	InvalidCredentials     = "invalid credentials"
	UserNotFoundError      = "user not found"
	InvalidTaxId           = "invalid tax id"
	OperationCreationError = "error creating operation"
	InvalidAddrError       = "invalid wallet address"
	DepositCreationError   = "error creating deposit"
	SignatureError         = "signature mismatch"
	UserAlreadyExists      = "user with this email already exists"
	UserWaitingPermit      = "burn or bridge operation still pending, must wait for it to finish"
)

const (
	MissingEmail    = "missing email field"
	MissingName     = "missing name field"
	MissingPassword = "missing password field"

	MissingPixKey = "missing pixKey field"

	MissingWalletAddress = "missing walletAddress field"
	MissingChain         = "missing chain field"
	MissingInputChain    = "missing inputChain field"
	MissingOutputChain   = "missing outputChain field"
	ChainsAreEqual       = "inputChain cannot be equal outputChain"

	MissingPermit      = "missing permit field"
	PasswordsDontMatch = "passwords don't match"

	InvalidAmount = "amount must be a positive integer"
	InvalidChain  = "invalid chain selected"
)
