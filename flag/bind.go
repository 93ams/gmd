package flag

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

type Binder struct {
	AutoEnv     bool
	EnvPrefix   string
	EnvReplacer viper.StringReplacer
}

func NewBinder() Binder {
	return Binder{}
}
func (b Binder) Flags(c ...any) func(*cobra.Command) {
	return func(cmd *cobra.Command) {
		for _, c := range c {
			for _, req := range b.scan(c, cmd.Flags()) {
				lo.Must0(cmd.MarkFlagRequired(req))
			}
		}
	}
}
func (b Binder) PFlags(c ...any) func(*cobra.Command) {
	return func(cmd *cobra.Command) {
		for _, c := range c {
			for _, req := range b.scan(c, cmd.PersistentFlags()) {
				lo.Must0(cmd.MarkPersistentFlagRequired(req))
			}
		}
	}
}
func (b Binder) scan(c any, fs *pflag.FlagSet) []string {
	t := reflect.TypeOf(c).Elem()
	val := reflect.ValueOf(c).Elem()
	var required []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := strings.ToLower(f.Name)
		flag, _ := f.Tag.Lookup("flag")
		if flag != "" {
			parts := strings.Split(flag, ",")
			name = parts[0]
		}
		fv := val.Field(i)
		fieldVal := flagPtr(fv)
		usage, _ := f.Tag.Lookup("usage")
		alias, _ := f.Tag.Lookup("alias")
		attrs, _ := f.Tag.Lookup("attrs")
		if strings.Contains(attrs, "required") {
			required = append(required, name)
		}
		switch f.Type.Kind() {
		case reflect.Uint8:
			addFlag(name, alias, fieldVal, usage, fs.Uint8Var, fs.Uint8VarP)
		case reflect.Uint16:
			addFlag(name, alias, fieldVal, usage, fs.Uint16Var, fs.Uint16VarP)
		case reflect.Uint32:
			addFlag(name, alias, fieldVal, usage, fs.Uint32Var, fs.Uint32VarP)
		case reflect.Uint64:
			addFlag(name, alias, fieldVal, usage, fs.Uint64Var, fs.Uint64VarP)
		case reflect.Uint:
			addFlag(name, alias, fieldVal, usage, fs.UintVar, fs.UintVarP)
		case reflect.Int8:
			addFlag(name, alias, fieldVal, usage, fs.Int8Var, fs.Int8VarP)
		case reflect.Int16:
			addFlag(name, alias, fieldVal, usage, fs.Int16Var, fs.Int16VarP)
		case reflect.Int32:
			addFlag(name, alias, fieldVal, usage, fs.Int32Var, fs.Int32VarP)
		case reflect.Int64:
			addFlag(name, alias, fieldVal, usage, fs.Int64Var, fs.Int64VarP)
		case reflect.Int:
			addFlag(name, alias, fieldVal, usage, fs.IntVar, fs.IntVarP)
		case reflect.String:
			addFlag(name, alias, fieldVal, usage, fs.StringVar, fs.StringVarP)
		case reflect.Bool:
			addFlag(name, alias, fieldVal, usage, fs.BoolVar, fs.BoolVarP)
		case reflect.Slice:
		case reflect.Array:
		case reflect.Map:
		case reflect.Ptr:
		case reflect.Interface:
		case reflect.Struct:
			required = append(required, b.scan(fieldVal, fs)...)
		default:
			log.Println("unsupported type", f.Type.Kind())
		}
	}
	return required
}

func BindSet(fs *pflag.FlagSet, v *viper.Viper) {
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			lo.Must0(fs.Set(f.Name, fmt.Sprintf("%v", v.Get(f.Name))))
		}
	})
}
func addFlag[T any](name, alias string, val any, desc string, fn func(*T, string, T, string), fnp func(*T, string, string, T, string)) {
	if alias != "" {
		fnp(val.(*T), name, alias, *val.(*T), desc)
	} else {
		fn(val.(*T), name, *val.(*T), desc)
	}
}

func flagPtr(fv reflect.Value) any {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv = reflect.New(fv.Type())
		}
		return fv.Interface()
	}
	return fv.Addr().Interface()
}
