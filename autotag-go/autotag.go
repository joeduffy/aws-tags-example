package main

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// registerAutoTags registers a global stack transformation that merges a set
// of tags with whatever was also explicitly added to the resource definition.
func registerAutoTags(ctx *pulumi.Context, autoTags map[string]string) {
	ctx.RegisterStackTransformation(
		func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			if isTaggable(args.Type) {
				ptr := reflect.ValueOf(args.Props)
				val := ptr.Elem()
				tags := val.FieldByName("Tags")

				var tagsMap pulumi.StringMap
				if !tags.IsZero() {
					tagsMap = tags.Interface().(pulumi.StringMap)
				} else {
					tagsMap = pulumi.StringMap(map[string]pulumi.StringInput{})
				}
				for k, v := range autoTags {
					tagsMap[k] = pulumi.String(v)
				}
				tags.Set(reflect.ValueOf(tagsMap))

				return &pulumi.ResourceTransformationResult{
					Props: args.Props,
					Opts:  args.Opts,
				}
			}
			return nil
		},
	)
}
