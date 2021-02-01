# N Nurses Problem

This problem is one of the more complicated examples, because it is closely related to real world employee scheduling problems. Also there is not really a "standard" problem. Here we are looking at the following problem setup:

- `n` is the number of nurses
- `d` is the number of days
- `s` is the number of shifts per day
- `shiftRequests` is a the dimensional array `[d, s, n]` with ones where nurse `n` requested to do shift `s` on day `d`
- `offRequests` is a two dimensional array `[d, n]` with ones where nurse `n` requested to have day `d` off
- `schedule` is a three dimensional array `[d, s, n]` with ones, if nurse `n` was assigned to shift `s` on day `d`

We are trying to minimize `alpha * not_respected_days_off - beta * respected_shift_requests`. In words: fullfill as many shift requests and as many off requests as possible.

Constrains:

- a nurse will have 1 day off after doing a night shift unless she is doing another night shift
- there is a maximum number of night shifts for the period
- a nurse can only do one shift per day
