package numint

// GaussLegendre returns a GaussLegendre rule of specified degree.
// The rules are based off of hard-coded constants, and are implemented up to n=12.
func GaussLegendre(degree int) *GaussLegendreRule {

	if degree > len(weightValues) || degree < 1 {
		return &GaussLegendreRule{}
	}

	weights := weightValues[degree-1]
	points := xValues[degree-1]

	var abscissae, weightCoeffs []float64

	if degree%2 == 0 {
		abscissae = expandEvenAbsiccae(points)
		weightCoeffs = expandEvenWeights(weights)
	} else {
		abscissae = expandOddAbsiccae(points)
		weightCoeffs = expandOddWeights(weights)
	}

	return &GaussLegendreRule{
		abscissae:    abscissae,
		weightCoeffs: weightCoeffs,
		weights:      make([]float64, len(weightCoeffs)),
		points:       make([]float64, len(abscissae)),
	}
}

// GaussLegendreRule provides Weights and Points functions for Gauss Legendre quadrature rules.
type GaussLegendreRule struct {
	abscissae, weightCoeffs []float64
	weights, points         []float64
}

// Weights returns the quadrature weights to use for the interval [a, b].
// The number of points returned depends on the degree of the rule.
func (g *GaussLegendreRule) Weights(a float64, b float64) []float64 {

	for i := range g.weightCoeffs {
		g.weights[i] = g.weightCoeffs[i] * (b - a) / 2
	}

	return g.weights
}

// Points returns the quadrature sampling points to use for the interval [a, b].
// The number of points returned depends on the degree of the rule.
func (g GaussLegendreRule) Points(a float64, b float64) []float64 {

	for i := range g.abscissae {
		g.points[i] = g.abscissae[i]*(b-a)/2 + (b+a)/2
	}

	return g.points
}

// source: http://www.holoborodko.com/pavel/numerical-methods/numerical-integration/
var xValues = [][]float64{
	{0},                           // n=1
	{0.5773502691896257645091488}, // n=2
	{0.0000000000000000000000000, 0.7745966692414833770358531},                              // n=3
	{0.3399810435848562648026658, 0.8611363115940525752239465},                              // n=4
	{0.0000000000000000000000000, 0.5384693101056830910363144, 0.9061798459386639927976269}, // n=5
	{0.2386191860831969086305017, 0.6612093864662645136613996, 0.9324695142031520278123016},
	{0.0000000000000000000000000, 0.4058451513773971669066064, 0.7415311855993944398638648, 0.9491079123427585245261897},
	{0.1834346424956498049394761, 0.5255324099163289858177390, 0.7966664774136267395915539, 0.9602898564975362316835609},
	{0.0000000000000000000000000, 0.3242534234038089290385380, 0.6133714327005903973087020, 0.8360311073266357942994298, 0.9681602395076260898355762},
	{0.1488743389816312108848260, 0.4333953941292471907992659, 0.6794095682990244062343274, 0.8650633666889845107320967, 0.9739065285171717200779640},
	{0.0000000000000000000000000, 0.2695431559523449723315320, 0.5190961292068118159257257, 0.7301520055740493240934163, 0.8870625997680952990751578, 0.9782286581460569928039380},
	{0.1252334085114689154724414, 0.3678314989981801937526915, 0.5873179542866174472967024, 0.7699026741943046870368938, 0.9041172563704748566784659, 0.9815606342467192506905491},
}

var weightValues = [][]float64{
	{2}, // n=1
	{1}, // n=2
	{0.8888888888888888888888889, 0.5555555555555555555555556},                              // n=3
	{0.6521451548625461426269361, 0.3478548451374538573730639},                              // n=4
	{0.5688888888888888888888889, 0.4786286704993664680412915, 0.2369268850561890875142640}, // n=5
	{0.4679139345726910473898703, 0.3607615730481386075698335, 0.1713244923791703450402961},
	{0.4179591836734693877551020, 0.3818300505051189449503698, 0.2797053914892766679014678, 0.1294849661688696932706114},
	{0.3626837833783619829651504, 0.3137066458778872873379622, 0.2223810344533744705443560, 0.1012285362903762591525314},
	{0.3302393550012597631645251, 0.3123470770400028400686304, 0.2606106964029354623187429, 0.1806481606948574040584720, 0.0812743883615744119718922},
	{0.2955242247147528701738930, 0.2692667193099963550912269, 0.2190863625159820439955349, 0.1494513491505805931457763, 0.0666713443086881375935688},
	{0.2729250867779006307144835, 0.2628045445102466621806889, 0.2331937645919904799185237, 0.1862902109277342514260976, 0.1255803694649046246346943, 0.0556685671161736664827537},
	{0.2491470458134027850005624, 0.2334925365383548087608499, 0.2031674267230659217490645, 0.1600783285433462263346525, 0.1069393259953184309602547, 0.0471753363865118271946160},
}

func expandEvenAbsiccae(abs []float64) []float64 {

	n := 2 * len(abs)
	result := make([]float64, n)

	for i := range result {
		if i < n/2 {
			result[i] = -1 * abs[(n-2)/2-i]
		} else {
			result[i] = abs[i-n/2]
		}
	}

	return result
}

func expandEvenWeights(ws []float64) []float64 {
	n := 2 * len(ws)
	result := make([]float64, n)

	for i := range result {
		if i < n/2 {
			result[i] = ws[(n-2)/2-i]
		} else {
			result[i] = ws[i-n/2]
		}
	}

	return result
}

func expandOddAbsiccae(abs []float64) []float64 {

	n := 2*len(abs) - 1
	result := make([]float64, n)

	for i := range result {
		if i < n/2 {
			result[i] = -1 * abs[(n-1)/2-i]
		} else {
			result[i] = abs[i-(n-1)/2]
		}
	}

	return result
}

func expandOddWeights(weights []float64) []float64 {
	n := 2*len(weights) - 1
	result := make([]float64, n)

	for i := range result {
		if i < n/2 {
			result[i] = weights[(n-1)/2-i]
		} else {
			result[i] = weights[i-(n-1)/2]
		}
	}

	return result
}
