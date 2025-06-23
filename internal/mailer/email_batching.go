package mailer

import (
	"sync"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type EmailBatchingJob struct {
	config               func() *com.Config
	service              *Mailer
	newNotifications     chan *batchedNotification
	pendingNotifications map[string][]*batchedNotification
	tasks                *models.ScheduledTask
	taskMutex            sync.Mutex
}

type batchedNotification struct{}

func NewEmailBatchingJob(s *Mailer, bufferSize int) *EmailBatchingJob {
	return &EmailBatchingJob{
		config:               s.config,
		service:              s,
		newNotifications:     make(chan *batchedNotification, bufferSize),
		pendingNotifications: make(map[string][]*batchedNotification),
	}
}

func (s *Mailer) InitEmailBatching() {
	if s.config().Email.GetEnableEmailBatching() {
		if s.EmailBatching == nil {
			s.EmailBatching = NewEmailBatchingJob(s, int(s.config().Email.GetEmailBatchingBufferSize()))
		}
		// s.EmailBatching.Start()
	}
}

// func (s  *Service ) AddNotificationEmailToBatch()  *models.AppError{
// 	if !s.config().Email.GetEnableEmailBatching() { }
// }
