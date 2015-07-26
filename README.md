# voyager

the main purpose of voyager is to provide a context aware file browser in a web page. 

voyager is a http server which can render file listings using the appropriate layout based
on the types of files found.

# features

- image resize 

# install

    
    export GOPATH=$HOME/go
    go get -u github.com/sigmonsays/voyager/cmd/voyager

create ~/.voyager file with the paths you wish to allow. They are relative to your $HOME directory. To allow 
~/Pictures to be browsable via HTTP, create ~/.voyager with this content:

    allow:
    - Pictures

# configuration

default configuration file is ~/.voyager 
