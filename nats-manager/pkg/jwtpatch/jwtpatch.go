package jwtpatch

import (
	"time"

	"github.com/nats-io/jwt/v2"
)

type StringListPatches struct {
	// Clear always happens first
	Clear bool
	// Remove happens second
	Remove jwt.StringList
	// Add happens third
	Add jwt.StringList
}

func PatchStringList(s *jwt.StringList, with *StringListPatches) {
	if with.Clear {
		*s = make(jwt.StringList, 0)
	}
	s.Remove(with.Remove...)
	s.Add(with.Add...)
}

type PermissionPatches struct {
	Allow StringListPatches
	Deny  StringListPatches
}

func PatchPermission(p *jwt.Permission, with *PermissionPatches) {
	PatchStringList(&p.Allow, &with.Allow)
	PatchStringList(&p.Deny, &with.Deny)
}

type ResponsePermissionPatches struct {
	MaxMsgs *int
	Expires *time.Duration
}

func patchResponsePermission(p *jwt.ResponsePermission, with *ResponsePermissionPatches) {
	if with.MaxMsgs != nil {
		p.MaxMsgs = *with.MaxMsgs
	}
	if with.Expires != nil {
		p.Expires = *with.Expires
	}
}

type PermissionsPatches struct {
	Pub  *PermissionPatches
	Sub  *PermissionPatches
	Resp *ResponsePermissionPatches
}

func PatchPermissions(p *jwt.Permissions, with *PermissionsPatches) {
	if with.Pub != nil {
		PatchPermission(&p.Pub, with.Pub)
	}
	if with.Sub != nil {
		PatchPermission(&p.Sub, with.Sub)
	}
	if with.Resp != nil {
		patchResponsePermission(p.Resp, with.Resp)
	}
}

type CIDRListPatches struct {
	// Clear always happens first
	Clear bool
	// Remove happens second
	Remove jwt.CIDRList
	// Add happens third
	Add jwt.CIDRList
}

func PatchCIDRList(s *jwt.CIDRList, with *CIDRListPatches) {
	if with.Clear {
		*s = make(jwt.CIDRList, 0)
	}
	s.Remove(with.Remove...)
	s.Add(with.Add...)
}

type TimeRangePatches struct {
	// Clear always happens first
	Clear bool
	// Remove happens second
	Remove []jwt.TimeRange
	// Add happens third
	Add []jwt.TimeRange
}

func PatchTimeRangeList(s *[]jwt.TimeRange, with *TimeRangePatches) {
	if with.Clear {
		*s = make([]jwt.TimeRange, 0)
	}

	j := 0
	for i := range *s {
		for _, rm := range with.Remove {
			if (*s)[i] != rm {
				(*s)[j] = (*s)[i]
				j++
			}
		}
	}
	*s = (*s)[:j]

	for _, add := range with.Add {
		exists := false
		for _, it := range *s {
			if it == add {
				exists = true
			}
		}
		if !exists {
			*s = append(*s, add)
		}
	}
}

type UserLimitsPatches struct {
	src    CIDRListPatches
	times  TimeRangePatches
	locale *string
}

func PatchUserLimits(l *jwt.UserLimits, with *UserLimitsPatches) {
	PatchCIDRList(&l.Src, &with.src)
	PatchTimeRangeList(&l.Times, &with.times)
	if with.locale != nil {
		l.Locale = *with.locale
	}
}

type NATSLimitsPatches struct {
	Subs    *int64
	Data    *int64
	Payload *int64
}

func PatchNATSLimits(l *jwt.NatsLimits, with *NATSLimitsPatches) {
	if with.Subs != nil {
		l.Subs = *with.Subs
	}
	if with.Data != nil {
		l.Data = *with.Data
	}
	if with.Payload != nil {
		l.Payload = *with.Payload
	}
}

type LimitsPatches struct {
	UserLimitsPatches
	NATSLimitsPatches
}

