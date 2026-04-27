package cloudfront

type CloudFront struct {
	Distributions []Distribution `tree:"distributions"`
	Functions     []Function     `tree:"functions"`
}
