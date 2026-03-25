# 백엔드 Go 충돌 해결 보고서

## 작업 개요
- 작업 일시: 2026-03-25
- 범위: `git merge upstream/main --no-commit` 상태에서 백엔드 Go 충돌 26개 해결
- 원칙: upstream 우선, Wrapped/Recommendations 관련 포크 기능 유지, add/add 충돌은 양쪽 내용 병합

## 해결 파일
- `engine/handlers/get_listen_activity.go`
- `engine/handlers/get_summary.go`
- `engine/handlers/handlers.go`
- `engine/handlers/lbz_submit_listen.go`
- `engine/handlers/stats.go`
- `engine/import_test.go`
- `internal/catalog/associate_album.go`
- `internal/cfg/cfg.go`
- `internal/db/period.go`
- `internal/db/psql/album.go`
- `internal/db/psql/artist.go`
- `internal/db/psql/counts.go`
- `internal/db/psql/counts_test.go`
- `internal/db/psql/listen.go`
- `internal/db/psql/listen_activity_test.go`
- `internal/db/psql/listen_test.go`
- `internal/db/psql/top_albums.go`
- `internal/db/psql/top_albums_test.go`
- `internal/db/psql/top_artists.go`
- `internal/db/psql/top_artists_test.go`
- `internal/db/psql/top_tracks.go`
- `internal/db/psql/top_tracks_test.go`
- `internal/db/psql/track.go`
- `internal/db/types.go`
- `internal/images/imagesrc.go`
- `internal/summary/summary.go`
- `internal/summary/summary_test.go`

## 핵심 병합 결정
- `get_summary.go`, `summary.go`, `summary_test.go` add/add 충돌은 두 버전을 합쳐 upstream의 summary 구조 위에 포크의 사용자 컨텍스트 처리와 테스트를 유지했다.
- `internal/db/types.go`는 upstream의 `InterestBucket`와 포크의 Wrapped/Recommendations 타입을 함께 유지했다.
- `internal/db/period.go`는 upstream의 timezone 기반 listen activity 계산을 유지하고, 포크가 의존하는 공용 timeframe 헬퍼와 중복 정의 충돌을 제거했다.
- `internal/images/imagesrc.go`는 upstream의 LastFM/Subsonic/CAA 흐름과 포크의 Spotify 경로를 함께 보존했다.
- `internal/cfg/cfg.go`는 upstream의 설정 로딩을 기준으로 두고 Spotify, Discogs, secure cookies, LastFM 관련 getter를 유지했다.
- `engine/handlers/get_listen_activity.go`는 upstream의 timezone 처리와 빈 bucket 보정 로직을 유지했다.
- `engine/handlers/lbz_submit_listen.go`는 upstream의 request body size 제한, timeout 기반 singleflight 처리, relay 재시도 로직을 유지했다.
- `internal/catalog/associate_album.go`는 upstream의 `GetAlbumWithNoMbzIDByTitles` 경로와 장르 저장 로직을 유지했다.
- `engine/import_test.go`는 upstream의 신규 ListenBrainz MBID mapping 테스트를 유지했다.

## 보조 정리
- `internal/db/psql/get_items_opts.go`는 충돌 해결 후 현재 `db.Timeframe` 구조와 맞지 않던 레거시 필드 변환을 정수/time 기반 구조로 맞췄다.

## 충돌 해소 확인
- `git diff --name-only --diff-filter=U` 결과: 없음
- `grep '^(<<<<<<<|=======|>>>>>>>)'` on `*.go` 결과: 없음

## 검증 시도
- 시도 명령: `go build ./...`
- 결과: 실패
- 실패 사유: 환경에 `go` 바이너리가 설치되어 있지 않아 `zsh:1: command not found: go`

## 진단 시도
- 시도 도구: LSP diagnostics (`gopls`)
- 결과: 실패
- 실패 사유: 환경에 `gopls`가 설치되어 있지 않음

## 현재 상태
- 백엔드 충돌 파일은 모두 해결 후 `git add` 처리됨
- 전체 머지 트리의 unmerged 파일은 없음
- 빌드 성공 여부는 Go 툴체인 설치 후 재검증 필요
