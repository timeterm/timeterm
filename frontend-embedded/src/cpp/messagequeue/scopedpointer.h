#ifndef SCOPEDPOINTER_H
#define SCOPEDPOINTER_H

namespace MessageQueue
{

template<typename T, typename R, R (*f)(T *)>
struct ScopedPointerDestroyerDeleter {
    static inline void cleanup(T *pointer)
    {
        f(pointer);
    }
};

} // namespace MessageQueue

#endif // SCOPEDPOINTER_H
