# Upstream 전체 머지: gabehf/Koito → kang-heewon/Koito

## TL;DR

> **Quick Summary**: upstream(gabehf/Koito)의 ~55개 커밋을 포크(kang-heewon/Koito)에 머지하면서, 포크의 핵심 기능 4개(Wrapped, Recommendations, backfill 개선, pg_bigm)를 유지하고 2개(GenreStats, 모바일 반응형)를 버린다.
> 
> **Deliverables**:
> - 70+ 충돌이 해결된 깨끗한 머지 커밋
> - upstream 신규 기능(InterestGraph, Timeframe, authenticate.go 등) 통합
> - 포크 고유 기능(Wrapped, Recommendations, backfill, pg_bigm) 정상 동작
> - GenreStats/모바일 반응형 제거 확인
> 
> **Estimated Effort**: Large
> **Parallel Execution**: YES - 7 waves
> **Critical Path**: 사전분석 → 머지시작 → 백엔드코어 → 백엔드핸들러 → 클라이언트 → 빌드검증 → 커밋

---

## Context

### Original Request
사용자가 gabehf/Koito(upstream)의 변경점을 포크(kang-heewon/Koito)에 받아 적용 요청. 인터뷰를 통해 "upstream 전체 머지 + 포크 핵심 기능 유지"로 범위 확정.

### Interview Summary
**Key Discussions**:
- 머지 전략: `git merge upstream/main` 직접 머지 선택
- 브랜치: main에 직접 머지 (별도 브랜치 없음)
- 유지할 포크 기능: Wrapped, Recommendations, backfill 개선, pg_bigm
- 버릴 포크 기능: GenreStats, 모바일 반응형

**Research Findings**:
- 70+ 파일 충돌 확인 (dry-run 완료)
- modify/delete 충돌 2건: astro.yml (HEAD 삭제), validate.go (upstream 삭제)
- add/add 충돌 7건: Rewind 관련 TSX 3개, RewindPage, get_summary.go, summary.go, summary_test.go
- upstream 변경: 클라이언트 51파일 +2899/-1558줄, 백엔드 87파일 +4565/-1429줄

### Metis Review
**Identified Gaps** (addressed):
- DB 마이그레이션 순서/호환성: 플랜에 사전 스키마 비교 태스크 추가
- Wrapped/Recommendations의 GenreStats/모바일 의존성: 사전 import 분석 태스크 추가
- validate.go→authenticate.go 전환 영향: 백엔드 코어 태스크에서 명시적 처리
- 롤백 전략: `git merge --abort` 또는 `git reset --hard HEAD` 명시

---

## Work Objectives

### Core Objective
upstream의 모든 변경사항을 포크에 통합하면서, 포크의 핵심 기능 4개를 깨뜨리지 않고 유지한다.

### Concrete Deliverables
- 70+ 충돌이 해결된 머지 완료 상태의 main 브랜치
- `make build`, `make test`, `yarn typecheck` 모두 통과

### Definition of Done
- [x] `make build` 성공 (exit code 0)
- [x] `make test` 모든 테스트 통과
- [x] `cd client && yarn typecheck` 타입 에러 없음
- [x] `git status` clean working tree

### Must Have
- upstream의 모든 커밋 통합 (InterestGraph, Timeframe, authenticate.go 등)
- Wrapped 페이지 정상 동작
- Recommendations 페이지 정상 동작
- backfill 개선사항 유지
- pg_bigm 마이그레이션 유지

### Must NOT Have (Guardrails)
- GenreStats 페이지 코드 유지 (버려야 함)
- 모바일 반응형 개선 코드 유지 (버려야 함)
- "혹시 몰라" 보존하는 버릴 기능 코드
- upstream 로직 수정 (포크 기능 유지 때문이 아닌 경우)
- Wrapped/Recommendations 기능 개선 (현 상태 유지만)
- 의존성 업그레이드 (머지 필수가 아닌 경우)
- authenticate.go 리팩토링 (upstream 그대로 채택)

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: YES (Go: `testing` + `testify`, Frontend: TypeScript typecheck)
- **Automated tests**: Tests-after (머지 완료 후 전체 테스트 실행)
- **Framework**: `go test` (backend), `yarn typecheck` (frontend)

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Backend**: Bash — `go build`, `go test`, `make test`
- **Frontend**: Bash — `yarn typecheck`, `yarn build`
- **Integration**: Bash (curl) — API 엔드포인트 호출, HTTP 상태 확인

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — 사전 분석, 병렬):
├── Task 1: 포크 기능 의존성 분석 [quick]
└── Task 2: DB 마이그레이션 스키마 비교 [quick]

