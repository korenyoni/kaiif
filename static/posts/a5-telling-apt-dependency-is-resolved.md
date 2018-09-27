Date: 2017-10-16
Title: Telling apt a dependency is resolved
cat: ops

In some cases you may install a package outside your package manager (e.g. using `dpkg -i` in Ubuntu), and this package may have dependencies unavailable to your system.

This will likely cause your package manager to prevent you from taking any actions unless you uninstall the said package.

An example of this occurs in my [post](https://yonatankoren.com/post/a1-synergy-odroid).

![apt broken after  manual resolve](/images/breaking-apt.png)

Building the dependencies from source won't let apt (or other package managers) that you resolved these dependencies.

![Building from source apt dependencies are not met](/images/fixing-apt-1.png)

One way to resolve this is by creating what is essentially a fake, empty package for whatever dependency you resolved manually.

For Ubuntu, `equivs` is the tool to create fake packages, but this theory can be extended to other distros.

###In Ubuntu:

```
$ sudo apt-get install equivs
$ equivs-control somepackage
```

Replace `somepackage` with whatever package you're faking. Use a text editor to edit the `somepackage` control-file and replace version info and other package data.

```
$ equivs-build somepackage
```

Will build a .deb file. You can now use `dpkg -i` to install this fake package.

![use equivs to build a fake package](/images/fixing-apt-2.png)
