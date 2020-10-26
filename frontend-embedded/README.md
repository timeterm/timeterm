# frontend-embedded

Embedded frontend for Timeterm.

## Compilation

### Compilation on the Raspberry Pi

```
frontend-embedded $ mkdir build && cd build
build $ cmake .. -G Ninja -DTIMETERMOS:BOOL=TRUE
build $ ninja
```