Wave 2 (After Wave 1 — 머지 시작 + 설정 충돌):
└── Task 3: git merge 실행 + 설정 파일 충돌 해결 [deep]

Wave 3 (After Wave 2 — 백엔드 코어, 병렬):
├── Task 4: validate.go→authenticate.go 충돌 해결 [deep]
└── Task 5: DB 쿼리/마이그레이션 충돌 해결 [quick]

Wave 4 (After Wave 3 — 백엔드 핸들러 + internal):
└── Task 6: 백엔드 핸들러/서비스 충돌 해결 [deep]

Wave 5 (After Wave 4 — 클라이언트 전체):
└── Task 7: 클라이언트 TSX/라우트 충돌 해결 [deep]

Wave 6 (After Wave 5 — 빌드 검증):
└── Task 8: 전체 빌드 + 테스트 검증 [quick]

Wave 7 (After Wave 6 — 커밋):
└── Task 9: 머지 커밋 완료 [quick]

Critical Path: T1 → T3 → T4 → T6 → T7 → T8 → T9
Parallel Speedup: Wave 1 (T1∥T2), Wave 3 (T4∥T5)
```

### Dependency Matrix

| Task | Depends On | Blocks |
|------|-----------|--------|
| 1 | — | 3 |
| 2 | — | 5 |
| 3 | 1 | 4, 5, 6, 7 |
| 4 | 3 | 6 |
| 5 | 2, 3 | 6 |
| 6 | 4, 5 | 7 |
| 7 | 6 | 8 |
| 8 | 7 | 9 |
| 9 | 8 | — |

### Agent Dispatch Summary

- **Wave 1**: **2** — T1 → `quick`, T2 → `quick`
- **Wave 2**: **1** — T3 → `deep` + `git-master`
- **Wave 3**: **2** — T4 → `deep` + `git-master`, T5 → `quick` + `git-master`
- **Wave 4**: **1** — T6 → `deep` + `git-master`
- **Wave 5**: **1** — T7 → `deep` + `git-master`
- **Wave 6**: **1** — T8 → `quick`
- **Wave 7**: **1** — T9 → `quick` + `git-master`

---

## TODOs

> Implementation + Verification = ONE Task.
> EVERY task MUST have: Recommended Agent Profile + Parallelization info + QA Scenarios.

- [x] 1. 포크 기능 의존성 분석 — Wrapped/Recommendations가 버릴 기능에 의존하는지 확인

  **What to do**:
  - Wrapped 페이지 관련 파일들에서 GenreStats, 모바일 반응형 관련 import/참조 검색
  - Recommendations 페이지 관련 파일들에서 동일 검색
  - backfill 관련 코드에서 GenreStats 의존성 검색
  - 결과를 `.sisyphus/evidence/task-1-dependency-analysis.md`에 기록

  **Must NOT do**:
  - 코드 수정 (분석만 수행)
  - 머지 시작

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 검색/분석만 수행하는 단순 태스크
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `git-master`: 코드 검색이지 git 작업 아님

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 2)
  - **Blocks**: Task 3
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `client/app/` — 클라이언트 전체 디렉토리에서 Wrapped/Recommendations 관련 파일 탐색
  - `internal/` — 백엔드에서 backfill 관련 코드 탐색

  **검색 대상 키워드**:
  - `GenreStats`, `genre-stats`, `genre_stats` — 버릴 기능 참조
  - `mobile`, `responsive`, `@media` — 모바일 반응형 관련 (포크 고유 추가분)
  - `Wrapped`, `Recap`, `Recommendations` — 유지할 기능 파일 위치 파악

  **WHY Each Reference Matters**:
  - 버릴 기능에 유지할 기능이 의존하면, 단순히 upstream 버전을 채택할 수 없고 추가 작업이 필요하므로 사전에 파악해야 함

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Wrapped/Recommendations의 GenreStats 의존성 없음 확인
    Tool: Bash (grep)
    Preconditions: 현재 HEAD 상태 (머지 전)
    Steps:
      1. grep -r "GenreStats\|genre-stats\|genre_stats" client/app/components/Wrapped/ client/app/components/Recommendations/ client/app/routes/RewindPage/ --include="*.tsx" --include="*.ts" -l
      2. grep -r "GenreStats\|genre-stats\|genre_stats" internal/ --include="*.go" -l
    Expected Result: 두 명령 모두 결과 없음 (exit code 1, 빈 출력)
    Failure Indicators: 파일 경로가 출력되면 의존성 존재
    Evidence: .sisyphus/evidence/task-1-dependency-analysis.md
  ```

  **Evidence to Capture:**
  - [ ] task-1-dependency-analysis.md — 의존성 분석 결과 전체

  **Commit**: NO

