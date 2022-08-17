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

package alertgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/prometheus/alertmanager/api/v2/models"
)

// GetAlertGroupsReader is a Reader for the GetAlertGroups structure.
type GetAlertGroupsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAlertGroupsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAlertGroupsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetAlertGroupsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAlertGroupsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAlertGroupsOK creates a GetAlertGroupsOK with default headers values
func NewGetAlertGroupsOK() *GetAlertGroupsOK {
	return &GetAlertGroupsOK{}
}

/*
GetAlertGroupsOK handles this case with default header values.

Get alert groups response
*/
type GetAlertGroupsOK struct {
	Payload models.AlertGroups
}

func (o *GetAlertGroupsOK) Error() string {
	return fmt.Sprintf("[GET /alerts/groups][%d] getAlertGroupsOK  %+v", 200, o.Payload)
}

func (o *GetAlertGroupsOK) GetPayload() models.AlertGroups {
	return o.Payload
}

func (o *GetAlertGroupsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlertGroupsBadRequest creates a GetAlertGroupsBadRequest with default headers values
func NewGetAlertGroupsBadRequest() *GetAlertGroupsBadRequest {
	return &GetAlertGroupsBadRequest{}
}

/*
GetAlertGroupsBadRequest handles this case with default header values.

Bad request
*/
type GetAlertGroupsBadRequest struct {
	Payload string
}

func (o *GetAlertGroupsBadRequest) Error() string {
	return fmt.Sprintf("[GET /alerts/groups][%d] getAlertGroupsBadRequest  %+v", 400, o.Payload)
}

func (o *GetAlertGroupsBadRequest) GetPayload() string {
	return o.Payload
}

func (o *GetAlertGroupsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlertGroupsInternalServerError creates a GetAlertGroupsInternalServerError with default headers values
func NewGetAlertGroupsInternalServerError() *GetAlertGroupsInternalServerError {
	return &GetAlertGroupsInternalServerError{}
}

/*
GetAlertGroupsInternalServerError handles this case with default header values.

Internal server error
*/
type GetAlertGroupsInternalServerError struct {
	Payload string
}

func (o *GetAlertGroupsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /alerts/groups][%d] getAlertGroupsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAlertGroupsInternalServerError) GetPayload() string {
	return o.Payload
}

func (o *GetAlertGroupsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
