package env

import "github.com/spf13/viper"

type Builder struct {
	StringReplacer viper.StringReplacer
	ListDelim      string
}

func (m *Builder) New(key string) Env {
	var vops []viper.Option
	if m.StringReplacer != nil {
		vops = append(vops, viper.EnvKeyReplacer(m.StringReplacer))
	}
	if m.ListDelim != "" {
		vops = append(vops, viper.KeyDelimiter(m.ListDelim))
	}
	var v *viper.Viper
	if len(vops) > 0 {
		v = viper.NewWithOptions(vops...)
	} else {
		v = viper.New()
	}
	return Env{key, v}
}