- [x] 2. DB 마이그레이션 스키마 비교 — 포크와 upstream 마이그레이션 충돌 여부 확인

  **What to do**:
  - 포크(HEAD)의 `db/migrations/` 파일 목록 확인
  - upstream의 `db/migrations/` 파일 목록 확인 (`git show upstream/main:db/migrations/` 사용)
  - 양쪽에만 있는 마이그레이션 파일 식별
  - 동일 번호 마이그레이션 충돌 여부 확인
  - pg_bigm 마이그레이션이 upstream 마이그레이션과 순서 충돌하는지 확인
  - 결과를 `.sisyphus/evidence/task-2-migration-diff.md`에 기록

  **Must NOT do**:
  - DB 마이그레이션 실행
  - 마이그레이션 파일 수정

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: git 명령어로 파일 비교만 수행
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocks**: Task 5
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `db/migrations/` — 마이그레이션 파일 디렉토리

  **WHY Each Reference Matters**:
  - 마이그레이션 번호가 겹치면 머지 후 goose가 실행 순서를 잘못 판단할 수 있음
  - pg_bigm 마이그레이션이 upstream 스키마 변경과 호환되는지 사전 확인

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: 마이그레이션 파일 비교 및 충돌 분석
    Tool: Bash (git)
    Preconditions: upstream 리모트 fetch 완료
    Steps:
      1. ls db/migrations/ > /tmp/fork-migrations.txt
      2. git ls-tree --name-only upstream/main:db/migrations/ > /tmp/upstream-migrations.txt
      3. diff /tmp/fork-migrations.txt /tmp/upstream-migrations.txt 또는 comm 명령으로 비교
    Expected Result: 차이점이 명확히 식별됨, 번호 충돌 여부 파악
    Failure Indicators: 동일 번호에 다른 내용의 마이그레이션 존재
    Evidence: .sisyphus/evidence/task-2-migration-diff.md
  ```

  **Evidence to Capture:**
  - [ ] task-2-migration-diff.md — 마이그레이션 비교 결과

  **Commit**: NO

- [x] 3. Git merge 실행 + 설정 파일 충돌 해결

  **What to do**:
  - `git merge upstream/main --no-commit` 실행하여 머지 시작 (커밋 없이)
  - 설정 파일 충돌 해결: `.gitignore`, `Makefile`, `go.mod`, `go.sum`, `client/package.json`, `client/yarn.lock`
  - `.github/workflows/astro.yml` modify/delete 충돌: HEAD에서 삭제했으므로 삭제 유지 (`git rm`)
  - 설정 파일은 upstream 기반으로 하되, 포크 고유 설정(pg_bigm 관련 등)은 유지
  - `go.sum`, `yarn.lock`은 나중에 재생성하므로 일단 upstream 버전 채택 후 최종 빌드에서 갱신

  **Must NOT do**:
  - 머지 커밋 생성 (--no-commit)
  - Go 코드나 TSX 코드 충돌 해결 (이 태스크는 설정 파일만)
  - 의존성 업그레이드 (머지에 필요한 것만)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: 70+ 충돌 머지 시작이므로 충분한 컨텍스트 필요
  - **Skills**: [`git-master`]
    - `git-master`: 복잡한 머지 충돌 해결에 git 전문성 필요

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (sequential)
  - **Blocks**: Tasks 4, 5, 6, 7
  - **Blocked By**: Task 1 (의존성 분석 결과 필요)

  **References**:

  **Pattern References**:
  - `.gitignore` — 포크에서 추가한 항목 확인 필요
  - `Makefile` — 포크에서 추가한 타겟 확인 필요
  - `go.mod` — 포크에서 추가한 의존성(pg_bigm 등) 확인 필요

  **External References**:
  - Task 1의 `.sisyphus/evidence/task-1-dependency-analysis.md` 결과 참조

  **WHY Each Reference Matters**:
  - 설정 파일에 포크 고유 설정이 포함되어 있을 수 있으므로 단순히 upstream 채택하면 안 됨
  - Task 1 결과를 참고하여 유지할 기능에 필요한 설정을 보존

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: 머지 시작 및 설정 파일 충돌 해결
    Tool: Bash (git)
    Preconditions: upstream/main fetch 완료, Task 1 완료
    Steps:
      1. git merge upstream/main --no-commit 실행
      2. git diff --name-only --diff-filter=U 로 충돌 파일 목록 확인
      3. 설정 파일(.gitignore, Makefile, go.mod, go.sum, package.json, yarn.lock) 충돌 해결
      4. .github/workflows/astro.yml: git rm .github/workflows/astro.yml
      5. git diff --name-only --diff-filter=U 로 남은 충돌 확인
    Expected Result: 설정 파일 충돌 모두 해결, 남은 충돌은 Go/TSX 코드 파일만
    Failure Indicators: 설정 파일이 여전히 unmerged 상태
    Evidence: .sisyphus/evidence/task-3-config-conflicts.md

  Scenario: 머지 실패 시 복구
    Tool: Bash (git)
    Preconditions: 머지 시도 중 심각한 오류 발생
    Steps:
      1. git merge --abort 실행
      2. git status 확인
    Expected Result: 머지 전 깨끗한 상태로 복귀
    Failure Indicators: git status에 unmerged 파일 남아있음
    Evidence: .sisyphus/evidence/task-3-merge-abort.md (실패 시에만)
  ```

  **Evidence to Capture:**
  - [ ] task-3-config-conflicts.md — 해결한 설정 파일 목록과 결정 사항

  **Commit**: NO (--no-commit 머지 진행 중)

