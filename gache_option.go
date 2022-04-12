package gache

import "time"

type Option[T addable] func(g *Gache[T])

func OptionCleanup[T addable](interval time.Duration) Option[T] {
	return func(g *Gache[T]) {
		g.cleanupInterval = interval
	}
}

func OptionTTLUnit[T addable](ttlUnit time.Duration) Option[T] {
	return func(g *Gache[T]) {
		g.ttlUnit = ttlUnit
	}
}
