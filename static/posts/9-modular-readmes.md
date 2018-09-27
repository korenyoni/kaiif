Date: 2017-10-07
Title: Creating modular READMEs
cat: dev

Github doesn't allow you to use include directives (reuse snippets of .md or .rst in a file without copy-paste):

![issue](/images/issuescreen.png)

So, if you want to reuse snippets of documentation in your repository, you have two options:

1. Copy and Paste (__BAD!__)
2. Compile your documentation before pushing to Github

I'm using the second technique in my [repo](https://github.com/yonkornilov/opus-api):

######README.recipe.rst
```
OPUS_ (opus.lingfil.uu.se) Python API

* Free software: MIT license
* Documentation: https://opus-api.readthedocs.io.

!requirements

!installation

!usage

!credits

```

Where docs/requirements.rst, docs/installation.rst ... are existing documentation files.

A command `diff -u replace replaceWith > path` followed by `patch README.recipe.rst patchfile -o README.rst` will replace all occurences of `replace` with `replaceWith` in `README.recipe.rst` and output it to `README.rst`. This is effectively our 'compilation' of `README.rst`.

With that, some redirection and piping, we can make a nice Makefile target:

######Makefile
```
readme: ## replace variables in README.recipe.rst and write README.rst
	rm README.rst
	bash -c "diff -u <(echo '!requirements') docs/requirements.rst | patch README.recipe.rst -o README.rst"
	bash -c "diff -u <(echo '!installation') docs/installation.rst | patch README.rst"
	bash -c "diff -u <(echo '!usage') docs/usage.rst | patch README.rst"
	bash -c "diff -u <(echo '!credits') docs/credits.rst | patch README.rst"
	rm README.rst.orig
```

Note that this target first patches the recipe and outputs to `README.rst`, then for each docfile we patch `README.rst` and overwrite it. After running `make readme`, we have our fully compiled `README.rst`, ready to be pushed:

######README.rst
```
OPUS_ (opus.lingfil.uu.se) Python API

* Free software: MIT license
* Documentation: https://opus-api.readthedocs.io.

.. _requirements:

.. highlight:: console
.. _PhantomJS: http://phantomjs.org/download.html

============
Requirements
============

... etc etc
```
