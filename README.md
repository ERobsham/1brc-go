#  One Billion Row Challenge - Go

Having completed the 1BRC in `Objective-C` ([more info here](https://github.com/ERobsham/1brc-objc/tree/main)), I decided to _give it a go_ in `go`!  I expect the 'baseline' / 'naive' implementation to be much faster than ObjC.  And I think for this round, I'm actually going to focus on getting some better performance squeezed out of the parsing logic / and optimizing GC / minimizing memory allocations etc, rather than just exploring how much overhead it creates (ie, what I was exploring for the objc challenge). 

## Rules:
* Input value ranges are as follows:
    * Station name: non null UTF-8 string of min length 1 character and max length 100 bytes, containing neither ; nor \n characters. (i.e. this could be 100 one-byte characters, or 50 two-byte characters, etc.)
    * Temperature value: non null double between -99.9 (inclusive) and 99.9 (inclusive), always with one fractional digit
* There is a maximum of 10,000 unique station names
* Line endings in the file are \n characters on all platforms
* Implementations must not rely on specifics of a given data set, e.g. any valid station name as per the constraints above and any data distribution (number of measurements per station) must be supported
* The rounding of output values must be done using the semantics of IEEE 754 rounding-direction "roundTowardPositive"
* Output should:
    * be sorted by 'station name'
    * follow the general format of: `{'station name'='min'/'mean'/'max', ...}`

Output example:
```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
```

Test runner specs:
2017 iMac / 3.6ghz Quad Core i7 / 8gb RAM.

# Initial Results

Go FTW!  A very naive implementation, using the standard `bufio` default line scanner, then parsing the resulting `string`s on a single thread is already more than twice as fast as the base `Objective-C` implementation:  `147.44 real       134.20 user         9.69 sys` (just under 2.5 mins)

This time, we can easily setup `pprof` and get some nice flame graphs to parse these results(`make profile` then `make view-prof` to open a browser with latest results).

And immediately the one standout bottle neck is just purely reading in the file: `bufio.(*Scanner).Scan() (93.78%)` 

This is outstanding, we already have a super clear direction of where to go to try speeding things up!

