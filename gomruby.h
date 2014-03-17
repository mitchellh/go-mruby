/*
 * This header exists to simplify the headers that are included within
 * the Go files. This header should include all the necessary headers
 * for the compilation of the Go library.
 * */

#ifndef _GOMRUBY_H_INCLUDED
#define _GOMRUBY_H_INCLUDED

#include <mruby.h>
#include <mruby/class.h>
#include <mruby/compile.h>
#include <mruby/irep.h>
#include <mruby/proc.h>
#include <mruby/string.h>
#include <mruby/value.h>
#include <mruby/variable.h>

//-------------------------------------------------------------------
// Helpers to deal with calling back into Go.
//-------------------------------------------------------------------
// This is declard in func.go and is a way for us to call back into
// Go to execute a method.
extern mrb_value *go_mrb_func_call(mrb_state*, mrb_value*, mrb_value*);

// This calls into Go with a similar signature to mrb_func_t. We have to
// change it slightly because cgo can't handle the union type of mrb_value,
// so we pass in a pointer instead. Additionally, the result is also a
// pointer to work around Go's confusion with unions.
static inline mrb_value _go_mrb_func_call(mrb_state *s, mrb_value self) {
    mrb_value exc = mrb_nil_value();
    mrb_value result = *go_mrb_func_call(s, &self, &exc);

    // We raise if we got an exception. We have to raise from here and
    // not from within Go because it messes with Go's calling conventions,
    // resulting in a broken stack.
    if (!mrb_nil_p(exc)) {
        mrb_exc_raise(s, exc);
    }

    return result;
}

// This method is used as a way to get a valid mrb_func_t that actually
// just calls back into Go.
static inline mrb_func_t _go_mrb_func_t() {
    return &_go_mrb_func_call;
}

//-------------------------------------------------------------------
// Helpers to deal with getting arguments
//-------------------------------------------------------------------
// This is declard in args.go
extern void go_get_arg_append(mrb_value*);

// This gets all arguments given to a function call and adds them to
// the accumulator in Go.
static inline int _go_mrb_get_args_all(mrb_state *s) {
    mrb_value *argv;
    mrb_value block;
    int argc, i, count;

    count = mrb_get_args(s, "*&", &argv, &argc, &block);
    for (i = 0; i < argc; i++) {
        go_get_arg_append(&argv[i]);
    }

    if (!mrb_nil_p(block)) {
        go_get_arg_append(&block);
    }

    return count;
}

//-------------------------------------------------------------------
// Misc. helpers
//-------------------------------------------------------------------

// This is used to help calculate the "send" value for the parser,
// since pointer arithmetic like this is hard in Go.
static inline const char *_go_mrb_calc_send(const char *s) {
    return s + strlen(s);
}

// Sets the capture_errors field on mrb_parser_state. Go can't access bit
// fields.
static inline void
_go_mrb_parser_set_capture_errors(struct mrb_parser_state *p, mrb_bool v) {
    p->capture_errors = v;
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

static inline short _go_mrb_is_dead(mrb_state *s, mrb_value o) {
    return is_dead(s, mrb_obj_ptr(o));
}

static inline struct RProc *_go_mrb_proc_ptr(mrb_value o) {
    return mrb_proc_ptr(o);
}

static inline enum mrb_vtype _go_mrb_type(mrb_value o) {
    return mrb_type(o);
}

#endif
