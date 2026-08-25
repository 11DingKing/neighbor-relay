package domain

import "time"

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleFieldWorker Role = "field_worker"
)

type CaseStatus string

const (
	CaseOpen       CaseStatus = "open"
	CaseInProgress CaseStatus = "in_progress"
	CaseResolved   CaseStatus = "resolved"
	CaseClosed     CaseStatus = "closed"
)

type VisitStatus string

const (
	VisitPlanned    VisitStatus = "planned"
	VisitInProgress VisitStatus = "in_progress"
	VisitCompleted  VisitStatus = "completed"
	VisitCancelled  VisitStatus = "cancelled"
)

type User struct {
	ID, Name, Email string
	Role            Role
	Active          bool
	CreatedAt       time.Time
}
type Session struct {
	ID, UserID, TokenHash           string
	ExpiresAt, CreatedAt, RevokedAt *time.Time
}
type Household struct {
	ID, Address, ContactName, Phone string
	Priority                        int
	CreatedAt, UpdatedAt            time.Time
}
type Case struct {
	ID, HouseholdID, Title, Summary string
	Status                          CaseStatus
	OwnerID                         string
	Version                         int64
	CreatedAt, UpdatedAt            time.Time
}
type Visit struct {
	ID, CaseID, HouseholdID, AssigneeID string
	Status                              VisitStatus
	ScheduledFor                        time.Time
	StartedAt, CompletedAt              *time.Time
	Version                             int64
	CreatedAt                           time.Time
}
type VisitNote struct {
	ID, VisitID, AuthorID, Body string
	CreatedAt                   time.Time
}
type AuditEvent struct {
	ID, ActorID, EntityType, EntityID, Action, Result, RequestID string
	CreatedAt                                                    time.Time
}
type OutboxJob struct {
	ID, Kind, AggregateID, Payload, Status string
	Attempts                               int
	AvailableAt, CreatedAt, UpdatedAt      time.Time
	LastError                              string
}

func (s CaseStatus) CanTransition(to CaseStatus) bool {
	switch s {
	case CaseOpen:
		return to == CaseInProgress || to == CaseClosed
	case CaseInProgress:
		return to == CaseResolved || to == CaseOpen
	case CaseResolved:
		return to == CaseClosed || to == CaseInProgress
	case CaseClosed:
		return false
	}
	return false
}
func (s VisitStatus) CanTransition(to VisitStatus) bool {
	switch s {
	case VisitPlanned:
		return to == VisitInProgress || to == VisitCancelled
	case VisitInProgress:
		return to == VisitCompleted || to == VisitCancelled
	case VisitCompleted, VisitCancelled:
		return false
	}
	return false
}
