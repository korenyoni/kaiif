Date: 2017-10-29
Title: Using Docker to build packages
cat: dev,ops

Have you ever really wanted to build some software?

Did it require dependencies not available for your system's release?

![Package dependencies break system](/images/unstable-package-building.png)

These unstable dependencies can cause your package manager to uninstall other packages. Their installation alone, even without uninstalls, can still break functionality of other packages.

Instead, we can use a docker container:

![Using docker for building packages](/images/docker-package-building.png)

### You can use another release or another distro

The package manager of whatever releaes or distro you choose may already have the packages you need to build your binaries.

### You won't break your system

Since its a docker container, whatever you do inside has no effect on the system running it. 

### You will save space

Imagine all the space you will save by deleting the container containing the dependencies required for building the binaries. Even if you do this manually, using a container is far more convenient.

### An example

Install `docker-ce` from an official repository, then:

Make sure you're in the docker group, so you don't need `sudo`:

```
$ sudo usermod -a -G docker $USER
```

Download and run the container:

```
$ docker run --security-opt seccomp:unconfined -i -t ubuntu:17.04 /bin/bash
```

If you exit, you can attach to it later:

```
$ docker container ls -a | awk '{print $1}'
CONTAINER
4d3f22cdabd5
$ docker start 4d3f22cdabd5
$ docker attach 4d3f22cdabd5
```

Then,

1. `apt-get update`
2. `apt-get upgrade`
3. Install your dependenices
4. Build

`cde` is a program which bundles all called libraries from your software into a package.

Download it from [here](http://www.pgbovine.net/cde.html)

```
cd /software/somepackage
$ chmod +x cde
./cde /software/somepackage/bin/somepackagebin
```

Once you're done, you can exit the container and copy the cde package from the container to your machine:

```
$ docker cp 4d3f22cdabd5:/software/cde-package ~/software/somepackage-cde
```

Now you can run it:

```
$ cd ~/software/somepackage-cde/cde-root
$ ../cde-exec software/somepackage/bin/somepackagebin
```

There you go!
