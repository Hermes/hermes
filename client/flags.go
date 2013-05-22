package client

import (
	"os"
)

func ValidateFlags(flags []string) bool {
	allFlags := []string{"generate", "update", "pull", "push", "lock", "load", "upgrade"}

	// Validating number of flags
	if len(flags) == 0 || len(flags) >=  3 {
		return false
	}

	// Validating one flag is passed
	count := 0
	for _, flag := range flags {
		for _, pflag := range allFlags {
			if (flag == pflag) {
            	count++
        	}
        	if count > 1 {
        		return false
        	}
		}
    }

    // Validating for specific cases
    switch flags[0] { 
		case "update":
			if len(flags) != 1 {
				return false
			}
		case "upgrade":
			if len(flags) != 1 {
				return false
			}
		case "generate":
			if len(flags) < 2 || len(flags) > 3 {
				return false
			} else if len(flags) == 2 {
				flags = append(flags, "")
			}
		case "pull":
			if len(flags) != 2 {
				return false
			}
		case "push":
			if len(flags) != 2 {
				return false
			} else if _, err := os.Stat(flags[1]); os.IsNotExist(err) {
				return false
			}
		case "lock":
			if len(flags) != 1 {
				return false
			}
		case "load":
			if len(flags) < 2 || len(flags) > 3 {
				return false
			} else if len(flags) == 2 {
				flags = append(flags, "")
			}
		default: return false
	}
	return true
}