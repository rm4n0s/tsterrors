# tsterrors


This package is not like the other error packages. For example, it will not save the stack trace of the “runtime” package inside an error. It is not like "xerrors".<br/>

This package will create its own types of stack traces that can be easily tested and validated with code.<br/>

Here is a list of things that the "tsterrors" can assist you with:<br/>

- Writing tests that will validate the input based on the stack trace of the error.

- Being able to see the stack trace of an error that has passed through microservices.

- Sysadmins can send the error as JSON back to developers to reproduce it and fix it with unit tests (as shown in the ["how_to_reproduce_error" example](https://github.com/rm4n0s/tsterrors/blob/main/v1/examples/how_to_reproduce_error/main.go)).

- Architect the software in code based on its constraints and limitations as it is explained in [Can't Driven Development](https://rm4n0s.github.io/posts/6-cant-driven-development/).


With this library, you can pass the [Error Handling Challenge](https://rm4n0s.github.io/posts/3-error-handling-challenge/) as shown in the ["error_handling_challenge" example](https://github.com/rm4n0s/tsterrors/blob/main/v1/examples/error_handling_challenge/main.go).<br/>

Lastly, I want to inform you that some of its functions contain panics to restrict the way it is used. It is very opinionated.