- [x] 4. validate.go→authenticate.go 충돌 해결 — 백엔드 인증 미들웨어 전환

  **What to do**:
  - `engine/middleware/validate.go` modify/delete 충돌 처리: upstream에서 삭제(authenticate.go로 대체)했으므로 upstream 결정 따름
  - `git rm engine/middleware/validate.go` 후 upstream의 `engine/middleware/authenticate.go` 채택
  - `engine/engine.go` 충돌 해결: upstream의 authenticate 미들웨어 사용 방식 채택
  - 포크에서 validate.go를 참조하는 코드가 있으면 authenticate.go로 교체
  - `engine/handlers/routes.go` 충돌 해결: upstream 라우트 구조 채택하되 포크 고유 라우트(Wrapped, Recommendations) 추가

  **Must NOT do**:
  - authenticate.go 로직 수정/리팩토링
  - 인증 방식 변경
  - 포크 고유 라우트 제거

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: 인증 미들웨어 전환은 다수 파일에 영향, 신중한 분석 필요
  - **Skills**: [`git-master`]
    - `git-master`: modify/delete 충돌 해결 전문성

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Task 5)
  - **Blocks**: Task 6
  - **Blocked By**: Task 3

  **References**:

  **Pattern References**:
  - `engine/middleware/validate.go` — 포크의 현재 인증 미들웨어 (삭제 예정)
  - `engine/middleware/authenticate.go` — upstream에서 이것이 validate.go를 대체 (upstream/main에서 확인: `git show upstream/main:engine/middleware/authenticate.go`)
  - `engine/engine.go` — 미들웨어 등록 코드
  - `engine/handlers/routes.go` — 라우트 등록 (포크 고유 라우트 보존 필요)

  **WHY Each Reference Matters**:
  - validate.go→authenticate.go 전환은 미들웨어 시그니처가 변경되었을 수 있으므로 호출부도 함께 수정 필요
  - routes.go에 포크 고유 라우트(Wrapped, Recommendations 엔드포인트)가 있으므로 단순 upstream 채택 불가

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: validate.go 제거 및 authenticate.go 채택
    Tool: Bash (git, grep)
    Preconditions: Task 3 완료 (머지 진행 중)
    Steps:
      1. git rm engine/middleware/validate.go (modify/delete 해결)
      2. 충돌 해결: engine/engine.go, engine/handlers/routes.go
      3. grep -r "validate" engine/ --include="*.go" -l 로 잔여 참조 확인
      4. go build ./engine/... 컴파일 확인
    Expected Result: validate.go 제거됨, authenticate.go 정상 채택, engine 패키지 컴파일 성공
    Failure Indicators: go build 실패 또는 validate 참조 남아있음
    Evidence: .sisyphus/evidence/task-4-auth-migration.md

  Scenario: 포크 고유 라우트 보존 확인
    Tool: Bash (grep)
    Steps:
      1. grep -n "wrapped\|recap\|recommendations\|backfill" engine/handlers/routes.go (대소문자 무시 -i)
    Expected Result: Wrapped/Recommendations/backfill 관련 라우트 존재
    Failure Indicators: 포크 고유 라우트가 누락됨
    Evidence: .sisyphus/evidence/task-4-routes-preserved.md
  ```

  **Evidence to Capture:**
  - [ ] task-4-auth-migration.md — validate→authenticate 전환 결과
  - [ ] task-4-routes-preserved.md — 포크 라우트 보존 확인

  **Commit**: NO (머지 진행 중)

- [x] 5. DB 쿼리/마이그레이션 충돌 해결

  **What to do**:
  - `db/queries/release.sql`, `db/queries/track.sql` 충돌 해결
  - `db/db.go`, `db/migrations/` 내 충돌 파일 해결
  - `internal/repository/psql/` 하위 파일 충돌 해결
  - pg_bigm 관련 마이그레이션 보존 (Task 2의 분석 결과 참조)
  - upstream의 새 쿼리/마이그레이션 채택

  **Must NOT do**:
  - 마이그레이션 실행 (해결만)
  - pg_bigm 마이그레이션 제거
  - 쿼리 로직 변경 (충돌 해결만)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: SQL/마이그레이션 파일은 비교적 단순한 충돌
  - **Skills**: [`git-master`]
    - `git-master`: 충돌 해결 전문성

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Task 4)
  - **Blocks**: Task 6
  - **Blocked By**: Tasks 2, 3

  **References**:

  **Pattern References**:
  - `db/queries/release.sql` — 릴리스 관련 쿼리
  - `db/queries/track.sql` — 트랙 관련 쿼리
  - `db/db.go` — DB 연결/설정 코드
  - `internal/repository/psql/` — PostgreSQL 리포지토리 구현체

  **External References**:
  - Task 2의 `.sisyphus/evidence/task-2-migration-diff.md` — 마이그레이션 비교 결과

  **WHY Each Reference Matters**:
  - pg_bigm 마이그레이션이 upstream 스키마 변경과 충돌하면 마이그레이션 순서 조정 필요
  - psql 리포지토리는 쿼리 변경에 따라 Go 코드도 함께 변경됨

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: DB 관련 충돌 해결 및 pg_bigm 보존
    Tool: Bash (git, grep)
    Preconditions: Task 3 완료, Task 2 결과 참조
    Steps:
      1. 충돌 파일 해결: db/queries/*.sql, db/db.go, db/migrations/*, internal/repository/psql/*
      2. grep -r "pg_bigm" db/ --include="*.sql" --include="*.go" -l
      3. git diff --name-only --diff-filter=U | grep "^db/" — 남은 DB 충돌 없음 확인
    Expected Result: DB 관련 충돌 모두 해결, pg_bigm 마이그레이션 파일 존재
    Failure Indicators: pg_bigm 관련 파일 누락 또는 DB 파일 unmerged 상태
    Evidence: .sisyphus/evidence/task-5-db-conflicts.md
  ```

  **Evidence to Capture:**
  - [ ] task-5-db-conflicts.md — DB 충돌 해결 목록

  **Commit**: NO (머지 진행 중)

