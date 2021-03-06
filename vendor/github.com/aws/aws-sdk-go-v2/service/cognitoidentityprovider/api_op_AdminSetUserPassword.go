// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cognitoidentityprovider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/cognito-idp-2016-04-18/AdminSetUserPasswordRequest
type AdminSetUserPasswordInput struct {
	_ struct{} `type:"structure"`

	// Password is a required field
	Password *string `min:"6" type:"string" required:"true"`

	Permanent *bool `type:"boolean"`

	// UserPoolId is a required field
	UserPoolId *string `min:"1" type:"string" required:"true"`

	// Username is a required field
	Username *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s AdminSetUserPasswordInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *AdminSetUserPasswordInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "AdminSetUserPasswordInput"}

	if s.Password == nil {
		invalidParams.Add(aws.NewErrParamRequired("Password"))
	}
	if s.Password != nil && len(*s.Password) < 6 {
		invalidParams.Add(aws.NewErrParamMinLen("Password", 6))
	}

	if s.UserPoolId == nil {
		invalidParams.Add(aws.NewErrParamRequired("UserPoolId"))
	}
	if s.UserPoolId != nil && len(*s.UserPoolId) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("UserPoolId", 1))
	}

	if s.Username == nil {
		invalidParams.Add(aws.NewErrParamRequired("Username"))
	}
	if s.Username != nil && len(*s.Username) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("Username", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/cognito-idp-2016-04-18/AdminSetUserPasswordResponse
type AdminSetUserPasswordOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s AdminSetUserPasswordOutput) String() string {
	return awsutil.Prettify(s)
}

const opAdminSetUserPassword = "AdminSetUserPassword"

// AdminSetUserPasswordRequest returns a request value for making API operation for
// Amazon Cognito Identity Provider.
//
//    // Example sending a request using AdminSetUserPasswordRequest.
//    req := client.AdminSetUserPasswordRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cognito-idp-2016-04-18/AdminSetUserPassword
func (c *Client) AdminSetUserPasswordRequest(input *AdminSetUserPasswordInput) AdminSetUserPasswordRequest {
	op := &aws.Operation{
		Name:       opAdminSetUserPassword,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &AdminSetUserPasswordInput{}
	}

	req := c.newRequest(op, input, &AdminSetUserPasswordOutput{})
	return AdminSetUserPasswordRequest{Request: req, Input: input, Copy: c.AdminSetUserPasswordRequest}
}

// AdminSetUserPasswordRequest is the request type for the
// AdminSetUserPassword API operation.
type AdminSetUserPasswordRequest struct {
	*aws.Request
	Input *AdminSetUserPasswordInput
	Copy  func(*AdminSetUserPasswordInput) AdminSetUserPasswordRequest
}

// Send marshals and sends the AdminSetUserPassword API request.
func (r AdminSetUserPasswordRequest) Send(ctx context.Context) (*AdminSetUserPasswordResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &AdminSetUserPasswordResponse{
		AdminSetUserPasswordOutput: r.Request.Data.(*AdminSetUserPasswordOutput),
		response:                   &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// AdminSetUserPasswordResponse is the response type for the
// AdminSetUserPassword API operation.
type AdminSetUserPasswordResponse struct {
	*AdminSetUserPasswordOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// AdminSetUserPassword request.
func (r *AdminSetUserPasswordResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
