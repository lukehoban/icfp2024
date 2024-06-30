# 3D

## 3d1
```
# Input
  `1 <= A <= 100`

# Output
  The factorial of `A`.
  `A! = 1 * 2 * .. * (A - 1) * A`

# Example
  * `A = 5`
    `Answer = 120`
```

```
loop [a, acc] {
  if a == 0 {
    return acc
  }
  loop(a-1, acc*a)
}(A, 1)

```

Solution?:




```
. . . . 1 . 0 .  
. A > . - . # .
. v 1 . . . . .
. . * . . . @ .
. . . . . . 3 .
```

Seems close-ish?
```
. . . . . . 0 . . . .
. A > . > . = . . . .
. v . . . . . > . . .
. . > . > . . . + S .
. v 1 . . v . ^ . . .
. . - . . . . 1 . . .
. . . . . v . v . . .
. . . . . . . . . . .
. . . . . . . * . . .
. 1 @ 6 . . . . . . .
. . 3 . . . . 0 @ 3 .
. . . . . . . 3 . . .
```



## 3d2

```
# Input
  `-100 <= A <= 100`

# Output
  The absolute value of `A`.

# Example
  * `A = 3`
    `Answer = 3`
  * `A = -6`
    `Answer = 6`
```

## 3d3
```
# Input
  `-100 <= A <= 100`

# Output
  The sign of `A`, which is `-1` for negative numbers, `0` for `0`, and `1` for positive numbers.

# Example
  * `A = 3`
    `Answer = 1`
  * `A = -6`
    `Answer = -1`
```

```
if x == 0 { return 0 }
a := (x + 1) / x
b := (x - 1) / x
return a-b
```

```
solve 3d3
. 0 . . . . . . . . . .
. = S . . . . . . . . .
^ . . . . . . . . . . .
A > . > . . . . . . . .
v . v . v . . . . . . .
. . . . . . . . . . . .
v 1 + . / . . . . . . .
. . . . . > . > . . . .
v . . . . . . . v . . .
. > . > . > . . . . . .
. . v 1 . . v . v . . .
. . . - . . . . . . . .
. . . . > . / . - . > .
. . . . . . . . . . 0 -
. . . . . . . . . . . S
```
