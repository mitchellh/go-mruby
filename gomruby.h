/*
 * This header exists to simplify the headers that are included within
 * the Go files. This header should include all the necessary headers
 * for the compilation of the Go library.
 * */

#include <mruby.h>
#include <mruby/class.h>
#include <mruby/compile.h>
#include <mruby/proc.h>
#include <mruby/string.h>
#include <mruby/value.h>

// This is declard in func.go and is a way for us to call back into
// Go to execute a method.
extern mrb_value go_mrb_func_call(mrb_state*, mrb_value*);

// This calls into Go with a similar signature to mrb_func_t. We have to
// change it slightly because cgo can't handle the union type of mrb_value,
// so we pass in a pointer instead.
static inline mrb_value _go_mrb_func_call(mrb_state *s, mrb_value self) {
    return go_mrb_func_call(s, &self);
}

// This method is used as a way to get a valid mrb_func_t that actually
// just calls back into Go.
static inline mrb_func_t _go_mrb_func_t() {
    return &_go_mrb_func_call;
}

//-------------------------------------------------------------------
// Functions below here expose defines or inline functions that were
// otherwise inaccessible to Go directly.
//-------------------------------------------------------------------

static inline mrb_aspec _go_MRB_ARGS_ANY() {
    return MRB_ARGS_ANY();
}

static inline mrb_aspec _go_MRB_ARGS_ARG(int r, int o) {
    return MRB_ARGS_ARG(r, o);
}

static inline mrb_aspec _go_MRB_ARGS_BLOCK() {
    return MRB_ARGS_BLOCK();
}

static inline mrb_aspec _go_MRB_ARGS_NONE() {
    return MRB_ARGS_NONE();
}

static inline mrb_aspec _go_MRB_ARGS_OPT(int n) {
    return MRB_ARGS_OPT(n);
}

static inline mrb_aspec _go_MRB_ARGS_REQ(int n) {
    return MRB_ARGS_REQ(n);
}

static inline int _go_mrb_fixnum(mrb_value o) {
    return mrb_fixnum(o);
}

static inline enum mrb_vtype _go_mrb_type(mrb_value o) {
    return mrb_type(o);
}
