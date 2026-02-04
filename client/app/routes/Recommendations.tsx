import { useQuery } from "@tanstack/react-query";
import { getRecommendations, imageUrl } from "api/api";
import { Link } from "react-router";
import ArtistLinks from "../components/ArtistLinks";
import TopListSkeleton from "../components/skeletons/TopListSkeleton";

function timeAgo(dateString: string) {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now.getTime() - date.getTime());
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) return "오늘";
  return `${diffDays}일 전`;
}

export default function Recommendations() {
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["recommendations"],
    queryFn: getRecommendations,
  });

  if (isPending) {
    return (
      <main className="flex flex-grow justify-center pb-4">
        <div className="flex-1 flex flex-col items-center gap-16 min-h-0 mt-20">
          <div className="w-full max-w-[600px]">
            <div className="mb-6">
              <h1 className="text-2xl font-bold mb-2">다시 들어볼 곡</h1>
              <p className="text-(--color-fg-secondary)">잠시 잊고 있었던 추억의 곡들을 다시 만나보세요.</p>
            </div>
            <TopListSkeleton numItems={10} />
          </div>
        </div>
      </main>
    );
  }

  if (isError) {
    return (
      <main className="flex flex-grow justify-center pb-4">
        <div className="flex-1 flex flex-col items-center gap-16 min-h-0 mt-20">
          <div className="w-full max-w-[600px]">
            <h1 className="text-2xl font-bold mb-2">다시 들어볼 곡</h1>
            <div className="p-4 border border-red-500/20 bg-red-500/10 rounded-lg">
              <p className="text-red-400">추천 목록을 불러오는 중 오류가 발생했습니다.</p>
              <p className="text-sm opacity-70 mt-1">{error.message}</p>
            </div>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="flex flex-grow justify-center pb-4">
      <div className="flex-1 flex flex-col items-center gap-16 min-h-0 mt-20">
        <div className="w-full max-w-[600px]">
          <div className="mb-8">
            <h1 className="text-2xl font-bold mb-2">
              다시 들어볼 곡
            </h1>
            <p className="text-(--color-fg-secondary)">
              과거에 즐겨 들었지만 최근에는 뜸했던 곡들을 모았습니다.
            </p>
          </div>

          <div className="space-y-2">
            {data.tracks.length === 0 ? (
              <div className="text-center py-12 border border-dashed border-(--color-fg-tertiary) rounded-xl">
                <p className="text-(--color-fg-secondary)">아직 추천할 만한 곡이 충분하지 않습니다.</p>
                <p className="text-sm text-(--color-fg-tertiary) mt-1">음악을 더 많이 들어보세요!</p>
              </div>
            ) : (
              <div className="flex flex-col gap-1">
                {data.tracks.map((track) => (
                  <div 
                    key={track.id} 
                    className="group flex items-center gap-3 p-2 rounded-lg hover:bg-(--color-bg-secondary) transition-colors"
                  >
                    <Link to={`/album/${track.album_id}`} className="shrink-0 relative overflow-hidden rounded-md">
                      <img 
                        src={imageUrl(track.image, "small")} 
                        alt={track.title}
                        className="w-12 h-12 object-cover"
                        loading="lazy"
                      />
                    </Link>

                    <div className="flex-grow min-w-0 flex flex-col justify-center">
                      <Link 
                        to={`/track/${track.id}`}
                        className="text-base font-medium truncate hover:underline block leading-tight mb-1"
                      >
                        {track.title}
                      </Link>
                      <div className="text-sm text-(--color-fg-secondary) truncate leading-tight">
                        <ArtistLinks artists={track.artists} />
                      </div>
                    </div>

                    <div className="flex flex-col items-end gap-1 shrink-0">
                      <span className="text-xs text-(--color-fg-secondary)">
                        이전 {track.past_listen_count}회
                      </span>
                      <span className="text-xs text-(--color-fg-tertiary)">
                        {timeAgo(track.last_listened_at)}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
