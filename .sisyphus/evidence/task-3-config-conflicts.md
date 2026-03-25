# Task 3 설정 충돌 해결 결과

## 수행 명령
- `GIT_MASTER=1 git fetch upstream`
- `GIT_MASTER=1 git merge upstream/main --no-commit`
- `GIT_MASTER=1 git diff --name-only --diff-filter=U`
- `GIT_MASTER=1 git checkout --theirs go.sum client/yarn.lock`
- `GIT_MASTER=1 git rm .github/workflows/astro.yml`
- `GIT_MASTER=1 git add .gitignore Makefile go.mod go.sum client/package.json client/yarn.lock`

## 해결한 충돌 파일
- `.gitignore`
- `Makefile`
- `go.mod`
- `go.sum`
- `client/package.json`
- `client/yarn.lock`
- `.github/workflows/astro.yml`

## 결정 사항

### `.gitignore`
- upstream의 `.env` 무시 규칙 유지
- 포크 고유 항목 `.claude/settings.local.json`, `koito` 유지
- 공통 항목 `test_config_dir` 유지

### `Makefile`
- upstream 구조 유지 (`.env` include, `api.debug: postgres.start`)
- 포크 고유 Postgres 이미지 `ghcr.io/kang-heewon/postgresql-local:18` 유지
- 포크 고유 `pg_bigm` preload 유지
- upstream의 named volume 패턴을 반영해 `koito_dev_db`, `koito_scratch_db` 추가

### `go.mod`
- 충돌 마커 없음
- upstream 기반 상태 유지
- 포크 고유로 보이는 `github.com/go-chi/httprate`, `golang.org/x/image` 항목이 이미 유지됨을 확인

### `go.sum`
- 요구사항대로 upstream 버전 채택
- 이후 빌드/의존성 재생성 단계에서 다시 정리 예정

### `client/package.json`
- upstream 기반 버전 채택
- 포크 고유 의존성 `motion` 유지
- 충돌 구간에서는 upstream의 `react-is@^19.2.3`, `recharts@^3.6.0` 채택

### `client/yarn.lock`
- 요구사항대로 upstream 버전 채택
- 이후 재생성 단계에서 갱신 예정

### `.github/workflows/astro.yml`
- modify/delete 충돌
- HEAD에서 삭제된 상태를 유지하도록 삭제 처리

## 검증 결과
- 설정 파일들은 더 이상 unmerged 상태가 아님
- 현재 남은 unmerged 파일은 Go/TSX/SQL 등 후속 태스크 범위의 코드 충돌만 존재
- `GIT_MASTER=1 git status --short -- .gitignore Makefile go.mod go.sum client/package.json client/yarn.lock .github/workflows/astro.yml` 결과상 대상 파일은 모두 `M` 상태이며 `UU`, `AA`, `UD` 없음

## 참고
- `lsp_diagnostics` 검증은 환경 제약이 있었음
  - `client/package.json`: biome LSP 미설치
  - `go.mod`: `.mod` 전용 LSP 미구성
- 따라서 이 태스크에서는 git 기준 충돌 해소 상태를 검증 근거로 사용
