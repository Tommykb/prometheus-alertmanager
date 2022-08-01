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

package general

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/prometheus/alertmanager/api/v2/models"
)

// GetMeOKCode is the HTTP code returned for type GetMeOK
const GetMeOKCode int = 200

/*GetMeOK Get user info

swagger:response getMeOK
*/
type GetMeOK struct {

	/*
	  In: Body
	*/
	Payload *models.User `json:"body,omitempty"`
}

// NewGetMeOK creates GetMeOK with default headers values
func NewGetMeOK() *GetMeOK {

	return &GetMeOK{}
}

// WithPayload adds the payload to the get me o k response
func (o *GetMeOK) WithPayload(payload *models.User) *GetMeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get me o k response
func (o *GetMeOK) SetPayload(payload *models.User) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetMeNoContentCode is the HTTP code returned for type GetMeNoContent
const GetMeNoContentCode int = 204

/*GetMeNoContent No user in basic authentication or a specified header

swagger:response getMeNoContent
*/
type GetMeNoContent struct {
}

// NewGetMeNoContent creates GetMeNoContent with default headers values
func NewGetMeNoContent() *GetMeNoContent {

	return &GetMeNoContent{}
}

// WriteResponse to the client
func (o *GetMeNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}
