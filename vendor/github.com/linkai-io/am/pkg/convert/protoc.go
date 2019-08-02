package convert

import (
	"github.com/linkai-io/am/am"
	"github.com/linkai-io/am/protocservices/prototypes"
)

// DomainToUser convert domain user type to protobuf user type
func DomainToUser(in *am.User) *prototypes.User {
	return &prototypes.User{
		OrgID:                      int32(in.OrgID),
		OrgCID:                     in.OrgCID,
		UserCID:                    in.UserCID,
		UserID:                     int32(in.UserID),
		UserEmail:                  in.UserEmail,
		FirstName:                  in.FirstName,
		LastName:                   in.LastName,
		StatusID:                   int32(in.StatusID),
		CreationTime:               in.CreationTime,
		Deleted:                    in.Deleted,
		AgreementAccepted:          in.AgreementAccepted,
		AgreementAcceptedTimestamp: in.AgreementAcceptedTimestamp,
		LastLoginTimestamp:         in.LastLoginTimestamp,
	}
}

// UserToDomain convert protobuf user type to domain user type
func UserToDomain(in *prototypes.User) *am.User {
	return &am.User{
		OrgID:                      int(in.OrgID),
		OrgCID:                     in.OrgCID,
		UserCID:                    in.UserCID,
		UserID:                     int(in.UserID),
		UserEmail:                  in.UserEmail,
		FirstName:                  in.FirstName,
		LastName:                   in.LastName,
		StatusID:                   int(in.StatusID),
		CreationTime:               in.CreationTime,
		Deleted:                    in.Deleted,
		AgreementAccepted:          in.AgreementAccepted,
		AgreementAcceptedTimestamp: in.AgreementAcceptedTimestamp,
		LastLoginTimestamp:         in.LastLoginTimestamp,
	}
}

func DomainToUserFilter(in *am.UserFilter) *prototypes.UserFilter {
	return &prototypes.UserFilter{
		Start:   int32(in.Start),
		Limit:   int32(in.Limit),
		OrgID:   int32(in.OrgID),
		Filters: DomainToFilterTypes(in.Filters),
	}
}

func UserFilterToDomain(in *prototypes.UserFilter) *am.UserFilter {
	return &am.UserFilter{
		Start:   int(in.Start),
		Limit:   int(in.Limit),
		OrgID:   int(in.OrgID),
		Filters: FilterTypesToDomain(in.Filters),
	}
}

// UserContextToDomain converts from a protoc usercontext to an am.usercontext
func UserContextToDomain(in *prototypes.UserContext) am.UserContext {
	return &am.UserContextData{
		TraceID:        in.TraceID,
		OrgID:          int(in.OrgID),
		OrgCID:         in.OrgCID,
		UserID:         int(in.UserID),
		UserCID:        in.UserCID,
		Roles:          in.Roles,
		IPAddress:      in.IPAddress,
		SubscriptionID: in.SubscriptionID,
		OrgStatusID:    int(in.OrgStatusID),
	}
}

// DomainToUserContext converts the domain usercontext to protobuf usercontext
func DomainToUserContext(in am.UserContext) *prototypes.UserContext {
	return &prototypes.UserContext{
		TraceID:        in.GetTraceID(),
		OrgID:          int32(in.GetOrgID()),
		OrgCID:         in.GetOrgCID(),
		UserCID:        in.GetUserCID(),
		UserID:         int32(in.GetUserID()),
		Roles:          in.GetRoles(),
		IPAddress:      in.GetIPAddress(),
		SubscriptionID: in.GetSubscriptionID(),
		OrgStatusID:    int32(in.GetOrgStatusID()),
	}
}