- [x] 6. 백엔드 핸들러/서비스 충돌 해결 — Go 코드 전체

  **What to do**:
  - `engine/handlers/` 충돌 해결: get_listen_activity.go, get_summary.go(add/add), handlers.go, lbz_submit_listen.go, stats.go, import_test.go
  - `internal/` 충돌 해결: catalog/associate_album.go, cfg/cfg.go, models/ 관련, images/imagesrc.go, summary/summary.go(add/add), summary_test.go(add/add), types.go
  - add/add 충돌(get_summary.go, summary.go, summary_test.go): 양쪽 코드를 비교하여 upstream 기반 + 포크의 Wrapped/Recommendations 로직 병합
  - `internal/cfg/cfg.go`, `opts.go`, `period.go` 충돌 해결: upstream 설정 구조 채택하되 포크 고유 설정 보존

  **Must NOT do**:
  - GenreStats 관련 핸들러/서비스 코드 유지
  - upstream 로직 수정 (포크 기능 통합 외)
  - 새 기능 추가

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: 20+ Go 파일 충돌 해결, add/add 충돌 병합은 코드 이해 필요
  - **Skills**: [`git-master`]
    - `git-master`: 복잡한 코드 충돌 해결

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (sequential)
  - **Blocks**: Task 7
  - **Blocked By**: Tasks 4, 5

  **References**:

  **Pattern References**:
  - `engine/handlers/get_summary.go` — add/add 충돌: 포크와 upstream 모두 새로 추가한 파일. 양쪽 로직 병합 필요
  - `internal/summary/summary.go` — add/add 충돌: summary 비즈니스 로직. Wrapped 기능에 핵심적
  - `internal/summary/summary_test.go` — add/add 충돌: 테스트 파일 병합
  - `internal/cfg/cfg.go` — 설정 구조체, 포크 고유 설정 키 보존 필요

  **WHY Each Reference Matters**:
  - add/add 충돌은 단순 --ours/--theirs로 해결 불가, 양쪽 코드를 이해하고 병합해야 함
  - summary 관련 코드는 Wrapped 페이지의 데이터 소스이므로 포크 로직 누락 시 Wrapped 작동 불가
  - cfg.go의 포크 고유 설정이 사라지면 포크 기능 설정 불가

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Go 백엔드 전체 컴파일 성공
    Tool: Bash (go)
    Preconditions: Tasks 4, 5 완료
    Steps:
      1. 모든 engine/handlers/, internal/ 충돌 해결
      2. go build ./... 실행
      3. git diff --name-only --diff-filter=U | grep -E "^(engine|internal)/" — 남은 충돌 없음
    Expected Result: go build 성공 (exit code 0), Go 관련 unmerged 파일 없음
    Failure Indicators: 컴파일 에러 또는 unmerged Go 파일 존재
    Evidence: .sisyphus/evidence/task-6-backend-conflicts.md

  Scenario: Wrapped/Recommendations 백엔드 로직 존재 확인
    Tool: Bash (grep)
    Steps:
      1. grep -rn "wrapped\|recap\|summary" internal/summary/ --include="*.go" | head -20
      2. grep -rn "recommend" internal/ --include="*.go" | head -20
    Expected Result: Wrapped/Recommendations 관련 함수/타입이 존재
    Failure Indicators: 관련 코드 누락
    Evidence: .sisyphus/evidence/task-6-fork-logic-preserved.md
  ```

  **Evidence to Capture:**
  - [ ] task-6-backend-conflicts.md — 백엔드 충돌 해결 목록
  - [ ] task-6-fork-logic-preserved.md — 포크 로직 보존 확인

  **Commit**: NO (머지 진행 중)

- [x] 7. 클라이언트 TSX/라우트 충돌 해결 — 프론트엔드 전체 (38+ 파일)

  **What to do**:
  - 컴포넌트 충돌 해결 (~38 파일): ActivityGrid, AlbumDisplay, AllTimeStats, ArtistAlbums, LastPlays, TopAlbums, TopArtists, TopItemList, TopThreeAlbums, TopTracks, 모달 전체, Sidebar, ThemeOption, Charts/*, MediaLayout, utils.ts
  - Rewind 관련 add/add 충돌: Rewind.tsx, RewindStatText.tsx, RewindTopItem.tsx, RewindPage.tsx — 포크의 Wrapped/Recap 코드와 upstream의 Rewind 코드 병합. 포크의 Wrapped 페이지 로직 보존
  - 라우트 충돌: root.tsx, routes.ts — upstream 라우트 구조 채택 + 포크 고유 라우트(Wrapped, Recommendations) 추가
  - package.json 의존성: upstream 기반 + 포크 고유 의존성(framer-motion 등 Wrapped에 필요한 것) 보존
  - yarn.lock: 모든 충돌 해결 후 `yarn install`로 재생성

  **GenreStats 처리**: upstream 버전 채택 (GenreStats 컴포넌트/라우트 제거)
  **모바일 반응형 처리**: upstream 버전 채택 (포크의 모바일 반응형 스타일 제거)

  **Must NOT do**:
  - GenreStats 관련 컴포넌트/라우트 유지
  - 모바일 반응형 코드 유지
  - Wrapped/Recommendations UI 개선
  - 새 컴포넌트 추가

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: 38+ 파일 충돌 해결, add/add 병합 포함, 가장 큰 태스크
  - **Skills**: [`git-master`]
    - `git-master`: 대량 충돌 해결

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (sequential)
  - **Blocks**: Task 8
  - **Blocked By**: Task 6

  **References**:

  **Pattern References**:
  - `client/app/components/rewind/Rewind.tsx` — add/add 충돌: 포크의 Wrapped/Recap UI 보존 핵심 파일
  - `client/app/routes/RewindPage.tsx` — add/add 충돌: 포크의 Wrapped 페이지 엔트리포인트
  - `client/app/routes.ts` — 라우트 정의: 포크 고유 라우트 보존 필요
  - `client/app/root.tsx` — 앱 레이아웃 루트

  **WHY Each Reference Matters**:
  - Rewind 관련 add/add 충돌은 포크의 Wrapped 기능 핵심이므로 포크 로직 우선 보존
  - routes.ts에서 GenreStats 라우트 제거, Wrapped/Recommendations 라우트 보존 필요
  - 38개 파일 중 대부분은 upstream 채택 가능하나, Wrapped/Recommendations 관련 파일은 주의 필요

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: 클라이언트 전체 충돌 해결 및 타입 체크
    Tool: Bash (git, yarn)
    Preconditions: Task 6 완료
    Steps:
      1. 모든 client/ 충돌 파일 해결
      2. git diff --name-only --diff-filter=U | grep "^client/" — 남은 충돌 없음
      3. cd client && yarn install (yarn.lock 재생성)
      4. cd client && yarn typecheck
    Expected Result: 충돌 없음, yarn install 성공, typecheck 통과
    Failure Indicators: unmerged 파일 존재 또는 타입 에러
    Evidence: .sisyphus/evidence/task-7-client-conflicts.md

  Scenario: GenreStats 라우트 제거 확인
    Tool: Bash (grep)
    Steps:
      1. grep -ri "genrestats\|genre-stats\|genre_stats" client/app/routes.ts client/app/routes/ -l
    Expected Result: 결과 없음 (GenreStats 관련 라우트/페이지 없음)
    Failure Indicators: GenreStats 관련 파일/라우트 발견
    Evidence: .sisyphus/evidence/task-7-genrestats-removed.md

  Scenario: Wrapped/Recommendations 라우트 존재 확인
    Tool: Bash (grep)
    Steps:
      1. grep -in "wrapped\|recap\|recommendations" client/app/routes.ts
    Expected Result: Wrapped/Recommendations 관련 라우트 정의 존재
    Failure Indicators: 라우트 누락
    Evidence: .sisyphus/evidence/task-7-fork-routes-preserved.md
  ```

  **Evidence to Capture:**
  - [ ] task-7-client-conflicts.md — 클라이언트 충돌 해결 목록
  - [ ] task-7-genrestats-removed.md — GenreStats 제거 확인
  - [ ] task-7-fork-routes-preserved.md — 포크 라우트 보존 확인

  **Commit**: NO (머지 진행 중)

