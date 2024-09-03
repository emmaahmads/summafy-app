package api

type awsConfig struct {
	s3_bucket string
	region    string
	creds     []string
}

func NewAwsConfig(s3_bucket string, region string, creds ...string) *awsConfig {
	return &awsConfig{
		s3_bucket: s3_bucket,
		region:    region,
		creds:     creds,
	}
}
