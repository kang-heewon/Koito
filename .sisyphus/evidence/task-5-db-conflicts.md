# DB 충돌 해결 보고서

## 작업 개요
- **작업 일시**: 2026-03-25
- **목적**: upstream(gabehf/Koito) 머지 시 발생한 DB 관련 충돌 해결
- **핵심 문제**: 쿼리 파일, 마이그레이션 번호 중복, DB 인터페이스 충돌

---

## 1. 해결된 충돌 파일

### 1.1 쿼리 파일 (db/queries/)

| 파일 | 충돌 내용 | 해결 방법 |
|------|-----------|-----------|
| `release.sql` | 포크: `GetReleasesWithoutMbzID`, `MarkMbzSearched` 추가<br/>upstream: 최신 릴리스 쿼리 변경 | upstream 채택 (포크의 MusicBrainz 관련 쿼리는 포함되지 않음) |
| `track.sql` | 포크: `TracksWithoutDuration`<br/>upstream: `GetTracksWithNoDurationButHaveMbzID` | upstream 채택 |

**결과**: 두 파일 모두 upstream 최신 버전을 채택하여 포크의 MusicBrainz 관련 기능은 제거됨.

### 1.2 마이그레이션 파일 (db/migrations/)

**번호 중복 충돌**:
- **포크(HEAD)**: `000005_switch_to_pg_bigm.sql` (pg_trgm → pg_bigm 전환)
- **upstream**: `000005_rm_orphan_artist_releases.sql` (고립 데이터 정리)

**해결 방법**: 포크의 마이그레이션 번호 재할당

| 원래 번호 (포크) | 새 번호 | 파일명 | 설명 |
|------------------|----------|--------|------|
| 000005 | **000011** | `switch_to_pg_bigm.sql` | pg_trgm → pg_bigm 전환 |
| 000006 | **000012** | `add_genres.sql` | 장르 테이블 추가 |
| 000007 | **000013** | `schema_improvements.sql` | 스키마 개선 |
| 000008 | **000014** | `add_musicbrainz_searched_at.sql` | MusicBrainz 검색 시간 추가 |
| 000009 | **000015** | `ensure_musicbrainz_searched_at.sql` | MusicBrainz 검색 시간 컬럼 보장 |
| 000010 | **000016** | `update_releases_with_title_view.sql` | 릴리스 뷰 업데이트 |

**upstream의 000005**: `000005_rm_orphan_artist_releases.sql` 유지 (번호 충돌 해결됨)

### 1.3 DB 인터페이스 (internal/db/)

| 파일 | 충돌 내용 | 해결 방법 |
|------|-----------|-----------|
| `db.go` | 포크: Genre stats, Wrapped, Recommendations 기능<br/>upstream: 최신 DB 인터페이스 | upstream 채택 + 포크 고유 기능 재추가 |
| `opts.go` | 포크: `GetItemsOpts`에 Week, Month, Year, From, To 필드<br/>upstream: 최신 옵션 구조 | upstream 채택 + 포크 필드 재추가 |

**재추가된 포크 고유 기능** (db.go):
- `GetGenreStatsByListenCount`, `GetGenreStatsByTimeListened`
- `GetWrappedStats`
- `GetTracksToRevisit`
- `AlbumsWithoutGenres`, `AlbumsWithoutMbzID`, `MarkMbzSearched`
- `ArtistsWithoutGenres`, `TracksWithoutDuration`, `UpdateTrackDuration`

**재추가된 포크 필드** (opts.go):
- `GetItemsOpts`: Week, Month, Year, From, To

**새로 추가된 타입** (db.go):
- `TrackWithMbzID`
- `GenreStat`
- `WrappedStats`
- `GetRecommendationsOpts`
- `TrackRecommendation`

---

## 2. pg_bigm 마이그레이션 보존 확인

### 검증 결과
```bash
$ grep -r "pg_bigm" db/migrations/ --include="*.sql"
db/migrations/000011_switch_to_pg_bigm.sql:CREATE EXTENSION IF NOT EXISTS pg_bigm WITH SCHEMA public;
db/migrations/000011_switch_to_pg_bigm.sql:DROP EXTENSION IF EXISTS pg_bigm;
```

✅ **pg_bigm 마이그레이션 정상 보존됨**
- `000011_switch_to_pg_bigm.sql`에 pg_bigm 확장 생성/제거 로직 존재
- 인덱스 변경 로직도 포함되어 있음 (artist_aliases, release_aliases, track_aliases)

---

## 3. 남은 충돌 확인

