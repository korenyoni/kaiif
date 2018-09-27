Date: 2015-09-07
Title: Lazy Prime Sequence in Clojure
cat: dev

##My curiosity on the Incremental Sieve
For a long time I've been aware of the [Sieve of Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes). I've also implemented the [Sieve of Atkin](https://en.wikipedia.org/wiki/Sieve_of_Atkin) in C++ (although I admit, I have not read the paper
explaining why it works), but these traditional sieve algorithms which work on vectors or arrays can only discover primes up to a fixed limit.

Project Euler's [7th Question](https://projecteuler.net/problem=7) was a challenge for me as it asks you to find the 10001st prime, where
I only knew how to use the two mentioned Sieves to compute all primes up to a certain limit.

I did not think much of this problem again as at the time I simply computed all the primes up to some high limit and took the 10001st prime from the vector. However after running into [the 67th 4clojure problem](http://www.4clojure.com/problem/67), wanting to use lazy
sequences and being unsatisfied with a lazy sequence created using trial division, I began researching.

I found a [paper by Melissa O'Neill written on this topic](https://www.cs.hmc.edu/~oneill/papers/Sieve-JFP.pdf)
and I also ran into a [blog post on Christophe Grand's blog](http://clj-me.cgrand.net/2009/07/30/everybody-loves-the-sieve-of-eratosthenes/).
And it turns out what I was curious about is called an *incremental sieve*.

##Lazy Sequences

If you are not already aware, a lazy sequence is a sequence whose values are evaluated only when needed. Thus, they can be infinite sequences and you can specify to take only the first n
amount of elements without having to compute the entire sequence (impossible if your lazy sequence is infinite).

You can create a lazy sequence in Clojure similar to `(range)` quite easily:

    (defn my-infinite-sequence
        []
        (letfn [(next-number
                    [x]
                    (lazy-seq (cons x (next-number (inc x)))))]
            (next-number 0)))

`lazy-seq` wraps around this recursive call to ensure we can access it safely (without having to compute the sequence to infinity).
`cons` is a function which (lazily) returns a collection with the specified element being at the front.

You can expand this lazy sequence like this:

    (cons 0 (next-number 1))
    (cons 0 (cons 1 (next-number 2)))
    (cons 0 (cons 1 (cons 2 (next-number 3))))

etc...

So we can take the first 10 elements of this infinite sequence:

    user=> (take 10 (my-infinite-sequence))
    (0 1 2 3 4 5 6 7 8 9)

##My incremental, Incremental Sieve implementations (;

After reading the paper by *O'Neill* and avoiding copying anything from Christophe Grand's implementation for the infinite lazy prime sequence,
my first solid attempt at the incremental sieve looked something like so:

    (defn prime-sieve
      []
      (letfn [(next-composite [composite-map x]
                (if-let [composite-prime-set (composite-map x)]
                  (merge-with clojure.set/union composite-map
                              (into {} (map #(vector (+ x %) #{%}) composite-prime-set)))
                  (assoc composite-map (+ x x) #{x})))
              (create-lazy [x composite-map]
                (let [composite-prime (composite-map x)
                      prime? (nil? composite-prime)
                      composite-map (next-composite composite-map x)
                      composite-map (dissoc composite-map x)]
                  (if prime?
                    (lazy-seq (cons x (create-lazy (+ x 1) composite-map)))
                    (recur (+ x 1) composite-map))))]
        (lazy-seq (create-lazy 2 {}))))

I was quite happy when this produced an infinite sequence of primes. I got the logic by looking at the paper briefly and doing some thinking on
my own.

In principle, just as there is a *crossing-off* of the next composite number in a sieve of Eratosthenes, the same crossing off is happening
here with the `next-composite` function. However since the limit for the sieve is infinite, we will only cross off each case of the specified number
added to one of its corresponding factors. And if a number did not exist in the map, i.e. it is a prime, we will associate `(+ number number)`
with `#{number}` (a set literal in Clojure: `#{}`). So if we arrive at 6 which was associated with `#{2 3}`, we will now associate 8 with `#{2}` and 9 with `#{3}`. In case this *crossing-off*
will lead to a number which already exists in the map, for example the case of 6 which was once associated with `#{3}` when x was 3,
once it got to 4, the sets of the two entries `6 #{3}` and `6 #{2}` were merged with `clojure.set/union` to create `6 #{2 3}`.

`create-lazy` calls this function and decides if a number is prime judging by its presence in the composite-map. If it exists, we will make a recursive call with
the number iterated by 1, otherwise we'll `cons` it into the sequence.

However, this implementation is slow. I knew this because Christophe Grand's solution could give you all the numbers up to 1,000,000 in 1.5s. This took 15 seconds...
There's a couple reasons why, and I realized it through some trial and error, refactoring, and another glance at Christophe Grand's lazy prime sequence implementation:

1. The merging of the sets within the map using `clojure.set/union` is most likely expensive.
2. We don't need to keep track of all the prime factors for each composite number.

Primarily, our map associations should take the form of `[composite prime]`, not `[composite #{(all prime factors)}]`. Then, the new implementation begins like so:

1. Start at 2, is a prime (got nothing), associate [4 2]
2. At 3, is a prime (got nothing), associate [6 3]
3. At 4, not a prime (got 2), 6 exists. Recur. Associate [8 2].
4. At 5, is a prime (got nothing), associate [10 5]
5. At 6, not a prime (got 3). Associate [9 3].
6. Repeat to infinity!

This implementation looked like this:

    (defn prime-sieve2
      []
      (letfn
          [(next-composite [composite-map n step]
             (if (composite-map n)
               (recur composite-map (+ n step) step)
               (let [final-value (if (= n step) (+ n step) n)]
                 (assoc composite-map final-value step))))
           (next-prime
             [n composite-map]
             (let [prime-factor (composite-map n)
                   prime? (nil? prime-factor)
                   step (if prime-factor prime-factor n) 
                   composite-map (next-composite composite-map n step)
                   composite-map (dissoc composite-map n)]
               (if prime?
                 (lazy-seq (cons n (next-prime (+ n 1) composite-map)))
                 (recur (+ n 1) composite-map))))]
      (lazy-seq (next-prime 2 {}))))

This was a lot faster! To get all the primes under 1,000,000 took 3.2 seconds. However for a long time I was wondering why it still was not as fast as Christophe Grand's implementation.
After playing around with things, trying to minimize branching (if statements) and not getting any faster code, I decided to look at his code and at the paper more thoroughly.

I missed a third idea:

* Only check for odd numbers! Iterate (+ n 2) instead of (+ n 1).

That's right, `cons` 2 into the lazy sequence, then start the algorithm at 3, making a recursive call for n + 2 instead of n + 1 in `next-prime`. After seeing this broke something,
and thinking about it again, since we are only dealing with odd numbers, that means we should *step-up* each composite number by (+ composite step step) instead of (+ composite step),
thus giving us composites in odd multiples, not even multiples which would never be reached when only dealing with odd numbers i.e. (+ 3 2 2 2 2 2... 2).

The refactored code now looks like this:

    (defn prime-sieve3
      []
      (letfn
          [(next-composite [composite-map n step]
             (if (composite-map n)
               (recur composite-map (+ n step step) step)
               (let [final-value (if (= n step) (+ n step step) n)]
                 (assoc composite-map final-value step))))
           (next-prime
             [n composite-map]
             (let [prime-factor (composite-map n)
                   prime? (nil? prime-factor)
                   step (if prime-factor prime-factor n) 
                   composite-map (next-composite composite-map n step)
                   composite-map (dissoc composite-map n)]
               (if prime?
                 (lazy-seq (cons n (next-prime (+ n 2) composite-map)))
                 (recur (+ n 2) composite-map))))]
      (lazy-seq (cons 2 (next-prime 3 {})))))

The algorithm gets all the primes under 1,000,000 in 1.55 seconds. So this literally cut its time in half, as it now only has to make half the initial calls to `next-prime`.

I guess it's still not as fast as Christophe Grand's implementation (1.5 seconds), although the logic is the same, I have one extra branch and I will figure that out soon.
It's been quite a learning experience for me and it was great to figure it out, debugging trying to get to the solution, and backwards when writing this blog post
while not having my old code to share my old implementations. What I care most about is understanding the concept and writing elegant code that's the best it could have been.
Of course I could use unboxed math until the primes will be auto-promoted to BigInteger when needed, which would cut down time, or try and remove that extra branch by making my
code exactly like the one in the blog, but I'm happy with the general simplicity of the solution which I built only from the logic I understood each time I created an implementation.

Again, check out the [paper by Melissa O'Neill](https://www.cs.hmc.edu/~oneill/papers/Sieve-JFP.pdf) and [Christophe Grand's blog post](http://clj-me.cgrand.net/2009/07/30/everybody-loves-the-sieve-of-eratosthenes/).

EDIT: Each implementation got slightly faster after following [bbatsov's style rule on if-let](https://github.com/bbatsov/clojure-style-guide#if-let) and removing unnecessary bindings.

The Clojure community is extremely passionate and kind when it comes to sharing ideas, and I hope someone finds this post to be of help or interest (:
