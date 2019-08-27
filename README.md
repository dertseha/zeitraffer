## Zeitraffer
A simple tool that grabs the content of the screen every second and saves a `.png` file. Written in Go.

Supported OS: MS Windows only.

### Motivation
I create digital sketchnotes at conferences and once another attendee wanted to see a time lapse video of creating them.
This was the basis for the idea of creating a screenshot every second to then convert these frames into a complete movie.

This tool creates the screenshots. The word `zeitraffer` is German for time lapse.

### Usage

(all from the console) 

* `go get` it with `go get github.com/dertseha/zeitraffer` (requires Go)
* Enter the directory you want to store the files in (should be empty)
* Start it `path\to\install\dir\zeitraffer.exe`
* Let it run until done, then press `Ctrl+C` to stop it.

### Convert to movie

This requires an `ffmpeg` installation (I do this via WSL):
```
ffmpeg -r 40 -f image2 -s 3840x2160 -i ~/path/to/images/%05d.png -vcodec libx264 -crf 25 -vf scale=1920:1080 -pix_fmt yuv420p output.mp4
```

This creates an `.mp4` file scaling it to HD resolution and 40 fps - the maximum parameters allowed for Twitter.

### License

The project is available under the terms of the **New BSD License** (see LICENSE file).
