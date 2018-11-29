package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"

	"github.com/AlexRudd/cognito-srp"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func main() {
	csrp, err := cognitosrp.NewCognitoSRP("isaac.dawson@linkai.io", "H0h0h0h0.",
		"us-east-1_gZiBStGnL",
		"68c9ta93d2mt6qo5jsbnpgdnpr", aws.String("1kb8c1ngs0lecjnn8dcsanbqkoqgm4sebfrikqalh176flc8379r")) // )
	if err != nil {
		log.Printf("failed creating CognitoSRP: %s", err.Error())
		return
	}

	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = endpoints.UsEast1RegionID
	cfg.Credentials = aws.AnonymousCredentials
	svc := cip.New(cfg)

	// initiate auth
	req := svc.InitiateAuthRequest(&cip.InitiateAuthInput{
		AuthFlow:       cip.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	resp, err := req.Send()
	if err != nil {
		checkError(err)
		log.Printf("%#v error sending initiate auth request\n", err)
		return
	}

	// respond to password verifier challenge
	if resp.ChallengeName == cip.ChallengeNameTypePasswordVerifier {
		challengeInput, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())
		chal := svc.RespondToAuthChallengeRequest(challengeInput)
		resp, err := chal.Send()
		if err != nil {
			checkError(err)
			log.Printf("Respond %#v error sending initiate auth request\n", err)
			return
		}

		if resp.ChallengeName == cip.ChallengeNameTypeNewPasswordRequired {
			log.Printf("New Password Required %#v\n", resp)
			secretHash, err := csrp.GetSecretHash("isaac.dawson@linkai.io")
			if err != nil {
				log.Printf("failed to get secret hash for user: %v\n", err)
				return
			}
			ch := svc.RespondToAuthChallengeRequest(&cip.RespondToAuthChallengeInput{
				ChallengeName: resp.ChallengeName,
				ChallengeResponses: map[string]string{
					"NEW_PASSWORD":                            "H0h0h0h02.",
					"USERNAME":                                "isaac.dawson@linkai.io",
					"SECRET_HASH":                             secretHash,
					"userAttributes.given_name":               "test",
					"userAttributes.family_name":              "test",
					"userAttributes.custom:organization_name": "test-org",
				},
				ClientId: aws.String(csrp.GetClientId()),
				Session:  resp.Session,
			})

			chpwdresp, err := ch.Send()
			if err != nil {
				checkError(err)
				log.Printf("Respond %#v error sending auth challenge response\n", err)
				return
			}
			log.Printf("TOKENS? %#v\n", chpwdresp)
			/*
				params := &cip.ChangePasswordInput{
					AccessToken:      chpwdresp.AuthenticationResult.AccessToken,
					PreviousPassword: aws.String("H0h0h0h0."),
					ProposedPassword: aws.String("H0h0h0h02."),
				}
				challengeInput := svc.ChangePasswordRequest(params)
				chpwdresp, err := challengeInput.Send()
				if err != nil {
					checkError(err)
					log.Printf("change password error: %#v\n", err)
					return
				}
				log.Printf("TOKENS? %#v\n", chpwdresp)
			*/
		}

		// print the tokens
		log.Printf("TOKENS? %#v\n", resp)
	} else {
		log.Printf("resp: %v %#v\n", resp.ChallengeName, resp)
	}
}

func checkError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		// Get error details
		log.Println("Respond Error:", awsErr.Code(), awsErr.Message())

		// Prints out full error message, including original error if there was one.
		log.Println("Respond Error:", awsErr.Error())

		// Get original error
		if origErr := awsErr.OrigErr(); origErr != nil {
			// operate on original error.
		}
	}
}
