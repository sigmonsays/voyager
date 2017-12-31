# voyager

the main purpose of voyager is to provide a context aware file browser in a web page.  voyager 
is a http server which can render file listings using the appropriate layout based on the 
types of files found.

A remote API is currently in development to allow connecting to remote servers.

# features
- context aware layout (photos, video, audio, etc)
- image resize for thumbnails with a local cache
- HTML5 video playback using jplayer
- HTML5 audio playback using jplayer
- remote API so multiple voyagers can talk to each other
- network ACL access control

# install
    
    export GOPATH=$HOME/go
    go get -u github.com/sigmonsays/voyager/cmd/voyager

create ~/.voyager.cfg file with the paths you wish to allow. They are relative to your $HOME directory. To allow 
~/Pictures to be browsable via HTTP, configure with this content:

    default:
       allow:
       - Pictures

The path to access files from your home directory is access via ~/username on the URL. So the url for Pictures in the previous example
would be /~user/Pictures.

# configuration

default configuration file is ~/.voyager.cfg

The configuration file has a section per host. The `default` section is applied to all host, so it is suitable for configuring parameters which
apply to all nodes. The node specific configuration section used is determined based on the hostname. For a given FQDN of `desktop.example.net` the string up to
the first dot is used as a configuration section; In this case `desktop` and `.example.net` would be ignored.

The following parameters are available

- startupbanner - string - printed when daemon starts
- autoupgrade - bool - if you run voyager with $GOPATH set the program will auto upgrade when changes are detected
- autorestart - bool - automatically restart daemon if the binary changes (useful with autoupgrade)
- http.bindaddr - string - http bind address. default is :8181
- allow - list - list of paths to allow from your home directory
- layouts - map - map of paths enforcing a layout. Most significant path is used.
- alias - map - map of top path name to alias local path


