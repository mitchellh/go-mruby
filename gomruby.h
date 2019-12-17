// vim: ft=c ts=2 sts=2 st=2
/*
 * This header exists to simplify the headers that are included within
 * the Go files. This header should include all the necessary headers
 * for the compilation of the Go library.
 * */

#ifndef _GOMRUBY_H_INCLUDED
#define _GOMRUBY_H_INCLUDED

#include <errno.h>
#include <mruby.h>
#include <mruby/array.h>
#include <mruby/class.h>
#include <mruby/compile.h>
#include <mruby/error.h>
#include <mruby/irep.h>
#include <mruby/gc.h>
#include <mruby/hash.h>
#include <mruby/proc.h>
#include <mruby/string.h>
#include <mruby/throw.h>
#include <mruby/value.h>
#include <mruby/variable.h>

// (erikh) this can be set in mruby/mrbconfig.h so we can default it here.
// XXX I don't know how this actually plays out when the config is modified.
// I'm taking a WAG here. Either way, the default is 16 in vm.c.
#ifndef MRB_FUNCALL_ARGC_MAX
  #define MRB_FUNCALL_ARGC_MAX 16
#endif // MRB_FUNCALL_ARGC_MAX

//-------------------------------------------------------------------
// Helpers to deal with calling back into Go.
//-------------------------------------------------------------------
// This is declard in func.go and is a way for us to call back into
// Go to execute a method.
extern mrb_value goMRBFuncCall(mrb_state*, mrb_value);

// This method is used as a way to get a valid mrb_func_t that actually
// just calls back into Go.
static inline mrb_func_t _go_mrb_func_t() {
    return &goMRBFuncCall;
}

//-------------------------------------------------------------------
// Helpers to deal with calling into Ruby (C)
//-------------------------------------------------------------------

static mrb_value load_string_cb(mrb_state *mrb, mrb_value in) {
  return mrb_load_string(mrb, (const char*)mrb_cptr(in));
}

static mrb_value _go_mrb_load_string(mrb_state *mrb, const char *s) {
  mrb_bool state;
  mrb_value result = mrb_protect(mrb, load_string_cb, mrb_cptr_value(mrb, (void*)s), &state);
  if (state) {
    mrb->exc = mrb_obj_ptr(result);
  }
  return result;
}

struct yield_data {
  mrb_value block;
  mrb_int argc;
  const mrb_value *argv;
};

static mrb_value yield_argv_cb(mrb_state *mrb, mrb_value in) {
  struct yield_data *d = (struct yield_data*)mrb_cptr(in);
  return mrb_yield_argv(mrb, d->block, d->argc, d->argv);
}

static mrb_value _go_mrb_yield_argv(mrb_state *mrb, mrb_value b, mrb_int argc, const mrb_value *argv) {
  struct yield_data d = { b, argc, argv };
  mrb_bool state;
  mrb_value result = mrb_protect(mrb, yield_argv_cb, mrb_cptr_value(mrb, &d), &state);
  if (state) {
    mrb->exc = mrb_obj_ptr(result);
  }
  return result;
}

struct call_data {
  mrb_value self;
  mrb_sym method;
  mrb_int argc;
  const mrb_value *argv;
  const mrb_value *block;
};

static mrb_value mrb_call_cb(mrb_state *mrb, mrb_value in) {
  struct call_data *d = (struct call_data*)mrb_cptr(in);
  if (d->block != NULL) {
    return mrb_funcall_with_block(mrb, d->self, d->method, d->argc, d->argv, *d->block);
  } else {
    return mrb_funcall_argv(mrb, d->self, d->method, d->argc, d->argv);
  }
}

static mrb_value _go_mrb_call(mrb_state *mrb, mrb_value self, mrb_sym method, mrb_int argc, const mrb_value *argv, const mrb_value *block) {
  struct call_data d = { self, method, argc, argv, block };
  mrb_bool state;
  mrb_value result = mrb_protect(mrb, mrb_call_cb, mrb_cptr_value(mrb, &d), &state);
  if (state) {
    mrb->exc = mrb_obj_ptr(result);
  }
  return result;
}

