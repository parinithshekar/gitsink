// Code generated by go-swagger; DO NOT EDIT.

package intrusion_settings

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

// NewGetOrganizationSecurityIntrusionSettingsParams creates a new GetOrganizationSecurityIntrusionSettingsParams object
// with the default values initialized.
func NewGetOrganizationSecurityIntrusionSettingsParams() *GetOrganizationSecurityIntrusionSettingsParams {
	var ()
	return &GetOrganizationSecurityIntrusionSettingsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetOrganizationSecurityIntrusionSettingsParamsWithTimeout creates a new GetOrganizationSecurityIntrusionSettingsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetOrganizationSecurityIntrusionSettingsParamsWithTimeout(timeout time.Duration) *GetOrganizationSecurityIntrusionSettingsParams {
	var ()
	return &GetOrganizationSecurityIntrusionSettingsParams{

		timeout: timeout,
	}
}

// NewGetOrganizationSecurityIntrusionSettingsParamsWithContext creates a new GetOrganizationSecurityIntrusionSettingsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetOrganizationSecurityIntrusionSettingsParamsWithContext(ctx context.Context) *GetOrganizationSecurityIntrusionSettingsParams {
	var ()
	return &GetOrganizationSecurityIntrusionSettingsParams{

		Context: ctx,
	}
}

// NewGetOrganizationSecurityIntrusionSettingsParamsWithHTTPClient creates a new GetOrganizationSecurityIntrusionSettingsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetOrganizationSecurityIntrusionSettingsParamsWithHTTPClient(client *http.Client) *GetOrganizationSecurityIntrusionSettingsParams {
	var ()
	return &GetOrganizationSecurityIntrusionSettingsParams{
		HTTPClient: client,
	}
}

/*GetOrganizationSecurityIntrusionSettingsParams contains all the parameters to send to the API endpoint
for the get organization security intrusion settings operation typically these are written to a http.Request
*/
type GetOrganizationSecurityIntrusionSettingsParams struct {

	/*OrganizationID*/
	OrganizationID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) WithTimeout(timeout time.Duration) *GetOrganizationSecurityIntrusionSettingsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) WithContext(ctx context.Context) *GetOrganizationSecurityIntrusionSettingsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) WithHTTPClient(client *http.Client) *GetOrganizationSecurityIntrusionSettingsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOrganizationID adds the organizationID to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) WithOrganizationID(organizationID string) *GetOrganizationSecurityIntrusionSettingsParams {
	o.SetOrganizationID(organizationID)
	return o
}

// SetOrganizationID adds the organizationId to the get organization security intrusion settings params
func (o *GetOrganizationSecurityIntrusionSettingsParams) SetOrganizationID(organizationID string) {
	o.OrganizationID = organizationID
}

// WriteToRequest writes these params to a swagger request
func (o *GetOrganizationSecurityIntrusionSettingsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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