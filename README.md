# findsecret
Find secret keys from JS file

## Installation
1. Install Go on your system
2. `go install github.com/burak0x01/findsecret@latest`

## Usage
Findsecret requires 2 parameters to run: -i (input), -o (output).

### For example 
`findsecret -i https://somedomain/something.js` </br> </br>
`findsecret -i main.js -o result.txt` </br> </br>
`findsecret -i domains.txt -o result.txt`
