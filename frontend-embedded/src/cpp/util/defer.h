#pragma once

#include <functional>

template<typename F>
struct ScopeGuard {
    F f;
    explicit ScopeGuard(F f)
        : f(f)
    {}
    ~ScopeGuard() { f(); }
};

template<typename F>
ScopeGuard<F> onScopeExit(F f)
{
    return ScopeGuard<F>(f);
}

template<typename F>
auto after(F f, const std::function<void()> &cleanup) -> decltype(f())
{
    auto _ = onScopeExit(cleanup);
    f();
}
