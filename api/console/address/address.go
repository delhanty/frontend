package address

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/linkai-io/am/pkg/convert"
	"github.com/linkai-io/am/pkg/inputlist"

	"github.com/go-chi/chi"
	"github.com/linkai-io/frontend/pkg/middleware"
	"github.com/rs/zerolog/log"

	"github.com/linkai-io/am/am"
)

type AddressResponse struct {
	Addrs     []*am.ScanGroupAddress `json:"addresses"`
	Status    string                 `json:"status"`
	LastIndex int64                  `json:"last_index"`
}

type AddressHandlers struct {
	addrClient       am.AddressService
	scanGroupClient  am.ScanGroupService
	ContextExtractor middleware.UserContextExtractor
}

func New(addrClient am.AddressService, scanGroupClient am.ScanGroupService) *AddressHandlers {
	return &AddressHandlers{
		addrClient:       addrClient,
		scanGroupClient:  scanGroupClient,
		ContextExtractor: middleware.ExtractUserContext,
	}
}

func (h *AddressHandlers) GetAddresses(w http.ResponseWriter, req *http.Request) {
	var err error

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}

	groupID, err := groupIDFromRequest(req)
	if err != nil {
		middleware.ReturnError(w, "invalid scangroup id supplied", 401)
		return
	}

	filter, err := h.ParseGetFilterQuery(req.URL.Query(), userContext.GetOrgID(), groupID)
	if err != nil {
		log.Error().Err(err).Int("OrgID", userContext.GetOrgID()).Int("GroupID", groupID).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("failed parse url query parameters")
		middleware.ReturnError(w, "invalid parameters supplied", 401)
		return
	}

	oid, addrs, err := h.addrClient.Get(req.Context(), userContext, filter)
	if err != nil {
		log.Error().Err(err).Int("OrgID", userContext.GetOrgID()).Int("GroupID", groupID).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("failed to get addresses")
		middleware.ReturnError(w, "failed to get addresses: "+err.Error(), 500)
		return
	}

	var lastAddr int64
	for _, addr := range addrs {
		if addr.AddressID > lastAddr {
			lastAddr = addr.AddressID
		}
		if oid != addr.OrgID {
			log.Error().Err(err).Int("OrgID", userContext.GetOrgID()).Int("GroupID", groupID).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("authorization failure")
			middleware.ReturnError(w, "failed to get addresses", 500)
			return
		}
	}
	response := &AddressResponse{
		Status:    "ok",
		LastIndex: lastAddr,
		Addrs:     addrs,
	}

	data, err := json.Marshal(response)
	if err != nil {
		middleware.ReturnError(w, "failed return addresses", 500)
		return
	}

	fmt.Fprintf(w, string(data))
}

