#pragma once

#include <QScopedPointer>

template<typename T, void (*f) (T*)>
struct ScopedPointerDestroyerDeleter {
    static void cleanup(T* t) {
        f(t);
    }
};

template<typename T, void (*cleanup) (T*)>
using ScopedPointer = QScopedPointer<T, ScopedPointerDestroyerDeleter<T, cleanup>>;
