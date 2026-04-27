package elasticbeanstalk

type ElasticBeanstalk struct {
	Environments []Environment `tree:"environments"`
}
