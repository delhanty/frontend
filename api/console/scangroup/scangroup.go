package scangroup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/linkai-io/frontend/pkg/serializers"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/go-chi/chi"
	"github.com/linkai-io/frontend/pkg/middleware"
	"github.com/rs/zerolog/log"

	"github.com/linkai-io/am/am"
)

// TODO: should probably make this better/find a good example.
func ValidateSubDomain(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return !strings.ContainsAny(str, " ~`!@#$%^&*()=+[]{};\"'|<>,.?\r\n\x00\x01\x02\x03\x03\x04\x05\x06\x07\x08\t")
}

type ScanGroupDetails struct {
	GroupName          string   `json:"group_name" validate:"required,gte=1,lte=128,excludesall=/"`
	PortScanEnabled    bool     `json:"port_scan_enabled"`
	CustomSubNames     []string `json:"custom_sub_names" validate:"omitempty,max=100,dive,gte=1,lte=128,subdomain"`
	CustomWebPorts     []int32  `json:"custom_web_ports" validate:"omitempty,max=10,dive,gte=1,lte=65535"`
	TCPPorts           []int32  `json:"tcp_ports" validate:"omitempty,max=50,dive,gte=1,lte=65535"`
	AllowedTLDs        []string `json:"allowed_tlds" validate:"omitempty,max=100"`
	AllowedHosts       []string `json:"allowed_hosts" validate:"omitempty,max=5000"`
	DisallowedTLDs     []string `json:"disallowed_tlds" validate:"omitempty,max=100"`
	DisallowedHosts    []string `json:"disallowed_hosts" validate:"omitempty,max=5000"`
	PortsPerSecond     int32    `json:"ports_per_second" validate:"omitempty,gte=1,lte=50"`
	ConcurrentRequests int32    `json:"concurrent_requests" validate:"required,gte=1,lte=20"` // lte 25, beta is limited to 10
	ArchiveAfterDays   int32    `json:"archive_after_days" validate:"required,gte=2,lte=14"`
}

type ScanGroupEnv struct {
	Env    string
	Region string
}

type ScanGroupHandlers struct {
	validate         *validator.Validate
	env              *ScanGroupEnv
	scanGroupClient  am.ScanGroupService
	userClient       am.UserService
	orgClient        am.OrganizationService
	ContextExtractor middleware.UserContextExtractor
}

func New(scanGroupClient am.ScanGroupService, userClient am.UserService, orgClient am.OrganizationService, env *ScanGroupEnv) *ScanGroupHandlers {
	validate := validator.New()
	validate.RegisterValidation("subdomain", ValidateSubDomain)
	return &ScanGroupHandlers{
		scanGroupClient:  scanGroupClient,
		userClient:       userClient,
		orgClient:        orgClient,
		env:              env,
		ContextExtractor: middleware.ExtractUserContext,
		validate:         validate,
	}
}

type statsResponse struct {
	Status     string                 `json:"status"`
	GroupStats map[int]*am.GroupStats `json:"group_stats"`
}

