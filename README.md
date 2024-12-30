# Hello YAML Difference Reporter! 👋

안녕하세요! 해당 CLI는 두 가지의 YAML 파일의 차이점을 비교하고, 간단한 리포트를 작성하는 도구입니다.

일반적인 YAML Diff 도구와는 다르게 검사 모드로 `type`, `value`, `key`, `index` 모드를 제공하기에 상황에 맞는 비교를 할 수 있습니다.

# Installation

```bash
$ go install github.com/YangTaeyoung/yaml-diff-reporter@v1.0.0
```

# Modes

## Type Mode

`type` 모드는 두 YAML 파일의 타입이 다른 경우를 검사합니다.

### Example

타입 불일치 케이스 (스칼라 값 vs 배열)

```yaml
# file1.yaml
key: value
``` 

```yaml
# file2.yaml
key:
  - value
```

## Value Mode

`value` 모드는 두 YAML 파일의 값이 다른 경우를 검사합니다

### Example

값 불일치 케이스 (`key`에 해당하는 값이 서로 다른 값임 `value1` vs `value2`)

```yaml
# file1.yaml
key: value1
```

```yaml
# file2.yaml
key: value2
```

## Key Mode

`key` 모드는 맵에서 두 파일 중 한쪽에만 존재하는 키가 있는지 검사합니다.

### Example

키 불일치 케이스 (`key1`이 `file2.yaml`에 없음, `key2`가 `file1.yaml`에 없음)

```yaml
# file1.yaml
key1: value
```

```yaml
# file2.yaml
key2: value
```

## Index Mode

`index` 모드는 배열에서 두 파일 중 한쪽에만 존재하는 인덱스가 있는지 검사합니다.

### Example

인덱스 불일치 케이스 (`file1.yaml`의 key 배열의 인덱스가 `file2`에 없음)

```yaml
# file1.yaml
key:
  - value1
  - value2
```

```yaml
# file2.yaml
key:
  - value3
```

# Error Codes

| Code              | Description         |
|-------------------|---------------------|
| `TYPE_UNMATCHED`  | 타입이 일치하지 않음         |
| `VALUE_UNMATCHED` | 값이 일치하지 않음          |
| `KEY_NOT_FOUND`   | 한쪽 파일에 키가 존재하지 않음   |
| `INDEX_NOT_FOUND` | 한쪽 파일에 인덱스가 존재하지 않음 |

# Flags

아래는 yaml-diff-reporter에서 사용 가능한 플래그 목록입니다. `--help` 플래그를 통해 보다 자세한 정보를 확인할 수도 있습니다.

| Flags                                      | Description                                              | Enums                          | Support Multiple Values | Required |
|--------------------------------------------|----------------------------------------------------------|--------------------------------|-------------------------|----------|
| `-M <value>`, <br>`--modes <value>`        | 비교 모드를 지정합니다. (default: `type`, `value`, `key`, `index`) | `type`, `value`,`key`, `index` | ✅                       | ❌        |
| `-l <value>`, <br>`--lhs-path <value>`     | 비교할 좌측 YAML 파일의 경로를 지정합니다.                               |                                | ❌                       | ✅        |
| `-r <value>` <br>`--rhs-path <value>`      | 비교할 우측 YAML 파일의 경로를 지정합니다.                               |                                | ❌                       | ✅        |
| `-la <value>`, <br>`--lhs-alias <value>`   | 좌측 YAML 파일의 별칭을 지정합니다. (default: `lhs`)                  |                                | ❌                       | ❌        |
| `-ra <value>`, <br>`--rhs-alias <value>`   | 우측 YAML 파일의 별칭을 지정합니다. (default: `rhs`)                  |                                | ❌                       | ❌        |
| `-ot <value>`, <br>`--output-type <value>` | 리포트를 출력할 방식을 지정합니다. (default: `stdout`)                  | `file`,`stdout`                | ❌                       | ❌        |
| `-o <value>`, <br>`--output-path <value>`  | 리포트를 저장할 경로를 지정합니다.                                      |                                | ❌                       | ❌        |
| `-f <value>`, <br>`--format <value>`       | 리포트 포맷을 지정합니다. (default: `json`)                         | `json`, `yaml`, `plain`        | ❌                       | ❌        |
| `-lang <value>`, <br>`--language <value>`  | 리포트 언어를 지정합니다. (default: `en`)                           | `en`, `ko`                     | ❌                       | ❌        |

# Simple Example

tests 디렉토리 내에 테스트를 위한 파일이 있습니다. 다음 명령어를 통해 해당 파일을 비교해보세요!

```bash
$ yaml-diff-reporter --lhs-path ./tests/A.yaml \ 
  --rhs-path ./tests/B.yaml \
  --output-path ./test_report.md \
  --output-type file \
  --format markdown \ 
  --language ko \ 
```

## Result

포맷(`--format`)과 출력 유형(`--output-type`)을 파일(`file`)로 지정하고 경로(`--output-path`)를 지정한 경우, 다양한 방식으로 파일을 생성하여 리포트를 조회할 수 있습니다.

- [markdown](./test_report.md)
- [json](./test_report.json)
- [plain](./test_report.txt)

# Trouble Shooting 👊

```bash
$ yaml-diff-reporter
> zsh: command not found: yaml-diff-reporter
```

go로 설치한 프로그램을 실행할 때 발생하는 에러입니다. `~/.zshrc` 파일(혹은 `~/.bashrc`)의 하단에 다음과 같이 환경변수를 추가합니다.

```bash
# ... (생략)
export PATH="$HOME/go/bin:$PATH"
```
