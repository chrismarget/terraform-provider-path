package tfpath

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithValidateConfig = &dataSourceExample{}

type dataSourceExample struct{}

type example struct {
	NumberList types.List `tfsdk:"number_list"`
}

type exampleNumber struct {
	Number     types.Int64 `tfsdk:"number"`
	MustBeEven types.Bool  `tfsdk:"must_be_even"`
}

func (o *dataSourceExample) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (o *dataSourceExample) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"number_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"number": schema.Int64Attribute{
						Required: true,
					},
					"must_be_even": schema.BoolAttribute{
						Required: true,
					},
				}},
				Required: true,
			},
		},
	}
}

func (o *dataSourceExample) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config example
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var numbers []exampleNumber
	resp.Diagnostics.Append(config.NumberList.ElementsAs(ctx, &numbers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, number := range numbers {
		if number.Number.ValueInt64()%2 != 0 && number.MustBeEven.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				// The following path directive produces an error which
				// identifies the line number (3) of the *list item* where the
				// problem value was encountered, but it doesn't point directly
				// at the problem value. The error looks like this:
				// ╷
				// │ Error: invalid attribute combination
				// │
				// │   with data.path_example.e,
				// │   on main.tf line 3, in data "path_example" "e":
				// │    3:     {
				// │    4:       number = "3"
				// │    5:       must_be_even = true
				// │    6:     }
				// │
				// │ number must be even when must_be_even is true
				// ╵
				// Pretty good, but I want the error to highlight line 4, rather
				// than line 3.
				//
				path.Root("number_list").AtListIndex(i),

				// Extending the `path.Path` with `AtName("number")` does not
				// produce the expected error message. Rather than highlighting
				// the `number` attribute on line 4, the error refers the user
				// to the entire resource block (line 1).
				// Why doesn't it point to the `number` attribute?
				// ╷
				// │ Error: invalid attribute combination
				// │
				// │   with data.path_example.e,
				// │   on main.tf line 1, in data "path_example" "e":
				// │    1: data "path_example" "e" {
				// │
				// │ number must be even when must_be_even is true
				// ╵
				//path.Root("number_list").AtListIndex(i).AtName("number"),

				"invalid attribute combination",
				"number must be even when must_be_even is true",
			)
		}
	}
}

func (o *dataSourceExample) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
}
