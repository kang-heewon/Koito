# Task 7 클라이언트 충돌 해결 결과

## 해결 원칙
- 기본적으로 upstream 버전을 채택했다.
- 포크 고유 기능인 `/wrapped`, `/recommendations`, Rewind UI는 유지했다.
- `GenreStats` 및 모바일 반응형 관련 포크 변경은 유지하지 않았다.
- add/add 충돌이었던 Rewind 관련 파일은 포크 버전을 기준으로 유지하되, upstream의 라우팅/응답 구조와 충돌하지 않도록 맞췄다.

## 파일별 처리

### upstream 버전 채택
- `client/app/components/ActivityGrid.tsx`
- `client/app/components/AlbumDisplay.tsx`
- `client/app/components/AllTimeStats.tsx`
- `client/app/components/ArtistAlbums.tsx`
- `client/app/components/LastPlays.tsx`
- `client/app/components/TopAlbums.tsx`
- `client/app/components/TopArtists.tsx`
- `client/app/components/TopItemList.tsx`
- `client/app/components/TopThreeAlbums.tsx`
- `client/app/components/TopTracks.tsx`
- `client/app/components/modals/Account.tsx`
- `client/app/components/modals/AddListenModal.tsx`
- `client/app/components/modals/ApiKeysModal.tsx`
- `client/app/components/modals/DeleteModal.tsx`
- `client/app/components/modals/EditModal/EditModal.tsx`
- `client/app/components/modals/ExportModal.tsx`
- `client/app/components/modals/LoginForm.tsx`
- `client/app/components/modals/MergeModal.tsx`
- `client/app/components/modals/SearchModal.tsx`
- `client/app/components/themeSwitcher/ThemeOption.tsx`
- `client/app/root.tsx`
- `client/app/routes/Charts/AlbumChart.tsx`
- `client/app/routes/Charts/ArtistChart.tsx`
- `client/app/routes/Charts/ChartLayout.tsx`
- `client/app/routes/Charts/TrackChart.tsx`
- `client/app/routes/MediaItems/MediaLayout.tsx`

### 포크 기능 보존을 위해 수동 병합
- `client/app/routes.ts`
  - upstream 라우트 구성을 기준으로 유지
  - 포크 고유 라우트 `/recommendations`, `/wrapped` 보존
  - Rewind 경로는 포크의 `/rewind/:year?/:month?` 형태 유지
  - `/chart/genres` 라우트는 제거
- `client/app/components/sidebar/Sidebar.tsx`
  - upstream 사이드바 구조를 기준으로 유지
  - 포크 고유 진입점 `Recommendations`, `Wrapped` 메뉴 보존
  - `Genre Stats` 메뉴는 제거
- `client/app/utils/utils.ts`
  - 포크의 `URLSearchParams` 기반 Rewind 파라미터 파싱 유지
  - upstream 유틸 구조와 공존하도록 정리
- `client/api/api.ts`
  - upstream의 `Ranked<T>` 데이터 구조와 대부분의 API 시그니처 채택
  - 포크 전용 `getWrapped`, `getRecommendations`, `getImageTier` 유지
  - `getRewindStats`는 포크/업스트림 모두와 호환되도록 쿼리 파라미터 빌더 방식으로 정리

### add/add 충돌 병합
- `client/app/routes/RewindPage.tsx`
  - 포크의 Rewind 전용 페이지 UI 유지
  - 경로형 `/rewind/:year?/:month?` 접근 유지
  - 기존 쿼리스트링 접근은 경로형 URL로 리다이렉트되도록 유지
- `client/app/components/rewind/Rewind.tsx`
- `client/app/components/rewind/RewindStatText.tsx`
- `client/app/components/rewind/RewindTopItem.tsx`
  - 포크의 Rewind 시각 디자인과 motion 기반 구현 유지

## 확인 사항
- 대상 클라이언트 파일에 충돌 마커가 남아 있지 않다.
- `git diff --name-only --diff-filter=U -- client` 결과가 비어 있다.
- 대상 파일들은 모두 스테이징되었다.

## 검증 메모
- 요구사항에 따라 `yarn install`, `yarn typecheck`는 실행하지 않았다.
- `lsp_diagnostics`로 TypeScript 오류 확인을 시도했으나 `typescript-language-server` 미설치로 실행할 수 없었다.
