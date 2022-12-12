# terraform-provider-sso
I experimented with writing a Terraform provider for SSO configuration. Originally I wanted to reduce the amount of logic
defined in the locals section as normally it would be used to prepare a number of maps that will be used later to create resources.

But it turned out Terraform SDKv2 doesn't allow passing information between a plugin and Terraform core
in a form of a map, so the output of a plugin cannot be used in for_each meta argument of a resource and I still had to so some mapping
in the locals section.

Here is an example of such mapping https://github.com/misha-slyusarev/terraform-provider-sso/blob/main/main.tf#L70

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
