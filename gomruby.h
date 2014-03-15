/*
 * This header exists to simplify the headers that are included within
 * the Go files. This header should include all the necessary headers
 * for the compilation of the Go library.
 * */

#include <mruby.h>
#include <mruby/compile.h>
#include <mruby/string.h>

inline static void *_go_mrb_ptr(mrb_value v) {
    return mrb_ptr(v);
}
