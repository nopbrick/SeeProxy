# SeeProxy
Golang reverse proxy with CobaltStrike malleable profile validation.  
The premise of this tool is to not open your teamserver to the world but to a single instance of SeeProxy instead.  
This way every request reaching your teamserver is a legitimate C2 traffic.

## Example deployment
Below you can find a very basic example deployment for a red team engagement. Only valid traffic from the intance of a SeeProxy is permitted to reach the C2. 

![Example Diagram](/demo/example_diagram.jpg)

<p>
<p>

![Demo](/demo/demo.gif)

<p>

## Usage: 

```bash
$ make
$ SeeProxy --teamserver <IP>:<PORT> --profile <path_to_malleable_profile> --port <local_port>
```

## Demo

<iframe width="800" height="450" src="https://www.youtube.com/embed/iWuphwQggxk" title="SeeProxy demo" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
