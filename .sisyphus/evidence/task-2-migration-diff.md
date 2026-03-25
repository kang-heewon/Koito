# DB 마이그레이션 스키마 비교 분석

## 분석 개요
- **분석 일시**: 2026-03-25
- **목적**: 포크(kang-heewon/Koito)와 upstream(gabehf/Koito)의 마이그레이션 파일 차이 분석
- **핵심 문제**: 동일한 번호(000005)의 마이그레이션이 서로 다른 내용을 담고 있어 충돌 가능성 확인

---

## 1. 포크(HEAD) 전용 마이그레이션

| 번호 | 파일명 | 설명 |
|------|--------|------|
| 000005 | `000005_switch_to_pg_bigm.sql` | pg_trgm에서 pg_bigm으로 전환 (전체 텍스트 검색 인덱스) |
| 000006 | `000006_add_genres.sql` | 장르(genres) 테이블 추가 |
| 000007 | `000007_schema_improvements.sql` | 스키마 개선 |
| 000008 | `000008_add_musicbrainz_searched_at.sql` | MusicBrainz 검색 시간 추가 |
| 000009 | `000009_ensure_musicbrainz_searched_at.sql` | MusicBrainz 검색 시간 컬럼 보장 |
| 000010 | `000010_update_releases_with_title_view.sql` | 릴리스 뷰 업데이트 |

**총 6개 파일** (000005-000010)

---

## 2. upstream(main) 전용 마이그레이션

| 번호 | 파일명 | 설명 |
|------|--------|------|
| 000005 | `000005_rm_orphan_artist_releases.sql` | 고립된 artist_releases 레코드 정리 |

**총 1개 파일** (000005만)

---

## 3. 번호 충돌 분석

### 🔴 심각한 충돌: 000005 마이그레이션

| 속성 | 포크 (kang-heewon) | upstream (gabehf) |
|------|-------------------|-------------------|
| **파일명** | `000005_switch_to_pg_bigm.sql` | `000005_rm_orphan_artist_releases.sql` |
| **목적** | 전체 텍스트 검색 엔진 교체 (pg_trgm → pg_bigm) | 고립된 artist_releases 데이터 정리 |
| **영향 범위** | 인덱스 재구성 (artist_aliases, release_aliases, track_aliases) | artist_releases 테이블 데이터 삭제 |
| **호환성** | 기존 pg_trgm 인덱스 삭제 후 pg_bigm 인덱스 생성 | orphan 데이터 삭제만 수행 |

**충돌 영향**:
1. **머지 시 goose 혼동**: 두 000005 파일 중 하나만 적용됨
2. **스키마 상태 불일치**:
   - 포크: pg_bgm 인덱스 존재
   - upstream: pg_trgm 인덱스 유지 (000005에서 인덱스 변경 없음)
3. **데이터 무결성**: orphan 데이터 정리 누락 가능성

---

## 4. pg_bigm 마이그레이션 상세

### 포크의 pg_bigm 전환 (000005_switch_to_pg_bigm.sql)

**Up Migration**:
```sql
-- 기존 pg_trgm 인덱스 삭제
DROP INDEX IF EXISTS idx_artist_aliases_alias_trgm;
DROP INDEX IF EXISTS idx_release_aliases_alias_trgm;
DROP INDEX IF EXISTS idx_track_aliases_alias_trgm;

-- pg_trgm 확장 제거 후 pg_bigm 설치
DROP EXTENSION IF EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pg_bigm WITH SCHEMA public;

-- pg_bigm 인덱스 생성
CREATE INDEX IF NOT EXISTS idx_artist_aliases_alias_bigm ON artist_aliases USING gin (alias gin_bigm_ops);
CREATE INDEX IF NOT EXISTS idx_release_aliases_alias_bigm ON release_aliases USING gin (alias gin_bigm_ops);
CREATE INDEX IF NOT EXISTS idx_track_aliases_alias_bigm ON track_aliases USING gin (alias gin_bigm_ops);
```

