
# ICFP 2024

Resources:
* https://icfpcontest2024.github.io/ 
* https://icfpcontest2024.github.io/task.html
* https://icfpcontest2024.github.io/scoreboard.html

This is work in progress for team __Team Cup&lt;T&gt;__.  

## Notes

The initial hint from the task:
```
Send: "S'%4}).$%8
Receive: SB%,,/}!.$}7%,#/-%}4/}4(%}M#(//,}/&}4(%}</5.$}P!2)!",%_~~<%&/2%}4!+).'}!}#/523%j}7%}35''%34}4(!4}9/5}(!6%}!},//+}!2/5.$l}S/5e2%}./7},//+).'}!4}4(%}u).$%8wl}N/}02!#4)#%}9/52}#/--5.)#!4)/.}3+),,3j}9/5}#!.}53%}/52}u%#(/w}3%26)#%l}@524(%2-/2%j}4/}+./7}(/7}9/5}!.$}/4(%2}345$%.43}!2%}$/).'j}9/5}#!.},//+}!4}4(%}u3#/2%"/!2$wl~~;&4%2},//+).'}!2/5.$j}9/5}-!9}"%}!$-)44%$}4/}9/52}&)234}#/523%3j}3/}-!+%}352%}4/}#(%#+}4()3}0!'%}&2/-}4)-%}4/}4)-%l}C.}4(%}-%!.4)-%j})&}9/5}7!.4}4/}02!#4)#%}-/2%}!$6!.#%$}#/--5.)#!4)/.}3+),,3j}9/5}-!9}!,3/}4!+%}/52}u,!.'5!'%y4%34wl~
```

That translates to sending `get index` and getting back:

```
% go run main.go "get index"                                     
Hello and welcome to the School of the Bound Variable!

Before taking a course, we suggest that you have a look around. You're now looking at the [index]. To practice your communication skills, you can use our [echo] service. Furthermore, to know how you and other students are doing, you can look at the [scoreboard].

Once you are ready, please progress to one of the courses that you are currently enrolled in:

 * [lambdaman]
 * [spaceship]
 * [3d]

After passing some tests, you may be admitted to other courses, so make sure to check this page from time to time. In the meantime, if you want to practice more advanced communication skills, you may also take our [language_test].
```

You can also send `get XXX` to:
* `get language_test` (runs diagnostics on our evaluator)
* `get lambdaman`
* `get spaceship`
* `get 3d` (was unlocked after making progress on above)

### Language Test

Evaluator was implemented in a fairly hacky way - doing true beta reduction in the term without any optimizations. Probably not ideal, but appears to work. The languge test appears to try every langauge feature and error if they don't work correctly. After a few fixes - got guidace to send `solve language_test 4w3s0m3` which results in:

```
% go run main.go "solve language_test 4w3s0m3"
Correct, you solved hello4!
```

### Spaceship

```
% go run main.go "get spaceship"
```

See the [task](./spaceship/spaceship.md).

Manually solved the first one (answer `31619`), visiting the points in order without real regard for controlling velocity.  Only solves about half.  Will need much better solutions.

Set up harness to pull each test, try to solve, then submit solution if we got one.

```
% go run main.go "get spaceship1"
1 -1
1 -3
2 -5
2 -8
3 -10

% go run main.go "solve spaceship1 31619"
Correct, you solved spaceship1 with a score of 5!
```

### Lambdaman

```
% go run main.go "get lambdaman"
```

See the [task](./lambdaman/lambdaman.md).

### 3d

```
% go run main.go "get 3d"
```

See the [task](./3d/3d.md).
