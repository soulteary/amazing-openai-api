package fn

import (
	"strings"

	"github.com/soulteary/amazing-openai-api/internal/define"
)

func ExtractModelAlias(alias string) define.ModelAlias {
	var result define.ModelAlias
	if alias == "" {
		return result
	}
	pairs := strings.Split(alias, ",")
	for _, pair := range pairs {
		alias := strings.Split(pair, ":")
		if len(alias) != 2 {
			continue
		}
		result = append(result, alias)
	}
	return result
}
