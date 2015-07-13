package mruby

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Decode converts the Ruby value to a Go value.
func Decode(out interface{}, v *MrbValue) error {
	// The out parameter must be a pointer since we must be
	// able to write to it.
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	var d decoder
	return d.decode("root", v, val.Elem())
}

type decoder struct {
	stack []reflect.Kind
}

func (d *decoder) decode(name string, v *MrbValue, result reflect.Value) error {
	k := result

	// If we have an interface with a valid value, we use that
	// for the check.
	if result.Kind() == reflect.Interface {
		elem := result.Elem()
		if elem.IsValid() {
			k = elem
		}
	}

	// Push current onto stack unless it is an interface.
	if k.Kind() != reflect.Interface {
		d.stack = append(d.stack, k.Kind())

		// Schedule a pop
		defer func() {
			d.stack = d.stack[:len(d.stack)-1]
		}()
	}

	switch k.Kind() {
	case reflect.Bool:
		return d.decodeBool(name, v, result)
	case reflect.Float64:
		return d.decodeFloat(name, v, result)
	case reflect.Int:
		return d.decodeInt(name, v, result)
	case reflect.Map:
		return d.decodeMap(name, v, result)
	case reflect.Ptr:
		return d.decodePtr(name, v, result)
	case reflect.Slice:
		return d.decodeSlice(name, v, result)
	case reflect.String:
		return d.decodeString(name, v, result)
	default:
		return fmt.Errorf(
			"%s: unknown kind to decode into: %s", name, k.Kind())
	}

	return nil
}

func (d *decoder) decodeBool(name string, v *MrbValue, result reflect.Value) error {
	switch t := v.Type(); t {
	case TypeFalse:
		result.Set(reflect.ValueOf(false))
	case TypeTrue:
		result.Set(reflect.ValueOf(true))
	default:
		return fmt.Errorf("%s: unknown type %v", name, t)
	}

	return nil
}

func (d *decoder) decodeFloat(name string, v *MrbValue, result reflect.Value) error {
	switch t := v.Type(); t {
	case TypeFloat:
		result.Set(reflect.ValueOf(v.Float()))
	default:
		return fmt.Errorf("%s: unknown type %v", name, t)
	}

	return nil
}

func (d *decoder) decodeInt(name string, v *MrbValue, result reflect.Value) error {
	switch t := v.Type(); t {
	case TypeFixnum:
		result.Set(reflect.ValueOf(v.Fixnum()))
	case TypeString:
		v, err := strconv.ParseInt(v.String(), 0, 0)
		if err != nil {
			return err
		}

		result.SetInt(int64(v))
	default:
		return fmt.Errorf("%s: unknown type %v", name, t)
	}

	return nil
}

func (d *decoder) decodeMap(name string, v *MrbValue, result reflect.Value) error {
	if v.Type() != TypeHash {
		return fmt.Errorf("%s: not a hash type for map (%v)", name, v.Type())
	}

	// If we have an interface, then we can address the interface,
	// but not the slice itself, so get the element but set the interface
	set := result
	if result.Kind() == reflect.Interface {
		result = result.Elem()
	}

	resultType := result.Type()
	resultElemType := resultType.Elem()
	resultKeyType := resultType.Key()
	if resultKeyType.Kind() != reflect.String {
		return fmt.Errorf(
			"%s: map must have string keys", name)
	}

	// Make a map if it is nil
	resultMap := result
	if result.IsNil() {
		resultMap = reflect.MakeMap(
			reflect.MapOf(resultKeyType, resultElemType))
	}

	// We're going to be allocating some garbage, so set the arena
	// so it is cleared properly.
	mrb := v.Mrb()
	defer mrb.ArenaRestore(mrb.ArenaSave())

	// Get the hash of the value
	hash := v.Hash()
	keysRaw, err := hash.Keys()
	if err != nil {
		return err
	}
	keys := keysRaw.Array()

	for i := 0; i < keys.Len(); i++ {
		// Get the key and value in Ruby. This should do no allocations.
		rbKey, err := keys.Get(i)
		if err != nil {
			return err
		}

		rbVal, err := hash.Get(rbKey)
		if err != nil {
			return err
		}

		// Make the field name
		fieldName := fmt.Sprintf("%s.<entry %d>", name, i)

		// Decode the key into the key type
		keyVal := reflect.Indirect(reflect.New(resultKeyType))
		if err := d.decode(fieldName, rbKey, keyVal); err != nil {
			return err
		}

		// Decode the value
		val := reflect.Indirect(reflect.New(resultElemType))
		if err := d.decode(fieldName, rbVal, val); err != nil {
			return err
		}

		// Set the value on the map
		resultMap.SetMapIndex(keyVal, val)

	}

	// Set the final map if we can
	set.Set(resultMap)
	return nil
}

func (d *decoder) decodePtr(name string, v *MrbValue, result reflect.Value) error {
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	resultType := result.Type()
	resultElemType := resultType.Elem()
	val := reflect.New(resultElemType)
	if err := d.decode(name, v, reflect.Indirect(val)); err != nil {
		return err
	}

	result.Set(val)
	return nil
}

func (d *decoder) decodeSlice(name string, v *MrbValue, result reflect.Value) error {
	// If we have an interface, then we can address the interface,
	// but not the slice itself, so get the element but set the interface
	set := result
	if result.Kind() == reflect.Interface {
		result = result.Elem()
	}

	// Create the slice if it isn't nil
	resultType := result.Type()
	resultElemType := resultType.Elem()
	if result.IsNil() {
		resultSliceType := reflect.SliceOf(resultElemType)
		result = reflect.MakeSlice(
			resultSliceType, 0, 0)
	}

	// Get the hash of the value
	array := v.Array()

	for i := 0; i < array.Len(); i++ {
		// Get the key and value in Ruby. This should do no allocations.
		rbVal, err := array.Get(i)
		if err != nil {
			return err
		}

		// Make the field name
		fieldName := fmt.Sprintf("%s[%d]", name, i)

		// Decode the value
		val := reflect.Indirect(reflect.New(resultElemType))
		if err := d.decode(fieldName, rbVal, val); err != nil {
			return err
		}

		// Append it onto the slice
		result = reflect.Append(result, val)
	}

	set.Set(result)
	return nil
}

func (d *decoder) decodeString(name string, v *MrbValue, result reflect.Value) error {
	switch t := v.Type(); t {
	case TypeFixnum:
		result.Set(reflect.ValueOf(
			strconv.FormatInt(int64(v.Fixnum()), 10)).Convert(result.Type()))
	case TypeString:
		result.Set(reflect.ValueOf(v.String()).Convert(result.Type()))
	default:
		return fmt.Errorf("%s: unknown type to string: %v", name, t)
	}

	return nil
}
