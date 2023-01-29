### Const
A `const` definition defines a single typed value. You can't refer to these values in the rest of the Bebop schema, but they come in handy when the parties using your schemas also need to agree on certain application-wide parameters.

Currently, valid `const` types are: boolean, integers, floats, strings, and GUIDs.

### Enum
An `enum` defines a type that acts as a wrapper around an integer type (defaults to `uint32`), with certain named constants, each having a corresponding underlying integer value. It is used much like an `enum` in C.

> The syntax is: `enum Flavor: uint8 { Vanilla = 1; Chocolate = 2; Mint = 3; }`.
> 
> * Unlike in C, all constants must be explicitly given an integer literal value.
>
> * You should never remove a constant from an `enum` definition. Instead, put `[deprecated("reason here")]` in front of the name.
>
> * You're free to add new constants to an `enum` at any point in the future.
#### Flags enum
By default, a Bebop enum type is _not_ supposed to represent any underlying values _outside_ of the ones listed.

If you want a more C-like, bitflags-like behavior, add a `[flags]` attribute before the enum:

```c
[flags]
enum Permissions {
    Read = 0x01;
    Write = 0x02;
    Comment = 0x04;
}
```

Defined this way, `Permissions` values like `0` (no permissions) and `3` (`Read` + `Write`) are valid too.

### Struct
A `struct` defines an aggregation of "fields", containing typed values in a fixed order. All values are always present. It is used much like a `struct` in C.

@@ -52,24 +87,39 @@ A `struct` defines an aggregation of "fields", containing typed values in a fixe
### Message
A `message` defines an indexed aggregation of fields containing typed values, each of which may be absent. It might correspond to something like a `class` in Java, or a JSON object.

> The syntax is: `message Song { string title = 1; uint16 year = 2; }` — note the indices.
> The syntax is: `message Song { 1 -> string title; 2 -> uint16 year; }` — note the indices before each field.
>
> * In the binary representation of a `message`, the message is prefixed with its length, and each field is prefixed with its index.
>
> * It's okay to add fields to a `message` with new indices later — in fact, this is the whole point of `message`. (When an unrecognized field index is encountered in the process of decoding a `message`, it is skipped over. This allows for compatibility with versions of your app that use an older version of the schema.)
### Union
A `union` defines a tagged union of one or more inline `struct` or `message` definitions. Each is preceded by a "discriminator" or "tag" value. This defines a type whose values may assume any _one_ of the aggregate layouts defined inside. It corresponds to something like C++'s [std::variant](https://en.cppreference.com/w/cpp/utility/variant).

> The syntax is: `union U { 1 -> message A { ... }; 2 -> struct B { ... } }`.
>
> * The binary representation of a `U` value is then: a length prefix, followed by either (a) a `01` byte followed by an encoding of an `A` message, or (b) a `02` byte followed by an encoding of a `B` struct.
>
> * Just like with messages, new branches may be added to a union later. When an unrecognized discriminator value is encountered, the length prefix is used to skip over the body, and decoding fails in a way your program may catch.
>
> * Nested types are not available globally but do reserve the identifier globally. E.g. in the above you cannot define `struct Other { A x; }` because `A` is private to `U` but you also cannot define `struct A { ... }` because `A` is reserved globally.

## Types
The following types are built-ins:

| Name | Description |
|---|---|
| `bool` | A Boolean value, true or false. |
| `byte` | An unsigned 8-bit integer. `uint8` is an alias. |
| `uint8` | An unsigned 8-bit integer. |
| `int8` | A signed 8-bit integer. |
| `uint16` | An unsigned 16-bit integer. |
| `int16` | A signed 16-bit integer. |
| `uint32` | An unsigned 32-bit integer. |
| `int32` | A signed 32-bit integer. |
| `uint64` | An unsigned 64-bit integer. |
| `int64` | A signed 64-bit integer. |
| `float32` | A 32-bit IEEE [single-precision floating point number](https://en.wikipedia.org/wiki/Single-precision_floating-point_format). |
| `float64` | A 64-bit IEEE [double-precision floating point number](https://en.wikipedia.org/wiki/Double-precision_floating-point_format).  |
| `string` | A length-prefixed UTF-8-encoded string. |
| `guid` | A [GUID](https://en.wikipedia.org/wiki/Universally_unique_identifier). |
| `date` | A [UTC](https://en.wikipedia.org/wiki/Coordinated_Universal_Time) date / timestamp. |
| `T[]` | A length-prefixed array of `T` values. `array[T]` is an alias. |
| `map[T1, T2]` | A map, as a length-prefixed array of (`T1`, `T2`) association pairs. |
You may also use user-defined types (`enum`s and other records) as field types.
A string is stored as a length-prefixed array of bytes. All length-prefixes are 32-bit unsigned integers, which means the maximum number of bytes in a string, or entries in an array or map, is about 4 billion (2^32).
A `guid` is stored as 16 bytes, in [Guid.ToByteArray](https://docs.microsoft.com/en-us/dotnet/api/system.guid.tobytearray?view=net-5.0) order.

A `date` is stored as a 64-bit integer amount of “ticks” since 00:00:00 UTC on January 1 of year 1 A.D. in the Gregorian calendar, where a “tick” is 100 nanoseconds.