- [x] 8. 전체 빌드 + 테스트 검증

  **What to do**:
  - `go build ./...` — Go 전체 컴파일 확인
  - `make test` — Go 테스트 전체 실행
  - `cd client && yarn typecheck` — TypeScript 타입 체크
  - `cd client && yarn build` — 클라이언트 빌드 확인
  - 실패 시 원인 파악 및 수정 (충돌 해결 오류, import 누락 등)
  - 모든 테스트 통과할 때까지 반복

  **Must NOT do**:
  - 새 기능 추가
  - 테스트 건너뛰기/비활성화
  - 컴파일 에러를 무시하고 진행

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 빌드/테스트 실행 및 에러 수정
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 6 (sequential)
  - **Blocks**: Task 9
  - **Blocked By**: Task 7

  **References**:

  **Pattern References**:
  - `Makefile` — `make test`, `make build` 타겟 정의
  - `client/package.json` — `yarn typecheck`, `yarn build` 스크립트 정의

  **WHY Each Reference Matters**:
  - 빌드/테스트 실패 시 어떤 명령어를 실행해야 하는지 확인
  - 에러 메시지에서 어떤 파일이 문제인지 추적

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Go 빌드 + 테스트 통과
    Tool: Bash
    Steps:
      1. go build ./...
      2. make test (또는 go test ./... -timeout 60s)
    Expected Result: 두 명령 모두 exit code 0
    Failure Indicators: 컴파일 에러 또는 테스트 실패
    Evidence: .sisyphus/evidence/task-8-go-build-test.md

  Scenario: 클라이언트 타입체크 + 빌드 통과
    Tool: Bash
    Steps:
      1. cd client && yarn typecheck
      2. cd client && yarn build
    Expected Result: 두 명령 모두 exit code 0
    Failure Indicators: 타입 에러 또는 빌드 실패
    Evidence: .sisyphus/evidence/task-8-client-build.md
  ```

  **Evidence to Capture:**
  - [ ] task-8-go-build-test.md — Go 빌드/테스트 결과
  - [ ] task-8-client-build.md — 클라이언트 빌드/타입체크 결과

  **Commit**: NO (검증만)

- [x] 9. 머지 커밋 완료

  **What to do**:
  - `git status`로 unmerged 파일 없음 최종 확인
  - `git add .` — 모든 해결된 파일 스테이징
  - `git commit` — 머지 커밋 생성 (커밋 메시지는 Commit Strategy 섹션 참조)
  - `git status` — clean working tree 확인

  **Must NOT do**:
  - unmerged 파일이 있는 상태에서 커밋
  - 커밋 후 push (사용자 요청 없이)
  - 커밋 메시지에서 포크 고유 기능 누락

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 단순 커밋 작업
  - **Skills**: [`git-master`]
    - `git-master`: 머지 커밋 형식

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 7 (final)
  - **Blocks**: None
  - **Blocked By**: Task 8

  **References**:

  **Pattern References**:
  - 이 플랜의 "Commit Strategy" 섹션 — 커밋 메시지 템플릿

  **WHY Each Reference Matters**:
  - 머지 커밋 메시지에 통합/보존/제거된 기능을 모두 기록하여 히스토리 추적 가능하게 함

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: 머지 커밋 완료 및 clean 상태
    Tool: Bash (git)
    Steps:
      1. git diff --name-only --diff-filter=U — unmerged 파일 없음 확인
      2. git add .
      3. git commit (머지 커밋 메시지 사용)
      4. git status
      5. git log -1 --oneline
    Expected Result: 커밋 성공, clean working tree, 머지 커밋 확인
    Failure Indicators: unmerged 파일 존재 또는 커밋 실패
    Evidence: .sisyphus/evidence/task-9-merge-commit.md
  ```

  **Evidence to Capture:**
  - [ ] task-9-merge-commit.md — 최종 커밋 정보 (hash, message, status)

  **Commit**: YES
  - Message: (Commit Strategy 섹션의 머지 커밋 메시지)
  - Files: 전체 (머지 충돌 해결 파일 70+)

