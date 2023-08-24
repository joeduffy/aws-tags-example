import pulumi
import pulumi_aws as aws
from autotag import register_auto_tags

# Automatically inject tags.
config = pulumi.Config()
register_auto_tags({
    'user:Project': pulumi.get_project(),
    'user:Stack': pulumi.get_stack(),
    'user:Cost Center': config.require('costCenter'),
})

# Create a bunch of AWS resources -- with auto-tags!

bucket = aws.s3.Bucket('my-bucket')

group = aws.ec2.SecurityGroup('web-secgrp',
    ingress=[
        { 'protocol': 'tcp', 'from_port': 22, 'to_port': 22, 'cidr_blocks': ['0.0.0.0/0']},
        { 'protocol': 'tcp', 'from_port': 80, 'to_port': 80, 'cidr_blocks': ['0.0.0.0/0']},
    ],
)

server = aws.ec2.Instance('web-server-www',
    instance_type='t2.micro',
    ami='ami-0a763bef4f952ec08',
    vpc_security_group_ids=[ group.id ],
)
