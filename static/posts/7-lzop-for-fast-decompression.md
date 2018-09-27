Date: 2017-09-14
Title: Switching to lzop for fast decompression
cat: ops

When I provision my instances, I need to download an archive of a translation engine and extract it. Some options that I've tried:

 - bz2: good compression ratio, but takes a long time to decompress
 - gz: also good compression, but not fast enough for my use

I was recommended to try **lzop** as the decompression speed is allegedly much faster.

I created an archive of the directory *default*:

`tar c default | lzop - > default.tar.lzo`

I compared the size to the *.gz* archive:

```
# stat default.gz | grep 'Size'
Size: 6648380497 Blocks: 12993623   IO Block: 131072 regular file
# stat default.tar.lzo | grep 'Size'
Size: 8810709804 Blocks: 17219473   IO Block: 131072 regular file
```

Clearly the *.lzo* archive is larger. But perhaps the decompression speed can make it worth it for my use case, so let's **time** it.

We also have to pipe it through tar in order to keep the directory structure, since it's not a single file:

```
# time lzop -d -c default.tar.lzo | tar xv
real    1m18.581s
user    0m50.349s
sys     0m33.539s
# time tar xvzf default.gz
real    2m23.931s
user    2m2.211s
sys     0m38.708s
```

Clearly it's significant. If you favour speed over compression ratio, you should definitely use lzop.

Apparently there's a way to make the decompression even faster. Let's try it out:

You need to add the `--fast` flag while compressing:

`tar c default | lzop --fast - > default.tar.lzo`

```
# time lzop -d -c default.tar.lzo | tar xv
real    1m21.701s
user    0m45.270s
sys     0m46.466s
# stat default.tar.lzo | grep 'Size' 
Size: 8855209044 Blocks: 17306449   IO Block: 131072 regular file
```

Surprisingly this run was even slower than the first lzop run. However the archive is slightly larger in size when you use the `--fast` flag. I recommend using the default, i.e. not using the flag.
