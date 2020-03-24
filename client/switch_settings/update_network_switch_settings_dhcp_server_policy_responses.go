// Code generated by go-swagger; DO NOT EDIT.

package switch_settings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// UpdateNetworkSwitchSettingsDhcpServerPolicyReader is a Reader for the UpdateNetworkSwitchSettingsDhcpServerPolicy structure.
type UpdateNetworkSwitchSettingsDhcpServerPolicyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateNetworkSwitchSettingsDhcpServerPolicyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateNetworkSwitchSettingsDhcpServerPolicyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateNetworkSwitchSettingsDhcpServerPolicyOK creates a UpdateNetworkSwitchSettingsDhcpServerPolicyOK with default headers values
func NewUpdateNetworkSwitchSettingsDhcpServerPolicyOK() *UpdateNetworkSwitchSettingsDhcpServerPolicyOK {
	return &UpdateNetworkSwitchSettingsDhcpServerPolicyOK{}
}

/*UpdateNetworkSwitchSettingsDhcpServerPolicyOK handles this case with default header values.

Successful operation
*/
type UpdateNetworkSwitchSettingsDhcpServerPolicyOK struct {
	Payload interface{}
}

func (o *UpdateNetworkSwitchSettingsDhcpServerPolicyOK) Error() string {
	return fmt.Sprintf("[PUT /networks/{networkId}/switch/settings/dhcpServerPolicy][%d] updateNetworkSwitchSettingsDhcpServerPolicyOK  %+v", 200, o.Payload)
}

func (o *UpdateNetworkSwitchSettingsDhcpServerPolicyOK) GetPayload() interface{} {
	return o.Payload
}

func (o *UpdateNetworkSwitchSettingsDhcpServerPolicyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}