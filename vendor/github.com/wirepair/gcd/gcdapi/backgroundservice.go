// AUTO-GENERATED Chrome Remote Debugger Protocol API Client
// This file contains BackgroundService functionality.
// API Version: 1.3

package gcdapi

import (
	"github.com/wirepair/gcd/gcdmessage"
)

// A key-value pair for additional event information to pass along.
type BackgroundServiceEventMetadata struct {
	Key   string `json:"key"`   //
	Value string `json:"value"` //
}

// No Description.
type BackgroundServiceBackgroundServiceEvent struct {
	Timestamp                   float64                           `json:"timestamp"`                   // Timestamp of the event (in seconds).
	Origin                      string                            `json:"origin"`                      // The origin this event belongs to.
	ServiceWorkerRegistrationId string                            `json:"serviceWorkerRegistrationId"` // The Service Worker ID that initiated the event.
	Service                     string                            `json:"service"`                     // The Background Service this event belongs to. enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	EventName                   string                            `json:"eventName"`                   // A description of the event.
	InstanceId                  string                            `json:"instanceId"`                  // An identifier that groups related events together.
	EventMetadata               []*BackgroundServiceEventMetadata `json:"eventMetadata"`               // A list of event-specific information.
}

// Called when the recording state for the service has been updated.
type BackgroundServiceRecordingStateChangedEvent struct {
	Method string `json:"method"`
	Params struct {
		IsRecording bool   `json:"isRecording"` //
		Service     string `json:"service"`     //  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	} `json:"Params,omitempty"`
}

// Called with all existing backgroundServiceEvents when enabled, and all new events afterwards if enabled and recording.
type BackgroundServiceBackgroundServiceEventReceivedEvent struct {
	Method string `json:"method"`
	Params struct {
		BackgroundServiceEvent *BackgroundServiceBackgroundServiceEvent `json:"backgroundServiceEvent"` //
	} `json:"Params,omitempty"`
}

type BackgroundService struct {
	target gcdmessage.ChromeTargeter
}

func NewBackgroundService(target gcdmessage.ChromeTargeter) *BackgroundService {
	c := &BackgroundService{target: target}
	return c
}

type BackgroundServiceStartObservingParams struct {
	//  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	Service string `json:"service"`
}

// StartObservingWithParams - Enables event updates for the service.
func (c *BackgroundService) StartObservingWithParams(v *BackgroundServiceStartObservingParams) (*gcdmessage.ChromeResponse, error) {
	return gcdmessage.SendDefaultRequest(c.target, c.target.GetSendCh(), &gcdmessage.ParamRequest{Id: c.target.GetId(), Method: "BackgroundService.startObserving", Params: v})
}

// StartObserving - Enables event updates for the service.
// service -  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
func (c *BackgroundService) StartObserving(service string) (*gcdmessage.ChromeResponse, error) {
	var v BackgroundServiceStartObservingParams
	v.Service = service
	return c.StartObservingWithParams(&v)
}

type BackgroundServiceStopObservingParams struct {
	//  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	Service string `json:"service"`
}

// StopObservingWithParams - Disables event updates for the service.
func (c *BackgroundService) StopObservingWithParams(v *BackgroundServiceStopObservingParams) (*gcdmessage.ChromeResponse, error) {
	return gcdmessage.SendDefaultRequest(c.target, c.target.GetSendCh(), &gcdmessage.ParamRequest{Id: c.target.GetId(), Method: "BackgroundService.stopObserving", Params: v})
}

// StopObserving - Disables event updates for the service.
// service -  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
func (c *BackgroundService) StopObserving(service string) (*gcdmessage.ChromeResponse, error) {
	var v BackgroundServiceStopObservingParams
	v.Service = service
	return c.StopObservingWithParams(&v)
}

type BackgroundServiceSetRecordingParams struct {
	//
	ShouldRecord bool `json:"shouldRecord"`
	//  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	Service string `json:"service"`
}

// SetRecordingWithParams - Set the recording state for the service.
func (c *BackgroundService) SetRecordingWithParams(v *BackgroundServiceSetRecordingParams) (*gcdmessage.ChromeResponse, error) {
	return gcdmessage.SendDefaultRequest(c.target, c.target.GetSendCh(), &gcdmessage.ParamRequest{Id: c.target.GetId(), Method: "BackgroundService.setRecording", Params: v})
}

// SetRecording - Set the recording state for the service.
// shouldRecord -
// service -  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
func (c *BackgroundService) SetRecording(shouldRecord bool, service string) (*gcdmessage.ChromeResponse, error) {
	var v BackgroundServiceSetRecordingParams
	v.ShouldRecord = shouldRecord
	v.Service = service
	return c.SetRecordingWithParams(&v)
}

type BackgroundServiceClearEventsParams struct {
	//  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
	Service string `json:"service"`
}

// ClearEventsWithParams - Clears all stored data for the service.
func (c *BackgroundService) ClearEventsWithParams(v *BackgroundServiceClearEventsParams) (*gcdmessage.ChromeResponse, error) {
	return gcdmessage.SendDefaultRequest(c.target, c.target.GetSendCh(), &gcdmessage.ParamRequest{Id: c.target.GetId(), Method: "BackgroundService.clearEvents", Params: v})
}

// ClearEvents - Clears all stored data for the service.
// service -  enum values: backgroundFetch, backgroundSync, pushMessaging, notifications, paymentHandler
func (c *BackgroundService) ClearEvents(service string) (*gcdmessage.ChromeResponse, error) {
	var v BackgroundServiceClearEventsParams
	v.Service = service
	return c.ClearEventsWithParams(&v)
}
