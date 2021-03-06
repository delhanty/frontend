package user

import (
	"context"
	"io"
	"time"

	"github.com/linkai-io/am/pkg/convert"
	"github.com/linkai-io/am/pkg/retrier"
	"github.com/pkg/errors"

	"github.com/linkai-io/am/am"
	service "github.com/linkai-io/am/protocservices/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type Client struct {
	client         service.UserServiceClient
	conn           *grpc.ClientConn
	defaultTimeout time.Duration
}

func New() *Client {
	return &Client{defaultTimeout: (time.Second * 10)}
}

func (c *Client) Init(config []byte) error {
	conn, err := grpc.DialContext(context.Background(), "srv://consul/"+am.UserServiceKey, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = service.NewUserServiceClient(conn)
	return nil
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.defaultTimeout = timeout
}

func (c *Client) get(ctx context.Context, userContext am.UserContext, in *service.UserRequest) (oid int, user *am.User, err error) {
	var resp *service.UserResponse

	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.Get(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to get users from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, nil, err
	}
	return int(resp.GetOrgID()), convert.UserToDomain(resp.User), nil
}

func (c *Client) Get(ctx context.Context, userContext am.UserContext, userEmail string) (oid int, user *am.User, err error) {
	in := &service.UserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		By:          service.UserRequest_USEREMAIL,
		UserEmail:   userEmail,
	}
	return c.get(ctx, userContext, in)
}

func (c *Client) GetWithOrgID(ctx context.Context, userContext am.UserContext, orgID int, userCID string) (oid int, user *am.User, err error) {
	in := &service.UserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		By:          service.UserRequest_USERWITHORGID,
		OrgID:       int32(orgID),
		UserCID:     userCID,
	}
	return c.get(ctx, userContext, in)
}

func (c *Client) GetByID(ctx context.Context, userContext am.UserContext, userID int) (oid int, user *am.User, err error) {
	in := &service.UserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		By:          service.UserRequest_USERID,
		UserID:      int32(userID),
	}
	return c.get(ctx, userContext, in)
}

func (c *Client) GetByCID(ctx context.Context, userContext am.UserContext, userCID string) (oid int, user *am.User, err error) {
	in := &service.UserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		By:          service.UserRequest_USERCID,
		UserCID:     userCID,
	}
	return c.get(ctx, userContext, in)
}

func (c *Client) List(ctx context.Context, userContext am.UserContext, filter *am.UserFilter) (oid int, users []*am.User, err error) {
	var resp service.UserService_ListClient
	oid = userContext.GetOrgID()

	in := &service.UserListRequest{
		UserContext: convert.DomainToUserContext(userContext),
		UserFilter:  convert.DomainToUserFilter(filter),
	}

	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.List(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to list users from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, nil, err
	}

	users = make([]*am.User, 0)
	for {
		userResp, err := resp.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, nil, err
		}
		users = append(users, convert.UserToDomain(userResp.User))
	}
	return oid, users, nil
}

func (c *Client) Create(ctx context.Context, userContext am.UserContext, user *am.User) (oid int, uid int, ucid string, err error) {
	var resp *service.UserCreatedResponse

	in := &service.CreateUserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		User:        convert.DomainToUser(user),
	}
	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.Create(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to create user from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, 0, "", err
	}

	return int(resp.GetOrgID()), int(resp.GetUserID()), resp.GetUserCID(), nil
}

func (c *Client) Update(ctx context.Context, userContext am.UserContext, user *am.User, userID int) (oid int, uid int, err error) {
	var resp *service.UserUpdatedResponse

	in := &service.UpdateUserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		User:        convert.DomainToUser(user),
		UserID:      int32(userID),
	}

	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.Update(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to update user from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, 0, err
	}
	return int(resp.GetOrgID()), int(resp.GetUserID()), nil
}

func (c *Client) Delete(ctx context.Context, userContext am.UserContext, userID int) (oid int, err error) {
	var resp *service.UserDeletedResponse

	in := &service.DeleteUserRequest{
		UserContext: convert.DomainToUserContext(userContext),
		UserID:      int32(userID),
	}

	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.Delete(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to delete user from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, err
	}
	return int(resp.GetOrgID()), nil
}

func (c *Client) AcceptAgreement(ctx context.Context, userContext am.UserContext, accepted bool) (oid int, uid int, err error) {
	var resp *service.UserAgreementResponse

	in := &service.UserAgreementRequest{
		UserContext: convert.DomainToUserContext(userContext),
		Agreement:   accepted,
	}

	ctxDeadline, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	err = retrier.RetryIfNot(func() error {
		var retryErr error

		resp, retryErr = c.client.AcceptAgreement(ctxDeadline, in)
		return errors.Wrap(retryErr, "unable to update user agreement from user client")
	}, "rpc error: code = Unavailable desc")

	if err != nil {
		return 0, 0, err
	}
	return int(resp.GetOrgID()), int(resp.GetUserID()), nil
}
