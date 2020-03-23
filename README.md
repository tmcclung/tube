# tube,

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

Currently only supports MP4 video files so you may need to re-encode your media to MP4 using something like [ffmpeg](https://ffmpeg.org/).

Since all of the video info comes from metadata it's also useful to have a metadata editor such as [EasyTAG](https://github.com/GNOME/easytag) (which supports attaching images as thumbnails too).

## Getting Started

### From Source

```#!sh
$ git clone https://github.com/proogic/tube
$ cd tube
$ make
$ ./tube
```

Open http://127.0.0.1:8000/ in your Browser!

### Using Docker

```#!sh
$ docker pull prologic/tube
$ docker run -p 8000:8000 -v /path/to/data:/data
```

## License

tube source code is available under the MIT [License](/LICENSE).

Previously based off of [tube](https://github.com/wybiral/tube) by [davy wybiral
](https://github.com/wybiral). (See [LICENSE.old](/LICENSE.old))
