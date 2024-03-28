package env

import (
	"errors"
	"github.com/93ams/gmd/util"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
)

const (
	current = "current"
)

type Env struct {
	name string
	*viper.Viper
}
type Envs map[string]Env
type Store struct {
	*viper.Viper
	key string
	Envs
}

var ErrRemSelected = errors.New("cannot remove selected environment")
var homeDir = lo.Must(user.Current()).HomeDir

func New(key string, opts ...func(Env)) Env { return util.Apply(Env{key, viper.New()}, opts) }
func (e Env) Name() string                  { return e.name }
func NewStore(key string) Store             { return Store{viper.New(), key, Envs{"": New("")}} }
func (s Store) Open() error {
	if err := util.EnsureFile(s.cfgPath()); err != nil {
		return err
	} else if err := util.EnsureDir(s.envsPath()); err != nil {
		return nil
	} else if dir, err := os.ReadDir(s.envsPath()); err != nil {
		return err
	} else if err := s.Load(lo.Map(dir, func(item os.DirEntry, _ int) string {
		ext := filepath.Ext(item.Name())
		if ext == "" || !slices.Contains(viper.SupportedExts, ext[1:]) {
			return ""
		}
		return strings.TrimSuffix(item.Name(), ext)
	})...); err != nil {
		return err
	}
	s.SetConfigFile(s.cfgPath())
	return s.ReadInConfig()
}
func (s Store) Close() error { return nil }

func (s Store) Keys() []string {
	return lo.Map(lo.Keys(s.Envs), func(i string, _ int) string { return i })
}
func (s Store) Rem(envs ...string) error {
	return s.forKeys(envs, func(v string, e Env) error {
		if e.Viper == nil {
			return nil
		} else if v == s.GetString(current) {
			return ErrRemSelected
		}
		delete(s.Envs, v)
		return s.removeFile(v)
	})
}
func (s Store) Add(envs ...string) error {
	return s.forKeys(envs, func(v string, e Env) error {
		s.setFile(v)
		return s.ensureFile(v)
	})
}
func (s Store) Save(envs ...string) error {
	return s.forKeys(envs, func(v string, e Env) error { return s.setFile(v).WriteConfig() })
}
func (s Store) Load(envs ...string) error {
	return s.forKeys(envs, func(v string, e Env) error {
		return s.setFile(v).ReadInConfig()
	})
}
func (s Store) Current() string { return s.GetString(current) }
func (s Store) GetCurrent() Env { return s.Envs[s.Current()] }
func (s Store) Select(key string) (Env, error) {
	if key == s.GetString(current) {
		return s.GetCurrent(), nil
	}
	s.Set(current, key)
	if key != "" {
		if err := s.Load(key); err != nil {
			return Env{}, err
		}
	}
	return s.Envs[key], s.WriteConfig()
}
func (s Store) forKeys(keys []string, fn func(string, Env) error) error {
	var errs []error
	for _, v := range keys {
		if v == "" {
			continue
		} else if err := fn(v, s.Envs[v]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
func (s Store) cfgPath() string           { return filepath.Join(homeDir, "."+s.key, "config.yml") }
func (s Store) envsPath() string          { return filepath.Join(homeDir, "."+s.key, "envs") }
func (s Store) envPath(v string) string   { return filepath.Join(s.envsPath(), v+".yml") }
func (s Store) removeFile(v string) error { return util.DesureFile(s.envPath(v)) }
func (s Store) ensureFile(v string) error { return util.EnsureFile(s.envPath(v)) }
func (s Store) setFile(v string) Env {
	e := s.Envs[v]
	if e.Viper == nil {
		e = New(v)
		e.SetConfigFile(s.envPath(v))
		s.Envs[v] = e
	}
	return e
}

func WithPrefix(p string) func(Env) { return func(v Env) { v.SetEnvPrefix(p) } }
func Merge(v *viper.Viper) func(Env) {
	return func(v Env) { lo.Must0(v.MergeConfigMap(v.AllSettings())) }
}
