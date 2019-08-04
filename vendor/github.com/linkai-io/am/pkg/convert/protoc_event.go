package convert

import (
	"github.com/linkai-io/am/am"
	"github.com/linkai-io/am/protocservices/prototypes"
)

func DomainToEvent(in *am.Event) *prototypes.EventData {
	return &prototypes.EventData{
		NotificationID: in.NotificationID,
		OrgID:          int32(in.OrgID),
		GroupID:        int32(in.GroupID),
		TypeID:         in.TypeID,
		EventTimestamp: in.EventTimestamp,
		Data:           in.Data,
		JsonData:       in.JSONData,
		Read:           false,
	}
}

func EventToDomain(in *prototypes.EventData) *am.Event {
	return &am.Event{
		OrgID:          int(in.OrgID),
		GroupID:        int(in.GroupID),
		NotificationID: in.NotificationID,
		TypeID:         in.TypeID,
		EventTimestamp: in.EventTimestamp,
		Data:           in.Data,
		JSONData:       in.JsonData,
		Read:           in.Read,
	}
}

func DomainToEventSubscriptions(in *am.EventSubscriptions) *prototypes.EventSubscriptions {
	return &prototypes.EventSubscriptions{
		TypeID:              in.TypeID,
		SubscribedTimestamp: in.SubscribedTimestamp,
		Subscribed:          in.Subscribed,
		WebhookEnabled:      in.WebhookEnabled,
		WebhookType:         in.WebhookType,
		WebhookURL:          in.WebhookURL,
		WebhookVersion:      in.WebhookVersion,
	}
}

func EventSubscriptionsToDomain(in *prototypes.EventSubscriptions) *am.EventSubscriptions {
	return &am.EventSubscriptions{
		TypeID:              in.TypeID,
		SubscribedTimestamp: in.SubscribedTimestamp,
		Subscribed:          in.Subscribed,
		WebhookEnabled:      in.WebhookEnabled,
		WebhookType:         in.WebhookType,
		WebhookURL:          in.WebhookURL,
		WebhookVersion:      in.WebhookVersion,
	}
}

func DomainToUserEventSettings(in *am.UserEventSettings) *prototypes.UserEventSettings {
	subs := make([]*prototypes.EventSubscriptions, 0)
	if in.Subscriptions != nil {
		for _, sub := range in.Subscriptions {
			subs = append(subs, DomainToEventSubscriptions(sub))
		}
	}
	return &prototypes.UserEventSettings{
		WeeklyReportSendDay: in.WeeklyReportSendDay,
		ShouldWeeklyEmail:   in.ShouldWeeklyEmail,
		DailyReportSendHour: in.DailyReportSendHour,
		ShouldDailyEmail:    in.ShouldDailyEmail,
		UserTimezone:        in.UserTimezone,
		Subscriptions:       subs,
		WebhookCurrentKey:   in.WebhookCurrentKey,
		WebhookPreviousKey:  in.WebhookPreviousKey,
	}
}

func UserEventSettingsToDomain(in *prototypes.UserEventSettings) *am.UserEventSettings {
	subs := make([]*am.EventSubscriptions, 0)
	if in.Subscriptions != nil {
		for _, sub := range in.Subscriptions {
			subs = append(subs, EventSubscriptionsToDomain(sub))
		}
	}
	return &am.UserEventSettings{
		WeeklyReportSendDay: in.WeeklyReportSendDay,
		ShouldWeeklyEmail:   in.ShouldWeeklyEmail,
		DailyReportSendHour: in.DailyReportSendHour,
		ShouldDailyEmail:    in.ShouldDailyEmail,
		UserTimezone:        in.UserTimezone,
		Subscriptions:       subs,
		WebhookCurrentKey:   in.WebhookCurrentKey,
		WebhookPreviousKey:  in.WebhookPreviousKey,
	}
}

func DomainToUserEvents(in []*am.Event) []*prototypes.EventData {
	events := make([]*prototypes.EventData, 0)
	if in != nil {
		for _, event := range in {
			events = append(events, DomainToEvent(event))
		}
	}
	return events
}

func EventsToDomain(in []*prototypes.EventData) []*am.Event {
	events := make([]*am.Event, 0)
	if in != nil {
		for _, event := range in {
			events = append(events, EventToDomain(event))
		}
	}
	return events
}

func DomainToEventFilter(in *am.EventFilter) *prototypes.EventFilter {
	return &prototypes.EventFilter{
		Start:   in.Start,
		Limit:   in.Limit,
		Filters: DomainToFilterTypes(in.Filters),
	}
}

func EventFilterToDomain(in *prototypes.EventFilter) *am.EventFilter {
	return &am.EventFilter{
		Start:   in.Start,
		Limit:   in.Limit,
		Filters: FilterTypesToDomain(in.Filters),
	}
}
