// Perform Bayesian hypothesis testing over a simple Binomial distribution
// Tested hypothesis: The success rate was random.


package main

import (
        &#34;fmt&#34;
        &#34;math&#34;
)



func main() {
        
        samples := 3000000.0                                    // number of samples to try
        count := float64(0.0)

        trials  := [16]float64{5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400,5400}
        picked  := [16]float64{2676,2759,2668,2618,2729,2766,2716,2695,2649,2663,2739,2723,2694,2692,2658,2709} // successes

        base := [16]float64{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}    // storage for baseline values
        avg := base						// precalculate them to save time
        for i, val := range trials { base[i] = math.Log(0.5) * val }


        offset := 0.5 * (math.Sqrt(5) - 1.0)                    // swap out random numbers for a quasi-random sequence,
        prob    := 0.5                                          //  which converges faster

        for count &lt; samples {
                
                slot    := int(count) % len(avg);               // cycle through the various tests
                acount  := float64(int(count) / len(avg))
                
                prob &#43;= offset					// generate a new quasi-random
                prob  = prob - float64(int(prob))
                
                sample := (math.Log(prob) * (picked[slot])) &#43; (math.Log(1-prob) * ((trials[slot]) - (picked[slot])))
                avg[slot] = (avg[slot]*acount &#43; math.Exp(sample-base[slot])) / (acount &#43; 1)
                count&#43;&#43;
                }
        
        fmt.Printf(&#34;Samples: %d\n&#34;, int(count))
        fmt.Printf(&#34;Success/Trials\tBayes Factor\t1 / Bayes Factor\n&#34;)

        over := float64(1.0)                                    // the culmulative Bayes Factor
        for i := 0; i &lt; len(avg); i&#43;&#43; {
                
                fmt.Printf(&#34;%4d/%4d\t%f\t%f\n&#34;, int(picked[i]), int(trials[i]), (1 / avg[i]), avg[i])
                over *= avg[i]
                }
        
        fmt.Printf(&#34;overall\t%f\t%f\n&#34;, (1/over), over)		// flip the order around

        }