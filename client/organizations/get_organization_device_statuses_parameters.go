// Code generated by go-swagger; DO NOT EDIT.

package organizations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetOrganizationDeviceStatusesParams creates a new GetOrganizationDeviceStatusesParams object
// with the default values initialized.
func NewGetOrganizationDeviceStatusesParams() *GetOrganizationDeviceStatusesParams {
	var ()
	return &GetOrganizationDeviceStatusesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetOrganizationDeviceStatusesParamsWithTimeout creates a new GetOrganizationDeviceStatusesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetOrganizationDeviceStatusesParamsWithTimeout(timeout time.Duration) *GetOrganizationDeviceStatusesParams {
	var ()
	return &GetOrganizationDeviceStatusesParams{

		timeout: timeout,
	}
}

// NewGetOrganizationDeviceStatusesParamsWithContext creates a new GetOrganizationDeviceStatusesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetOrganizationDeviceStatusesParamsWithContext(ctx context.Context) *GetOrganizationDeviceStatusesParams {
	var ()
	return &GetOrganizationDeviceStatusesParams{

		Context: ctx,
	}
}

// NewGetOrganizationDeviceStatusesParamsWithHTTPClient creates a new GetOrganizationDeviceStatusesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetOrganizationDeviceStatusesParamsWithHTTPClient(client *http.Client) *GetOrganizationDeviceStatusesParams {
	var ()
	return &GetOrganizationDeviceStatusesParams{
		HTTPClient: client,
	}
}

/*GetOrganizationDeviceStatusesParams contains all the parameters to send to the API endpoint
for the get organization device statuses operation typically these are written to a http.Request
*/
type GetOrganizationDeviceStatusesParams struct {

	/*OrganizationID*/
	OrganizationID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) WithTimeout(timeout time.Duration) *GetOrganizationDeviceStatusesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) WithContext(ctx context.Context) *GetOrganizationDeviceStatusesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) WithHTTPClient(client *http.Client) *GetOrganizationDeviceStatusesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOrganizationID adds the organizationID to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) WithOrganizationID(organizationID string) *GetOrganizationDeviceStatusesParams {
	o.SetOrganizationID(organizationID)
	return o
}

// SetOrganizationID adds the organizationId to the get organization device statuses params
func (o *GetOrganizationDeviceStatusesParams) SetOrganizationID(organizationID string) {
	o.OrganizationID = organizationID
}

// WriteToRequest writes these params to a swagger request
func (o *GetOrganizationDeviceStatusesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param organizationId
	if err := r.SetPathParam("organizationId", o.OrganizationID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}