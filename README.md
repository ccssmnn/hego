# hego

hego aims to be an efficient, zero dependency library for metaheuristic algorithms in Go. The advantage of Go is the combination of understandable syntax, a static type system and efficiency.

Most real world optimization problems have high costs when it comes to objective function evaluation. This package allows you to define your objectives in Go. Moreover Go's approach to concurrency makes it easy to take advantages of higher cpu counts.

Finding the right algorithm and parameters for your problem is not trivial. hego aims to provide a rich set of helper functions to parametrize the algorithms for your needs as well as allowing you to provide your own functions.

## Algorithms

- Simulated Annealing
  - [x] Binary
  - [ ] Integer
  - [ ] Continuous
- Genetic Algorithms
  - [ ] Binary
  - [ ] Integer
  - [ ] Continuous
- Ant Colony Optimization
  - [ ] Integer
