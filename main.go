package main

import (
  "fmt"
  "math/rand"
  "time"
  "sync"
)

// global state for the size of the matrices
// in C I would use #define but I looked into it and go doesn't appear to have
// a direct equivalent, so this will do
const SIZE = 5

// opting for arrays instead of slices because I do not need dynamically sized
// memory, I always know it will be 100x100 (SIZE x SIZE)
type Matrix [SIZE][SIZE]float32

// main creates two matrices, records the amount of time their multiplication
// takes, and prints that time to stdout
func main() {
  A := construct_matrix()
  B := construct_matrix()
  
  // begin timing
  start := time.Now()
  naive_mult(A, B)
  
  fmt.Printf(
    "%dx%d Matrix Multiplication WITHOUT goroutines accomplished in: %s \n",
    SIZE, 
    SIZE,
    time.Since(start))
  
  // I am going to make new matrices to avoid any caching advantages that 
  // the goroutine code may have over the naive implementation
  A = construct_matrix()
  B = construct_matrix()
  
  // begin timing for optimized multiplication
  start = time.Now()
  optimized_mult(A, B)

  fmt.Printf(
    "%dx%d Matrix Multiplication WITH    goroutines accomplished in: %s \n",
    SIZE, 
    SIZE,
    time.Since(start))
}

// constructs a matrix of SIZE x SIZE filled with floats from 0-1
func construct_matrix() Matrix {
  var M Matrix 
  for i := 0; i < len(M); i++ {
    for j := 0; j < len(M[i]); j++ {
      M[i][j] = rand.Float32() 
    }
  }
  return M
}

// naive matrix multiplication without goroutines
func naive_mult(A Matrix, B Matrix) Matrix {
  var C Matrix
  for i := 0; i < len(A); i++ {
    for j := 0; j < len(B[i]); j++ {
      C[i][j] = 0;
      for k := 0; k < len(B); k++ {
        C[i][j] += A[i][k] * B[k][j]
      }
    }
  }
  return C
}

// matrix multiplication optimized with goroutines 
func optimized_mult(A Matrix, B Matrix) Matrix {
  var C Matrix
  var wg sync.WaitGroup

  for i := 0; i < len(A); i++ {
    // each row of the matrix multiplication will be done in parallel
    // so we add 1 to the wait group per row
    wg.Add(1)
    
    go func(row int) {
      // guarantee goroutine finishes with Done
      defer wg.Done()

      for j := 0; j < len(B[row]); j++ {
        C[row][j] = 0;
        for k := 0; k < len(B); k++ {
          C[row][j] += A[row][k] * B[k][j]
        }
      }
    }(i)
  }
  
  // wait for all goroutines to be Done
  wg.Wait()
  return C
}

// for testing purposes, this function takes a matrix and prints out each element
// all of the brackets are to make it easy to copy-paste in wulfram alpha to
// see if my matrix multiplications actually work 
// it did not occur to me that there might be better formatting already
// available in go 
func print_matrix(M Matrix) {
  for i := 0; i < len(M); i++ {
    if i == 0 { 
      fmt.Printf("[[")
    } else {
      fmt.Printf("\n[")
    }
    for j := 0; j < len(M[i]); j++ {
      fmt.Printf("%f", M[i][j])
      if j != len(M[i])-1 {
        fmt.Printf(", ")
      }
    }
    fmt.Printf("]")
    if i != len(M)-1 {
      fmt.Printf(",")
    }
  }
  fmt.Printf("]\n")
}
