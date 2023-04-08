package consts

type MirrorStatusFlag rune

const (
	ReverseCachedProxy = MirrorStatusFlag('C')
	ReverseProxy       = MirrorStatusFlag('R')
	Unknown            = MirrorStatusFlag('U')
)