**Down Migration**:
```sql
-- pg_bgm 인덱스 삭제 후 pg_trgm 복원
DROP INDEX IF EXISTS idx_artist_aliases_alias_bigm;
DROP INDEX IF EXISTS idx_release_aliases_alias_bigm;
DROP INDEX IF EXISTS idx_track_aliases_alias_bigm;

DROP EXTENSION IF EXISTS pg_bigm;
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

CREATE INDEX IF NOT EXISTS idx_artist_aliases_alias_trgm ON artist_aliases USING gin (alias gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_release_aliases_alias_trgm ON release_aliases USING gin (alias gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_track_aliases_alias_trgm ON track_aliases USING gin (alias gin_trgm_ops);
```

**upstream 확인 결과**:
- upstream에는 pg_bigm 관련 마이그레이션 없음
- upstream은 기본 pg_trgm 확장 유지

---

## 5. 공통 마이그레이션 (000001-000004)

| 번호 | 파일명 | 포크 | upstream | 상태 |
|------|--------|------|----------|------|
| 000001 | `000001_initial_schema.sql` | ✓ | ✓ | 동일 (가정) |
| 000002 | `000002_fix_api_key_fkey.sql` | ✓ | ✓ | 동일 (가정) |
| 000003 | `000003_add_primary_artist.sql` | ✓ | ✓ | 동일 (가정) |
| 000004 | `000004_fix_usernames.sql` | ✓ | ✓ | 동일 (가정) |

---

## 6. 결론 및 권장 사항

### 핵심 문제
1. **000005 번호 충돌**: 포크와 upstream이 서로 다른 000005 마이그레이션을 가짐
2. **머지 시 스키마 분기**: pg_bigm vs pg_trgm 인덱스 차이로 인한 검색 기능 불일치
3. **누락된 upstream 변경**: upstream의 000005(orphan 정리)가 포크에 적용되지 않음

### 권장 해결 방법

**옵션 A: 포크의 마이그레이션 번호 재할당 (권장)**
1. 포크의 000005-000010을 000011-000016으로 재번호 할당
2. upstream의 000005를 머지 후 적용
3. 포크의 pg_bigm 마이그레이션을 새 번호(000011)로 적용
4. 장점: upstream 변경사항 유지, 번호 충돌 해결
5. 단점: 기존 DB가 있는 경우 마이그레이션 재적용 필요

**옵션 B: upstream의 000005를 포크에 수동 적용**
1. upstream의 000005 내용을 포크의 000006에 병합
2. 포크의 기존 000006-000010을 000007-000011로 이동
3. 장점: 포크의 pg_bigm 마이그레이션 번호 유지
4. 단점: 번호 이동으로 인한 혼란 가능

**옵션 C: 새로운 통합 마이그레이션 생성**
1. 포크와 upstream의 000005 내용을 모두 포함하는 새 마이그레이션 생성
2. 두 기존 000005 파일 모두 삭제
3. 장점: 한 번에 두 변경사항 적용
4. 단점: 마이그레이션 역사 불명확

### 추가 고려사항
- **pg_bigm 의존성**: pg_bigm 확장이 PostgreSQL 설치에 포함되어 있는지 확인 필요 (기본 설치엔 없을 수 있음)
- **인덱스 재구성 비용**: 대용량 DB에서 pg_trgm ↔ pg_bigm 전환 시 인덱스 재구성 시간 고려
- **검색 기능 차이**: pg_bigm은 2-gram 기반으로 pg_trgm(3-gram)과 검색 결과가 다를 수 있음

---

## 7. 검증 방법

머지 후 다음을 확인:
```bash
# 1. 마이그레이션 상태 확인
goose status

# 2. 인덱스 확인
psql -c "\d artist_aliases"
psql -c "\d release_aliases"
psql -c "\d track_aliases"

# 3. 확장 확인
psql -c "SELECT extname FROM pg_extension WHERE extname IN ('pg_trgm', 'pg_bigm');"

# 4. orphan 데이터 확인
psql -c "SELECT COUNT(*) FROM artist_releases ar WHERE NOT EXISTS (SELECT 1 FROM artist_tracks at JOIN tracks t ON at.track_id = t.id WHERE at.artist_id = ar.artist_id AND t.release_id = ar.release_id);"
```
