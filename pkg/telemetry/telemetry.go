/*
2021 © Postgres.ai
*/

// Package telemetry contains tools to collect Database Lab Engine data.
package telemetry

import (
	"context"
	"time"

	"gitlab.com/postgres-ai/database-lab/v2/pkg/client/platform"
	"gitlab.com/postgres-ai/database-lab/v2/pkg/log"
)

const (
	// EngineStartedEvent defines the engine start event.
	EngineStartedEvent = "engine_started"

	// EngineStoppedEvent describes the engine stop event.
	EngineStoppedEvent = "engine_stopped"

	// CloneCreatedEvent describes the clone creation event.
	CloneCreatedEvent = "clone_created"

	// CloneResetEvent describes the clone reset event.
	CloneResetEvent = "clone_reset"

	// CloneDestroyedEvent describes a clone destruction event.
	CloneDestroyedEvent = "clone_destroyed"

	// SnapshotCreatedEvent describes a snapshot creation event.
	SnapshotCreatedEvent = "snapshot_created"

	// AlertEvent describes alert events.
	AlertEvent = "alert"
)

// Agent represent a telemetry agent to collect engine data.
type Agent struct {
	instanceID string
	platform   *platform.Client
}

// New creates a new agent.
func New(instanceID string, platform *platform.Client) *Agent {
	return &Agent{instanceID: instanceID, platform: platform}
}

// SendEvent sends a telemetry event.
func (a *Agent) SendEvent(ctx context.Context, eventType string, payload interface{}) {
	_, err := a.platform.SendTelemetryEvent(ctx, platform.TelemetryEvent{
		InstanceID: a.instanceID,
		EventType:  eventType,
		Timestamp:  time.Now(),
		Payload:    payload,
	})

	if err != nil {
		log.Err("Failed to send telemetry event", err)
	}
}
