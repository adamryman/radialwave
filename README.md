# visualizer-experiment

![Harder Better Faster Stronger by Daft Punk](./screenshots/harder.png)

This project is converting a music visualizer I wrote in matlab and I hope to
rewrite in go.

See [_matlab/matlab.m](_matlab/matlab.m) for the original matlab code with
comments I had to write to understand it.

It partitions the song, uses an fft to determine the maximum frequency of that
partition, and stores the max frequency of each partition in order.

It then creates a mapping of unique dominate frequencies to wavelengths of
color.

Then it then can render the song as sections of color based on an mapping the
dominate frequency of each partition to a wavelength of light using the
previous mapping and then mapping that wavelength to RGB with a static
csv file that in golang would be a `[][]float32`. The matlab code draws this
onto a circle from the start of the music to the end of the music.

## TODO

1. Do the problem first by breaking the song by seconds as that defines the sampling rate.
2. Don't worry about optimizing. Do things at the end. Make them simple. Do them fast. You can always break off more parts and make them faster later.
3. After that worry about upsampling so you can have more arcs than one per second.

// TODO:
// - read in data chunks
// - possibly assign chunks numbers to recombine later
// - pass groups of chunks to something to do fft on them.
// - Find dominate frequency of group of chunks
// - Pass dominate frequency info to both something to create a slice of the frequency data in
//   original music AND something to create a sorted unique slice of frequencies,
//   this could be done effectily with insertion sort
// - After all pieces have been collected:
// - Create mapping from wavelength to color
// - Render some somehow

