Date: 2017-10-9
Title: Installing Synergy on Odroid
cat: ops

If you install synergy via apt-get:

```
$ sudo apt-get install synergy
```

You may be disappointed in that you're not installing the latest version.
Furthermore, if you're using an armhf device such as Odroid which only currently supports Ubuntu 14.04, you will not find a PPA that provides the latest synergy package.

###Finding the latest release

Find the latest release of Synergy for armhf [here](http://ftp.debian.org/debian/pool/main/s/synergy/?C=N;O=D)

###Why you can't dpkg -i then apt-get install -f

If you `dpkg -i`, you will be further dissapointed because your Ubuntu 14.04 does not provide libssl1.1. So there's no way to `dpkg -i` and then `apt-get install -f`.

###Compiling OpenSSL 1.1.0 from source

You need to get OpenSSL 1.1.0 without the help of any repository:

```
$ wget https://www.openssl.org/source/openssl-1.1.0.tar.gz
$ tar xvzf openssl-1.1.0.tar.gz
$ cd openssl-1.1.0
$ sudo ./configure.sh
$ sudo make
$ sudo make install
```

Now, make sure that the newly installed library is in the ldpath:

```
$ sudo ldconfig /usr/local/lib/
```

Or wherever libssl.1.1.so was installed.  
You can now start the synergy server:

```
$ synergys
```

If you're having any errors, run synergys in the foreground to see the output:

```
$ synergys -f
```
