package main

const (
	eax uint8 = iota
	ecx
	edx
	ebx
	esp
	ebp
	esi
	edi
	registersSize
)

func registerStr(r uint8) string {
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
