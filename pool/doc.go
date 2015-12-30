// Package pool manages a pool of routines to perform work. It does so my providing
// a Do function that will block when the pool is busy. This also allows the pool
// to monitor and report pushback. The pool also supports the dynamic re-sizing
// of the number of routines in the pool.
//
// Worker
//
//     type Worker interface {
//         Work(context interface{}, id int)
//     }
//
// The Worker interface is how you can provide work to the pool. A user-defined type
// implements this interface, then values of that type can be passed into the Do
// function.
//
// Sample Application
//
// The following is a sample application using the work pool.
//
//     // theWork is the customer work type for using the pool.
//     type theWork struct{}
//
//     // Work implements the DoWorker interface.
//     func (*theWork) Work(context string, id int) {
//         fmt.Printf("%s : Performing Work\n", context)
//     }
//
//     // ExampleNewDoPool provides a basic example for using a DoPool.
//     func ExampleNewDoPool() {
//         // Create a new do pool.
//         p, err := pool.New(context, "TheWork", 3, func() time.Duration { return time.Minute })
//         if err != nil {
//             fmt.Println(err)
//             return
//         }
//
//         // Pass in some work to be performed.
//         p.Do("TEST", &theWork{})
//         p.Do("TEST", &theWork{})
//         p.Do("TEST", &theWork{})
//
//         // Wait to the work to be processed.
//         time.Sleep(1 * time.Second)
//
//         // Shutdown the pool.
//         p.Shutdown(context)
//     }
package pool
