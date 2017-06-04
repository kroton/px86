package main

type register int

const (
	eax register = iota
	ecx
	edx
	ebx
	esp
	ebp
	esi
	edi
	registersSize int = iota
)

func (r register) String() string {
	switch r {
	case eax:
		return "EAX"
	case ecx:
		return "ECX"
	case edx:
		return "EDX"
	case ebx:
		return "EBX"
	case esp:
		return "ESP"
	case ebp:
		return "EBP"
	case esi:
		return "ESI"
	case edi:
		return "EDI"
	}
	return ""
}
