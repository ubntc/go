# TypeScript Things

Now and then I need to hack some TS to assist friends and family.
Some of the learnings will emerge here.

## TS 1: Optional and Named Parameters with Default Values

### Example
Directly go to the example: [optional-params.ts](optional-params.ts)

### Explanation

TypeScript does not have named (non-positional) parameters out of the box. Also named parameters are most useful if you can make them optional and assign a default value. This is also a bit complicated in TS.

Here is what I mean (Python):
```python
class Thing:
    # constructor with named, optional params and default values
    def __init__(self, firstname:str='The', lastname:str='Thing'):
        self.name = f'{firstname} {lastname}'

# call with any param omitted as desired:
t0 = Thing()
t1 = Thing(firstname="A")
t2 = Thing(lastname="Thinglehooper")
t3 = Thing(firstname="What a", lastname="Thing!")
```

In Python (and other languages) it is easy to define
named and optional paramaters and assign them default values.

## Implementation
Achieving this in TS is a bit harder:

1. First define an additional `interface` type with optional field names using `?` as field name suffix.
   ```typescript
   interface ThingOptions {
       firstname?:string
       lastname?:string
   }
   ```
2. Your constructor needs to accept this interface as a parameter object and implicitly deconstruct this object with the desired default values.
   ```typescript
   constructor({firstname="The", lastname="Thing"}: ThingOptions) {}
   ```
3. When calling the optionalized function you need to pass the params as object, but can now omit or add any params you like.
   ```typescript
   let t1 = new Thing({firstname:"A"})
   ```

This solution is not beauityful but at least it is type-safe and gets the job done.

See the full example [here](optional-params.ts).