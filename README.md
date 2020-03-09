# argparser
A Powerful Argument Parser for Go


```
type:   opt/pos
name:   "name without punctuation etc."
short:  "string with one charater"
long:   "override long option name string"
nargs:  number or special characters line +/*
mutex: string label
usage: usage string

type cmdArgs struct {
    // --- positional arguments ---
    posA argparse.Int     `type:"pos",help:"posA help message"`
    // custom name to display
    posB argparse.String  `type:"pos",name:"POS_B",help:"posB help message"`
    //groups multiple positional arguments into a list
    posC argparse.Uint64List `type:"pos",nargs:"5/5+/5-10"`
    //user defined type
    posD cutomType       `type:"pos",help:"posD help message"`
    
    // --- optional arguments ---
    optE argparse.String `type:"opt",short:"-e",mutex:"grp1",help:"help for optE"`
    optF argparse.String `type:"opt",short:"-f",long:"--option-f",mutex:"grp1",help:"help for optF"`

    // --- sub commands --- ???
    subCommands []interface{}
}
```