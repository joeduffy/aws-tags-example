package myproject;

import com.pulumi.Pulumi;

import java.util.List;
import java.util.Map;

import com.pulumi.Context;
import com.pulumi.aws.ec2.InstanceArgs;
import com.pulumi.aws.ec2.SecurityGroupArgs;
import com.pulumi.aws.ec2.inputs.SecurityGroupIngressArgs;
import com.pulumi.aws.inputs.ProviderDefaultTagsArgs;
import com.pulumi.resources.StackOptions;
import com.pulumi.resources.CustomResourceOptions;
import com.pulumi.aws.s3.Bucket;
import com.pulumi.aws.ec2.SecurityGroup;
import com.pulumi.aws.Provider;
import com.pulumi.aws.ProviderArgs;
import com.pulumi.aws.ec2.Instance;

public class App {
    public static void main(String[] args) {
        Pulumi.withOptions(StackOptions.builder()
                //.resourceTransformations(App::stackTransformation)
                .build())
                .run(App::stack);
    }

    private static void stack(Context ctx) {

        var provider = new Provider("Provider", ProviderArgs.builder()
            .defaultTags(ProviderDefaultTagsArgs.builder()
                .tags(Map.of(
                    "user:Project", ctx.projectName(),
                    "user:Stack", ctx.stackName(),
                    "user:Cost Center", ctx.config().require("costCenter")
                )).build()
            )
            .build()
        );

        // Create a bunch of AWS resources -- with auto-tags!
        var bucket = new Bucket("my-bucket", null, CustomResourceOptions.builder().provider(provider).build());

        var group = new SecurityGroup("web-secgrp", SecurityGroupArgs.builder()
                .ingress(SecurityGroupIngressArgs.builder()
                        .protocol("tcp")
                        .fromPort(22)
                        .toPort(22)
                        .cidrBlocks("0.0.0.0/0")
                        .build())
                .ingress(SecurityGroupIngressArgs.builder()
                        .protocol("tcp")
                        .fromPort(80)
                        .toPort(80)
                        .cidrBlocks("0.0.0.0/0")
                        .build())
                .build(),
                CustomResourceOptions.builder().provider(provider).build()
            );

        var server = new Instance("web-server-www", InstanceArgs.builder()
                .instanceType("t2.micro")
                .ami("ami-0a763bef4f952ec08")
                .vpcSecurityGroupIds(group.id().applyValue(List::of))
                .build(),
                CustomResourceOptions.builder().provider(provider).build()
            );

    }

}
