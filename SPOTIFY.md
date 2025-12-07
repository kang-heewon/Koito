## Spotify 앨범 이미지 공급자 추가 사항

- 새 Spotify 클라이언트(`internal/images/spotify.go`)를 추가해 앨범 커버를 가져옵니다. `client_credentials` 플로우로 토큰을 받아 캐시하며 만료 전에 자동 갱신합니다.
- 이미지 조회 순서가 Spotify → Subsonic → Cover Art Archive → Deezer 순으로 변경되었습니다(`internal/images/imagesrc.go`).
- 설정: `KOITO_SPOTIFY_CLIENT_ID`, `KOITO_SPOTIFY_CLIENT_SECRET`를 모두 지정하면 Spotify 공급자가 활성화됩니다. 둘 중 하나만 있으면 기동 시 설정 오류가 발생합니다(`internal/cfg/cfg.go`, `engine/engine.go`).
- 종료 시 각 클라이언트의 요청 큐를 안전하게 닫도록 `Shutdown` 처리가 보강되었습니다(Spotify/Subsonic/Deezer).
- 간단 빌드 확인: `go test ./internal/images -count=1` (패키지 빌드 및 포맷 정상).
