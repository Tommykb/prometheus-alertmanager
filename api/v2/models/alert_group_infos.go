// Code generated by go-swagger; DO NOT EDIT.

// Copyright Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AlertGroupInfos alert group infos
//
// swagger:model alertGroupInfos
type AlertGroupInfos struct {

	// alert group infos
	AlertGroupInfos []*AlertGroupInfo `json:"alertGroupInfos"`

	// next token
	NextToken string `json:"nextToken,omitempty"`
}

// Validate validates this alert group infos
func (m *AlertGroupInfos) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlertGroupInfos(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertGroupInfos) validateAlertGroupInfos(formats strfmt.Registry) error {
	if swag.IsZero(m.AlertGroupInfos) { // not required
		return nil
	}

	for i := 0; i < len(m.AlertGroupInfos); i++ {
		if swag.IsZero(m.AlertGroupInfos[i]) { // not required
			continue
		}

		if m.AlertGroupInfos[i] != nil {
			if err := m.AlertGroupInfos[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("alertGroupInfos" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("alertGroupInfos" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this alert group infos based on the context it is used
func (m *AlertGroupInfos) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAlertGroupInfos(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertGroupInfos) contextValidateAlertGroupInfos(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.AlertGroupInfos); i++ {

		if m.AlertGroupInfos[i] != nil {
			if err := m.AlertGroupInfos[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("alertGroupInfos" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("alertGroupInfos" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertGroupInfos) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertGroupInfos) UnmarshalBinary(b []byte) error {
	var res AlertGroupInfos
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