// DomainToOrganization converts the domain organization to protobuf organization
func DomainToOrganization(in *am.Organization) *prototypes.Org {
	return &prototypes.Org{
		OrgID:                      int32(in.OrgID),
		OrgCID:                     in.OrgCID,
		OrgName:                    in.OrgName,
		OwnerEmail:                 in.OwnerEmail,
		UserPoolID:                 in.UserPoolID,
		UserPoolAppClientID:        in.UserPoolAppClientID,
		UserPoolAppClientSecret:    in.UserPoolAppClientSecret,
		IdentityPoolID:             in.IdentityPoolID,
		FirstName:                  in.FirstName,
		LastName:                   in.LastName,
		Phone:                      in.Phone,
		Country:                    in.Country,
		StatePrefecture:            in.StatePrefecture,
		Street:                     in.Street,
		Address1:                   in.Address1,
		Address2:                   in.Address2,
		City:                       in.City,
		PostalCode:                 in.PostalCode,
		CreationTime:               in.CreationTime,
		StatusID:                   int32(in.StatusID),
		Deleted:                    in.Deleted,
		SubscriptionID:             in.SubscriptionID,
		UserPoolJWK:                in.UserPoolJWK,
		LimitTLD:                   in.LimitTLD,
		LimitTLDReached:            in.LimitTLDReached,
		LimitHosts:                 in.LimitHosts,
		LimitHostsReached:          in.LimitHostsReached,
		LimitCustomWebFlows:        in.LimitCustomWebFlows,
		LimitCustomWebFlowsReached: in.LimitCustomWebFlowsReached,
		PortScanEnabled:            in.PortScanEnabled,
	}
}

// OrganizationToDomain converts the protobuf organization to domain organization
func OrganizationToDomain(in *prototypes.Org) *am.Organization {
	return &am.Organization{
		OrgID:                      int(in.OrgID),
		OrgCID:                     in.OrgCID,
		OrgName:                    in.OrgName,
		OwnerEmail:                 in.OwnerEmail,
		UserPoolID:                 in.UserPoolID,
		UserPoolAppClientID:        in.UserPoolAppClientID,
		UserPoolAppClientSecret:    in.UserPoolAppClientSecret,
		IdentityPoolID:             in.IdentityPoolID,
		UserPoolJWK:                in.UserPoolJWK,
		FirstName:                  in.FirstName,
		LastName:                   in.LastName,
		Phone:                      in.Phone,
		Country:                    in.Country,
		StatePrefecture:            in.StatePrefecture,
		Street:                     in.Street,
		Address1:                   in.Address1,
		Address2:                   in.Address2,
		City:                       in.City,
		PostalCode:                 in.PostalCode,
		CreationTime:               in.CreationTime,
		StatusID:                   int(in.StatusID),
		Deleted:                    in.Deleted,
		SubscriptionID:             in.SubscriptionID,
		LimitTLD:                   in.LimitTLD,
		LimitTLDReached:            in.LimitTLDReached,
		LimitHosts:                 in.LimitHosts,
		LimitHostsReached:          in.LimitHostsReached,
		LimitCustomWebFlows:        in.LimitCustomWebFlows,
		LimitCustomWebFlowsReached: in.LimitCustomWebFlowsReached,
		PortScanEnabled:            in.PortScanEnabled,
	}
}

func DomainToOrgFilter(in *am.OrgFilter) *prototypes.OrgFilter {
	return &prototypes.OrgFilter{
		Start:   int32(in.Start),
		Limit:   int32(in.Limit),
		Filters: DomainToFilterTypes(in.Filters),
	}
}

// OrgFilterToDomain convert org filter protobuf to orgfilter domain
func OrgFilterToDomain(in *prototypes.OrgFilter) *am.OrgFilter {
	return &am.OrgFilter{
		Start:   int(in.Start),
		Limit:   int(in.Limit),
		Filters: FilterTypesToDomain(in.Filters),
	}
}

