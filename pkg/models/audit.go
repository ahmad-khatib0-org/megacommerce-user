package models

import (
	"log"
	"sync"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type EventName string

const (
	EventNameSupplierCreate = "supplier_create"
)

type EventStatus string

const (
	EventStatusFail    EventStatus = "fail"
	EventStatusSuccess EventStatus = "success"
	EventStatusAttempt EventStatus = "attempt"
)

var (
	cachedAuditEventDataFieldValue *anypb.Any
	cachedAuditEventDataFieldOnce  sync.Once
)

func getCachedAuditEventDataFieldValue() *anypb.Any {
	cachedAuditEventDataFieldOnce.Do(func() {
		value, err := anypb.New(cachedAuditEventDataFieldValue)
		if err != nil {
			log.Fatalf("failed to marshal an empty map[string]any value, err: %v", &InternalError{
				Err:  err,
				Path: "user.models.audit.getCachedAuditEventDataFieldValue",
				Msg:  "failed to marshal an empty map[string]any value",
			})
		}
		cachedAuditEventDataFieldValue = value
	})

	return cachedAuditEventDataFieldValue
}

func AuditRecordNew(ctx *Context, event EventName, initialStatus EventStatus) *common.AuditRecord {
	return &common.AuditRecord{
		EventName: string(event),
		Status:    string(initialStatus),
		Actor: &common.AuditEventActor{
			UserId:        ctx.Session.UserId,
			SessionId:     ctx.Session.Id,
			IpAddress:     ctx.IPAddress,
			Client:        ctx.UserAgent,
			XForwardedFor: ctx.XForwardedFor,
		},
		Meta: &common.AuditRecordMeta{Path: ctx.Path},
		EventData: &common.AuditEventData{
			Parameters:     getCachedAuditEventDataFieldValue(),
			PriorState:     getCachedAuditEventDataFieldValue(),
			ResultingState: getCachedAuditEventDataFieldValue(),
		},
	}
}
