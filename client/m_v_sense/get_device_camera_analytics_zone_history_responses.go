// Code generated by go-swagger; DO NOT EDIT.

package m_v_sense

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// GetDeviceCameraAnalyticsZoneHistoryReader is a Reader for the GetDeviceCameraAnalyticsZoneHistory structure.
type GetDeviceCameraAnalyticsZoneHistoryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceCameraAnalyticsZoneHistoryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDeviceCameraAnalyticsZoneHistoryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetDeviceCameraAnalyticsZoneHistoryOK creates a GetDeviceCameraAnalyticsZoneHistoryOK with default headers values
func NewGetDeviceCameraAnalyticsZoneHistoryOK() *GetDeviceCameraAnalyticsZoneHistoryOK {
	return &GetDeviceCameraAnalyticsZoneHistoryOK{}
}

/*GetDeviceCameraAnalyticsZoneHistoryOK handles this case with default header values.

Successful operation
*/
type GetDeviceCameraAnalyticsZoneHistoryOK struct {
	Payload interface{}
}

func (o *GetDeviceCameraAnalyticsZoneHistoryOK) Error() string {
	return fmt.Sprintf("[GET /devices/{serial}/camera/analytics/zones/{zoneId}/history][%d] getDeviceCameraAnalyticsZoneHistoryOK  %+v", 200, o.Payload)
}

func (o *GetDeviceCameraAnalyticsZoneHistoryOK) GetPayload() interface{} {
	return o.Payload
}

func (o *GetDeviceCameraAnalyticsZoneHistoryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}