// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package sts

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

var _ aws.Config
var _ = awsutil.Prettify

// The identifiers for the temporary security credentials that the operation
// returns.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/sts-2011-06-15/AssumedRoleUser
type AssumedRoleUser struct {
	_ struct{} `type:"structure"`

	// The ARN of the temporary security credentials that are returned from the
	// AssumeRole action. For more information about ARNs and how to use them in
	// policies, see IAM Identifiers (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html)
	// in Using IAM.
	//
	// Arn is a required field
	Arn *string `min:"20" type:"string" required:"true"`

	// A unique identifier that contains the role ID and the role session name of
	// the role that is being assumed. The role ID is generated by AWS when the
	// role is created.
	//
	// AssumedRoleId is a required field
	AssumedRoleId *string `min:"2" type:"string" required:"true"`
}

// String returns the string representation
func (s AssumedRoleUser) String() string {
	return awsutil.Prettify(s)
}

// AWS credentials for API authentication.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/sts-2011-06-15/Credentials
type Credentials struct {
	_ struct{} `type:"structure"`

	// The access key ID that identifies the temporary security credentials.
	//
	// AccessKeyId is a required field
	AccessKeyId *string `min:"16" type:"string" required:"true"`

	// The date on which the current credentials expire.
	//
	// Expiration is a required field
	Expiration *time.Time `type:"timestamp" required:"true"`

	// The secret access key that can be used to sign requests.
	//
	// SecretAccessKey is a required field
	SecretAccessKey *string `type:"string" required:"true"`

	// The token that users must pass to the service API to use the temporary credentials.
	//
	// SessionToken is a required field
	SessionToken *string `type:"string" required:"true"`
}

// String returns the string representation
func (s Credentials) String() string {
	return awsutil.Prettify(s)
}

// Identifiers for the federated user that is associated with the credentials.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/sts-2011-06-15/FederatedUser
type FederatedUser struct {
	_ struct{} `type:"structure"`

	// The ARN that specifies the federated user that is associated with the credentials.
	// For more information about ARNs and how to use them in policies, see IAM
	// Identifiers (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html)
	// in Using IAM.
	//
	// Arn is a required field
	Arn *string `min:"20" type:"string" required:"true"`

	// The string that identifies the federated user associated with the credentials,
	// similar to the unique ID of an IAM user.
	//
	// FederatedUserId is a required field
	FederatedUserId *string `min:"2" type:"string" required:"true"`
}

// String returns the string representation
func (s FederatedUser) String() string {
	return awsutil.Prettify(s)
}

// A reference to the IAM managed policy that is passed as a session policy
// for a role session or a federated user session.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/sts-2011-06-15/PolicyDescriptorType
type PolicyDescriptorType struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) of the IAM managed policy to use as a session
	// policy for the role. For more information about ARNs, see Amazon Resource
	// Names (ARNs) and AWS Service Namespaces (https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html)
	// in the AWS General Reference.
	Arn *string `locationName:"arn" min:"20" type:"string"`
}

// String returns the string representation
func (s PolicyDescriptorType) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *PolicyDescriptorType) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "PolicyDescriptorType"}
	if s.Arn != nil && len(*s.Arn) < 20 {
		invalidParams.Add(aws.NewErrParamMinLen("Arn", 20))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}