func PatchLimits(l *jwt.Limits, with *LimitsPatches) {
	PatchUserLimits(&l.UserLimits, &with.UserLimitsPatches)
	PatchNATSLimits(&l.NatsLimits, &with.NATSLimitsPatches)
}

type TagListPatches struct {
	// Clear always happens first
	Clear bool
	// Remove happens second
	Remove jwt.TagList
	// Add happens third
	Add jwt.TagList
}

func PatchTagList(s *jwt.TagList, with *TagListPatches) {
	if with.Clear {
		*s = make(jwt.TagList, 0)
	}
	s.Remove(with.Remove...)
	s.Add(with.Add...)
}

type GenericFieldsPatches struct {
	Tags    TagListPatches
	Type    *jwt.ClaimType
	Version *int
}

func PatchGenericFields(f *jwt.GenericFields, with *GenericFieldsPatches) {
	PatchTagList(&f.Tags, &with.Tags)
	if with.Type != nil {
		f.Type = *with.Type
	}
	if with.Version != nil {
		f.Version = *with.Version
	}
}

type ClaimsDataPatches struct {
	Audience  *string
	Expires   *int64
	ID        *string
	IssuedAt  *int64
	Issuer    *string
	Name      *string
	NotBefore *int64
	Subject   *string
}

func PatchClaimsData(c *jwt.ClaimsData, with *ClaimsDataPatches) {
	if with.Audience != nil {
		c.Audience = *with.Audience
	}
	if with.Expires != nil {
		c.Expires = *with.Expires
	}
	if with.ID != nil {
		c.ID = *with.ID
	}
	if with.IssuedAt != nil {
		c.IssuedAt = *with.IssuedAt
	}
	if with.Issuer != nil {
		c.Issuer = *with.Issuer
	}
	if with.Name != nil {
		c.Name = *with.Name
	}
	if with.NotBefore != nil {
		c.NotBefore = *with.NotBefore
	}
	if with.Subject != nil {
		c.Subject = *with.Subject
	}
}

type UserPatches struct {
	PermissionsPatches
	LimitsPatches
	BearerToken            *bool
	AllowedConnectionTypes StringListPatches
	IssuerAccount          *string
	GenericFieldsPatches
}

func PatchUser(u *jwt.User, with *UserPatches) {
	PatchPermissions(&u.Permissions, &with.PermissionsPatches)
	PatchLimits(&u.Limits, &with.LimitsPatches)
	if with.BearerToken != nil {
		u.BearerToken = *with.BearerToken
	}
	PatchStringList(&u.AllowedConnectionTypes, &with.AllowedConnectionTypes)
	if with.IssuerAccount != nil {
		u.IssuerAccount = *with.IssuerAccount
	}
	PatchGenericFields(&u.GenericFields, &with.GenericFieldsPatches)
}

type UserClaimsPatches struct {
	ClaimsDataPatches
	UserPatches
}

func PatchUserClaims(c *jwt.UserClaims, with *UserClaimsPatches) {
	if with != nil {
		PatchClaimsData(&c.ClaimsData, &with.ClaimsDataPatches)
		PatchUser(&c.User, &with.UserPatches)
	}
}

type RevocationListPatches struct {
	// Clear always happens first
	Clear bool
	// Remove happens second
	Remove jwt.StringList
	// Add happens third
	Add map[string]time.Time
}

func PatchRevocationList(l *jwt.RevocationList, with *RevocationListPatches) {
	if with.Clear {
		*l = make(jwt.RevocationList)
	}
	for _, pubKey := range with.Remove {
		l.ClearRevocation(pubKey)
	}
	for pubKey, ts := range with.Add {
		l.Revoke(pubKey, ts)
	}
}

type ImportPatches struct {
	Name    *string
	Subject *jwt.Subject
	Account *string
	Token   *string
	To      *jwt.Subject
	Type    *jwt.ExportType
}

