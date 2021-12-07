	/**
		Still not sure why this is not working :/


		(x1-t)^2+(x2-t)^2+...+(xn-t)^2 = 0
		(x1^2 + ... + xn^2) - 2*(sum(x))*t + t^2 = 0

		a = 1
		b = - 2 * sum(x)
		c = (x1^2 + ... + xn^2)

		a*t^2 + b*t + c = 0

		t = (-b +/- sqrt(b^2 - 4*a*c))/2*a
	**/

	// a := 1.0
	// b := 0.0
	// c := 0.0

	// for i := 0; i < len(positions); i++ {
	// 	b += -2.0 * float64(positions[i])
	// 	c += float64(positions[i] * positions[i])
	// }

	// min_move := math.MaxFloat64
	// target := 0.0
	// for t := 0.0; t < 100.0; t += 1.0 {
	// 	move := t*t + b*t + c
	// 	if move < min_move {
	// 		target = t
	// 		min_move = move
	// 		fmt.Printf("t = %v, move = %v\n", target, move)
	// 		fmt.Printf("actual move = %v\n\n", compute_move(positions, int(target)))
	// 	}
	// }

	// t := (-float64(b) - math.Sqrt(b*b-4*a*c)) / (2 * a)

	