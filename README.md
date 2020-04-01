# tube

`tube` is a Youtube-like (_without censorship and features you don't need!_)
Video Sharing App written in Go which also supports automatic transcoding to
MP4 H.265 AAC, multiple collections and RSS feed.

## Features

- Easy to add videos (just move a file into the folder)
- Easy to upload videos (just use the builtin uploader and automatic transcoder!)
- Builtin ffmpeg-based Transcoder that automatically converts your uploaded content to MP4 H.264 / AAC
- Builtin automatic thumbnail generator
- No database (video info pulled from file metadata)
- No JavaScript (the player UI is entirely HTML, except for the uploader which degrades!))
- Easy to customize CSS and HTML template
- Automatically generates RSS feed (at `/feed.xml`)
- Clean, simple, familiar UI

### Screenshots

![screenshot-1](screenshot-1.png?raw=true "Main Screen and Video Player")
![screenshot-2](screenshot-2.png?raw=true "Video Upload Screen")

## Getting Started

### Using Homebrew

```#!sh
$ brew tap prologic/tube
$ brew install tube
$ tube
```

Open http://127.0.0.1:8000/ in your Browser!

### Using a Binary

1. Go grab the latest binary from the
   [Releases](https://github.com/prologic/tube/releases) page for your
   platform / operating system.
2. Extract the archive.
3. Run `./tube`

Open http://127.0.0.1:8000/ in your Browser!

### Using Docker

```#!sh
$ docker pull prologic/tube
$ docker run -p 8000:8000 -v /path/to/data:/data
```

Open http://DOCKER_MACHINE_IP:8000/ in your Browser!

Where `DOCKER_MACHINE_IP` is the IP Address of your Docker Node.

### From Source

```#!sh
$ git clone https://github.com/proogic/tube
$ cd tube
$ make
$ ./tube
```

Open http://127.0.0.1:8000/ in your Browser!

### A Production Deployment

A Production Deployment can be found at https://tube.mills.io/ -- This is run
and managed via a Docker Swarm cluster with a Docker-Compose / Stack very
similiar to the one you can find in the repo: [docker-compose.yml](docker-compose.yml)

Beyond this a "Production Deployment" is out-of-scope at this time for the
documentation being provided here. Please don't hesitate to file an
[Issue](https://github.com/prologic/tube/issues/new) however for ask for help
or advice or contact the author directly!

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/prologic/tube.svg)](https://starcharts.herokuapp.com/prologic/tube)

## Support

Support the ongoing development of Tube!

**Sponser**

- Become a [Sponsor](https://www.patreon.com/prologic)

## Contributors

Thank you to all those that have contributed to this project, battle-tested it,
used it in their own projects or products, fixed bugs, improved performance
and even fix tiny typos in documentation! Thank you and keep contributing!

You can find an [AUTHORS](/AUTHORS) file where we keep a list of contributors
to the project. If you contriibute a PR please consider adding your name there.
There is also Github's own [Contributors](https://github.com/prologic/tube/graphs/contributors) statistics.

[![](https://sourcerer.io/fame/prologic/prologic/tube/images/0)](https://sourcerer.io/fame/prologic/prologic/tube/links/0)
[![](https://sourcerer.io/fame/prologic/prologic/tube/images/1)](https://sourcerer.io/fame/prologic/prologic/tube/links/1)
[![](https://sourcerer.io/fame/prologic/prologic/tube/images/2)](https://sourcerer.io/fame/prologic/prologic/tube/links/2)
[![](https://sourcerer.io/fame/prologic/prologic/tube/images/3)](https://sourcerer.io/fame/prologic/prologic/tube/links/3)

## License

tube source code is available under the MIT [License](/LICENSE).

Previously based off of [tube](https://github.com/wybiral/tube) by [davy wybiral
](https://github.com/wybiral). (See [LICENSE.old](/LICENSE.old))
