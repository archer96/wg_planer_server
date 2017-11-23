// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ListItem list item
// swagger:model ListItem

type ListItem struct {

	// bought at
	// Read Only: true
	BoughtAt strfmt.DateTime `json:"boughtAt,omitempty"`

	// category
	// Required: true
	Category *string `json:"category"`

	// count
	// Required: true
	Count *int64 `json:"count"`

	// created at
	// Read Only: true
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// group Uid
	// Read Only: true
	GroupUID strfmt.UUID `json:"groupUid,omitempty"`

	// id
	// Read Only: true
	ID strfmt.UUID `json:"id,omitempty"`

	// price
	Price int64 `json:"price,omitempty"`

	// requested by
	// Read Only: true
	RequestedBy string `json:"requestedBy,omitempty"`

	// requested for
	RequestedFor []string `json:"requestedFor"`

	// title
	// Required: true
	// Max Length: 150
	Title *string `json:"title"`

	// updated at
	// Read Only: true
	UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
}

/* polymorph ListItem boughtAt false */

/* polymorph ListItem category false */

/* polymorph ListItem count false */

/* polymorph ListItem createdAt false */

/* polymorph ListItem groupUid false */

/* polymorph ListItem id false */

/* polymorph ListItem price false */

/* polymorph ListItem requestedBy false */

/* polymorph ListItem requestedFor false */

/* polymorph ListItem title false */

/* polymorph ListItem updatedAt false */

// Validate validates this list item
func (m *ListItem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCategory(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateRequestedFor(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ListItem) validateCategory(formats strfmt.Registry) error {

	if err := validate.Required("category", "body", m.Category); err != nil {
		return err
	}

	return nil
}

func (m *ListItem) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("count", "body", m.Count); err != nil {
		return err
	}

	return nil
}

func (m *ListItem) validateRequestedFor(formats strfmt.Registry) error {

	if swag.IsZero(m.RequestedFor) { // not required
		return nil
	}

	return nil
}

func (m *ListItem) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	if err := validate.MaxLength("title", "body", string(*m.Title), 150); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ListItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ListItem) UnmarshalBinary(b []byte) error {
	var res ListItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
