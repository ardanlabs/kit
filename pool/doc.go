// Package pool manages a pool of routines to perform work. It does so my providing
// a Do function that will block when the pool is busy. This also allows the pool
// to monitor and report pushback. The pool also supports the dynamic re-sizing
// of the number of routines in the pool.
//
// Worker
//
//     type Worker interface {
//         Work(ctx context.Context, id int)
//     }
//
// The Worker interface is how you can provide work to the pool. A user-defined type
// implements this interface, then values of that type can be passed into the Do
// function.
//
// Sample Application
//
// https://github.com/ardanlabs/kit/blob/master/examples/pool/main.go
//
package pool
