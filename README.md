# argparser
A Powerful Argument Parser for Go

## Struct Tags Syntax
```
`argparser:"<key=value>|<key=value>|..."`
```

## Valid Tag keys and Values

| Tag | Mandatory | Go Type | Possible Values | Default | Description |
| :---: | :---: | --- | :---: | :---: | :--- |
| `pos` | no | string | yes | N/A | create a positional argument if pos=yes is given otherwise create an optional argument |
| `name` | yes | string | a valid string | "" |the name to identify the argument with |
| `nargs` | no | int | a valid int | 1 | number of values required by the argument, if set to 0 then the argument will be treated an optional `switch` argument |
| `help` | no | string | a valid string | "" | help message for the user |

## Example

For more details please refer to `examples/`.
