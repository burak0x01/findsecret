# findsecret
Find secret keys from JS file

## Installation
### From Source
1. Install Go on your system
2. `go install github.com/burak0x01/findsecret@latest`

### From Binary
You can download the pre-built binaries from the releases page and run. For example:
`unzip findsecret_linux_amd64.zip` </br>
`mv findsecret /usr/bin` </br>
`findsecret -h`

## Usage
Findsecret requires 2 parameters to run: -i (input), -o (output).

### For example 
`findsecret -i https://somedomain/something.js` </br> </br>
`findsecret -i main.js -o result.txt` </br> </br>
`findsecret -i domains.txt -o result.txt`
