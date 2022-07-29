# findsecret
Find secret keys from JS file

## Installation
### From Source
1. Install Go on your system (min v17.0.0 and GO111MODULE env should be "on")
2. `go install github.com/burak0x01/findsecret@latest`

### From Binary
You can download the pre-built binaries from the releases page and run. For example: </br> </br>
`unzip findsecret_linux_amd64.zip` </br>
`mv findsecret /usr/bin` </br>
`findsecret -h`

## Usage
Findsecret requires 2 parameters to run: -i (input), -o (output).

#### Input
`findsecret -i https://domain.tld` <br>
`findsecret -i ./local/path/main.js` <br>
`findsecret -i https://somedomain/something.js` <br>

#### Output
`findsecret -i https://domain.tld -o cli (default)` <br>
`findsecret -i https://domain.tld -o result.txt` <br><br>

### For example 
`findsecret -i https://findsecret.s3.eu-central-1.amazonaws.com/index.html`

![Screenshot 2022-07-20 102542](https://user-images.githubusercontent.com/60436338/179922445-09657d4a-ebba-4395-85bf-6cf904b43b2a.png)

### ToDo
- [X] New input - Domain
- [X] New input - Domain list
- [ ] Relative url for domain input  
- [ ] New input - JS list
- [ ] User-Agent option
