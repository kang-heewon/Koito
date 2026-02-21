import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getWrapped, imageUrl, type WrappedTrack, type WrappedArtist, type WrappedAlbum } from "api/api";
import { Link } from "react-router";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';

export default function Wrapped() {
  const [year, setYear] = useState(new Date().getFullYear());
  const years = [2022, 2023, 2024, 2025, 2026];

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["wrapped", year],
    queryFn: () => getWrapped(year),
  });

  if (isPending) return <div className="flex justify-center mt-20 text-xl font-medium">로딩 중...</div>;
  if (isError) return <div className="flex justify-center mt-20 text-red-500 font-medium">오류: {error.message}</div>;

  const totalHours = Math.floor(data.total_seconds_listened / 3600);
  const totalMinutes = Math.floor((data.total_seconds_listened % 3600) / 60);

  return (
    <main className="flex flex-grow justify-center pb-4">
      <div className="flex-1 flex flex-col items-center gap-16 min-h-0 mt-20">
        <div className="w-full max-w-[1400px] px-5 flex flex-col gap-12">
        
          <div className="flex flex-col md:flex-row justify-between items-end md:items-center border-b border-[var(--color-fg-tertiary)] pb-6">
            <h1 className="text-4xl md:text-6xl font-black tracking-tight text-[var(--color-fg)] mb-4 md:mb-0">
              {year} Wrapped
            </h1>
            <select 
              value={year} 
              onChange={(e) => setYear(Number(e.target.value))}
              className="px-4 py-2 rounded-lg bg-[var(--color-bg-secondary)] border border-[var(--color-fg-tertiary)] text-[var(--color-fg)] focus:outline-none focus:border-[var(--color-primary)] font-medium cursor-pointer"
            >
              {years.map(y => <option key={y} value={y}>{y}</option>)}
            </select>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard title="총 재생 횟수" value={data.total_listens.toLocaleString()} />
            <StatCard title="총 청취 시간" value={<>{totalHours}<span className="text-xl ml-1 mr-2 opacity-60">시간</span>{totalMinutes}<span className="text-xl ml-1 opacity-60">분</span></>} />
            <StatCard title="아티스트" value={data.unique_artists.toLocaleString()} />
            <StatCard title="트랙" value={data.unique_tracks.toLocaleString()} />
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="col-span-1 lg:col-span-2 bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)]">
               <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
                  <span className="text-[var(--color-accent)]">●</span> 시간대별 청취
               </h2>
               <div className="h-[300px]">
                 {(data.listening_hours || []).length > 0 ? (
                 <ResponsiveContainer width="100%" height="100%">
                   <BarChart data={data.listening_hours}>
                     <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.1} stroke="var(--color-fg)" />
                     <XAxis 
                       dataKey="hour" 
                       tickFormatter={(val) => `${val}`} 
                       stroke="var(--color-fg-secondary)" 
                       tick={{fontSize: 12}}
                       axisLine={false}
                       tickLine={false}
                     />
                     <YAxis hide />
                     <Tooltip 
                       contentStyle={{
                         backgroundColor: 'var(--color-bg)', 
                         border: '1px solid var(--color-fg-tertiary)',
                         borderRadius: '8px',
                         color: 'var(--color-fg)'
                       }} 
                       itemStyle={{color: 'var(--color-fg)'}}
                       cursor={{fill: 'var(--color-fg-tertiary)', opacity: 0.1}}
                     />
                     <Bar 
                       dataKey="listen_count" 
                       fill="var(--color-primary)" 
                       radius={[4, 4, 0, 0]} 
                     />
                    </BarChart>
                  </ResponsiveContainer>
                  ) : (
                    <div className="flex items-center justify-center h-full opacity-50">데이터 없음</div>
                  )}
               </div>
            </div>

            <div className="flex flex-col gap-6 h-full">
               <div className="flex-1 bg-[var(--color-bg-secondary)] p-6 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col justify-center">
                  <h2 className="text-lg font-bold mb-1 opacity-80">가장 바빴던 주</h2>
                   {data.busiest_week ? (
                      <>
                          <div className="text-4xl font-black text-[var(--color-fg)]">{data.busiest_week.listen_count}회</div>
                          <div className="text-sm opacity-60 mt-1">{new Date(data.busiest_week.week_start).toLocaleDateString()} 주간</div>
                      </>
                   ) : <div className="opacity-50">데이터 없음</div>}
               </div>
               <div className="flex-1 bg-[var(--color-bg-secondary)] p-6 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col justify-center">
                  <h2 className="text-lg font-bold mb-1 opacity-80">아티스트 집중도</h2>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-black text-[var(--color-primary)]">{data.artist_concentration}%</div>
                  </div>
                  <p className="text-sm opacity-60 mt-2">청취량 중 상위 아티스트 비율</p>
               </div>
               <div className="flex-1 bg-[var(--color-bg-secondary)] p-6 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col justify-center">
                  <h2 className="text-lg font-bold mb-1 opacity-80">트랙 집중도</h2>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-black text-[var(--color-accent)]">{data.track_concentration}%</div>
                  </div>
                  <p className="text-sm opacity-60 mt-2">청취량 중 상위 트랙 비율</p>
               </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <TopList title="최다 재생 트랙" items={data.top_tracks} type="track" />
            <TopList title="최다 재생 아티스트" items={data.top_artists} type="artist" />
            <TopList title="최다 재생 앨범" items={data.top_albums} type="album" />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)]">
               <h2 className="text-2xl font-bold mb-8 flex items-center gap-2">
                  <span className="text-[var(--color-primary)]">★</span> 새로 발견한 아티스트
               </h2>
               {(data.top_new_artists || []).length === 0 ? (
                 <div className="flex items-center justify-center h-32 opacity-50">데이터 없음</div>
               ) : (
               <div className="space-y-6">
                  {(data.top_new_artists || []).slice(0, 5).map((artist) => (
                   <Link to={`/artist/${artist.id}`} key={artist.id} className="flex items-center gap-4 group hover:bg-[var(--color-fg-tertiary)] p-2 rounded-lg transition-colors -mx-2">
                     <div className="w-14 h-14 rounded-full overflow-hidden shadow-sm">
                       <img src={imageUrl(artist.image, "small")} alt={artist.name} className="w-full h-full object-cover" />
                     </div>
                     <div className="flex-1 min-w-0">
                       <div className="font-bold truncate group-hover:text-[var(--color-primary)] transition-colors">{artist.name}</div>
                       <div className="text-sm opacity-60">{artist.listen_count}회</div>
                     </div>
                   </Link>
                  ))}
                </div>
               )}
            </div>

            <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col">
               <h2 className="text-2xl font-bold mb-8 flex items-center gap-2">
                  <span className="text-[var(--color-accent)]">↻</span> 가장 많이 반복한 곡
               </h2>
               {data.most_replayed_track ? (
                 <div className="flex-1 flex flex-col items-center justify-center text-center">
                    <div className="w-48 h-48 rounded-lg overflow-hidden mb-6 shadow-md relative group">
                       <img src={imageUrl(data.most_replayed_track.track.image, "medium")} alt={data.most_replayed_track.track.title} className="w-full h-full object-cover" />
                       <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                          <Link to={`/track/${data.most_replayed_track.track.id}`} className="px-4 py-2 bg-white text-black rounded-full font-bold text-sm">트랙 보기</Link>
                       </div>
                    </div>
                    <h3 className="text-2xl font-bold mb-1">{data.most_replayed_track.track.title}</h3>
                    <p className="text-lg opacity-70 mb-6">{data.most_replayed_track.track.artists.map(a => a.name).join(", ")}</p>
                    <div className="inline-block bg-[var(--color-primary)]/20 text-[var(--color-primary)] px-6 py-3 rounded-full font-bold border border-[var(--color-primary)]/30">
                      {data.most_replayed_track.streak_count}회 연속 재생
                    </div>
                 </div>
               ) : (
                 <div className="flex-1 flex items-center justify-center opacity-50">
                   데이터 없음
                 </div>
               )}
            </div>
          </div>

        </div>
      </div>
    </main>
  );
}

