package main

import (
        &#34;fmt&#34;
        &#34;math&#34;

        )

// a list of all the variables
var input       = []float64     { 0.0 , 0.0 , 0.0 , 0.08, 0.06, 0.29, 0.05, 0.0 ,-1.0 ,28.1109445277361, 9.37031484257871,
                                  0.0 , 0.5 , 0.5 , 0.25, 0.04, 0.0 , 0.25, 0.3 , 0.0 , 0.66, 0.33,      4.68515742128936,
                                  0.4 , 0.8 , 9.37031484257871, 0.8 , 0.25, 0.8 , 0.95, 0.1 };

// evaluate the likelihood of a nesting or attempt, given the input
func evaluate( values []float64 ) float64 {

        // fill in the pre-calculated value first
        values[8]       = values[5] * values[6] * values[3]*(1.0 - values[3])

        // now begin the calculations
        result := 1.0

        result *= values[3]             // Douglas Hugh
        result *= values[3]             // Pauline Gray
        result *= values[3]             // Ali Smyth
                                        
                                        // maryann
        result *= (values[8] * values[9]) / ( (values[8] * values[9]) &#43; (values[4] * values[17]) )
                                        
                                        // daufnie_odie
        result *= 1.0 - ( values[4] * values[18] * (1.0 - values[3]) ) / (
                        ( values[4] * values[18] * (1.0 - values[3]) ) &#43; (values[8] * values[10] * values[12]) &#43;
                        (( 1.0 - (values[8] * values[10]) ) * values[8] * values[9]) );
                                        
                                        // Biff Jag
        result *= 1.0 - ( values[4] * values[20] * (1.0 - values[3]) ) / (
                        ( values[4] * values[20] * (1.0 - values[3]) ) &#43; (values[8] * values[10] * values[15]) &#43;
                        (( 1.0 - (values[8] * values[10]) ) * values[8] *values[10]) );
                                        
                                        // Puppy/Tompsin
        result *= 1.0 - ( values[4] * values[21] * (1.0 - values[3]) ) / (
                        ( values[4] * values[21] * (1.0 - values[3]) ) &#43; (values[8] * values[10] * values[14]) &#43;
                        (( 1.0 - (values[8] * values[10]) ) * values[8] * values[22]) );
                                        
                                        // Puppy/Myerson
        result *= 1.0 - ( values[4] * values[23] * values[24] * (1.0 - values[3]) ) / (
                        ( values[4] * values[23] * values[24] * (1.0 - values[3]) ) &#43;
                        (values[8] * values[10] * values[14]) &#43; (( 1.0 - (values[8]*values[10]) ) * values[8] * values[25]) );
                                        
                                        // unnamed via Myerson
        result *= 1.0 - ( values[4] * values[26] * (1.0 - values[3]) ) / (
                        ( values[4] * values[26] * (1.0 - values[3]) ) &#43; (values[8] * values[10] * values[13]) &#43;
                        (( 1.0 - (values[8] * values[10]) ) * values[8] * values[9]) );
                                        
                                        // skippingthem
        result *= 1.0 - ( values[4] * values[18] * (1.0 - values[3]) ) / (
                        ( values[4] * values[18] * (1.0 - values[3]) ) &#43; (values[8] * values[10] * values[12]) &#43;
                        (( 1.0 - (values[8]*values[10]) ) * values[8] * values[9]) );
                                        
                                        // Grandie
        result *= ( (values[8] * values[10]) &#43; values[27] ) / (
                        (values[4] * values[18] * values[28] * (1.0 - values[3])) &#43; (values[8]*values[10]&#43;values[27]) );

                                        // disturbed nest
        result *= ( (1.0 - values[4]) * values[30] ) / (
                        (( 1.0 - values[4]) * values[30]) &#43; (values[4] * values[29]) );

        return (1.0 - result)
        }

func main() {
        
        var slopes [31]float64

        fmt.Printf( &#34;Confidence in a nest or attempted nest at 84744 M.S.\n&#34; )
        fmt.Printf( &#34;Baseline:\t%.16f\tLog (1-Baseline):\t%.16f\n&#34;, evaluate(input), math.Log(1.0-evaluate(input)) )
        fmt.Printf( &#34;\n&#34; )

        fmt.Printf( &#34;Weightings:\n&#34; )
        fmt.Printf( &#34;Row&#34; )

        maxVal := 0.0
        for i, val := range input {

                if val &lt;= 0.0 { continue }      // skip blank values
                if i == 8 { continue }          // skip calculated values

                fmt.Printf( &#34;\tE%d&#34;, i )

                trial           := input
                trial[i]        = val * 1.01
                high            := evaluate( trial )

                trial[i]        = val * 0.99
                low             := evaluate( trial )

                slopes[i]       = (low - high) / (val * (1.01 - 0.99))

                if math.Abs(slopes[i]) &gt; maxVal { maxVal = math.Abs(slopes[i]) }

                trial[i]        = val
                }

        fmt.Printf( &#34;\n&#34; )
        fmt.Printf( &#34;Slope (raw)&#34; )

        for i, val := range input {

                if val &lt;= 0.0 { continue }      // skip blank values
                if i == 8 { continue }          // skip calculated values
                fmt.Printf( &#34;\t%.10f&#34;, slopes[i] )
                }

        fmt.Printf( &#34;\n&#34; )
        fmt.Printf( &#34;Slope (normed)&#34; )

        for i, val := range input {

                if val &lt;= 0.0 { continue }      // skip blank values
                if i == 8 { continue }          // skip calculated values
                fmt.Printf( &#34;\t%f&#34;, slopes[i] / maxVal )
                }

        fmt.Printf( &#34;\n&#34; )
        fmt.Printf( &#34;\n&#34; )

        fmt.Printf( &#34;Greedy variable reduction:\n&#34; )
        fmt.Printf( &#34;Row&#34; )
        for i, val := range input {

                if val &lt;= 0.0 { continue }      // skip blank values
                if i == 8 { continue }          // skip calculated values

                fmt.Printf( &#34;\tE%d&#34;, i )
                }

        fmt.Printf( &#34;\n&#34; )
                                                // when do we print values? How much do we step?
        cutoffs := []float64 { 0.999, 0.99, 0.95, 0.9, 0.8, 0.75, 0.666, 0.5 }
        step    := 0.0001

        bestResult      := evaluate(input)
        altered         := make( []float64, 31 )        // manually copy by value, just in case
        for i := range input {  altered[i]      = input[i] }

        for cutIndex    := range cutoffs {

                for bestResult &gt; cutoffs[cutIndex] {

                        bestIndex       := -1
                        bestVar         := -1.0

                        for i, val := range altered {

                                if val &lt;= 0.0 { continue }      // skip blank values
                                if i == 8 { continue }          // skip calculated values

                                if slopes[i] &lt; 0 {      altered[i]      *= (1 - step)
                                } else          {       altered[i]      *= (1 &#43; step) }

                                if evaluate( altered ) &lt; bestResult {

                                        bestResult      = evaluate( altered )
                                        bestIndex       = i
                                        bestVar         = altered[i]
                                        }

                                altered[i]      = val           // put this back, in case of copy by ref
                                }

                        altered[bestIndex]      = bestVar       // only once sure, switch in the new value

                        } // while (we haven&#39;t hit the cutoff)

                fmt.Printf( &#34;%.3f&#34;, cutoffs[cutIndex] )
                for i, val := range altered {

                        if val &lt;= 0.0 { continue }      // skip blank values
                        if i == 8 { continue }          // skip calculated values

                        fmt.Printf( &#34;\t%.3f&#34;, val )
                        }
                fmt.Printf( &#34;\n&#34; )

                }       // for (each cutoff)

        fmt.Printf( &#34;\n&#34; )
        fmt.Printf( &#34;All-variable reduction:\n&#34; )
        fmt.Printf( &#34;Row&#34; )
        for i, val := range input {

                if val &lt;= 0.0 { continue }      // skip blank values
                if i == 8 { continue }          // skip calculated values

                fmt.Printf( &#34;\tE%d&#34;, i )
                }

        fmt.Printf( &#34;\n&#34; )

        bestResult      = evaluate(input)	// manually copy by value, just in case
        for i := range input {  altered[i]      = input[i] }

        for cutIndex    := range cutoffs {

                for bestResult &gt; cutoffs[cutIndex] {

						// for each variable, reduce it by the step times the slope
                        for i, val := range altered {

                                if val &lt;= 0.0 { continue }      // skip blank values
                                if i == 8 { continue }          // skip calculated values

				altered[i]	*= 1.0 &#43; slopes[i]/maxVal*step
				}	
						// recalculate the best result
			newResult	:= evaluate( altered )
                        if newResult &lt; bestResult { bestResult = newResult }

						// recalculate the slopes
			maxVal	= 0.0
			for i, val := range altered {

			        if val &lt;= 0.0 { continue }      // skip blank values
			        if i == 8 { continue }          // skip calculated values

		                trial           := altered
		                trial[i]        = val * 1.01
		                high            := evaluate( altered )

		                trial[i]        = val * 0.99
		                low             := evaluate( trial )

		                slopes[i]       = (low - high) / (val * (1.01 - 0.99))

		                if math.Abs(slopes[i]) &gt; maxVal { maxVal = math.Abs(slopes[i]) }

                		trial[i]        = val
				}
				
			} // for (less than cutoff)
			
                fmt.Printf( &#34;%.4f&#34;, cutoffs[cutIndex] )
                for i, val := range altered {

                        if val &lt;= 0.0 { continue }      // skip blank values
                        if i == 8 { continue }          // skip calculated values

                        fmt.Printf( &#34;\t%.4f&#34;, val )
                        }
                fmt.Printf( &#34;\n&#34; )

                }       // for (each cutoff)

        } // main