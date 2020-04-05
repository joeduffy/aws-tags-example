import * as aws from "@pulumi/aws";
import * as pulumi from "@pulumi/pulumi";
import { registerAutoTags } from "./autotag";

// Automatically inject tags.
const config = new pulumi.Config();
registerAutoTags({
    "user:Project": pulumi.getProject(),
    "user:Stack": pulumi.getStack(),
    "user:Cost Center": config.require("costCenter"),
});

// Create a bunch of AWS resources -- with auto-tags!

const bucket = new aws.s3.Bucket("my-bucket");

const group = new aws.ec2.SecurityGroup("web-secgrp", {
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
    ],
});

const server = new aws.ec2.Instance("web-server-www", {
    instanceType: "t2.micro",
    ami: "ami-0c55b159cbfafe1f0",
    vpcSecurityGroupIds: [ group.id ],
});