func (h *AddressHandlers) ParseGetFilterQuery(values url.Values, orgID, groupID int) (*am.ScanGroupAddressFilter, error) {
	var err error
	filter := &am.ScanGroupAddressFilter{
		OrgID:               orgID,
		GroupID:             groupID,
		WithIgnored:         false,
		IgnoredValue:        false,
		WithLastScannedTime: false,
		SinceScannedTime:    0,
		WithLastSeenTime:    false,
		SinceSeenTime:       0,
		Start:               0,
		Limit:               0,
	}

	ignored := values.Get("ignored")
	if ignored == "true" {
		filter.WithIgnored = true
		filter.IgnoredValue = true
	} else if ignored == "false" {
		filter.WithIgnored = true
		filter.IgnoredValue = false
	}

	sinceScanned := values.Get("since_scanned")
	if sinceScanned != "" {
		filter.WithLastScannedTime = true
		filter.SinceScannedTime, err = strconv.ParseInt(sinceScanned, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	sinceSeen := values.Get("since_seen")
	if sinceSeen != "" {
		filter.WithLastSeenTime = true
		filter.SinceSeenTime, err = strconv.ParseInt(sinceSeen, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	start := values.Get("start")
	if start == "" {
		filter.Start = 0
	} else {
		filter.Start, err = strconv.ParseInt(start, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	limit := values.Get("limit")
	if limit == "" {
		filter.Limit = 0
	} else {
		filter.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		if filter.Limit > 1000 {
			return nil, errors.New("limit max size exceeded (1000)")
		}
	}
	return filter, nil
}

type PutResponse struct {
	Status       string                  `json:"status"`
	ParserErrors []*inputlist.ParseError `json:"errors,omitempty"`
	Count        int                     `json:"count,omitempty"`
}

func (h *AddressHandlers) PutInitialAddresses(w http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	putResponse := &PutResponse{}

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}

	groupID, err := groupIDFromRequest(req)
	if err != nil {
		middleware.ReturnError(w, "invalid scangroup id supplied", 401)
		return
	}

	addrs, parserErrors := inputlist.ParseList(req.Body, 100000)
	if len(parserErrors) != 0 {
		putResponse.ParserErrors = parserErrors
		putResponse.Status = "NG"
		data, err = json.Marshal(putResponse)
		if err != nil {
			middleware.ReturnError(w, "internal marshal error", 500)
			return
		}
		w.WriteHeader(400)
		fmt.Fprint(w, string(data))
		return
	}

	sgAddrs := makeAddrs(addrs, userContext.GetOrgID(), userContext.GetUserID(), groupID)

	oid, count, err := h.addrClient.Update(req.Context(), userContext, sgAddrs)
	if err != nil {
		log.Error().Err(err).Int("OrgID", userContext.GetOrgID()).Int("GroupID", groupID).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("failed to add addresses")
		middleware.ReturnError(w, "failed to add addresses to scangroup", 500)
		return
	}

	if oid != userContext.GetOrgID() {
		log.Error().Err(am.ErrOrgIDMismatch).Int("OrgID", userContext.GetOrgID()).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	putResponse.Count = count
	putResponse.Status = "OK"

	data, err = json.Marshal(putResponse)
	if err != nil {
		middleware.ReturnError(w, "internal marshal error", 500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func makeAddrs(in map[string]struct{}, orgID, userID, groupID int) map[string]*am.ScanGroupAddress {
	addrs := make(map[string]*am.ScanGroupAddress, len(in))
	i := 0
	for addr := range in {
		sgAddr := &am.ScanGroupAddress{
			OrgID:               orgID,
			GroupID:             groupID,
			DiscoveredBy:        "input_list",
			DiscoveryTime:       time.Now().UnixNano(),
			ConfidenceScore:     100.0,
			UserConfidenceScore: 0.0,
		}

		if inputlist.IsIP(addr) {
			sgAddr.IPAddress = addr
		} else {
			sgAddr.HostAddress = addr
		}
		sgAddr.AddressHash = convert.HashAddress(sgAddr.IPAddress, sgAddr.HostAddress)
		addrs[sgAddr.AddressHash] = sgAddr
		i++
	}
	return addrs
}

type countResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (h *AddressHandlers) GetGroupCount(w http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	userContext, ok := h.ContextExtractor(req.Context())
	if !ok {
		middleware.ReturnError(w, "missing user context", 401)
		return
	}

	id, err := groupIDFromRequest(req)
	if err != nil {
		middleware.ReturnError(w, "invalid scangroup id supplied", 401)
		return
	}

	oid, count, err := h.addrClient.Count(req.Context(), userContext, id)
	if oid != userContext.GetOrgID() {
		log.Error().Err(am.ErrOrgIDMismatch).Int("OrgID", userContext.GetOrgID()).Int("UserID", userContext.GetUserID()).Str("TraceID", userContext.GetTraceID()).Msg("authorization failure")
		middleware.ReturnError(w, "internal error", 500)
		return
	}

	data, err = json.Marshal(&countResponse{Status: "OK", Count: count})
	if err != nil {
		middleware.ReturnError(w, "internal marshal error", 500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, string(data))
}

func groupIDFromRequest(req *http.Request) (int, error) {
	param := chi.URLParam(req, "id")
	id, err := strconv.Atoi(param)
	return id, err
}