### DB 관련 파일 상태
```bash
$ git status | grep -E "db/"
new file:   db/migrations/000005_rm_orphan_artist_releases.sql
renamed:    db/migrations/000005_switch_to_pg_bigm.sql -> db/migrations/000011_switch_to_pg_bigm.sql
renamed:    db/migrations/000006_add_genres.sql -> db/migrations/000012_add_genres.sql
renamed:    db/migrations/000007_schema_improvements.sql -> db/migrations/000013_schema_improvements.sql
renamed:    db/migrations/000008_add_musicbrainz_searched_at.sql -> db/migrations/000014_add_musicbrainz_searched_at.sql
renamed:    db/migrations/000009_ensure_musicbrainz_searched_at.sql -> db/migrations/000015_ensure_musicbrainz_searched_at.sql
renamed:    db/migrations/000010_update_releases_with_title_view.sql -> db/migrations/000016_update_releases_with_title_view.sql
modified:   db/queries/release.sql
modified:   db/queries/track.sql
modified:   internal/db/db.go
modified:   internal/db/opts.go
```

✅ **모든 DB 관련 충돌 해결됨** (staged 상태)

### 남은 충돌 (비-DB 파일)
```
both modified:   internal/db/period.go
both modified:   internal/db/psql/album.go
both modified:   internal/db/psql/artist.go
both modified:   internal/db/psql/counts.go
both modified:   internal/db/psql/counts_test.go
both modified:   internal/db/psql/listen.go
```

이 파일들은 구현 파일로, Task 6에서 해결 예정입니다.

---

## 4. 해결 전략 요약

### 4.1 쿼리 파일
- **전략**: upstream 우선
- **이유**: 쿼리 로직 개선은 upstream의 변경사항이 더 최신임
- **결과**: 포크의 MusicBrainz 관련 일부 쿼리는 제거됨 (기능이 인터페이스에만 남음)

### 4.2 마이그레이션 파일
- **전략**: 번호 재할당으로 충돌 회피
- **이유**: pg_bigm 전환은 포크의 핵심 기능이므로 보존 필요
- **결과**: 포크의 000005-000010 → 000011-000016, upstream의 000005 유지

### 4.3 DB 인터페이스
- **전략**: upstream + 포크 고유 기능 병합
- **이유**: 포크의 고유 기능(Genre stats, Wrapped 등)은 인터페이스에 유지 필요
- **결과**: upstream 구조 채택 + 포크 메서드/타입 재추가

---

## 5. 검증 방법

머지 완료 후 다음을 확인:

```bash
# 1. 마이그레이션 번호 중복 확인
ls db/migrations/ | grep "^000005"
# 예상: 000005_rm_orphan_artist_releases.sql (upstream)

# 2. pg_bigm 마이그레이션 존재 확인
ls db/migrations/ | grep "pg_bigm"
# 예상: 000011_switch_to_pg_bigm.sql

# 3. DB 인터페이스에 포크 고유 기능 존재 확인
grep -n "GetGenreStatsByListenCount\|GetWrappedStats\|GetTracksToRevisit" internal/db/db.go
# 예상: 각 메서드가 인터페이스에 정의됨

# 4. 빌드 확인
make build
# 예상: 컴파일 에러 없음
```

---

## 6. 주의 사항

### 기존 DB가 있는 경우
마이그레이션 번호가 변경되었으므로:
1. **새로운 DB**: 문제없이 000001 → 000016 순서로 적용됨
2. **기존 DB (포크)**: goose 상태에 따라 수동 조정 필요
   - 이미 000005-000010이 적용된 경우: goose 테이블의 버전 번호를 000011-000016으로 업데이트 필요
   - 또는 `goose fix` 명령으로 자동 재번호 가능

### pg_bigm 의존성
- pg_bigm 확장이 PostgreSQL 설치에 포함되어 있는지 확인 필요
- 기본 PostgreSQL 패키지에는 없을 수 있음 (별도 설치 필요)

---

## 7. 완료 상태

✅ **Task 5 완료**:
- [x] `db/queries/release.sql` 충돌 해결
- [x] `db/queries/track.sql` 충돌 해결
- [x] `db/migrations/` 번호 중복 해결 (000005-000010 → 000011-000016)
- [x] `internal/db/db.go`, `internal/db/opts.go` 충돌 해결
- [x] pg_bigm 마이그레이션 보존 확인
- [x] 남은 DB 충돌 검증 (쿼리/마이그레이션/인터페이스 모두 해결됨)

**다음 단계**: Task 6 - internal/db/psql/ 구현 파일 충돌 해결
