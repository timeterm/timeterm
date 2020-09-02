#pragma once

#include <QScopedPointer>

template<typename T, void (*f)(T *)>
struct ScopedPointerDestroyFnDeleter {
    static void cleanup(T *t)
    {
        f(t);
    }
};

template<typename T, void (*cleanup)(T *)>
using ScopedPointer = QScopedPointer<T, ScopedPointerDestroyFnDeleter<T, cleanup>>;
