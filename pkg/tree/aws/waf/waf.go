package waf

type WAF struct {
	WebACLs   []WebACL   `tree:"web_acls"`
	WebACLv2s []WebACLv2 `tree:"web_acl_v2s"`
}
