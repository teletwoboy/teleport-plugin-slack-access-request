package telemetry

const (
	OpenModal              = "open-modal"
	OpenModalAccessRequest = "open-access-request-modal"
	OpenModalAccessReview  = "open-access-review-modal"
	OpenModalAccessPolicy  = "open-access-policy-modal"

	ARequest                              = "access-request"
	ARequestRoleSelection                 = "select-access-request-role"
	ARequestChannelSelection              = "select-access-request-channel"
	ARequestStartDateOptionSelection      = "select-access-request-start-date-option"
	ARequestStartDateSelection            = "select-access-request-start-date"
	ARequestStartTimeSelection            = "select-access-request-start-time"
	ARequestAccessDurationOptionSelection = "select-access-request-access-duration-option"
	ARequestAccessDurationDateSelection   = "select-access-request-access-duration-date"
	ARequestAccessDurationTimeSelection   = "select-access-request-access-duration-time"
	ARequestRequestTTLOptionSelection     = "select-access-request-request-ttl-option"
	ARequestRequestTTLDateSelection       = "select-access-request-request-ttl-date"
	ARequestRequestTTLTimeSelection       = "select-access-request-request-ttl-time"
	ARequestModalSubmission               = "submit-access-request-modal"

	AReview                = "access-review"
	AReviewModalSubmission = "submit-access-review-modal"

	APolicy                   = "access-policy"
	APolicyChannelSelection   = "select-access-policy-channel"
	APolicyRoleSelection      = "select-access-policy-role"
	APolicyUserSelection      = "select-access-policy-user"
	APolicyStartDateSelection = "select-access-policy-start-date"
	APolicyStartTimeSelection = "select-access-policy-start-time"
	APolicyEndDateSelection   = "select-access-policy-end-date"
	APolicyEndTimeSelection   = "select-access-policy-end-time"
	APolicyEffectSelection    = "select-access-policy-effect"
	APolicyModalSubmission    = "submit-access-policy-modal"

	SlackService    = "slack-service"
	TeleportService = "teleport-service"
)
