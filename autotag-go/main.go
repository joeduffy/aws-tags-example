package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		registerAutoTags(ctx, map[string]string{
			"user:Project":     ctx.Project(),
			"user:Stack":       ctx.Stack(),
			"user:Cost Center": config.Require(ctx, "gopol:costCenter"),
		})

		// Create a bunch of AWS resources -- with auto-tags!

		_, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{})
		if err != nil {
			return err
		}

		grp, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewInstance(ctx, "web-server-www", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			Ami:                 pulumi.String("ami-0a763bef4f952ec08"),
			VpcSecurityGroupIds: pulumi.StringArray{grp.ID()},
		})
		return err
	})
}