func (h *ScanGroupHandlers) GetGroupStats(w http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	oid, stats, err := h.scanGroupClient.GroupStats(req.Context(), userContext)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group statistics")
		middleware.ReturnError(w, "internal error getting group statistics", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	data, err = json.Marshal(&statsResponse{Status: "OK", GroupStats: stats})
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling response")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func (h *ScanGroupHandlers) GetScanGroups(w http.ResponseWriter, req *http.Request) {
	log.Info().Msg("getting scan groups for user")
	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	oid, groups, err := h.scanGroupClient.Groups(req.Context(), userContext)
	if err != nil {
		logger.Error().Err(err).Msg("error getting groups for user")
		middleware.ReturnError(w, "error listing groups: "+err.Error(), 500)
		return
	}

	logger.Info().Msgf("groups: %#v", groups)

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Int("org_id", oid).Msg("authorization failure")
		middleware.ReturnError(w, "internal authorization error", 500)
		return
	}

	groupsForUser := make([]*serializers.ScanGroupForUser, len(groups))
	for i, g := range groups {
		groupsForUser[i] = &serializers.ScanGroupForUser{g}
	}

	data, _ := json.Marshal(groupsForUser)
	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func (h *ScanGroupHandlers) GetScanGroupByID(w http.ResponseWriter, req *http.Request) {
	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	param := chi.URLParam(req, "id")
	groupID, err := strconv.Atoi(param)
	if err != nil {
		logger.Error().Err(err).Msg("invalid group_id parameter")
		middleware.ReturnError(w, "invalid parameter", 403)
		return
	}

	oid, group, err := h.scanGroupClient.Get(req.Context(), userContext, groupID)
	if err != nil {
		logger.Error().Err(err).Msg("error getting group")
		middleware.ReturnError(w, "error getting group", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	groupForUser := &serializers.ScanGroupForUser{group}

	data, _ := json.Marshal(groupForUser)
	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func (h *ScanGroupHandlers) GetScanGroupByName(w http.ResponseWriter, req *http.Request) {
	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	param := chi.URLParam(req, "name")

	oid, group, err := h.scanGroupClient.GetByName(req.Context(), userContext, param)
	if err != nil {
		logger.Error().Err(err).Msg("error listing groups")
		middleware.ReturnError(w, "error listing groups", 400)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	groupForUser := &serializers.ScanGroupForUser{group}

	data, _ := json.Marshal(groupForUser)
	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

type groupCreated struct {
	Status           string `json:"status"`
	GroupID          int    `json:"group_id"`
	UploadAddressURI string `json:"upload_address_uri"`
}

func (h *ScanGroupHandlers) CreateScanGroup(w http.ResponseWriter, req *http.Request) {
	var err error
	var body []byte
	var gid int
	var orgPortScanEnabled bool

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	if middleware.AccountDisabled(userContext) {
		middleware.ReturnError(w, "user account disabled", 401)
		return
	}

	_, user, err := h.userClient.GetByCID(req.Context(), userContext, userContext.GetUserCID())
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user details")
		middleware.ReturnError(w, "unable to retrieve user details", 500)
		return
	}

	if user.AgreementAccepted == false {
		logger.Warn().Msg("user has not accepted agreement")
		middleware.ReturnError(w, "user has not accepted agreement, unable to create scan group.", 401)
		return
	}

	_, org, err := h.orgClient.GetByCID(req.Context(), userContext, userContext.GetOrgCID())
	if err != nil {
		logger.Warn().Msg("unable to retrieve organization for checking features")
	} else {
		orgPortScanEnabled = org.PortScanEnabled
	}

	param := chi.URLParam(req, "name")
	if strings.Contains(param, "/") {
		logger.Error().Msg("invalid character '/' in group name")
		middleware.ReturnError(w, "'/' is not allowed in the group name", 401)
		return
	}

	_, exists, _ := h.scanGroupClient.GetByName(req.Context(), userContext, param)
	if exists != nil {
		middleware.ReturnError(w, "group name already exists", 400)
		return
	}

	switch userContext.GetSubscriptionID() {
	case am.SubscriptionMonthlySmall:
		_, g, err := h.scanGroupClient.Groups(req.Context(), userContext)
		if err != nil {
			middleware.ReturnError(w, "error listing current scan groups", 400)
			return
		}
		if len(g) != 0 {
			middleware.ReturnError(w, "this pricing plan only allows one scan group", 400)
			return
		}
	case am.SubscriptionMonthlyMedium:
		_, g, err := h.scanGroupClient.Groups(req.Context(), userContext)
		if err != nil {
			middleware.ReturnError(w, "error listing current scan groups", 400)
			return
		}
		if len(g) >= 3 {
			middleware.ReturnError(w, "this pricing plan only allows three scan groups", 400)
			return
		}
	}

	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		middleware.ReturnError(w, "error reading scangroup from body", 400)
		return
	}
	defer req.Body.Close()

	groupDetails := &ScanGroupDetails{}
	if err := json.Unmarshal(body, groupDetails); err != nil {
		middleware.ReturnError(w, "error reading scangroup", 400)
		return
	}

	if err := h.validate.Struct(groupDetails); err != nil {
		logger.Error().Err(err).Msg("invalid data passed")
		middleware.ReturnError(w, err.Error(), 401) // TODO: don't expose internal errors
		return
	}

	now := time.Now().UnixNano()
	group := &am.ScanGroup{}
	group.GroupName = groupDetails.GroupName
	group.OrgID = userContext.GetOrgID()
	group.CreatedByID = userContext.GetUserID()
	group.CreationTime = now
	group.ModifiedByID = userContext.GetUserID()
	group.ModifiedTime = now
	group.OriginalInputS3URL = "s3://empty"
	group.Paused = true
	group.ArchiveAfterDays = groupDetails.ArchiveAfterDays

	portConfig := &am.PortScanModuleConfig{
		RequestsPerSecond: groupDetails.ConcurrentRequests,
		PortScanEnabled:   false,
		CustomWebPorts:    groupDetails.CustomWebPorts,
	}

	if orgPortScanEnabled && groupDetails.PortScanEnabled {
		portConfig, err = createPortConfig(groupDetails)
		if err != nil {
			middleware.ReturnError(w, err.Error(), 401)
			return
		}
	}

	group.ModuleConfigurations = &am.ModuleConfiguration{
		NSModule: &am.NSModuleConfig{
			RequestsPerSecond: groupDetails.ConcurrentRequests,
		},
		BruteModule: &am.BruteModuleConfig{
			CustomSubNames:    groupDetails.CustomSubNames,
			RequestsPerSecond: groupDetails.ConcurrentRequests,
			MaxDepth:          2,
		},
		PortModule: portConfig,
		WebModule: &am.WebModuleConfig{
			TakeScreenShots:       true,
			RequestsPerSecond:     groupDetails.ConcurrentRequests,
			MaxLinks:              10,
			ExtractJS:             true,
			FingerprintFrameworks: true,
		},
		KeywordModule: &am.KeywordModuleConfig{
			Keywords: []string{""},
		},
	}

	oid, gid, err := h.scanGroupClient.Create(req.Context(), userContext, group)
	if err != nil {
		logger.Error().Err(err).Msg("error creating scangroup")
		middleware.ReturnError(w, "error creating scangroup", 400)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	created := &groupCreated{
		Status:           "OK",
		GroupID:          gid,
		UploadAddressURI: fmt.Sprintf("/address/%d/initial", gid),
	}

	data, err := json.Marshal(created)
	if err != nil {
		middleware.ReturnError(w, "failed to create response", 500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func (h *ScanGroupHandlers) UpdateScanGroup(w http.ResponseWriter, req *http.Request) {
	var err error
	var body []byte
	var orgPortScanEnabled bool

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}

	if middleware.AccountDisabled(userContext) {
		middleware.ReturnError(w, "user account disabled", 401)
		return
	}

	logger := middleware.UserContextLogger(userContext)

	param := chi.URLParam(req, "name")

	oid, original, err := h.scanGroupClient.GetByName(req.Context(), userContext, param)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group by name")
		middleware.ReturnError(w, "group not found", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Int("org_id", oid).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	_, org, err := h.orgClient.GetByCID(req.Context(), userContext, userContext.GetOrgCID())
	if err != nil {
		logger.Warn().Msg("unable to retrieve organization for checking features")
	} else {
		orgPortScanEnabled = org.PortScanEnabled
	}

	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read body")
		middleware.ReturnError(w, "error reading group details", 400)
		return
	}
	defer req.Body.Close()

	updatedGroup := &ScanGroupDetails{}
	if err := json.Unmarshal(body, updatedGroup); err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal")
		middleware.ReturnError(w, "error reading group", 400)
		return
	}

	if err := h.validate.Struct(updatedGroup); err != nil {
		logger.Error().Err(err).Msg("invalid data passed")
		middleware.ReturnError(w, err.Error(), 401) // TODO: don't expose internal errors
		return
	}

	if strings.Contains(updatedGroup.GroupName, "/") {
		middleware.ReturnError(w, "'/' is not allowed in the group name", 401)
		return
	}

	original.GroupName = updatedGroup.GroupName
	original.ModifiedByID = userContext.GetUserID()
	original.ArchiveAfterDays = updatedGroup.ArchiveAfterDays

	portConfig := &am.PortScanModuleConfig{
		RequestsPerSecond: updatedGroup.ConcurrentRequests,
		PortScanEnabled:   false,
		CustomWebPorts:    updatedGroup.CustomWebPorts,
	}

	if orgPortScanEnabled && updatedGroup.PortScanEnabled {
		portConfig, err = createPortConfig(updatedGroup)
		if err != nil {
			middleware.ReturnError(w, err.Error(), 401)
			return
		}
		log.Info().Msgf("%#v", portConfig)
	}

	original.ModuleConfigurations = &am.ModuleConfiguration{
		NSModule: &am.NSModuleConfig{
			RequestsPerSecond: updatedGroup.ConcurrentRequests,
		},
		BruteModule: &am.BruteModuleConfig{
			CustomSubNames:    updatedGroup.CustomSubNames,
			RequestsPerSecond: updatedGroup.ConcurrentRequests,
			MaxDepth:          2,
		},
		PortModule: portConfig,
		WebModule: &am.WebModuleConfig{
			TakeScreenShots:       true,
			RequestsPerSecond:     updatedGroup.ConcurrentRequests,
			MaxLinks:              10,
			ExtractJS:             true,
			FingerprintFrameworks: true,
		},
		KeywordModule: &am.KeywordModuleConfig{
			Keywords: []string{""},
		},
	}

	_, _, err = h.scanGroupClient.Update(req.Context(), userContext, original)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update scangroup")
		middleware.ReturnError(w, "internal error updating scangroup", 500)
		return
	}

	middleware.ReturnSuccess(w, "group updated", 200)
}

func (h *ScanGroupHandlers) DeleteScanGroup(w http.ResponseWriter, req *http.Request) {
	var err error

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}
	logger := middleware.UserContextLogger(userContext)

	if userContext.GetSubscriptionID() == am.SubscriptionMonthlySmall {
		middleware.ReturnError(w, "this pricing plan does not allow for deleting scan groups", 400)
		return
	}

	param := chi.URLParam(req, "name")

	oid, group, err := h.scanGroupClient.GetByName(req.Context(), userContext, param)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group by name")
		middleware.ReturnError(w, "failure retrieving group", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	_, _, err = h.scanGroupClient.Delete(req.Context(), userContext, group.GroupID)
	if err != nil {
		logger.Error().Err(err).Msg("deletion failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	middleware.ReturnSuccess(w, "OK", 200)
}

type groupStatus struct {
	Status string `json:"status"`
}

func (h *ScanGroupHandlers) UpdateScanGroupStatus(w http.ResponseWriter, req *http.Request) {
	var err error
	var body []byte

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}

	if middleware.AccountDisabled(userContext) {
		middleware.ReturnError(w, "user account disabled", 401)
		return
	}

	logger := middleware.UserContextLogger(userContext)

	param := chi.URLParam(req, "name")
	logger.Info().Str("group_name", param).Msg("looking up group")
	oid, group, err := h.scanGroupClient.GetByName(req.Context(), userContext, param)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group by name")
		middleware.ReturnError(w, "group not found", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		logger.Error().Err(am.ErrOrgIDMismatch).Int("org_id", oid).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read body")
		middleware.ReturnError(w, "error reading group details", 400)
		return
	}
	defer req.Body.Close()

	status := &groupStatus{}
	if err := json.Unmarshal(body, status); err != nil {
		logger.Error().Err(err).Msg("failed to read status")
		middleware.ReturnError(w, "error reading status", 400)
		return
	}

	if status.Status == "pause" {
		_, _, err = h.scanGroupClient.Pause(req.Context(), userContext, group.GroupID)
	} else if status.Status == "resume" {
		_, _, err = h.scanGroupClient.Resume(req.Context(), userContext, group.GroupID)
	} else {
		middleware.ReturnError(w, "unknown status supplied, must be pause or resume", 400)
		return
	}

	if err != nil {
		logger.Error().Err(err).Msg("update failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	middleware.ReturnSuccess(w, "OK", 200)
}

func createPortConfig(details *ScanGroupDetails) (*am.PortScanModuleConfig, error) {
	pps := details.PortsPerSecond
	if pps == 0 {
		pps = 5
	}

	webPorts, tcpPorts, valid := VerifyWebScanPorts(details.CustomWebPorts, details.TCPPorts)
	if !valid {
		return nil, errors.New("not all webports exist in tcp ports")
	}

	if len(details.AllowedTLDs) == 0 && len(details.AllowedHosts) == 0 {
		return nil, errors.New("you have not specified any allowed hosts or TLDs")
	}

	portConfig := &am.PortScanModuleConfig{
		RequestsPerSecond: pps,
		PortScanEnabled:   details.PortScanEnabled,
		CustomWebPorts:    webPorts,
		TCPPorts:          tcpPorts,
		UDPPorts:          nil,
		AllowedTLDs:       details.AllowedTLDs,
		AllowedHosts:      details.AllowedHosts,
		DisallowedTLDs:    details.DisallowedTLDs,
		DisallowedHosts:   details.DisallowedHosts,
	}
	return portConfig, nil
}

// VerifyWebScanPorts a terribly inefficient method of validating webports are a subset of tcpports
// and removes duplicates
func VerifyWebScanPorts(webPorts, tcpPorts []int32) ([]int32, []int32, bool) {
	web := make(map[int32]struct{})
	for _, port := range webPorts {
		web[port] = struct{}{}
	}

	tcp := make(map[int32]struct{})
	for _, port := range tcpPorts {
		tcp[port] = struct{}{}
	}

	// manually add 80 and 443 if not in there
	if _, ok := tcp[80]; !ok {
		tcp[80] = struct{}{}
	}

	if _, ok := tcp[443]; !ok {
		tcp[443] = struct{}{}
	}

	// validate web ports are in tcp ports
	for port := range web {
		if _, ok := tcp[port]; !ok {
			return nil, nil, false
		}
	}

	w := make([]int32, len(web))
	i := 0
	for port := range web {
		w[i] = port
		i++
	}

	t := make([]int32, len(tcp))
	i = 0
	for port := range tcp {
		t[i] = port
		i++
	}

	return w, t, true
}
