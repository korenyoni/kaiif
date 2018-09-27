Date: 2017-10-25
Title: Fixing only dummy output showing
cat: ops

Sometimes you may lose sound on your linux machine after installing a package.

### A DANGEROUS SOLUTION
The usual suggestion is to `apt-get purge alsa-base pulseaudio`.

If your distribution has core packages that depend on alsa-base and pulseaudio, this command will remove them too and break your system.


### An alternative solution

Instead, make sure your user is in the `audio` group:

```
$ groups

username adm cdrom sudo dip plugdev lpadmin sambashare
```

Add yourself to the `audio` group:

```
$ sudo usermod -a -G audio username
```

Check your groups:

```
$ groups

username adm cdrom sudo audio dip plugdev lpadmin sambashare
```

Reboot.
