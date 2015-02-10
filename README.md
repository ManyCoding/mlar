# mlp - media library parser
Traverses through the input folder and renames albums' covers and artwork folders. Considers an album folder as a folder with one or more audio files (m4a, flac).


## Installation
	
	Get GO compiler if you still haven't - http://golang.org/


## Usage

	$ go run mlp.go -i="path"
	-h/-help
	-i/-inputDir - input directory
	Optional arguments:
	-n/-coverName - name for covers, default is "cover"
	-afn/-artworkDirName - name for artwork folders, default is "artwork"
	

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request
