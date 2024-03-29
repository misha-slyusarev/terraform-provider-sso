# terraform-provider-sso

I experimented with writing a Terraform provider for SSO configuration. Originally I wanted to reduce the amount of logic defined in the locals section. The locals section is sometimes used to perform imperative data manipulations under the guise of defining a variable which is an ugly way of getting around the declarative nature of HCL.

This is an excellent example of what I'm talking about https://github.com/prodapt-cloud/TerraformRepo/blob/main/AWS/ApplicationIntegration/terraform-sso/modules/permission-sets/main.tf#L44. Several variables are defined just to allow running several for loops to shape the data.

But it turned out Terraform SDKv2 doesn't allow passing information between a plugin and Terraform core
in the form of a map of objects ([for example](https://github.com/hashicorp/terraform-plugin-sdk/issues/616)) , so I cannot use the output of a plugin in the `for_each` meta argument of a resource, and I still had to do some mapping in the locals section (https://github.com/misha-slyusarev/terraform-provider-sso/blob/main/main.tf#L70).

I also wanted to try to implment the "based on" functionality when one permission set could be based on another one. The code currently has a naive version implemntation of such functionality.

Not everything worked out, but it was fun to play around with building my own Terraform provider.

## Development
You will need to set up the development overrides configuration to be able to test locally. Create `.terraformrc` in your HOME folder and use the following example. Replace the override path with the path to wherever you cloned the repository.
```
provider_installation {
	dev_overrides {
		"mysso" = "/Users/you/terraform-provider-sso/"
	}

	direct {}
}
```

Then build this project, creating the plugin binary. *Note that the naming convention is important.*
```
make build
```

You can finally run `terraform`. You don't have to initialize the folder, just use plan.
```
TF_LOG=INFO terraform plan
```
You should be able to see Terraform plans to add several new resources. Check `main.tf` for usage examples.
