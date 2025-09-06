# lumberyard_dirmap (Data Source)

Traverses a nested directory of YAML and JSON files and constructs a nested object.

## Example Usage

```hcl
data "lumberyard_dirmap" "example" {
  path = "./environments"
}

output "dirmap" {
  value = data.lumberyard_dirmap.example.result
}
```

## Schema

### Required

- `path` (String) The base directory to traverse.

### Optional

- `filter` (String) A glob pattern to filter files.

### Read-Only

- `result` (Map of String) The constructed object from the directory structure.
