# TODO

Make the raycast program compile correctly

At a minimum:

* Fix the interpretation of `(asm 16) sub bl, [0x46c]`.
* make `0x9ff6 -> stack` and then `stack -> es` mean the same as `0x9ff6 -> es`. Make the latter possible.
