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
extern mrb_value go_mrb_func_call(mrb_state*, mrb_value);

static mrb_value poop(mrb_state *s, mrb_value self) {
    mrb_value result = go_mrb_func_call(s, self);
    printf("BANG: %d\n", mrb_type(result));
    return result;
}


// This method is used as a way to get a valid mrb_func_t that actually
// just calls back into Go.
static inline mrb_func_t _go_mrb_func_call() {
    //return &go_mrb_func_call;
    return &poop;
}

//-------------------------------------------------------------------
// Functions below here expose defines or inline functions that were
// otherwise inaccessible to Go directly.
//-------------------------------------------------------------------

static inline mrb_aspec _go_MRB_ARGS_ANY() {
    return MRB_ARGS_ANY();
}

static inline enum mrb_vtype _go_mrb_type(mrb_value o) {
    return mrb_type(o);
}
