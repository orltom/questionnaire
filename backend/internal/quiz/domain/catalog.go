package domain

import (
	"errors"
	"slices"
	"uuid"
)

type CatalogID uuid.UUID

type Catalog struct {
	id          CatalogID
	title       string
	description string
	questionIDs []QuestionID
}

func NewCatalog(title string, description string, refs ...QuestionID) (Catalog, error) {
	if len(title) == 0 {
		return Catalog{}, errors.New("title is required")
	}

	return Catalog{
		id:          CatalogID(uuid.New()),
		title:       title,
		description: description,
		questionIDs: refs,
	}, nil
}

func (c *Catalog) ID() CatalogID {
	return c.id
}

func (c *Catalog) Title() string {
	return c.title
}

func (c *Catalog) Description() string {
	return c.description
}

func (c *Catalog) QuestionIDs() []QuestionID {
	return slices.Clone(c.questionIDs)
}

func (c *Catalog) AddQuestion(ref QuestionID) error {
	if slices.Contains(c.questionIDs, ref) {
		return errors.New("question already exists")
	}

	c.questionIDs = append(c.questionIDs, ref)

	return nil
}

func (c *Catalog) Rename(title string) error {
	if len(title) == 0 {
		return errors.New("title is required")
	}

	c.title = title

	return nil
}

func (c *Catalog) ChangeDescription(description string) error {
	if len(description) == 0 {
		return errors.New("description is required")
	}

	c.description = description

	return nil
}
