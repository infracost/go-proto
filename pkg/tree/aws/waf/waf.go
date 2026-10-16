package waf

type WAF struct {
	WebACLs                     []WebACL                     `tree:"web_acls"`
	RegionalWebACLs             []WebACL                     `tree:"regional_web_acls"`
	WebACLv2s                   []WebACLv2                   `tree:"web_acl_v2s"`
	WebACLLoggingConfigurations []WebACLLoggingConfiguration `tree:"web_acl_logging_configurations"`
}
