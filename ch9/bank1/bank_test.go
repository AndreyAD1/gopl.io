// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package main

import (
	"fmt"
	"testing"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		Deposit(200)
		fmt.Println("=", Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
	var withdrawSuccess bool
	go func() {
		withdrawSuccess = Withdraw(301)
		done <- struct{}{}
	}()
	<-done
	if withdrawSuccess != false {
		t.Errorf(
			"Withdraw result is %v, want %v",
			withdrawSuccess,
			false,
		)
	}
	if got, want := Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
	go func() {
		withdrawSuccess = Withdraw(300)
		done <- struct{}{}
	}()
	<-done
	if withdrawSuccess != true {
		t.Errorf(
			"Withdraw result is %v, want %v",
			withdrawSuccess,
			true,
		)
	}
	if got, want := Balance(), 0; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
