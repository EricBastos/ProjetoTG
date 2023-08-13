package utils

const (
	BadRequest             = "bad request"
	InternalError          = "internal error"
	PixToUsdAlreadyExists  = "only one pix to usd order can exist simultaneously"
	CantSwapAtomically     = "can't ensure atomic swap"
	BigSwapNotAtomic       = "can't ensure atomic big swap"
	InvalidCredentials     = "invalid credentials"
	EmailUnverified        = "verify your email"
	UserNotFoundError      = "user not found"
	AlreadyPassedKyc       = "user already passed kyc"
	InvalidKycUpgrade      = "kyc document must be the same when upgrading level"
	InvalidTaxId           = "invalid tax id"
	TransferLogError       = "error fetching transfer logs for user"
	DepositLogError        = "error fetching deposit logs for user"
	KycStatusHistoryError  = "error fetching kyc history for user"
	KycLinkError           = "error fetching kycs for user"
	AccountCreationError   = "error creating account"
	TaxesIdDontMatch       = "tax id of the bank account associated to this key doesn't match user's tax id"
	UserAlreadyVerified    = "user was verified already"
	OperationCreationError = "error creating operation"
	InvalidAddrError       = "invalid wallet address"
	DepositCreationError   = "error creating deposit"
	SignatureError         = "signature mismatch"
	NotOptedIn             = "account has not opted in"
	Frozen                 = "account is frozen"
	BlackListed            = "account is blacklisted"
	MissingKYC             = "user must pass KYC first"
	KycConflic             = "kyc level in DB and token is different, login again to renew token"
	MustBeSuperuser        = "user must be superuser"
	MustBePj               = "user must be PJ"
	MustBePf               = "user must be PF"
	UserAlreadyExists      = "user with this email already exists"
	ErrorParsingWebhookUrl = "error parsing webhook url"
	WebhookMustBeHttps     = "webhook url must be https"
	TooManyApiKeys         = "maximum of 1 API keys allowed"
	TooManyWebhooks        = "maximum of 1 webhooks allowed"
	TooManyWallets         = "maximum of 20 wallets allowed"
	TooManyTaxIdRequests   = "too many requests for this taxId, try again in 5 minutes"
	QuotingError           = "error quoting pair"
	OnlyInSandbox          = "available only in sandbox"
	OnlyInProduction       = "available only in production"
	WalletAlreadyAdded     = "this wallet was already added"
	WalletNameRepeated     = "this wallet name was already used"
	AccountNameRepeated    = "this account name was already used"
	AccountAlreadyAdded    = "this account was already added"
	AccountNotFound        = "account not found"
	KycStatusNotFound      = "kyc status not found"
	KycLinkNotFound        = "kyc link not found"
	KycStillPending        = "kyc still pending"
	WalletNotFound         = "wallet not found"
	WalletWrongChain       = "wallet chain doesn't match input"
	InvalidPixKey          = "invalid pix key"
	SseNotSupported        = "connection does not support streaming"
	ErrorQuoting           = "error registering user (already quoting?)"
	TokenNotYours          = "token does not belong to user"
	TokenUsed              = "this swap token has already been used"
	UserWaitingPermit      = "burn or swap operation still pending, must wait for it to finish"

	KycAlreadyUsed       = "this taxId is used by another account. contact us if you think that's a mistake"
	KycPendingThirdParty = "this taxId is being processed by another account. contact us if you think that's a mistake"
	KycMustBeLevel1      = "must have kyc level 1 before attempting level 2"

	UserAlreadyUsing = "user is already consuming the endpoint"
)

const (
	MissingEmail          = "missing email field"
	MissingPassword       = "missing password field"
	MissingTaxId          = "missing taxId field"
	MissingCPF            = "missing cpf field"
	MissingCNPJ           = "missing cnpj field"
	WalletAddressRequired = "walletAddress must be passed for quoting above 100 qty"
	QuoteSurpassBalance   = "quoted amount must not surpass wallet balance"
	QuoteSurpassAllowance = "quoted amount must not surpass wallet allowance"

	MissingCurrentPassword    = "missing currentPassword field"
	MissingNewPassword        = "missing newPassword field"
	MissingNewPasswordConfirm = "missing newPasswordConfirm field"
	WrongCurrentPassword      = "current password is wrong"

	InvalidCpfField    = "invalid cpf field"
	InvalidCnpjField   = "invalid cnpj field"
	MissingSdkToken    = "missing SDKTokenn field"
	MissingMotherName  = "missing motherName field"
	MissingFullName    = "missing fullName field"
	MissingCompanyName = "missing companyName field"
	InvalidBirthDate   = "birthDate field is invalid"
	InvalidStartDate   = "startDate field is invalid"

	MissingFirstName   = "missing firstName field"
	MissingLastName    = "missing lastName field"
	MissingFantasyName = "missing fantasyName field"
	MissingPhone       = "missing phone field"
	MissingCountry     = "missing country field"
	MissingState       = "missing state field"
	MissingWorkspace   = "missing workspaceId field"
	MissingIdentifier  = "missing identifier field"
	MissingWebhookUrl  = "missing webhookUrl field"

	MissingAccountNickname   = "missing accountNickname field"
	MissingWalletName        = "missing name field"
	MissingAccountNumber     = "missing accountNumber field"
	MissingPixKey            = "missing pixKey field"
	MissingReceiverTaxId     = "missing receiverTaxId field"
	MissingBranchCode        = "missing branchCode field"
	MissingBankCodeTED       = "missing bankCodeTED field"
	MissingBankCodePIX       = "missing bankCodePIX field"
	MissingWalletAddress     = "missing walletAddress field"
	MissingFromWalletAddress = "missing from field"
	MissingToWalletAddress   = "missing to field"
	MissingValueAddress      = "missing value field"
	MissingPermitAddress     = "must send either usdcPermit or usdtPermit"
	MissingChain             = "missing chain field"

	MissingFaceData     = "missing faceData field"
	MissingDocumentData = "must either fill CNHData field or both RGFrontData and RGBackData fields"
	InvalidImageType    = "imageType must be one of [image/jpg, image/jpeg, image/png, application/pdf]"

	MissingToken             = "missing token field"
	MissingPermit            = "missing permit field"
	MissingAccountId         = "missing accountId field"
	MissingWalletId          = "missing walletId field"
	MissingUserTaxId         = "missing userTaxId field"
	MissingUserName          = "missing userName field"
	MissingAccountBankCode   = "missing accountBankCode field"
	MissingAccountBranchCode = "missing accountBranchCode field"
	MissingSignature         = "missing signature field"
	MissingMessageTime       = "missing messageTime field"
	TokenNotFound            = "reset password token not found"
	TokenExpired             = "reset password token expired"
	TokenAlreadyUsed         = "token has been used already"
	PasswordsDontMatch       = "passwords don't match"
	TokenInvalid             = "invalid token"
	WontSwapZero             = "wont swap zero quantity of either side"

	InvalidAmount          = "amount must be a positive integer"
	InvalidMintToUsdAmount = "amount must be an integer greater than 40000 (R$400.00)"
	InvalidBurnAmount1     = "amount must be an integer greater than 5000 (R$50.00)"
	InvalidBurnAmount2     = "amount must be an integer greater than 75 (R$0.75)"
	InvalidChain           = "invalid chain selected"
	InvalidCoin            = "coin must be either USDT or USDC"
	InvalidQuoteAmount     = "amount must be a positive float with at most 2 decimal places"
	InvalidMarkupAmount    = "markup must be a positive float that represents a percentage with at most 4 decimal places"
)
