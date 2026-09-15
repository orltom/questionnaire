package domain

import (
	"errors"
	"slices"
	"time"
	"uuid"
)

type ChallengeID uuid.UUID

type Access string

const (
	AccessPublic  Access = "public"
	AccessInvited Access = "invited"
)

func (a Access) Valid() bool {
	return a == AccessPublic || a == AccessInvited
}

type State string

const (
	Scheduled State = "scheduled"
	Running   State = "running"
	Finished  State = "finished"
)

type Challenge struct {
	id       ChallengeID
	owner    UserID
	quiz     QuizSnapshot
	startsAt time.Time
	endsAt   *time.Time
	closedAt *time.Time
	access   Access
	invitees []UserID
}

func NewChallenge(owner UserID, snapshot QuizSnapshot, startsAt time.Time, endsAt *time.Time, access Access, invitees ...UserID) (Challenge, error) {
	if !access.Valid() {
		return Challenge{}, errors.New("unknown access")
	}

	if startsAt.IsZero() {
		return Challenge{}, errors.New("start is required")
	}

	if endsAt != nil && !endsAt.After(startsAt) {
		return Challenge{}, errors.New("end must be after the start")
	}

	return Challenge{
		id:       ChallengeID(uuid.NewV7()),
		owner:    owner,
		quiz:     snapshot,
		startsAt: startsAt,
		endsAt:   endsAt,
		closedAt: nil,
		access:   access,
		invitees: invitees,
	}, nil
}

func (c *Challenge) ID() ChallengeID {
	return c.id
}

func (c *Challenge) Owner() UserID {
	return c.owner
}

func (c *Challenge) Quiz() QuizSnapshot {
	return c.quiz
}

func (c *Challenge) StartsAt() time.Time {
	return c.startsAt
}

func (c *Challenge) EndsAt() *time.Time {
	return c.endsAt
}

func (c *Challenge) ClosedAt() *time.Time {
	return c.closedAt
}

func (c *Challenge) Access() Access {
	return c.access
}

func (c *Challenge) Invitees() []UserID {
	return slices.Clone(c.invitees)
}

func (c *Challenge) State(now time.Time) State {
	if c.closedAt != nil && !now.Before(*c.closedAt) {
		return Finished
	}

	if now.Before(c.startsAt) {
		return Scheduled
	}

	if c.endsAt != nil && !now.Before(*c.endsAt) {
		return Finished
	}

	return Running
}

func (c *Challenge) CanEdit(actor UserID) bool {
	return c.owner == actor
}

func (c *Challenge) CanView(actor UserID) bool {
	if c.owner == actor || c.access == AccessPublic {
		return true
	}

	return slices.Contains(c.invitees, actor)
}

func (c *Challenge) CanParticipate(actor UserID, now time.Time) bool {
	return c.CanView(actor) && c.State(now) == Running
}

func (c *Challenge) Close(now time.Time) error {
	if c.closedAt != nil {
		return errors.New("challenge is already closed")
	}

	c.closedAt = &now

	return nil
}

func (c *Challenge) Invite(actor UserID) error {
	if c.access != AccessInvited {
		return errors.New("only an invited challenge can have invitees")
	}

	if slices.Contains(c.invitees, actor) {
		return nil
	}

	c.invitees = append(c.invitees, actor)

	return nil
}