func PortResultsToDomain(in *prototypes.PortResults) *am.PortResults {
	if in == nil {
		return &am.PortResults{}
	}
	ports := &am.Ports{Current: &am.PortData{}, Previous: &am.PortData{}}
	if in != nil && in.Ports != nil && in.Ports.Current != nil {
		ports.Current = &am.PortData{
			IPAddress:  in.Ports.Current.IPAddress,
			TCPPorts:   in.Ports.Current.TCPPorts,
			UDPPorts:   in.Ports.Current.UDPPorts,
			TCPBanners: in.Ports.Current.TCPBanners,
			UDPBanners: in.Ports.Current.UDPBanners,
		}
	}

	if in != nil && in.Ports != nil && in.Ports.Previous != nil {
		ports.Previous = &am.PortData{
			IPAddress:  in.Ports.Previous.IPAddress,
			TCPPorts:   in.Ports.Previous.TCPPorts,
			UDPPorts:   in.Ports.Previous.UDPPorts,
			TCPBanners: in.Ports.Previous.TCPBanners,
			UDPBanners: in.Ports.Previous.UDPBanners,
		}
	}

	return &am.PortResults{
		PortID:                   in.PortID,
		OrgID:                    int(in.OrgID),
		GroupID:                  int(in.GroupID),
		HostAddress:              in.HostAddress,
		Ports:                    ports,
		ScannedTimestamp:         in.ScannedTimestamp,
		PreviousScannedTimestamp: in.PreviousScannedTimestamp,
	}
}

func DomainToPortResults(in *am.PortResults) *prototypes.PortResults {
	if in == nil {
		return &prototypes.PortResults{}
	}

	ports := &prototypes.Ports{Current: &prototypes.PortData{}, Previous: &prototypes.PortData{}}
	if in != nil && in.Ports != nil && in.Ports.Current != nil {
		ports.Current = &prototypes.PortData{
			IPAddress:  in.Ports.Current.IPAddress,
			TCPPorts:   in.Ports.Current.TCPPorts,
			UDPPorts:   in.Ports.Current.UDPPorts,
			TCPBanners: in.Ports.Current.TCPBanners,
			UDPBanners: in.Ports.Current.UDPBanners,
		}
	}

	if in != nil && in.Ports != nil && in.Ports.Previous != nil {
		ports.Previous = &prototypes.PortData{
			IPAddress:  in.Ports.Previous.IPAddress,
			TCPPorts:   in.Ports.Previous.TCPPorts,
			UDPPorts:   in.Ports.Previous.UDPPorts,
			TCPBanners: in.Ports.Previous.TCPBanners,
			UDPBanners: in.Ports.Previous.UDPBanners,
		}
	}

	return &prototypes.PortResults{
		PortID:                   in.PortID,
		OrgID:                    int32(in.OrgID),
		GroupID:                  int32(in.GroupID),
		HostAddress:              in.HostAddress,
		Ports:                    ports,
		ScannedTimestamp:         in.ScannedTimestamp,
		PreviousScannedTimestamp: in.PreviousScannedTimestamp,
	}
}

func AddressToDomain(in *prototypes.AddressData) *am.ScanGroupAddress {
	if in == nil {
		return nil
	}
	return &am.ScanGroupAddress{
		AddressID:           in.AddressID,
		OrgID:               int(in.OrgID),
		GroupID:             int(in.GroupID),
		HostAddress:         in.HostAddress,
		IPAddress:           in.IPAddress,
		DiscoveryTime:       in.DiscoveryTime,
		DiscoveredBy:        in.DiscoveredBy,
		LastScannedTime:     in.LastScannedTime,
		LastSeenTime:        in.LastSeenTime,
		ConfidenceScore:     in.ConfidenceScore,
		UserConfidenceScore: in.UserConfidenceScore,
		IsSOA:               in.IsSOA,
		IsWildcardZone:      in.IsWildcardZone,
		IsHostedService:     in.IsHostedService,
		Ignored:             in.Ignored,
		FoundFrom:           in.FoundFrom,
		NSRecord:            in.NSRecord,
		AddressHash:         in.AddressHash,
		Deleted:             in.Deleted,
	}
}

func DomainToAddress(in *am.ScanGroupAddress) *prototypes.AddressData {
	if in == nil {
		return nil
	}

	return &prototypes.AddressData{
		OrgID:               int32(in.OrgID),
		AddressID:           in.AddressID,
		GroupID:             int32(in.GroupID),
		HostAddress:         in.HostAddress,
		IPAddress:           in.IPAddress,
		DiscoveryTime:       in.DiscoveryTime,
		DiscoveredBy:        in.DiscoveredBy,
		LastScannedTime:     in.LastScannedTime,
		LastSeenTime:        in.LastSeenTime,
		ConfidenceScore:     in.ConfidenceScore,
		UserConfidenceScore: in.UserConfidenceScore,
		IsSOA:               in.IsSOA,
		IsWildcardZone:      in.IsWildcardZone,
		IsHostedService:     in.IsHostedService,
		Ignored:             in.Ignored,
		Deleted:             in.Deleted,
		FoundFrom:           in.FoundFrom,
		NSRecord:            in.NSRecord,
		AddressHash:         in.AddressHash,
	}
}

