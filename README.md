# terraform-provider-path

This is a quick-and-dirty terraform provider which demonstrates an error path problem I've been unable to figure out.

## Configuration schema

The example configuration below is expected to produce a validation error because the value `3` is supplied when `must_be_even` is asserted.
```hcl
1  data "path_example" "e" {
2    number_list = [
3      {
4        number = "3"
5        must_be_even = true
6      }
7    ]
8  }
```

## The Validation Error

The validation error is implemented with `AddAttributeError()`, which includes a `path.Path` argument.

#### Close:

When creating the error, I can point the `path.Path` argument argument at the specific list entry like this:

```go
path.Root("number_list").AtListIndex(i)
```

The error cites "line 3" (the beginning of the list item with the problem attribute value):
```text
╷
│ Error: invalid attribute combination
│ 
│   with data.path_example.e,
│   on main.tf line 3, in data "path_example" "e":
│   3:     {
│   4:       number = "3"
│   5:       must_be_even = true
│   6:     }
│ 
│ number must be even when must_be_even is true
╵
```
#### Big Miss:

But I want to cite the `number` attribute in the error, so I extend the `path.Path` with an additional `path.PathStep` to the `number` attribute:

```go
path.Root("number_list").AtListIndex(i).AtName("number"),
```

Instead of pointing more specifically at the problem (line 4 - the `number` attribute), the error now cites *the whole resource block* (line 1):

```text
╷
│ Error: invalid attribute combination
│ 
│   with data.path_example.e,
│   on main.tf line 1, in data "path_example" "e":
│   1: data "path_example" "e" {
│ 
│ number must be even when must_be_even is true
╵
```

## Reproduction

* Clone this repo and build the provider
* run the configuration in the `test` directory
