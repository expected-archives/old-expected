package docker

func convertEnv(env map[string]string) []string {
	var converted []string
	for k, v := range env {
		converted = append(converted, k+"="+v)
	}
	return converted
}
