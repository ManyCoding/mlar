# mlp - media library parser
Traverses through input folder and renames albums' artwork. Considers an album folder as a folder with one or more audio files (m4a, flac). Skips renaming if there is more than one image.


## Installation
	
	Get GO compiler if you still haven't - http://golang.org/


## Usage

	$ go run mlp.go -i="path" (-n="cover")
	-i/-inputDir - input directory
	-n/coverName - optional, new name for artwork. Default is "folder"
	

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request