func PatchImport(i *jwt.Import, with *ImportPatches) {
	if with.Name != nil {
		i.Name = *with.Name
	}
	if with.Subject != nil {
		i.Subject = *with.Subject
	}
	if with.Account != nil {
		i.Account = *with.Account
	}
	if with.Token != nil {
		i.Token = *with.Token
	}
	if with.To != nil {
		i.To = *with.To
	}
	if with.Type != nil {
		i.Type = *with.Type
	}
}

type ImportsPatches struct {
	// Clear always happens first
	Clear bool
	// Edit happens second
	Edit map[string]ImportPatches
	// Remove happens third
	Remove jwt.StringList
	// Add happens fourth
	Add jwt.Imports
}

func PatchImports(i *jwt.Imports, with *ImportsPatches) {
	if with.Clear {
		*i = make(jwt.Imports, 0)
	}
	for name, patches := range with.Edit {
		patches := patches
		for _, imp := range *i {
			if imp.Name == name {
				PatchImport(imp, &patches)
			}
		}
	}
	for _, removeName := range with.Remove {
		j := 0
		for _, imp := range *i {
			if imp.Name != removeName {
				(*i)[j] = imp
				j++
			}
		}
		*i = (*i)[:j]
	}
	i.Add(with.Add...)
}

type ServiceLatencyPatches struct {
	Sampling *int
	Results  *jwt.Subject
}

func PatchServiceLatency(l *jwt.ServiceLatency, with *ServiceLatencyPatches) {
	if with.Sampling != nil {
		l.Sampling = *with.Sampling
	}
	if with.Results != nil {
		l.Results = *with.Results
	}
}

type ExportPatches struct {
	Name                 *string
	Subject              *jwt.Subject
	Type                 *jwt.ExportType
	TokenReq             *bool
	Revocations          RevocationListPatches
	ResponseType         *jwt.ResponseType
	ResponseThreshold    *time.Duration
	Latency              *ServiceLatencyPatches
	AccountTokenPosition *uint
}

func PatchExport(e *jwt.Export, with *ExportPatches) {
	if with.Name != nil {
		e.Name = *with.Name
	}
	if with.Subject != nil {
		e.Subject = *with.Subject
	}
	if with.Type != nil {
		e.Type = *with.Type
	}
	if with.TokenReq != nil {
		e.TokenReq = *with.TokenReq
	}
	PatchRevocationList(&e.Revocations, &with.Revocations)
	if with.ResponseType != nil {
		e.ResponseType = *with.ResponseType
	}
	if with.ResponseThreshold != nil {
		e.ResponseThreshold = *with.ResponseThreshold
	}
	if with.Latency != nil {
		if e.Latency == nil {
			e.Latency = new(jwt.ServiceLatency)
		}
		PatchServiceLatency(e.Latency, with.Latency)
	}
	if with.AccountTokenPosition != nil {
		e.AccountTokenPosition = *with.AccountTokenPosition
	}
}

type ExportsPatches struct {
	// Clear always happens first
	Clear bool
	// Edit happens second
	Edit map[string]ExportPatches
	// Remove happens third
	Remove jwt.StringList
	// Add happens fourth
	Add jwt.Exports
}

func PatchExports(e *jwt.Exports, with *ExportsPatches) {
	if with.Clear {
		*e = make(jwt.Exports, 0)
	}
	for name, patches := range with.Edit {
		patches := patches
		for _, exp := range *e {
			if exp.Name == name {
				PatchExport(exp, &patches)
			}
		}
	}
	for _, removeName := range with.Remove {
		j := 0
		for _, exp := range *e {
			if exp.Name != removeName {
				(*e)[j] = exp
				j++
			}
		}
		*e = (*e)[:j]
	}
	e.Add(with.Add...)
}

type AccountLimitsPatches struct {
	Imports         *int64
	Exports         *int64
	WildcardExports *bool
	Conn            *int64
	LeafNodeConn    *int64
}

