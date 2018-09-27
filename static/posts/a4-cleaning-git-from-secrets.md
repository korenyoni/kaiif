Date: 2017-10-12
Title: Removing secrets from a git repository
cat: dev

So, you just realized your git repository had some sensitive information, i.e. passwords, private keys at some point in its lifetime.

If you're in a public repository, consider this information compromised. **Delete the remote repo**. Replace all involved passwords and disable your private keys immediately.

In a private repository, it's a good idea to clean your history anyways.

![git repository with secrets actions to take](/images/git-clean-flow-intro.png)

###Your dirty history

Removing these files from your git repository in a commit won't remove them from your history.

![git repository with secrets removed](/images/cleaning-git-bad.png)

You can check this using `git log` (the `p` flag shows insertions / deletions):

```
$ git log -p | grep 'secret' -B 2 -A 2
+}
+
+variable "secret_key" {
+  default = "mypass123"
+}
-}
-
-variable "secret_key" {
-  default = "mypass123"
-}
```

The options `B` and `A` set how many lines grep will print before and after the match, respectively. This will show you the value you're looking for in your git log and show you some context.

###Rewriting the repository history

![git repository with visible secrets](/images/cleaning-git.png)

GitHub provides a command that can be used to rewrite the repository's history. In [this](https://help.github.com/articles/removing-sensitive-data-from-a-repository/) post however, GitHub recommends using [bfg](https://rtyley.github.io/bfg-repo-cleaner/) to more easily and efficiently perform specific prunes from your repository history.

Download the latest `bfg.jar` and put it somewhere in your machine. Now create a `replacements.txt` file which contains the sensitive string:

```
$ echo mypass123 >> replacements.txt
```

Which will replace the sensitive string with `***REMOVED***`. You can replace it with whatever you like (credits: [this](https://stackoverflow.com/a/15730571/4650776) StackOverflow answer):

```
$ echo 'mypass123==>WhateverYouLike' >> replacements.txt 
```

Now you can run `bfg`:

```
$ java -jar ~/bfg.jar --replace-text replacements.txt
```

Now, two things might happen:

1. `bfg` will tell you that the offending files are protected
    * Get rid of those files in the head commit
    * Re-run
2. `bfg` will tell you to prune your history:
    *  do `$ git reflog expire --expire=now --all && git gc --prune=now --aggressive`
    
You can check if the sensitive string is anywhere in your repository, including your .git files:

```
$ rm replacements.txt
$ grep -rnw ./ -e 'mypass123'
```

***Use the above command to look for any other sensitive strings you may have forgotten about.***

Now you can update your repositories:

![git repository with secrets actions to take after cleaning](/images/git-clean-flow-end.png)
