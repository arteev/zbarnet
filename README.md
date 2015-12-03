# zbarnet
Wrapper over zbarcam (httpclient)


[![Build Status](https://travis-ci.org/arteev/zbarnet.svg)](https://travis-ci.org/arteev/zbarnet)

Welcome to ZBarNet, a wrapper over zbarcam ([ZBar](http://zbar.sourceforge.net/) bar code reader)

**Features:**

- Analysis of the bar code from the output zbarcode
- The output to the console in the format: "raw", "json"
- Sending barcode protocol http get or post methods
- A single execution and closing zbarcam
- Work as a service / daemon or cli

**Menu:**

- [Installation](#installation)
- [Quick start](#quick-start)
- [Service/daemon](#servicedaemon)
- [Configure](#configure)

Installation
------------

1. Install [ZBar](http://zbar.sourceforge.net)
2. $ `go get github.com/arteev/zbarnet`

Quick start
-----------

Configure $HOME/".zbarnet.json"

The file can be located in the same directory with the program or in your home folder.

```

{
        "source": "zbar",
        "output": "json",
        "once": false,
        "zbar": {
                "enabled": true,
                "location": "/usr/bin/zbarcam",         
                "device": "/dev/video0",
                "args": [
                  "-q",
                  "--xml"
                ]
        },
        "http": {
          "enabled": true,
          "url": "http://httpbin.org/post?barcode=${barCode}&type=${barCodeType}&quality=${quality}&api=${apikey}",
          "method": "POST",
          "apikey": "THIS_IS_API_KEY",
          "apikeyhdr": true
        }
}

```

$ zbarcam

Output:
```
Hit Ctrl+C to [EXIT]
{
"type": "QR-Code",
"quality": 1,
"data": "d3d3LmRlZmVuZGVyLnJ1L2NvbnRlbnQvcXIvP2FydD02MzExMA=="
}
```

Service/daemon
--------------

1. Configure $GOPATH/bin/.zbarnet.conf
2. $ sudo $GOPATH/bin/zbarnet service install
3. $ sudo $GOPATH/bin/zbarnet service start

Also see usage: $ $GOPATH/bin/zbarnet service -h

Configure
---------

TODO

License
-------

MIT

Author
------

Arteev Aleksey
