import * as policy from "@pulumi/policy";
import { isTaggable } from "../autotag-ts/taggable";

new policy.PolicyPack("aws-tags-policies", {
    policies: [{
        name: "check-required-tags",
        description: "Ensure required tags are present on all taggable resources.",
        configSchema: {
            properties: {
                requiredTags: {
                    type: "array",
                    items: { type: "string" },
                },
            },
        },
        validateResource: (args, reportViolation) => {
            const config = args.getConfig<AwsTagsPolicyConfig>();
            const requiredTags = config.requiredTags;
            if (requiredTags && isTaggable(args.type)) {
                const ts = args.props["tags"];
                for (const rt of requiredTags) {
                    if (!ts || !ts[rt]) {
                        reportViolation(`Taggable resource '${args.urn}' is missing required tag '${rt}'`);
                    }
                }
            }
        },
    }],
});

interface AwsTagsPolicyConfig {
    requiredTags?: string[];
}
