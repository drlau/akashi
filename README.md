# Akashi

Akashi is a Go tool that can be used to parse `terraform plan` outputs and validate the changes.

Still a WIP project. Supports `terraform` v0.12 currently.

- [x] Created and destroyed resources
- [x] JSON Output
- [x] stdout output
- [x] Map / Array rules
- [ ] Updated resources
- [ ] Index matching
- [ ] Module matching
- [ ] Multiple rule matching
- [ ] Other validations(regex, int in range, etc)
- [ ] Customizable output

## Installation

```bash
go get -u github.com/drlau/akashi
```

## Usage

```
Usage:
  akashi <path to ruleset> [flags]

Flags:
  -e, --error-on-fail   for non-quiet runs, make akashi return exit code 1 on fails
      --failed-only     only output failing lines
  -f, --file string     read plan output from file
  -h, --help            help for akashi
  -j, --json            read the contents as the output from 'terraform state show -json'
      --no-color        disable color output
  -q, --quiet           compare only, and error if there is a failing rule
  -s, --strict          require all resources to match a comparer
```

By default, `akashi` will read a `terraform plan` output from `stdin`, so you should pipe the result of `terraform plan`:

```bash
terraform plan | akashi <path to ruleset>
```

If you produced a plan output with `terraform plan -out=<file>`, `akashi` can parse it by specifying `--json`. Note that this file is encoded and needs to be decoded by `terraform` before `akashi` can parse it properly:

```bash
terraform state show -json <file> | akashi <path to ruleset> --json
```

If the `terraform plan` output or the decoded json is in a file, you can read directly from the file by specifying the path with `-f`.

## Ruleset schema

**NOTE**: Ruleset schema is in the early stages and is subject to change in later versions.

Rulesets are written in YAML and have the following schema:

```yaml
# Rules to apply to created resources.
createdResources:
  # Set to true if you want all created resources to match a rule.
  # Default is false.
  strict: true

  # Default compare options to apply to all resources.
  # If a resource specifies the same option, the resource's value will be used.
  default:
    # Set to true if you want every enforcedValue to be present in a resource's plan.
    # If set to true, if an enforcedValue is missing, it will count as a failure.
    # Default is false.
    enforceAll: true

    # Set to true if you want to ignore extra arguments in a resource's plan
    # that don't match any enforcedValue or ignoredArg.
    # If set to false, any extra arguments will count as a failure.
    # Default is false.
    ignoreExtraArgs: true

    # Set to true if you want to ignore computed arguments in a resource's plan.
    # Default is false.
    ignoreComputed: true

    # Set to true if you want every enforcedValue and ignoredArg to be present in a resource's plan.
    # If set to true, if any enforcedValue or ignoredArg is missing, it will count as a failure.
    # Default is false.
    requireAll: true

    # Set to true if you want to automatically fail before comparison
    # if a matching resource is found.
    # Default is false.
    autoFail: true

  # List of rules.
  resources:
  - resource:
      # Resource name to match on.
      # At least one of "name" or "type" must be set.
      # Defaults is empty.
      name: resource-name

      # Resource type to match on.
      # At least one of "name" or "type" must be set.
      # Defaults is empty.
      type: resource-type

      # The same compare options from "default" can be specified per resource.
      # The resource level option will take priority over the option specified in "default"
      # If omitted, the option specified in "default" is used.
      # If not specified in "default", the value is false.
      enforceAll: true
      ignoreExtraArgs: true
      ignoreComputed: true
      requireAll: true
      autoFail: true

      # List of arguments to ignore.
      # Default is empty.
      ignored:
        - ignored-arg-1
        - ignored-arg-2

      # List of arguments to enforce.
      # Default is empty.
      enforced:
        stringEnforced: string
        intEnforced: 1
        boolEnforced: true
        mapEnforced:
          mapKey: mapValue
        arrayEnforced:
          - array1
          - array2

# Rules to apply to destroyed resources.
# Has the exact same schema as createdResources.
destroyedResources:
```

### Example

Say you provision `google_compute_instance` and you want to validate that all new instances are created in zone `us-central1-a`, and you don't care about any other argument. To validate that, you would create the following ruleset:

```yaml
createdResources:
  default:
    enforceAll: true
    ignoreExtraArgs: true
  resources:
    - type: google_compute_instance
      enforced:
        zone: us-central1-a
```