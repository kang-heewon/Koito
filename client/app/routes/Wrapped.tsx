import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getWrapped, imageUrl, type WrappedTrack, type WrappedArtist, type WrappedAlbum } from "api/api";
import { Link } from "react-router";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';

export default function Wrapped() {
  const currentYear = new Date().getFullYear();
  const [year, setYear] = useState(currentYear);
  const years = Array.from({ length: currentYear - 2021 }, (_, i) => 2022 + i);

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["wrapped", year],
    queryFn: () => getWrapped(year),
  });

  if (isPending) return <div className="flex justify-center mt-20 text-xl font-medium">Loading...</div>;
  if (isError) return <div className="flex justify-center mt-20 text-red-500 font-medium">Error: {error.message}</div>;

  const totalHours = Math.floor(data.total_seconds_listened / 3600);
  const totalMinutes = Math.floor((data.total_seconds_listened % 3600) / 60);

  return (
    <main className="flex min-h-screen flex-grow justify-center px-0 pb-8 pt-6 sm:pb-10 sm:pt-12">
      <div className="flex w-full flex-1 justify-center">
        <div className="w-19/20 sm:17/20 flex max-w-[1400px] flex-col gap-8 sm:gap-10 lg:gap-12">
          <div className="border-b border-[var(--color-fg-tertiary)] pb-6 sm:pb-8">
            <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
              <h1 className="text-4xl font-black tracking-tight text-[var(--color-fg)] sm:text-5xl md:text-6xl">
                {year} Wrapped
              </h1>
              <div className="flex w-full flex-wrap gap-2 md:w-auto md:justify-end">
                {years.map((optionYear) => {
                  const isActive = optionYear === year;

                  return (
                    <button
                      key={optionYear}
                      type="button"
                      onClick={() => setYear(optionYear)}
                      disabled={isActive}
                      className={`min-w-[84px] flex-1 rounded-lg border px-4 py-2 text-sm font-semibold transition-colors cursor-pointer sm:flex-none ${
                        isActive
                          ? "border-[var(--color-primary)] bg-[var(--color-primary)] text-[var(--color-bg)]"
                          : "border-[var(--color-accent)] bg-[var(--color-bg-secondary)] text-[var(--color-fg)] hover:bg-[var(--color-accent)]/10 hover:border-[var(--color-primary)]"
                      }`}
                    >
                      {optionYear}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4 sm:gap-6">
            <StatCard title="Total Plays" value={data.total_listens.toLocaleString()} />
            <StatCard title="Total Listening Time" value={<>{totalHours}<span className="text-xl ml-1 mr-2 opacity-60">hrs</span>{totalMinutes}<span className="text-xl ml-1 opacity-60">min</span></>} />
            <StatCard title="Artists" value={data.unique_artists.toLocaleString()} />
            <StatCard title="Tracks" value={data.unique_tracks.toLocaleString()} />
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3 lg:gap-8">
            <div className="col-span-1 lg:col-span-2 bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)]">
               <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
                  <span className="text-[var(--color-accent)]">●</span> Listening by Hour
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
                    <div className="flex items-center justify-center h-full opacity-50">No data</div>
                  )}
               </div>
            </div>

            <div className="flex h-full flex-col gap-6">
               <div className="flex flex-1 flex-col justify-between rounded-2xl border border-[var(--color-primary)]/15 bg-[var(--color-bg)] px-6 py-6 md:px-7 md:py-7">
                  <h2 className="mb-4 text-xs font-bold tracking-[0.3em] text-[var(--color-primary)]">Busiest Week</h2>
                   {data.busiest_week ? (
                      <>
                          <div className="text-4xl font-black tracking-tight text-[var(--color-fg)]">{data.busiest_week.listen_count} plays</div>
                          <div className="mt-2 text-sm text-[var(--color-fg-secondary)]">week of {new Date(data.busiest_week.week_start).toLocaleDateString()}</div>
                      </>
                   ) : <div className="text-sm text-[var(--color-fg-secondary)]/80">No data</div>}
               </div>
               <div className="flex flex-1 flex-col justify-between rounded-2xl border border-[var(--color-primary)]/15 bg-[var(--color-bg)] px-6 py-6 md:px-7 md:py-7">
                  <h2 className="mb-4 text-xs font-bold tracking-[0.3em] text-[var(--color-primary)]">Artist Concentration</h2>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-black tracking-tight text-[var(--color-primary)]">{data.artist_concentration}%</div>
                  </div>
                  <p className="mt-2 text-sm text-[var(--color-fg-secondary)]">Top artist share of total listens</p>
               </div>
               <div className="flex flex-1 flex-col justify-between rounded-2xl border border-[var(--color-primary)]/15 bg-[var(--color-bg)] px-6 py-6 md:px-7 md:py-7">
                  <h2 className="mb-4 text-xs font-bold tracking-[0.3em] text-[var(--color-primary)]">Track Concentration</h2>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-black tracking-tight text-[var(--color-accent)]">{data.track_concentration}%</div>
                  </div>
                  <p className="mt-2 text-sm text-[var(--color-fg-secondary)]">Top track share of total listens</p>
               </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-3 md:gap-8">
            <TopList title="Top Tracks" items={data.top_tracks} type="track" />
            <TopList title="Top Artists" items={data.top_artists} type="artist" />
            <TopList title="Top Albums" items={data.top_albums} type="album" />
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 md:gap-8">
            <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)]">
               <h2 className="text-2xl font-bold mb-8 flex items-center gap-2">
                  <span className="text-[var(--color-primary)]">★</span> New Artists Discovered
               </h2>
               {(data.top_new_artists || []).length === 0 ? (
                 <div className="flex items-center justify-center h-32 opacity-50">No data</div>
               ) : (
               <div className="space-y-6">
                  {(data.top_new_artists || []).slice(0, 5).map((artist) => (
                   <Link to={`/artist/${artist.id}`} key={artist.id} className="flex items-center gap-4 group hover:bg-[var(--color-fg-tertiary)] p-2 rounded-lg transition-colors -mx-2">
                     <div className="w-14 h-14 rounded-full overflow-hidden shadow-sm">
                       <img src={imageUrl(artist.image, "small")} alt={artist.name} className="w-full h-full object-cover" />
                     </div>
                     <div className="flex-1 min-w-0">
                       <div className="font-bold truncate group-hover:text-[var(--color-primary)] transition-colors">{artist.name}</div>
                      <div className="text-sm opacity-60">{artist.listen_count} plays</div>
                     </div>
                   </Link>
                  ))}
                </div>
               )}
            </div>

            <div className="bg-[var(--color-bg-secondary)] p-8 rounded-lg border border-[var(--color-fg-tertiary)] flex flex-col">
               <h2 className="text-2xl font-bold mb-8 flex items-center gap-2">
                  <span className="text-[var(--color-accent)]">↻</span> Most Replayed Track
               </h2>
               {data.most_replayed_track ? (
                 <div className="flex-1 flex flex-col items-center justify-center text-center">
                    <div className="w-48 h-48 rounded-lg overflow-hidden mb-6 shadow-md relative group">
                       <img src={imageUrl(data.most_replayed_track.track.image, "medium")} alt={data.most_replayed_track.track.title} className="w-full h-full object-cover" />
                       <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                          <Link to={`/track/${data.most_replayed_track.track.id}`} className="px-4 py-2 bg-white text-black rounded-full font-bold text-sm">View Track</Link>
                       </div>
                    </div>
                    <h3 className="text-2xl font-bold mb-1">{data.most_replayed_track.track.title}</h3>
                    <p className="text-lg opacity-70 mb-6">{data.most_replayed_track.track.artists.map(a => a.name).join(", ")}</p>
                    <div className="inline-block bg-[var(--color-primary)]/20 text-[var(--color-primary)] px-6 py-3 rounded-full font-bold border border-[var(--color-primary)]/30">
                      {data.most_replayed_track.streak_count} consecutive plays
                    </div>
                 </div>
               ) : (
                 <div className="flex-1 flex items-center justify-center opacity-50">
                  No data
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
    <div className="flex h-full flex-col justify-between rounded-2xl border border-[var(--color-primary)]/15 bg-[var(--color-bg)] px-6 py-7 md:px-8 md:py-8">
      <div className="mb-5 text-xs font-bold uppercase tracking-[0.3em] text-[var(--color-primary)]">{title}</div>
      <div className="text-4xl font-black tracking-tight text-[var(--color-fg)] lg:text-5xl">{value}</div>
    </div>
  );
}

type TopListProps =
  | { title: string; type: "track"; items: WrappedTrack[] }
  | { title: string; type: "artist"; items: WrappedArtist[] }
  | { title: string; type: "album"; items: WrappedAlbum[] };

function TopList(props: TopListProps) {
  return (
    <div className="rounded-2xl border border-[var(--color-primary)]/15 bg-[var(--color-bg)] px-6 py-7 shadow-[0_24px_80px_-56px_rgba(0,0,0,0.9)] md:px-8 md:py-8">
      <h2 className="mb-6 text-sm font-bold tracking-[0.24em] text-[var(--color-primary)]">{props.title}</h2>
      {props.items.length === 0 ? (
        <div className="flex h-32 items-center justify-center text-sm text-[var(--color-fg-secondary)]/80">No data</div>
      ) : (
      <div className="space-y-3">
        {props.type === "track" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="group flex items-center gap-4 rounded-2xl border border-transparent bg-[var(--color-bg-secondary)]/45 px-3 py-3 transition-all duration-200 hover:border-[var(--color-primary)]/15 hover:bg-[var(--color-bg-secondary)]/80">
            <div className={`flex h-10 w-8 flex-shrink-0 items-center justify-center text-xl font-black tracking-tight 
                ${i === 0 ? 'text-[var(--color-accent)]' : 
                  i === 1 ? 'text-[var(--color-fg)]/85' : 
                  i === 2 ? 'text-[var(--color-fg)]/65' : 'text-[var(--color-fg)]/45 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/track/${item.id}`} className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-xl border border-white/8 bg-black/20 shadow-[0_18px_40px_-24px_rgba(0,0,0,0.95)] transition-transform duration-200 group-hover:scale-[1.03]">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.title}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/track/${item.id}`} className="block truncate text-base font-black tracking-tight text-[var(--color-fg)] transition-colors hover:text-[var(--color-primary)] md:text-lg">
                {item.title}
              </Link>
              <div className="mt-1 truncate text-sm text-[var(--color-fg-secondary)]">
                {item.artists.map((a) => a.name).join(", ")}
                <span className="mx-1">•</span>
                {item.listen_count} plays
              </div>
            </div>
          </div>
        )) : null}

        {props.type === "artist" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="group flex items-center gap-4 rounded-2xl border border-transparent bg-[var(--color-bg-secondary)]/45 px-3 py-3 transition-all duration-200 hover:border-[var(--color-primary)]/15 hover:bg-[var(--color-bg-secondary)]/80">
            <div className={`flex h-10 w-8 flex-shrink-0 items-center justify-center text-xl font-black tracking-tight 
                ${i === 0 ? 'text-[var(--color-accent)]' : 
                  i === 1 ? 'text-[var(--color-fg)]/85' : 
                  i === 2 ? 'text-[var(--color-fg)]/65' : 'text-[var(--color-fg)]/45 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/artist/${item.id}`} className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-xl border border-white/8 bg-black/20 shadow-[0_18px_40px_-24px_rgba(0,0,0,0.95)] transition-transform duration-200 group-hover:scale-[1.03]">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.name}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/artist/${item.id}`} className="block truncate text-base font-black tracking-tight text-[var(--color-fg)] transition-colors hover:text-[var(--color-primary)] md:text-lg">
                {item.name}
              </Link>
              <div className="mt-1 truncate text-sm text-[var(--color-fg-secondary)]">{item.listen_count} plays</div>
            </div>
          </div>
        )) : null}

        {props.type === "album" ? props.items.slice(0, 5).map((item, i) => (
          <div key={item.id} className="group flex items-center gap-4 rounded-2xl border border-transparent bg-[var(--color-bg-secondary)]/45 px-3 py-3 transition-all duration-200 hover:border-[var(--color-primary)]/15 hover:bg-[var(--color-bg-secondary)]/80">
            <div className={`flex h-10 w-8 flex-shrink-0 items-center justify-center text-xl font-black tracking-tight 
                ${i === 0 ? 'text-[var(--color-accent)]' : 
                  i === 1 ? 'text-[var(--color-fg)]/85' : 
                  i === 2 ? 'text-[var(--color-fg)]/65' : 'text-[var(--color-fg)]/45 text-lg'}`}>
              {i + 1}
            </div>
            <Link to={`/album/${item.id}`} className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-xl border border-white/8 bg-black/20 shadow-[0_18px_40px_-24px_rgba(0,0,0,0.95)] transition-transform duration-200 group-hover:scale-[1.03]">
              <img
                src={imageUrl(item.image, "small")}
                alt={item.title}
                className="w-full h-full object-cover"
              />
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/album/${item.id}`} className="block truncate text-base font-black tracking-tight text-[var(--color-fg)] transition-colors hover:text-[var(--color-primary)] md:text-lg">
                {item.title}
              </Link>
              <div className="mt-1 truncate text-sm text-[var(--color-fg-secondary)]">{item.listen_count} plays</div>
            </div>
          </div>
        )) : null}
      </div>
      )}
    </div>
  )
}
