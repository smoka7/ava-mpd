# Ava
Ava is a web client for MPD.
![web screenshot](https://raw.githubusercontent.com/smoka7/ava-mpd/main/assets/screen.webp)

![mobile screenshot](https://raw.githubusercontent.com/smoka7/ava-mpd/main/assets/screen-mob.webp)

## Installation

Download the binaries from release page; or you can build it like this.

```bash
git clone https://github.com/smoka7/ava-mpd.git
cd ava-mpd
chmod +x install.sh
./install.sh
```

## Usage
```bash
# it reads MPD_HOST and MPD_PORT from env or you can use these flags
#  -address string
#   	address of mpd server host:port
#  -password string
#    	password of mpd server
#  -port string
#   	The port to run this app on it (default "3001")

ava -address 127.0.0.1:6602  -password xxxx -port=3333
```

## Contributing
Pull requests are welcome.

## License
[MIT](https://choosealicense.com/licenses/mit/)
