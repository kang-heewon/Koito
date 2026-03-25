# 포크 기능 의존성 분석 결과

## 분석 개요
**목적**: 유지할 포크 기능(Wrapped, Recommendations, backfill 개선)이 버릴 포크 기능(GenreStats, 모바일 반응형)에 의존하는지 확인

**분석 일시**: 2026-03-25

---

## 1. 유지할 기능 파일 위치

### 1.1 Wrapped/Recap
- **경로**: `client/app/routes/Wrapped.tsx`
- **용도**: 연간 청취 요약 페이지

### 1.2 Recommendations
- **경로**: `client/app/routes/Recommendations.tsx`
- **용도**: 음악 추천 페이지

### 1.3 Backfill 개선
- **경로**: `internal/catalog/backfill_genres.go`
- **용도**: 하이브리드 장르 백필 시스템 (MusicBrainz → Discogs → Last.fm → Spotify)

---

## 2. GenreStats 의존성 분석

### 2.1 Wrapped/Recap 페이지
**검색 패턴**: `GenreStats|genre-stats|genre_stats`

**결과**: **의존성 없음**
- `Wrapped.tsx`에서 GenreStats 관련 코드 발견 안 됨
- GenreStats 페이지는 `client/app/routes/Charts/GenreStats.tsx`로 독립적

### 2.2 Recommendations 페이지
**검색 패턴**: `GenreStats|genre-stats|genre_stats`

**결과**: **의존성 없음**
- `Recommendations.tsx`에서 GenreStats 관련 코드 발견 안 됨

### 2.3 Backfill 기능
**검색 결과**: **의존성 있음 (데이터베이스 계층)**

**발견 위치**:
```go
// internal/catalog/backfill_genres.go:342
stats, err := store.GetGenreStatsByListenCount(ctx, db.PeriodToTimeframe(db.PeriodAllTime))
```

**상세 분석**:
- `countTotalGenres()` 함수에서 `GetGenreStatsByListenCount()` 호출
- 목적: 백필 완료 후 전체 장르 수 확인용
- 구현: `internal/db/psql/store_genre.go`에 정의된 DB 쿼리 함수

**의존성 유형**: **데이터베이스 계층 의존성 (삭제 불가)**
- 이 함수는 GenreStats 페이지가 아닌 **데이터베이스 저장소 계층**에 속함
- GenreStats 페이지는 이 DB 함수를 사용하지만, 역 의존성은 아님
- Backfill이 GenreStats 페이지에 의존하는 것이 아니라, **둘 다 같은 DB 함수를 사용**하는 구조

**영향도**: **낮음**
- `GetGenreStatsByListenCount()`는 GenreStats 페이지뿐만 아니라 다른 기능에서도 사용 가능한 공유 DB 함수
- GenreStats 페이지를 삭제해도 DB 계층 함수는 유지 필요

### 2.4 GenreStats 관련 Go 코드
**Engine 레이어**:
- `engine/handlers/genre_stats.go` — GenreStats HTTP 핸들러
- `engine/routes.go` — GenreStats 라우팅

**DB 레이어**:
- `internal/db/psql/store_genre.go` — `GetGenreStatsByListenCount()`, `GetGenreStatsByTimeListened()`
- `internal/db/db.go` — GenreStats 인터페이스 정의
- `internal/repository/genre.sql.go` — SQL 쿼리 정의

---

## 3. 모바일 반응형 의존성 분석

### 3.1 Wrapped 페이지
**검색 패턴**: `@media|useMediaQuery|isMobile`

**결과**: **의존성 없음**
- `Wrapped.tsx`에서 모바일 반응형 관련 코드 발견 안 됨

### 3.2 Recommendations 페이지
**검색 패턴**: `@media|useMediaQuery|isMobile`

**결과**: **의존성 없음**
- `Recommendations.tsx`에서 모바일 반응형 관련 코드 발견 안 됨

---

## 4. 결론

### 4.1 GenreStats 의존성
| 유지할 기능 | GenreStats 의존성 | 영향도 | 조치 필요 |
|------------|-------------------|--------|----------|
| Wrapped/Recap | 없음 | - | 없음 |
| Recommendations | 없음 | - | 없음 |
| Backfill | 있음 (DB 계층) | 낮음 | **DB 함수 유지 필요** |

**핵심 발견**:
- Backfill의 `countTotalGenres()`는 GenreStats 페이지가 아닌 **데이터베이스 계층 함수**에 의존
- GenreStats 페이지를 삭제해도 `GetGenreStatsByListenCount()` 등 DB 함수는 **반드시 유지**해야 함
- 이는 GenreStats 페이지 기능이 아니라 **공유 DB 저장소 기능**이므로 upstream 병합 시 영향 없음

### 4.2 모바일 반응형 의존성
| 유지할 기능 | 모바일 반응형 의존성 | 영향도 | 조치 필요 |
|------------|---------------------|--------|----------|
| Wrapped/Recap | 없음 | - | 없음 |
| Recommendations | 없음 | - | 없음 |

**핵심 발견**:
- 모바일 반응형은 포크에서 독립적으로 추가된 UI 개선사항
- 유지할 기능(Wrapped, Recommendations)은 이에 의존하지 않음

---

## 5. 업스트림 병합 시 영향도 평가

### 5.1 GenreStats 페이지 삭제
- **Wrapped/Recap**: 영향 없음
- **Recommendations**: 영향 없음
- **Backfill**: DB 계층 함수(`GetGenreStatsByListenCount`)는 유지되어야 함
  - 이 함수는 GenreStats 페이지 전용이 아니므로 upstream에 이미 존재할 가능성 높음
  - 만약 upstream에 없다면 백필 기능 유지를 위해 별도 보존 필요

### 5.2 모바일 반응형 삭제
- **Wrapped/Recap**: 영향 없음
- **Recommendations**: 영향 없음
- **Backfill**: 영향 없음

---

## 6. 권장 사항

1. **Backfill 기능 유지**: `GetGenreStatsByListenCount()` DB 함수는 GenreStats 페이지와 무관하게 유지
2. **업스트림 확인**: upstream에 `GetGenreStatsByListenCount()` 함수가 있는지 확인
   - 있음: 그대로 병합
   - 없음: 백필 기능을 위해 해당 함수를 포크 브랜치에 유지
3. **GenreStats 페이지 안전 삭제**: UI 레벨에서만 삭제되므로 다른 기능에 영향 없음
4. **모바일 반응형 안전 삭제**: 유지할 기능과 독립적이므로 안전하게 제거 가능

---

## 7. 추가 조사 항목 (선택)

- upstream에 `GetGenreStatsByListenCount()` 함수 존재 여부 확인
- Backfill 기능이 upstream에 이미 포함되어 있는지 확인
- 포크의 backfill 개선사항이 upstream에 병합되었는지 확인
