## Task 4: validate.go → authenticate.go 충돌 해결

### 수행 내용
- `engine/middleware/validate.go` 삭제 처리
- `engine/middleware/authenticate.go` 채택 확인
- `engine/engine.go`에서 인증 전환 관련 충돌 마커 제거
- `engine/routes.go`에서 upstream 인증 구조를 채택하면서 포크 고유 라우트(`/wrapped`, `/recommendations`, `/admin/backfill-genres`) 유지

### 파일별 결과
- `engine/middleware/validate.go`: 삭제 상태 반영
- `engine/middleware/authenticate.go`: 추가 상태 유지
- `engine/engine.go`: `bindRoutes(..., discogsC, lastfmC, images.GetSpotifyClient(), backfillController)` 호출 유지
- `engine/routes.go`:
  - 공개 조회 그룹: `middleware.Authenticate(db, middleware.AuthModeLoginGate)` 사용
  - 수정 그룹: `middleware.Authenticate(db, middleware.AuthModeSessionOrAPIKey)` 사용
  - 관리자 backfill 그룹: `middleware.Authenticate(db, middleware.AuthModeSessionCookie)` 사용
  - ListenBrainz 라우트: `middleware.Authenticate(db, middleware.AuthModeAPIKey)` 사용

### 포크 고유 라우트 보존 확인
- `/wrapped`
- `/recommendations`
- `/admin/backfill-genres`

### validate 참조 확인
- `engine/` 범위에서 `ValidateSession(`, `ValidateApiKey(` 실사용 참조는 `engine/middleware/validate_test.go`만 남음
- `engine/routes.go`의 실제 미들웨어 참조는 모두 `Authenticate(...)`로 전환됨

### 충돌 상태 확인
- `engine/engine.go`, `engine/routes.go`, `engine/middleware/authenticate.go`, `engine/middleware/validate.go`에 대해 `git diff --name-only --diff-filter=U` 결과 없음

### 검증 메모
- `gopls` 미설치로 LSP diagnostics 수행 불가
- `go`, `gofmt` 바이너리가 환경 PATH 및 일반 설치 경로에서 확인되지 않아 빌드/테스트 실행 불가

### 현재 상태
- 인증 마이그레이션 충돌은 해소됨
- 정적 검증 일부는 환경 의존성 부재로 미실행
