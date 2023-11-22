package validate

type (

	// Validatable is the interface that must be implemented to support (recursive) validations of structs
	Validatable interface {
		Validate(string, *Validator) *Validator
	}
)

func (v *Validator) SaveContext(ctx string) {
	v.ctxStack = append(v.ctxStack, ctx)
}

func (v *Validator) RestoreContext() {
	n := len(v.ctxStack)
	if n > 0 {
		v.ctxStack = v.ctxStack[:n-1]
	}
}

func (v *Validator) Context() string {
	n := len(v.ctxStack)
	if n == 0 {
		return "root"
	}
	return v.ctxStack[n-1]
}

/*
// Validate starts the chain of validations. Attribute root sets the context for the validations and
// should be passed along when sub structures are validated.
func (v *Validator) Validate(root string, src interface{}) *Validator {
	fmt.Println("push " + root)

	vv := src.(Validatable).Validate(root, v)

	fmt.Println("pop")
	return vv
}
*/
