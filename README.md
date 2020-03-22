# A Powerful Argument Parser for Go

## Struct Tags Syntax
```
type <struct name> struct {
    Field1    <field type>                                               // 'argparser' struct tag not given hence ignored
    Field2    <field type>    `argparser:"<key=value>,<key=value>,..."   // 'argparser' struct tag given hence parsed
    Field3    <field type>    `argparser:"<key=value>,<key=value>,..."   // 'argparser' struct tag given hence parsed
    ...
}
```
**PS:** The fields must be public otherwise the `reflect` package will fail to parse the struct.

## Valid Tag keys and Values

| Key | Mandatory | Value Type (Go) | Possible Values | Default | Description |
| :---: | :---: | --- | :---: | :---: | :--- |
| `pos` | no | string | exact string `yes` | N/A | create a positional argument if given otherwise create an optional argument |
| `name` | yes | string | a valid string containing alphanumeric charaters and/or '-' | N/A |the name to identify the argument with |
| `nargs` | no | int | a valid int | 1 | number of values required by the argument, if set to 0 then the argument will be treated as an optional `switch` argument |
| `help` | no | string | any valid string, escape `,` as `\\,`  | "" | help message for the user |

## Example

For full examples please refer to `examples/`.