func HostListToDomain(in *prototypes.HostListData) *am.ScanGroupHostList {
	return &am.ScanGroupHostList{
		AddressIDs:  in.AddressIDs,
		OrgID:       int(in.OrgID),
		GroupID:     int(in.GroupID),
		ETLD:        in.ETLD,
		HostAddress: in.HostAddress,
		IPAddresses: in.IPAddresses,
	}
}

func DomainToHostList(in *am.ScanGroupHostList) *prototypes.HostListData {
	return &prototypes.HostListData{
		AddressIDs:  in.AddressIDs,
		OrgID:       int32(in.OrgID),
		GroupID:     int32(in.GroupID),
		ETLD:        in.ETLD,
		HostAddress: in.HostAddress,
		IPAddresses: in.IPAddresses,
	}
}

func AddressFilterToDomain(in *prototypes.AddressFilter) *am.ScanGroupAddressFilter {
	return &am.ScanGroupAddressFilter{
		OrgID:   int(in.OrgID),
		GroupID: int(in.GroupID),
		Start:   in.Start,
		Limit:   int(in.Limit),
		Filters: FilterTypesToDomain(in.Filters),
	}
}

func DomainToAddressFilter(in *am.ScanGroupAddressFilter) *prototypes.AddressFilter {
	return &prototypes.AddressFilter{
		OrgID:   int32(in.OrgID),
		GroupID: int32(in.GroupID),
		Start:   in.Start,
		Limit:   int32(in.Limit),
		Filters: DomainToFilterTypes(in.Filters),
	}
}

// ModuleToDomain converts protoc ModuleConfiguration to am.ModuleConfiguration
func ModuleToDomain(in *prototypes.ModuleConfiguration) *am.ModuleConfiguration {
	return &am.ModuleConfiguration{
		NSModule: &am.NSModuleConfig{
			RequestsPerSecond: in.NSConfig.RequestsPerSecond,
		},
		BruteModule: &am.BruteModuleConfig{
			RequestsPerSecond: in.BruteConfig.RequestsPerSecond,
			CustomSubNames:    in.BruteConfig.CustomSubNames,
			MaxDepth:          in.BruteConfig.MaxDepth,
		},
		PortModule: &am.PortScanModuleConfig{
			RequestsPerSecond: in.PortConfig.RequestsPerSecond,
			PortScanEnabled:   in.PortConfig.PortScanEnabled,
			CustomWebPorts:    in.PortConfig.CustomWebPorts,
			TCPPorts:          in.PortConfig.TCPPorts,
			UDPPorts:          in.PortConfig.UDPPorts,
			AllowedTLDs:       in.PortConfig.AllowedTLDs,
			AllowedHosts:      in.PortConfig.AllowedHosts,
			DisallowedTLDs:    in.PortConfig.DisallowedTLDs,
			DisallowedHosts:   in.PortConfig.DisallowedHosts,
		},
		WebModule: &am.WebModuleConfig{
			RequestsPerSecond:     in.WebModuleConfig.RequestsPerSecond,
			TakeScreenShots:       in.WebModuleConfig.TakeScreenShots,
			MaxLinks:              in.WebModuleConfig.MaxLinks,
			ExtractJS:             in.WebModuleConfig.ExtractJS,
			FingerprintFrameworks: in.WebModuleConfig.FingerprintFrameworks,
		},
		KeywordModule: &am.KeywordModuleConfig{
			Keywords: in.KeywordModuleConfig.Keywords,
		},
	}
}

