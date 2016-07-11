package main

import (
	breaker &#34;.&#34;
	&#34;fmt&#34;
	&#34;time&#34;
)

// RemoteService is a fictitious interface to an encapsulation of some
// unreliable subsystem that is outside of our domain of control.  It should
// be treated for this example as merely a recording of behavior.
type RemoteService struct {
	// The subsystem&#39;s circuit breaker.
	Breaker breaker.Consecutive
	// The pre-recorded success or failure results that are used to drive
	// the behavior of this example.
	Results []bool
}

// ConductRequest is the supposed interface point for user&#39;s of this subsystem.
// They call it as necessary to perform whatever work to yield the result they
// want.
func (s *RemoteService) ConductRequest() {
	// For purposes of not convoluting the example, we use real sleep operations
	// here.
	time.Sleep(time.Second/2 &#43; time.Second/4)

	// If the circuit is broken, merely bail.
	if s.Breaker.Open() {
		fmt.Println(&#34;Unavailable; Trying Again Later...&#34;)

		return
	}

	// Emulate the actual remote interface here that is supposedly unreliable.
	err := s.performRequest()
	// WARNING: We make an implicit assumption that any err value is retryable
	// and not a permanent error.
	if err != nil {
		fmt.Println(&#34;Operation Failed&#34;)
		s.Breaker.Fail()
	} else {
		fmt.Println(&#34;Operation Succeeded&#34;)
		s.Breaker.Succeed()
	}
}

// performRequest models an interaction with an unreliable external
// system---e.g., a remote API server.
func (s *RemoteService) performRequest() error {
	result := s.Results[0]
	s.Results = s.Results[1:]

	if !result {
		return fmt.Errorf(&#34;Temporary Unavailable&#34;)
	}

	return nil
}

func main() {
	subsystem := &amp;RemoteService{
		Breaker: breaker.Consecutive{
			FailureAllowance: 2,
			RetryTimeout:     time.Second,
		},
		Results: []bool{
			true,  // Success.
			true,  // Success.
			false, // One-off failure; do not trip circuit.
			true,  // This success negates past failures.
			false, // String of contiguous failures to create open circuit.
			false, // Open.  :-(
			false, // Open; however, we&#39;ve timed out.
			true,  // We have a success here.
		},
	}

	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
}