//-------------------------------------------------------------------
// Helpers to deal with getting arguments
//-------------------------------------------------------------------
// This is declard in args.go
extern void goGetArgAppend(mrb_value);

// This gets all arguments given to a function call and adds them to
// the accumulator in Go.
static inline int _go_mrb_get_args_all(mrb_state *s) {
  mrb_value *argv;
  mrb_value block;
  mrb_bool append;
  mrb_int argc, i;

  mrb_get_args(s, "*&?", &argv, &argc, &block, &append);

  for (i = 0; i < argc; i++) {
    goGetArgAppend(argv[i]);
  }

  if (append == FALSE || mrb_type(block) == MRB_TT_FALSE) {
    return argc;
  }

  argc++;
  goGetArgAppend(block);

  return argc;
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

static inline float _go_mrb_float(mrb_value o) {
  return mrb_float(o);
}

static inline int _go_mrb_fixnum(mrb_value o) {
  return mrb_fixnum(o);
}

static inline struct RBasic *_go_mrb_basic_ptr(mrb_value o) {
  return mrb_basic_ptr(o);
}

static inline struct RProc *_go_mrb_proc_ptr(mrb_value o) {
  return mrb_proc_ptr(o);
}

static inline enum mrb_vtype _go_mrb_type(mrb_value o) {
  return mrb_type(o);
}

static inline mrb_bool _go_mrb_nil_p(mrb_value o) {
  return mrb_nil_p(o);
}

static inline struct RClass *_go_mrb_class_ptr(mrb_value o) {
  return mrb_class_ptr(o);
}

static inline void _go_set_gc(mrb_state *m, int val) {
  mrb_gc *gc = &m->gc;
  gc->disabled = val;
}

static inline void _go_disable_gc(mrb_state *m) {
  _go_set_gc(m, 1);
}

static inline void _go_enable_gc(mrb_state *m) {
  _go_set_gc(m, 0);
}

static inline int _go_get_max_funcall_args() {
  return MRB_FUNCALL_ARGC_MAX;
}

// this function returns 1 if the value is dead, aka reaped or otherwise
// terminated by the GC.
static inline int _go_isdead(mrb_state *m, mrb_value o) {
  // immediate values such as Fixnums and symbols are never to be garbage
  // collected, so converting them to a basic pointer yields an invalid one.
  // This pattern is seen in the mruby source's gc.c.
  if mrb_immediate_p(o) {
    return 0;
  }

  struct RBasic *ptr = mrb_basic_ptr(o);

  // I don't actually know this is a potential condition but better safe than sorry.
  if (ptr == NULL) {
    return 1;
  }

  return mrb_object_dead_p(m, ptr);
}

static inline int _go_gc_live(mrb_state *m) {
  mrb_gc *gc = &m->gc;
  return gc->live;
}

static inline void _go_mrb_context_set_capture_errors(struct mrbc_context *ctx, int state) {
  ctx->capture_errors = FALSE;

  if (state != 0) {
    ctx->capture_errors = TRUE;
  }
}

static inline mrb_value _go_mrb_context_run(mrb_state *m, struct RProc *proc, mrb_value self, int *stack_keep) {
  mrb_value result = mrb_context_run(m, proc, self, *stack_keep);
  *stack_keep = proc->body.irep->nlocals;
  return result;
}

static inline struct RObject* _go_mrb_getobj(mrb_value v) {
  return mrb_obj_ptr(v);
}

static inline void _go_mrb_iv_set(mrb_state *m, mrb_value self, mrb_sym sym, mrb_value v) {
  mrb_iv_set(m, self, sym, v);
}

static inline mrb_value _go_mrb_iv_get(mrb_state *m, mrb_value self, mrb_sym sym) {
  return mrb_iv_get(m, self, sym);
}

static inline void _go_mrb_gv_set(mrb_state *m, mrb_sym sym, mrb_value v) {
  mrb_gv_set(m, sym, v);
}

static inline mrb_value _go_mrb_gv_get(mrb_state *m, mrb_sym sym) {
  return mrb_gv_get(m, sym);
}

static inline mrb_int _mrb_ary_len(mrb_value ary) {
  return RARRAY_LEN(ary);
}

#endif