func PatchAccountLimits(l *jwt.AccountLimits, with *AccountLimitsPatches) {
	if with.Imports != nil {
		l.Imports = *with.Imports
	}
	if with.Exports != nil {
		l.Exports = *with.Exports
	}
	if with.WildcardExports != nil {
		l.WildcardExports = *with.WildcardExports
	}
	if with.Conn != nil {
		l.Conn = *with.Conn
	}
	if with.LeafNodeConn != nil {
		l.LeafNodeConn = *with.LeafNodeConn
	}
}

type JetStreamLimitsPatches struct {
	MemoryStorage *int64
	DiskStorage   *int64
	Streams       *int64
	Consumer      *int64
}

func PatchJetStreamLimits(l *jwt.JetStreamLimits, with *JetStreamLimitsPatches) {
	if with.MemoryStorage != nil {
		l.MemoryStorage = *with.MemoryStorage
	}
	if with.DiskStorage != nil {
		l.DiskStorage = *with.DiskStorage
	}
	if with.Streams != nil {
		l.Streams = *with.Streams
	}
	if with.Consumer != nil {
		l.Consumer = *with.Consumer
	}
}

type OperatorLimitsPatches struct {
	NATSLimitsPatches
	AccountLimitsPatches
	JetStreamLimitsPatches
}

func PatchOperatorLimits(l *jwt.OperatorLimits, with *OperatorLimitsPatches) {
	PatchNATSLimits(&l.NatsLimits, &with.NATSLimitsPatches)
	PatchAccountLimits(&l.AccountLimits, &with.AccountLimitsPatches)
	PatchJetStreamLimits(&l.JetStreamLimits, &with.JetStreamLimitsPatches)
}

type AccountPatches struct {
	Imports            ImportsPatches
	Exports            ExportsPatches
	Limits             OperatorLimitsPatches
	SigningKeys        StringListPatches
	Revocations        RevocationListPatches
	DefaultPermissions PermissionsPatches
	GenericFieldsPatches
}

func PatchAccount(a *jwt.Account, with *AccountPatches) {
	PatchImports(&a.Imports, &with.Imports)
	PatchExports(&a.Exports, &with.Exports)
	PatchOperatorLimits(&a.Limits, &with.Limits)
	PatchStringList(&a.SigningKeys, &with.SigningKeys)
	PatchRevocationList(&a.Revocations, &with.Revocations)
	PatchPermissions(&a.DefaultPermissions, &with.DefaultPermissions)
	PatchGenericFields(&a.GenericFields, &with.GenericFieldsPatches)
}

type AccountClaimsPatches struct {
	ClaimsDataPatches
	AccountPatches
}

func PatchAccountClaims(c *jwt.AccountClaims, with *AccountClaimsPatches) {
	if with != nil {
		PatchClaimsData(&c.ClaimsData, &with.ClaimsDataPatches)
		PatchAccount(&c.Account, &with.AccountPatches)
	}
}

type OperatorPatches struct {
	SigningKeys         StringListPatches
	AccountServerURL    *string
	OperatorServiceURLs StringListPatches
	SystemAccount       *string
	AssertServerVersion *string
	GenericFieldsPatches
}

func PatchOperator(o *jwt.Operator, with *OperatorPatches) {
	PatchStringList(&o.SigningKeys, &with.SigningKeys)
	if with.AccountServerURL != nil {
		o.AccountServerURL = *with.AccountServerURL
	}
	PatchStringList(&o.OperatorServiceURLs, &with.OperatorServiceURLs)
	if with.SystemAccount != nil {
		o.SystemAccount = *with.SystemAccount
	}
	if with.AssertServerVersion != nil {
		o.AccountServerURL = *with.AssertServerVersion
	}
	PatchGenericFields(&o.GenericFields, &with.GenericFieldsPatches)
}

type OperatorClaimsPatches struct {
	ClaimsDataPatches
	OperatorPatches
}

func PatchOperatorClaims(c *jwt.OperatorClaims, with *OperatorClaimsPatches) {
	if with != nil {
		PatchClaimsData(&c.ClaimsData, &with.ClaimsDataPatches)
		PatchOperator(&c.Operator, &with.OperatorPatches)
	}
}