---

## Final Verification Wave

> 머지 작업의 특성상 별도의 Final Verification Wave 대신, Task 8(빌드+테스트)이 최종 검증 역할을 수행한다.
> Task 8에서 `make build`, `make test`, `yarn typecheck`가 모두 통과해야만 Task 9(커밋)로 진행한다.

---

## Commit Strategy

이 작업은 `git merge upstream/main` 머지이므로 단일 머지 커밋으로 완료한다.
충돌 해결 후 `git commit` (merge commit)으로 마무리.

커밋 메시지:
```
merge: upstream/main - preserve Wrapped, Recommendations, backfill, pg_bigm

Upstream features integrated:
- InterestGraph system
- Timeframe system
- authenticate.go (replacing validate.go)
- Rewind improvements
- MBZ ID updates
- LastFM images
- API key authentication
- Duration backfill
- Interest queries
- All-time rank display

Fork features preserved:
- Wrapped/Recap page
- Recommendations page
- Backfill improvements
- pg_bigm full-text search migration

Fork features discarded (upstream version adopted):
- GenreStats page
- Mobile responsive improvements

Conflicts resolved: 70+ files
```

---

## Success Criteria

### Verification Commands
```bash
make build          # Expected: exit code 0
make test           # Expected: all tests pass
cd client && yarn typecheck  # Expected: no type errors
git status          # Expected: clean working tree
```

### Final Checklist
- [x] All "Must Have" present (upstream 통합, Wrapped, Recommendations, backfill, pg_bigm)
- [x] All "Must NOT Have" absent (GenreStats 코드 없음, 모바일 반응형 코드 없음)
- [x] All tests pass (`make test`, `yarn typecheck`)
- [x] Clean git status (no uncommitted changes)