func DomainToModule(in *am.ModuleConfiguration) *prototypes.ModuleConfiguration {
	return &prototypes.ModuleConfiguration{
		NSConfig: &prototypes.NSModuleConfig{
			RequestsPerSecond: in.NSModule.RequestsPerSecond,
		},
		BruteConfig: &prototypes.BruteModuleConfig{
			RequestsPerSecond: in.BruteModule.RequestsPerSecond,
			CustomSubNames:    in.BruteModule.CustomSubNames,
			MaxDepth:          in.BruteModule.MaxDepth,
		},
		PortConfig: &prototypes.PortModuleConfig{
			RequestsPerSecond: in.PortModule.RequestsPerSecond,
			PortScanEnabled:   in.PortModule.PortScanEnabled,
			CustomWebPorts:    in.PortModule.CustomWebPorts,
			TCPPorts:          in.PortModule.TCPPorts,
			UDPPorts:          in.PortModule.UDPPorts,
			AllowedTLDs:       in.PortModule.AllowedTLDs,
			AllowedHosts:      in.PortModule.AllowedHosts,
			DisallowedTLDs:    in.PortModule.DisallowedTLDs,
			DisallowedHosts:   in.PortModule.DisallowedHosts,
		},
		WebModuleConfig: &prototypes.WebModuleConfig{
			RequestsPerSecond:     in.WebModule.RequestsPerSecond,
			TakeScreenShots:       in.WebModule.TakeScreenShots,
			MaxLinks:              in.WebModule.MaxLinks,
			ExtractJS:             in.WebModule.ExtractJS,
			FingerprintFrameworks: in.WebModule.FingerprintFrameworks,
		},
		KeywordModuleConfig: &prototypes.KeywordModuleConfig{
			Keywords: in.KeywordModule.Keywords,
		},
	}
}

// ScanGroupToDomain convert protoc group to domain type ScanGroup
func ScanGroupToDomain(in *prototypes.Group) *am.ScanGroup {
	return &am.ScanGroup{
		OrgID:                int(in.OrgID),
		GroupID:              int(in.GroupID),
		GroupName:            in.GroupName,
		CreationTime:         in.CreationTime,
		CreatedBy:            in.CreatedBy,
		CreatedByID:          int(in.CreatedByID),
		OriginalInputS3URL:   in.OriginalInputS3URL,
		ModifiedBy:           in.ModifiedBy,
		ModifiedByID:         int(in.ModifiedByID),
		ModifiedTime:         in.ModifiedTime,
		ModuleConfigurations: ModuleToDomain(in.ModuleConfiguration),
		Paused:               in.Paused,
		Deleted:              in.Deleted,
		LastPausedTime:       in.LastPausedTime,
		ArchiveAfterDays:     in.ArchiveAfterDays,
	}
}

// DomainToScanGroup convert domain type SdcanGroup to protoc Group
func DomainToScanGroup(in *am.ScanGroup) *prototypes.Group {
	return &prototypes.Group{
		OrgID:               int32(in.OrgID),
		GroupID:             int32(in.GroupID),
		GroupName:           in.GroupName,
		CreationTime:        in.CreationTime,
		CreatedBy:           in.CreatedBy,
		CreatedByID:         int32(in.CreatedByID),
		OriginalInputS3URL:  in.OriginalInputS3URL,
		ModifiedBy:          in.ModifiedBy,
		ModifiedByID:        int32(in.ModifiedByID),
		ModifiedTime:        in.ModifiedTime,
		ModuleConfiguration: DomainToModule(in.ModuleConfigurations),
		Paused:              in.Paused,
		Deleted:             in.Deleted,
		LastPausedTime:      in.LastPausedTime,
		ArchiveAfterDays:    in.ArchiveAfterDays,
	}
}

func DomainToScanGroupFilter(in *am.ScanGroupFilter) *prototypes.ScanGroupFilter {
	return &prototypes.ScanGroupFilter{
		Filters: DomainToFilterTypes(in.Filters),
	}
}

func ScanGroupFilterToDomain(in *prototypes.ScanGroupFilter) *am.ScanGroupFilter {
	return &am.ScanGroupFilter{
		Filters: FilterTypesToDomain(in.Filters),
	}
}

