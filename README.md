# Gonads

Gonads is an experimental library/package for implementing monad like data structures in Go. 

_Be warned, many Gophers will likely find this code non-idiomatic and maybe even a blight on the Go community! This library was inspired by my experience in other languages, Java/Kotlin/C# and I've been getting into Rust. In particular, I'm a big fan of Option[T] and Result[T, E] types from Rust. Because of the limitations of the Go implementation of generics and the weak type system we can't mimic Rust exactly, but this an attempt to at least build something useful._