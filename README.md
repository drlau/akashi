# Akashi / 証

Akashi is a Go tool that can be used to parse `terraform plan` outputs and validate the changes.

Still a WIP project. Supports `terraform` v0.12 currently.

- [x] Created and destroyed resources
- [x] JSON Output
- [x] stdout output
- [x] Map / Array rules
- [x] Updated resources
- [ ] Index matching
- [ ] Module matching
- [ ] Multiple rule matching
- [ ] Other validations(regex, int in range, etc)
- [ ] Combining multiple rulesets
- [ ] Customizable output

## Why

By parsing the `terraform plan` result, you can validate changes shown in the `plan` output before actually applying. Also, it can be used to track unexpected changes during routine operations that may have occured from misconfiguration or infrastructure drift.

For example, say you had a `google_container_node_pool` with `autoscaling` configured, and due to some incident you had to modify the `autoscaling` configuration outside of `terraform`. Some time later, you update the `version` of your `google_container_node_pool` as a routine operation. Due to the declarative nature of `terraform`, the previous configuration of `autoscaling` would be applied and overwrite the modified configuration set during the incident. This might be unintentional and would not be caught from parsing the `terraform` configuration, but can be caught by validating the `terraform plan` output.

## Installation

```bash
go get -u github.com/drlau/akashi
```

## Usage

```
Usage:
  akashi <compare | diff> <path to ruleset> [flags]
```

By default, `akashi` will read a `terraform plan` output from `stdin`, so you should pipe the result of `terraform plan`:

```bash
terraform plan | akashi <compare | diff> <path to ruleset>
```

If you produced a plan output with `terraform plan -out=<file>`, `akashi` can parse it by specifying `--json`. Note that this file is encoded and needs to be decoded by `terraform` before `akashi` can parse it properly:

```bash
terraform state show -json <file> | akashi <compare | diff> <path to ruleset> --json
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

  # Set to true if you want Akashi to validate the ruleset specifies names for
  # created resources.
  # Default is false.
  requireName: true

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
    -
      # Resource name to match on.
      # At least one of "name" or "type" must be set.
      # Default is empty.
      name: resource-name

      # Resource type to match on.
      # At least one of "name" or "type" must be set.
      # Default is empty.
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
        stringEnforced:
          value: string
        intEnforced:
          value: 1
        boolEnforced:
          value: true
        mapEnforced:
          value:
            mapKey: mapValue
        arrayEnforced:
          value:
          - array1
          - array2
        stringMatchAny:
          matchAny:
          - validValue1
          - validValue2

# Rules to apply to destroyed resources.
# Has the exact same schema as createdResources.
destroyedResources:

# Rules to apply to updated resources.
updatedResources:
  # Set to true if you want all updated resources to match a rule.
  # Default is false.
  strict: true

  # Set to true if you want Akashi to validate the ruleset specifies names for
  # updated resources.
  # Default is false.
  requireName: true

  # Default compare options to apply to all resources.
  # If a resource specifies the same option, the resource's value will be used.
  # All options for created and destroyed resources work here, but also has a few additional options that can be enabled
  default:
    # Set to true if you want to ignore all unchanged attributes
    # Default is false.
    ignoreNoOp: true

  # List of rules.
  resources:
  -
    # Similar to created and destroyed resources, matching can be done on a name and/or type of resource
    name: resource-name
    type: resource-type

    # Every compare option can also be specified at resource level, overriding the top level default
    ignoreNoOp: true

    # Rules to enforce on the attributes before the planned changes
    # Consists of ignored and enforced, with the same behaviour as created and destroyed resources
    before:
      # List of arguments to ignore.
      # Default is empty.
      ignored:
        - ignored-arg-1
        - ignored-arg-2

      # List of arguments to enforce.
      # Default is empty.
      enforced:
        stringEnforced:
          value: string

    # Rules to enforce on the attributes after the planned changes
    # Same schema as before.
    after:
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
        zone:
          value: us-central1-a
```
