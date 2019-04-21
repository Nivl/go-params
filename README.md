# Params options

[![Actions Status](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/badge/Nivl/go-params)](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/results/Nivl/go-params)
[![Go Report Card](https://goreportcard.com/badge/github.com/nivl/go-params)](https://goreportcard.com/report/github.com/nivl/go-params)
[![codecov](https://codecov.io/gh/Nivl/go-params/branch/master/graph/badge.svg)](https://codecov.io/gh/Nivl/go-params)
[![GoDoc](https://godoc.org/github.com/Nivl/go-params?status.svg)](https://godoc.org/github.com/Nivl/go-params)

## Source types (`from:""`)

- `url`: The param is part of the URL as in `/item/your-param`.
- `query`: The param is part of the query string as in `item?id=your-param`.
- `form`: The param is part of the body of the request. It can be from a JSON payload or a basic form-urlencoded payload.
- `file`: The param is a file sent using `multipart/form-data`. The param type MUST be a `*formfile.FormFile`.

## Params type (`params:""`)

### global options (works on most of the types)

- `required`: The field is required and an error will be returned if the field is empty (or contains less than 1 element for an array).

### String specific params

- `trim`: The value will be trimmed of its trailing spaces.
- `uuid`: The value is required to be a valid UUIDv4.
- `url`: The value is required to be a valid http(s) url.
- `email`: The value is required to be a valid email.
- `slug`: The value is required to be a valid slug.
- `slugOrUuid`: The value is required to be either a valid slug or a valid UUIDv4.

**If used on an array, those params will be applied on each values of the array**

### Pointers specific params

- `noempty`: The value cannot be empty. The pointer can be nil, but if a value is provided it cannot be an empty string. the difference with `required` is that `required` does not accept nil pointer (works on array too).

### Files specific params

- `image`: The provided file is required to be an image.

### Array specific params

- `no_empty_items`: If the array contains an empty item, an error will be thrown.

## Ignoring and naming `json:""`

- Use `json:"_"` to prevent a field to be altered or checked.
- Use `json:"field_name"` to name a field.

## Default value

Use `default:"my_value"` to set a default value. The default value will be used
if nothing is found in the payload or if the provided value is an empty string.

For an array, you can use a comma (`,`) to separate a default value of multiple
data: `default:"1,2,3"` on a `[]string` type will output `[]string{1,2,3}`.

## Enum values

You can set a comma separated list of valid values using
`enum:"value1,value2,value3"`. Empty and nil values are accepted.

## Min/Max values for integers

You can set a min value or a max value for an integer using
`min_int:"0" max_int:"10"`.

**If used on an array, those params will be applied on each values of the array**

## Maxlen of a string

Use `maxlen:"255"` to make sure the len of a string is not bigger than 255 char. Any invalid values (including `0`) will be ignored.

**If used on an array, those params will be applied on each values of the array**

## Custom Validator

You can add a custom validator by implementing `params.CustomValidation`.

## Examples

```golang
type UpdateParams struct {
  ID        string  `from:"url" json:"id" params:"required,uuid"`
  Name      *string `from:"form" json:"name" params:"trim,noempty"`
  ShortName *string `from:"form" json:"short_name" params:"trim"`
  Website   *string `from:"form" json:"website" params:"url,trim"`
  InTrash   *bool   `from:"form" json:"in_trash"`
  Items     []int   `from:"form" json:"items" max_items="3" params:"trim,empty"`
}
```
