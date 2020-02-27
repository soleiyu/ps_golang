package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("PS")

	p := MkParticle(10, 30, 0, 0, 9.8, 0.0)
	f1 := false
	f2 := false

	l := MkLine(-5, 10, 15, -10, 0.8)
	l2 := MkLine(20, -15, 50, -20, 0.8)

	fmt.Printf("%4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f\n", 
		p.T, p.Px, p.Py, p.Vx, p.Vy, l.X1, l.Y1)
	fmt.Printf("%4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f\n", 
		p.T, p.Px, p.Py, p.Vx, p.Vy, l.X2, l.Y2)
	fmt.Printf("%4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f\n", 
		p.T, p.Px, p.Py, p.Vx, p.Vy, l2.X1, l2.Y1)
	fmt.Printf("%4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f, %4.4f\n", 
		p.T, p.Px, p.Py, p.Vx, p.Vy, l2.X2, l2.Y2)

	p.Show()
	for i := 0; i < 150; i ++ {
		np := MkZParticle()
		np, f1 = lineHit(l, 0.03125, p, f1)
		if !f1 {
			np, f2 = lineHit(l2, 0.03125, p, f2)
		}

		p = np

		p.Showh()
	}
}

type Particle struct {
	Px, Py, Vx, Vy, G, T float64
}

type Line struct {
	X1, Y1, X2, Y2, K float64
}

func HitLine(l Line, x3, y3, px, py float64) bool {
	vx1 := l.X2 - l.X1
	vy1 := l.Y2 - l.Y1

	vx2 := x3 - l.X2
	vy2 := y3 - l.Y2

	vx3 := l.X1 - x3
	vy3 := l.Y1 - y3

	pvx1 := px - l.X1
	pvy1 := py - l.Y1

	pvx2 := px - l.X2
	pvy2 := py - l.Y2

	pvx3 := px - x3
	pvy3 := py - y3

	return BrokenLine (
		vx1 * pvy1 - pvx1 * vy1,
		vx2 * pvy2 - pvx2 * vy2,
		vx3 * pvy3 - pvx3 * vy3)
}

func BrokenLine(v1, v2, v3 float64) bool {

	if v1 < 0 && 0 < v2 && 0 < v3 {
		return true
	} else if 0 < v1 && v2 < 0 && v3 < 0 {
		return true
	} else {
		return false
	}
}

func lineHit(l Line, t float64, p Particle, flag bool) (Particle, bool) {
	np := Renew(t, p)
	hit := HitLine(l, p.Px, p.Py, np.Px, np.Py) 

	if flag {
	} else if hit{
		np = BinHitLine(l, 8, p, np)

		sita := math.Pi - qsita(l, p, np)
		np.Vx = np.Vx * math.Cos(sita) - np.Vy * math.Sin(sita)
		np.Vy = np.Vx * math.Sin(sita) + np.Vy * math.Cos(sita)
	}

	return np, hit
}

func qsita(l Line, p, np Particle) float64 {
	lvx := l.X2 - l.X1
	lvy := l.Y2 - l.Y1
	_, _, lvs := uv(lvx, lvy)

	pvx := np.Px - p.Px
	pvy := np.Py - p.Py
	_, _, pvs := uv(pvx, pvy)

	inner := lvx * pvy + lvy * pvx
	inner /= lvs
	inner /= pvs

	return math.Acos(inner)
}

func uv (xs, ys float64) (float64, float64, float64) {

	s := math.Sqrt(xs * xs + ys * ys)

	return xs / s, ys / s, s
}

func BinHitLine(l Line, time int, p, np Particle) Particle {
	dt := np.T - p.T
	cp := Renew(0.5 * dt, p)

	if 0 < time {
		if HitLine(l, p.Px, p.Py, np.Px, np.Py) {
			return BinHitLine(l, time - 1, p, cp)
		} else {
			return BinHitLine(l, time - 1, cp, np)
		}
	} else {
		return cp
	}
}


func yukabehit(t, h, v, k float64, p Particle) Particle {
	np := Renew(t, p)

	if np.Py < h {
		np = BinHitYuka(h, 8, p, np)
		np.Vy *= -k
	} else if v < np.Px {
		np = BinHitYuka(v, 8, p, np)
		np.Vx *= -k
	}
	return np
}

func yukahit(t, h, k float64, p Particle) Particle {
	np := Renew(t, p)

	if np.Py < h {
		np = BinHitYuka(h, 4, p, np)
		np.Vy *= k
	}
	return np
}

func BinHitKabe(v float64, time int, p, np Particle) Particle {
	dt := np.T - p.T
	cp := Renew(0.5 * dt, p)

	if 0 < time {
		if v < cp.Px {
			return BinHitYuka(v, time - 1, p, cp)
		} else {
			return BinHitYuka(v, time - 1, cp, np)
		}
	} else {
		return cp
	}
}

func BinHitYuka(h float64, time int, p, np Particle) Particle {
	dt := np.T - p.T
	cp := Renew(0.5 * dt, p)

	if 0 < time {
		if cp.Py < h {
			return BinHitYuka(h, time - 1, p, cp)
		} else {
			return BinHitYuka(h, time - 1, cp, np)
		}
	} else {
		return cp
	}
}

func MkLine(x1, y1, x2, y2, k float64) Line {
	var p Line

	p.X1 = x1
	p.Y1 = y1
	p.X2 = x2
	p.Y2 = y2
	p.K = k

	return p
}

func MkParticle (px, py, vx, vy, g, t float64) Particle {
	var p Particle

	p.Px = px
	p.Py = py
	p.Vx = vx
	p.Vy = vy
	p.G = g
	p.T = t

	return p
}

func MkZParticle () Particle {
	return MkParticle(0, 0, 0, 0, 9.8, 0)
}

func (this *Particle) ShowPos () {
	fmt.Println(this.T, this.Px, this.Py)
}

func (this *Particle) Show () {
	fmt.Println(this.T, this.Px, this.Py, this.Vx, this.Vy)
}

func (this *Particle) Showh () {
	fmt.Printf("%4.2f, %4.2f, %4.2f, %4.2f, %4.2f\n", this.T, this.Px, this.Py, this.Vx, this.Vy)
}

func Renew(t float64, p Particle) Particle {
	res := MkParticle(
		p.Px + p.Vx * t,
		p.Py + p.Vy * t,
		p.Vx,
		p.Vy - p.G * t,
		p.G, p.T + t)

	return res
}

