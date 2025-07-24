package enums

const (
	UNUSED  int = 0
	USED    int = 1
	EXPIRED int = 2
)

func String(e int) string {
	name := []string{
		"UNUSED",
		"USED",
		"EXPIRED",
	}

	if e < UNUSED || e > EXPIRED {
		return "Unknown"
	}

	return name[e]
}
