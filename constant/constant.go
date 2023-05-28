package constant

const (
	// role user
	RoleAdmin     = "admin"
	RoleValidator = "validator"
	RoleFarmer    = "farmer"
	RoleBuyer     = "buyer"

	// status proposal
	ProposalStatusPending  = "pending"
	ProposalStatusApproved = "approved"
	ProposalStatusRejected = "rejected"

	// transaction type
	TransactionTypePerennials = "perennials"
	TransactionTypeAnnuals    = "annuals"

	// status transaction
	TransactionStatusPending  = "pending"
	TransactionStatusAccepted = "accepted"
	TransactionStatusRejected = "rejected"
	TransactionStatusCancel   = "cancelled"

	// status batch
	BatchStatusPlanting = "planting"
	BatchStatusHarvest  = "harvest"
	BatchStatusCancel   = "cancel"

	// status harvest
	HarvestStatusPending  = "pending"
	HarvestStatusApproved = "approved"
	HarvestStatusRevision = "revision"

	// status treatment record
	TreatmentRecordStatusWaitingResponse = "waitingResponse"
	TreatmentRecordStatusPending         = "pending"
	TreatmentRecordStatusApproved        = "approved"
	TreatmentRecordStatusRevision        = "revision"

	// folder cloudinary
	CloudinaryFolderCommodities      = "commodities"
	CloudinaryFolderTreatmentRecords = "treatmentRecords"
	CloudinaryFolderHarvests         = "harvests"

	// template mailgun
	MailgunForgotPasswordTemplate = "forgot_password"
)
