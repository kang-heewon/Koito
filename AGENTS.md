# Repository Guidelines

## 프로젝트 구조와 모듈
- `cmd/api/main.go`: API 서버 엔트리포인트로 엔진과 구성요소를 초기화합니다.
- `engine/`: 라우팅, 미들웨어, 핸들러 계층을 담고 있어 HTTP 동작을 조립합니다.
- `internal/`: 핵심 도메인 로직(`repository`, `catalog`, `importer`, `mbz`, `images`, `models`, `cfg`)이 위치하며 테스트가 집중되는 영역입니다.
- `db/`: Goose 마이그레이션과 sqlc 입력용 `queries/`가 있습니다. `test_assets/`는 테스트 픽스처, `assets/`는 정적 자원, `docs/`는 개발자 문서 사이트, `client/`는 React Router + Vite 기반 UI입니다.
- `queue/`, `romanizer/`: API 보조 유틸리티 모듈로 필요 시 참고하세요.

## 빌드·테스트·개발 명령
- `make postgres.run` / `make postgres.start` / `make postgres.stop`: 로컬 Postgres 컨테이너 수명주기 관리(기본 비밀번호 `secret`).
- `make api.debug`: 로컬 DB(5432)로 API를 디버그 모드로 기동합니다.
- `make client.dev`: `client/`에서 프론트엔드 개발 서버를 실행합니다.
- `make docs.dev`: 개발 문서 사이트를 로컬에서 미리보기 합니다.
- `make test` 또는 `go test ./... -timeout 60s`: Go 단위·통합 테스트 실행.
- `make build`: API(`make api.build`)와 클라이언트(`make client.build`)를 모두 빌드합니다. 프런트는 `yarn install` 선행 필요.

## 코드 스타일과 네이밍
- Go 코드는 `gofmt` 기본 탭 들여쓰기를 따르고, 핸들러는 `engine/handlers`, 비즈니스 로직은 `internal/*`로 분리합니다.
- 에러는 문맥이 드러나게 감싸고(`fmt.Errorf`), 로깅은 `zerolog` 패턴을 유지합니다.
- TypeScript는 컴포넌트/라우트 파일을 PascalCase, 훅/유틸을 camelCase로 명명하고, 스타일은 기존 Vanilla Extract/Tailwind 유틸 조합을 따라갑니다.
- 설정 키와 환경 변수는 대문자 스네이크케이스(`KOITO_ALLOWED_HOSTS`, `KOITO_DATABASE_URL`)를 유지합니다.

## 테스트 가이드
- 표준 `testing` + `testify`를 사용하며, 테이블 주도 테스트로 입력/출력 경계를 명확히 합니다.
- DB 의존 테스트는 `make postgres.run`으로 컨테이너를 띄운 뒤 수행하고, 픽스처는 `test_assets/`에 추가합니다.
- 새 핸들러나 리포지토리를 추가할 때는 성공·실패 경로를 모두 커버하고, 외부 호출은 인터페이스로 주입해 모킹합니다.
- 프런트 변경 시 `yarn typecheck`로 타입 안정성을 확인하고, 주요 UI 변경은 스크린샷을 PR에 첨부합니다.

## 커밋·PR 지침
- 히스토리와 동일하게 Conventional Commits(`feat:`, `fix:`, `chore:` 등)를 사용하고, 가능한 한 변경 범위를 작게 유지합니다.
- PR에는 문제 배경, 해결 내용, 검증 방법(예: `make test`, `yarn typecheck`), UI 변경 시 스크린샷을 포함하고 관련 이슈를 링크합니다.
- CI가 없으므로 로컬 테스트 결과를 본문에 명시하고, 마이그레이션이나 공개 API 변경 시 요약을 강조합니다.

## 보안·구성 팁
- 비밀값을 리포지터리에 두지 말고 환경 변수(`KOITO_DATABASE_URL`, `KOITO_ALLOWED_HOSTS`, `KOITO_CONFIG_DIR`, `KOITO_LOG_LEVEL`)로 주입합니다.
- 로컬 기본 DB 비밀번호(`secret`)는 개발 전용입니다. 배포 시 즉시 변경하고 방화벽/네트워크 정책을 적용하세요.
- 마이그레이션은 `db/migrations/`를 수정한 뒤 Goose로 적용하고, 쿼리 변경 시 `sqlc generate`로 Go 코드를 재생성합니다.