function StatCard({ title, value }: { title: string, value: React.ReactNode }) {
  return (
    <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col justify-between h-full">
      <div className="text-sm uppercase tracking-wider opacity-60 font-bold mb-4">{title}</div>
      <div className="text-4xl lg:text-5xl font-black text-[var(--color-fg)] tracking-tight">{value}</div>
    </div>
  );
}

type TopListProps =
  | { title: string; type: "track"; items: WrappedTrack[] }
  | { title: string; type: "artist"; items: WrappedArtist[] }
  | { title: string; type: "album"; items: WrappedAlbum[] };

function TopList(props: TopListProps) {
  return (
    <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)]">
      <h2 className="text-2xl font-bold mb-8">{props.title}</h2>
      {props.items.length === 0 ? (
        <div className="flex items-center justify-center h-32 opacity-50">데이터 없음</div>
      ) : (
      <div className="space-y-6">
        {props.type === "track" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="flex items-center gap-4 group">
            <div className={`flex-shrink-0 w-8 h-8 flex items-center justify-center font-black text-xl 
                ${i === 0 ? 'text-[var(--color-accent)] scale-110' : 
                  i === 1 ? 'text-[var(--color-fg)] opacity-80' : 
                  i === 2 ? 'text-[var(--color-fg)] opacity-60' : 'text-[var(--color-fg)] opacity-40 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/track/${item.id}`} className="w-14 h-14 flex-shrink-0 rounded-md overflow-hidden bg-black/20 shadow-sm relative group-hover:scale-105 transition-transform duration-200">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.title}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/track/${item.id}`} className="font-bold text-lg truncate block hover:text-[var(--color-primary)] transition-colors">
                {item.title}
              </Link>
              <div className="text-sm opacity-60 truncate">
                {item.artists.map((a) => a.name).join(", ")}
                <span className="mx-1">•</span>
                {item.listen_count}회
              </div>
            </div>
          </div>
        )) : null}

        {props.type === "artist" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="flex items-center gap-4 group">
            <div className={`flex-shrink-0 w-8 h-8 flex items-center justify-center font-black text-xl 
                ${i === 0 ? 'text-[var(--color-accent)] scale-110' : 
                  i === 1 ? 'text-[var(--color-fg)] opacity-80' : 
                  i === 2 ? 'text-[var(--color-fg)] opacity-60' : 'text-[var(--color-fg)] opacity-40 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/artist/${item.id}`} className="w-14 h-14 flex-shrink-0 rounded-md overflow-hidden bg-black/20 shadow-sm relative group-hover:scale-105 transition-transform duration-200">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.name}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/artist/${item.id}`} className="font-bold text-lg truncate block hover:text-[var(--color-primary)] transition-colors">
                {item.name}
              </Link>
              <div className="text-sm opacity-60 truncate">{item.listen_count}회</div>
            </div>
          </div>
        )) : null}

        {props.type === "album" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="flex items-center gap-4 group">
            <div className={`flex-shrink-0 w-8 h-8 flex items-center justify-center font-black text-xl 
                ${i === 0 ? 'text-[var(--color-accent)] scale-110' : 
                  i === 1 ? 'text-[var(--color-fg)] opacity-80' : 
                  i === 2 ? 'text-[var(--color-fg)] opacity-60' : 'text-[var(--color-fg)] opacity-40 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/album/${item.id}`} className="w-14 h-14 flex-shrink-0 rounded-md overflow-hidden bg-black/20 shadow-sm relative group-hover:scale-105 transition-transform duration-200">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.title}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/album/${item.id}`} className="font-bold text-lg truncate block hover:text-[var(--color-primary)] transition-colors">
                {item.title}
              </Link>
              <div className="text-sm opacity-60 truncate">{item.listen_count}회</div>
            </div>
          </div>
        )) : null}
      </div>
      )}
    </div>
  )
}
