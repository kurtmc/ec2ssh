package search

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func ListRegions() []*ec2.Region {
	sess := session.Must(session.NewSession())

	svc := ec2.New(sess, &aws.Config{Region: aws.String("ap-southeast-2")})

	resp, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})

	if err != nil {
		panic(err)
	}

	return resp.Regions
}

func ListInstances(name string) ([]string, error) {
	sess := session.Must(session.NewSession())

	// Get all regions
	in := make(chan *ec2.Region, 100)
	regions := ListRegions()
	for _, r := range regions {
		in <- r
	}
	close(in)

	out := make(chan string, 100)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for region := range in {
				svc := ec2.New(sess, &aws.Config{Region: region.RegionName})
				resp, err := svc.DescribeInstances(params)

				if err != nil {
					panic(err)
				}

				for _, r := range resp.Reservations {
					for _, t := range r.Instances[0].Tags {
						if *r.Instances[0].State.Name == ec2.InstanceStateNameRunning && *t.Key == "Name" {
							out <- *t.Value
						}
					}
				}
			}
		}(&wg)
	}

	wg.Wait()
	close(out)

	var instances []string
	for instance := range out {
		instances = append(instances, instance)
	}

	return instances, nil

}
