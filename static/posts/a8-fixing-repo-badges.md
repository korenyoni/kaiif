Date: 2017-10-29
Title: Fixing repo badges
cat: dev

So you have a public repo. You're happy with your code coverage, your test status, and your documentation status, but one or more badges is telling people something is wrong:

![failing docs badge](https://media.readthedocs.org/static/projects/badges/failing.svg)

GitHub caches your badges, but sometimes it doesn't work as it should. One possilble solution is to change your badge url from:

```
https://readthedocs.org/projects/someproject/badge/
```

to

```
https://img.shields.io/readthedocs/someproject.svg
```

You're now using `shields.io` for your badges, which does the caching on their end --and properly, at that.

This will fix your problem.

![passing docs badge](https://media.readthedocs.org/static/projects/badges/passing.svg)
