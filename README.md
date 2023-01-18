# SeeProxy
Golang reverse proxy with CobaltStrike malleable profile validation. 
The premise of this tool is to not open your teamserver to the world but to a single instance of SeeProxy instead.
This way every request reaching your teamserver is a legitimate C2 traffic. 

![Demo](/demo/demo.gif)

## Usage: 

```bash
$ SeeProxy --teamserver <IP>:<PORT> --profile <path_to_malleable_profile>
```