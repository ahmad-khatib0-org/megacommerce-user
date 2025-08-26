package models

type EventName string

const (
	EventNameSupplierCreate    = "supplier_create"
	EventNameEmailConfirmation = "email_confirmation"
	EventNamePasswordForgot    = "password_forgot"
)

type EventStatus string

const (
	EventStatusFail    EventStatus = "fail"
	EventStatusSuccess EventStatus = "success"
	EventStatusAttempt EventStatus = "attempt"
)

// AuditRecord provides a consistent set of fields used for all audit logging.
type AuditRecord struct {
	EventName EventName       `json:"event_name"`
	Status    EventStatus     `json:"status"`
	EventData AuditEventData  `json:"event"`
	Actor     AuditEventActor `json:"actor"`
	Meta      map[string]any  `json:"meta"`
	Error     EventError      `json:"error,omitempty"`
}

// AuditEventData contains all event specific data about the modified entity
type AuditEventData struct {
	Parameters  map[string]any `json:"parameters"`      // Payload and parameters being processed as part of the request
	PriorState  map[string]any `json:"prior_state"`     // Prior state of the object being modified, nil if no prior state
	ResultState map[string]any `json:"resulting_state"` // Resulting object after creating or modifying it
	ObjectType  string         `json:"object_type"`     // String representation of the object type. eg. "post"
}

// AuditEventActor is the subject triggering the event
type AuditEventActor struct {
	UserID        string `json:"user_id"`
	SessionID     string `json:"session_id"`
	Client        string `json:"client"`
	IPAddress     string `json:"ip_address"`
	XForwardedFor string `json:"x_forwarded_for"`
}

// EventError contains error information in case of failure of the event
type EventError struct {
	Description string `json:"description,omitempty"`
	Code        int    `json:"status_code,omitempty"`
}

// Success marks the audit record status as successful.
func (ar *AuditRecord) Success() {
	ar.Status = EventStatusSuccess
}

// Fail marks the audit record status as failed.
func (ar *AuditRecord) Fail() {
	ar.Status = EventStatusFail
}

func AuditRecordNew(ctx *Context, event EventName, initialStatus EventStatus) *AuditRecord {
	return &AuditRecord{
		EventName: event,
		Status:    initialStatus,
		Actor: AuditEventActor{
			UserID:        ctx.Session.UserID,
			SessionID:     ctx.Session.ID,
			IPAddress:     ctx.IPAddress,
			Client:        ctx.UserAgent,
			XForwardedFor: ctx.XForwardedFor,
		},
		Meta: map[string]any{},
		EventData: AuditEventData{
			Parameters:  map[string]any{},
			PriorState:  map[string]any{},
			ResultState: map[string]any{},
		},
	}
}

func AuditEventDataParameter(ar *AuditRecord, key string, val any) {
	if ar.EventData.Parameters == nil {
		ar.EventData.Parameters = make(map[string]any)
	}
	ar.EventData.Parameters[key] = val
}

func (ar *AuditRecord) AuditEventDataPriorState(data map[string]any) {
	if ar.EventData.PriorState == nil {
		ar.EventData.PriorState = make(map[string]any)
	}
	ar.EventData.PriorState = data
}

func (ar *AuditRecord) AuditEventDataResultState(data map[string]any) {
	if ar.EventData.ResultState == nil {
		ar.EventData.ResultState = make(map[string]any)
	}
	ar.EventData.ResultState = data
}