func DomainToGroupStats(in *am.GroupStats) *prototypes.GroupStats {
	return &prototypes.GroupStats{
		OrgID:           int32(in.OrgID),
		GroupID:         int32(in.GroupID),
		ActiveAddresses: in.ActiveAddresses,
		BatchSize:       in.BatchSize,
		LastUpdated:     in.LastUpdated,
		BatchStart:      in.BatchStart,
		BatchEnd:        in.BatchEnd,
	}
}

func DomainToGroupsStats(in map[int]*am.GroupStats) map[int32]*prototypes.GroupStats {
	stats := make(map[int32]*prototypes.GroupStats, len(in))
	for groupID, stat := range in {
		stats[int32(groupID)] = DomainToGroupStats(stat)
	}
	return stats
}

func GroupStatsToDomain(in *prototypes.GroupStats) *am.GroupStats {
	return &am.GroupStats{
		OrgID:           int(in.OrgID),
		GroupID:         int(in.GroupID),
		ActiveAddresses: in.ActiveAddresses,
		BatchSize:       in.BatchSize,
		LastUpdated:     in.LastUpdated,
		BatchStart:      in.BatchStart,
		BatchEnd:        in.BatchEnd,
	}
}

func GroupsStatsToDomain(in map[int32]*prototypes.GroupStats) map[int]*am.GroupStats {
	stats := make(map[int]*am.GroupStats, len(in))
	for groupID, stat := range in {
		stats[int(groupID)] = GroupStatsToDomain(stat)
	}

	return stats
}

func DomainToScanGroupAggregates(in map[string]*am.ScanGroupAggregates) map[string]*prototypes.ScanGroupAggregates {
	if in == nil {
		return nil
	}
	agg := make(map[string]*prototypes.ScanGroupAggregates, len(in))
	for k, v := range in {
		agg[k] = &prototypes.ScanGroupAggregates{Time: v.Time, Count: v.Count}
	}
	return agg
}

func DomainToScanGroupAddressStats(in *am.ScanGroupAddressStats) *prototypes.ScanGroupAddressStats {
	return &prototypes.ScanGroupAddressStats{
		OrgID:             int32(in.OrgID),
		GroupID:           int32(in.GroupID),
		DiscoveredBy:      in.DiscoveredBy,
		DiscoveredByCount: in.DiscoveredByCount,
		Aggregates:        DomainToScanGroupAggregates(in.Aggregates),
		Total:             in.Total,
		ConfidentTotal:    in.ConfidentTotal,
	}
}

func DomainToScanGroupsAddressStats(in []*am.ScanGroupAddressStats) []*prototypes.ScanGroupAddressStats {
	stats := make([]*prototypes.ScanGroupAddressStats, 0)
	for _, v := range in {
		stats = append(stats, DomainToScanGroupAddressStats(v))
	}
	return stats
}

func ScanGroupAggregatesToDomain(in map[string]*prototypes.ScanGroupAggregates) map[string]*am.ScanGroupAggregates {
	if in == nil {
		return nil
	}
	agg := make(map[string]*am.ScanGroupAggregates, len(in))
	for k, v := range in {
		agg[k] = &am.ScanGroupAggregates{Time: v.Time, Count: v.Count}
	}
	return agg
}

func ScanGroupAddressStatsToDomain(in *prototypes.ScanGroupAddressStats) *am.ScanGroupAddressStats {
	return &am.ScanGroupAddressStats{
		OrgID:             int(in.OrgID),
		GroupID:           int(in.GroupID),
		DiscoveredBy:      in.DiscoveredBy,
		DiscoveredByCount: in.DiscoveredByCount,
		Aggregates:        ScanGroupAggregatesToDomain(in.Aggregates),
		Total:             in.Total,
		ConfidentTotal:    in.ConfidentTotal,
	}
}

func ScanGroupsAddressStatsToDomain(in []*prototypes.ScanGroupAddressStats) []*am.ScanGroupAddressStats {
	stats := make([]*am.ScanGroupAddressStats, 0)
	for _, v := range in {
		stats = append(stats, ScanGroupAddressStatsToDomain(v))
	}
	return stats
}
