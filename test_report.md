## Difference Report

| Key | Error Code | lhs | rhs | Description |
| --- | --- | --- | --- | --- |
| `d[2]` | `INDEX_NOT_FOUND` | `(null)null` | `(string)c` | 인덱스가 존재하지 않습니다. |
| `c.b` | `KEY_NOT_FOUND` | `(string)b` | `(null)null` | 키가 존재하지 않습니다. |
| `b` | `TYPE_UNMATCHED` | `(array)[a b c]` | `(map)map[a:a b:b c:c]` | 타입이 일치하지 않습니다.  |
| `a` | `VALUE_UNMATCHED` | `(string)b` | `(string)a` | 값이 일치하지 않습니다. |
