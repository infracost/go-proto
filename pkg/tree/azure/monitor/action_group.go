package monitor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ActionGroup struct {
	resource.Resource         `tree:"-"`
	EmailReceivers            value.Int       `tree:"email_receivers"`
	ITSMEventReceivers        value.Int       `tree:"itsm_event_receivers"`
	PushNotificationReceivers value.Int       `tree:"push_notification_receivers"`
	SecureWebHookReceivers    value.Int       `tree:"secure_webhook_receivers"`
	WebHookReceivers          value.Int       `tree:"webhook_receivers"`
	SMSReceivers              []SMSReceiver   `tree:"sms_receivers"`
	VoiceCallReceivers        []VoiceReceiver `tree:"voice_call_receivers"`
}

type SMSReceiver struct {
	CountryCode int `tree:"country_code"`
	Count       int `tree:"count"`
}

type VoiceReceiver struct {
	CountryCode int `tree:"country_code"`
	Count       int `tree:"count"`
}
