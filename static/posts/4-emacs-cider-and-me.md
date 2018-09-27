Date: 2015-09-30
Title: Emacs, CIDER, and clojure-pretty-lambda-mode
cat: dev

#Emacs

Dual wield experience and theory as you do [Vim and Emacs](http://briancarper.net/page/422/about). One is trusty, sharp, and is ready to swing.
The latter takes longer to sharpen but will be leading the former in the end.

Well, both of these take some practice to get a grip on. But over the past few months I've gotten quite proficient with the Vim bindings and Vim's important features (macros, buffers).

Historically, Emacs provided a good lisp environment because it itself has a lisp dialect, Emacs Lisp, which is used for implementing most of its editing features. Its REPL, debugger, interactive expression evaluation and other features helped inspire [SLIME](https://common-lisp.net/project/slime/), an emacs mode for Common Lisp. For Clojure, the modern Emacs environment is [CIDER](https://github.com/clojure-emacs/cider). Bozhidar Batsov, the creator of CIDER, talks about the history of emacs as a lisp environment and what inspired CIDER in [this video](https://www.youtube.com/watch?v=4X-1fJm25Ww).

I never got a chance to use Emacs until I became more involved with Clojure and its community. CIDER, along with [Clojure-mode](https://github.com/clojure-emacs/clojure-mode) within Emacs is a far better experience than Vim + some REPL in another window (multiplexed or not) as you have powerful interaction ocurring between the Clojure code buffer and the REPL buffer. Now CIDER even has a debugger for your Clojure functions.

Give it a try, Clojure-mode makes for some nice highlighting right away, but the environment is super sexy when you tweak out your emacs setup:

![Emacs setup](/images/emacs-setup.png)

In particular, this setup has [Rainbow Delimeters](https://github.com/Fanael/rainbow-delimiters), [Rainbow Identifiers](https://github.com/Fanael/rainbow-identifiers), and [Relative Line Numbers](https://github.com/Fanael/relative-line-numbers) (all of these are by [Fanael](https://github.com/Fanael)).

This catches many eyes, not in the same way that having terminals open in tmux does (people will think you're a "hacker"), but in a "whoa, what is that trippy editor?" way.

#Emacs Lisp and clojure-pretty-lambda-mode.el

I thought that my second lisp language to start hacking in would be Common Lisp. But it ended up being Emacs Lisp as I very much wanted a Clojure-specific pretty lambda mode for my setup. At first, all this involved was taking [pretty-lambdada.el](http://www.emacswiki.org/emacs/pretty-lambdada.el) (pretty lambda for Emacs Lisp) and replacing the `"\<lambda\>"` regular expression with `"\<fn\>"`, but I was immediately unhappy with the fact that `fn?` became a lambda as well. After all, when you call `fn?` you're checking if something is a function, not if it is an anonymous function.

This was more difficult than I anticipated.

When I first tackled this problem, I tried to exclude matches of the length of fn that contain a `?` symbol at the end. However, the regular expression I built, `"\<fn\>[^\?]"` means that a third character, any character aside from `?` must be typed out. So, without changing any of the parameters of `font-lock-add-keywords`, particularly the index parameters for `(match-beginning)` and `(match-end)`, the entire match, including the character following `fn` is changed to a lambda.

Looking about how to create such a match involved something called negative lookahead. This isn't supported by emacs.

More importantly, if the `fn` in `fn?` is already matched, it will turn into a lambda and the composed region for the lambda will need to be decomposed. This involved a lot of re-writing for the functions and without any prior practice with Emacs Lisp, it was quite a learning experience. It was cool to see the similarities in the syntax and functionality of Clojure and another Lisp language.

I also tried to replace the matching region with a string `"<lambda><last-char>"` for the regex which only matches after `fn` is followed by a non-`?` character, where `<last-char>` is that particular character. But I was having trouble fitting the logic into `font-lock-add-keywords`, and for some reason `let` was not working for me. I'm sure there was a way of doing it, but after hours of playing around with it and not being able to get the last character of the match properly, I wasn't seeing much progress with this method.

I then pursued the idea of not using strings, but rather just decomposing the text region once the match contains a `?` symbol at the end. After a lot of hacking around with the logic, not being able decompose text-regions exactly the way I wanted to, and having some difficulty working with the core functions in Emacs, I decided to check online for a pretty lambda mode for Clojure.

Very quickly, I stumbled into a [config file](https://github.com/cemerick/.emacs.d#pretty-lambda-and-co) created back when Swank (a SLIME extension for Clojure) was popular. Basically, the regular expression matches for whenever the user types `(fn ` followed by a space. And, using `(match-begin 1)` and `(match-end 1)`, whose parameters I was not aware of how to use, (each number represents an index to skip, namely skip one cell forward in `(match-begin)` and one cell backwards in `(match-end)`). Then, I could revert all the changes to the functions' structures, and replaced the old regular expressions with the one in the config.
I could have went back to my solution with strings, but I've read somewhere that using strings is discouraged as it may lead to performance issues. But ultimately, I had a solution to the problem, and I understood the problem.

And then I asked myself, how come I didn't think of this simple method?

But it's okay, it works wonderfully (:
You can find the code [here](https://github.com/yonkornilov/clojure-pretty-lambda.el), and what it looks like in action below:

![Clojure Pretty lambda works!](/images/clojure-pretty-lambda-works.png)

I also got to dabble around with Emacs Lisp, and had the experience of taking someone elses code, particularly code in a language I am not very familiar with, and trying to solve a problem within the code as best as I possibly can. It was a very valuable experience.

##Bonus setup:

And, while on the subject of sexy setups, I am also using [spf13's Ultimate Vim Distribution](http://vim.spf13.com/). It works right off the bat and has some very useful features. However I changed the default colorscheme.

![My vim setup (;](/images/vim-setup.png)

The bar at the very bottom is the tmux status bar.
