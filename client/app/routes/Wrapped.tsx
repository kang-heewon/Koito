import { useState, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { motion } from "motion/react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import {
  getWrapped,
  imageUrl,
  type WrappedAlbum,
  type WrappedArtist,
  type WrappedTrack,
} from "api/api";
import { Link } from "react-router";
import RecapLayout from "~/components/recap/RecapLayout";
import AnimatedCounter from "~/components/recap/AnimatedCounter";
import GradientBackground from "~/components/recap/GradientBackground";
import StatReveal from "~/components/recap/StatReveal";
import TopItemCard from "~/components/recap/TopItemCard";
import { getWrappedGradient } from "~/components/recap/colors";

type WrappedData = Awaited<ReturnType<typeof getWrapped>> & {
  first_listen?: WrappedFirstListen | null;
  tracks_played_every_month?: WrappedTrack[];
};

type WrappedFirstListen = {
  time: string;
  track: WrappedTrack;
};

type TopItem = WrappedTrack | WrappedArtist | WrappedAlbum;

const sectionIds = [
  "intro",
  "totals",
  "top-tracks",
  "top-artists",
  "top-albums",
  "listening-hours",
  "discovery",
  "busiest-week",
  "most-replayed",
  "concentration",
] as const;

const introEasing = [0.22, 1, 0.36, 1] as const;

export default function Wrapped() {
  const currentYear = new Date().getFullYear();
  const [year, setYear] = useState(currentYear);
  const years = Array.from({ length: currentYear - 2021 }, (_, index) => currentYear - index);

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["wrapped", year],
    queryFn: () => getWrapped(year),
  });

  if (isPending) {
    return <WrappedState text="Loading your year in music..." tone="default" />;
  }

  if (isError) {
    return <WrappedState text={`Error: ${error.message}`} tone="error" />;
  }

  const wrappedData = data as WrappedData;

  const sections = [
      {
        id: sectionIds[0],
        gradient: getWrappedGradient(0),
        component: <IntroSection year={year} years={years} selectedYear={year} onSelectYear={setYear} />,
      },
      {
        id: sectionIds[1],
        gradient: getWrappedGradient(1),
        component: <TotalsSection data={wrappedData} />,
      },
      {
        id: sectionIds[2],
        gradient: getWrappedGradient(2),
        component: (
          <TopItemsSection
            eyebrow="Top Tracks"
            title="These songs owned your year."
            description="Your five most-played tracks, ready for another replay."
            items={wrappedData.top_tracks}
            type="track"
          />
        ),
      },
      {
        id: sectionIds[3],
        gradient: getWrappedGradient(3),
        component: (
          <TopItemsSection
            eyebrow="Top Artists"
            title="The voices that stayed closest."
            description="The artists that kept showing up whenever you pressed play."
            items={wrappedData.top_artists}
            type="artist"
          />
        ),
      },
      {
        id: sectionIds[4],
        gradient: getWrappedGradient(4),
        component: (
          <TopItemsSection
            eyebrow="Top Albums"
            title="Front-to-back favorites."
            description="Your heaviest rotation, from first track to last."
            items={wrappedData.top_albums}
            type="album"
          />
        ),
      },
      {
        id: sectionIds[5],
        gradient: getWrappedGradient(5),
        component: <ListeningHoursSection hours={wrappedData.listening_hours ?? []} />,
      },
      {
        id: sectionIds[6],
        gradient: getWrappedGradient(6),
        component: (
          <DiscoverySection
            topNewArtists={wrappedData.top_new_artists ?? []}
            firstListen={wrappedData.first_listen ?? null}
            tracksPlayedEveryMonth={wrappedData.tracks_played_every_month ?? []}
          />
        ),
      },
      {
        id: sectionIds[7],
        gradient: getWrappedGradient(7),
        component: <BusiestWeekSection busiestWeek={wrappedData.busiest_week ?? null} />,
      },
      {
        id: sectionIds[8],
        gradient: getWrappedGradient(8),
        component: <MostReplayedSection mostReplayedTrack={wrappedData.most_replayed_track ?? null} />,
      },
      {
        id: sectionIds[9],
        gradient: getWrappedGradient(9),
        component: (
          <ConcentrationSection
            artistConcentration={wrappedData.artist_concentration}
            trackConcentration={wrappedData.track_concentration}
            artistCount={wrappedData.top_artists.length}
            trackCount={wrappedData.top_tracks.length}
          />
        ),
      },
    ];

  return <RecapLayout sections={sections} title={`${year} Wrapped`} />;
}

function WrappedState({ text, tone }: { text: string; tone: "default" | "error" }) {
  return (
    <main className="flex min-h-screen items-center justify-center px-6 py-16 text-center">
      <div
        className={`rounded-[32px] border px-8 py-10 backdrop-blur-sm ${
          tone === "error"
            ? "border-red-400/25 bg-red-500/10 text-red-200"
            : "border-[var(--color-primary)]/15 bg-[var(--color-bg)]/70 text-[var(--color-fg)]"
        }`}
      >
        <p className="header-font text-3xl font-semibold tracking-[-0.03em] sm:text-4xl">{text}</p>
      </div>
    </main>
  );
}

function IntroSection({
  year,
  years,
  selectedYear,
  onSelectYear,
}: {
  year: number;
  years: number[];
  selectedYear: number;
  onSelectYear: (year: number) => void;
}) {
  return (
    <div className="flex min-h-[70vh] flex-col justify-between gap-10 lg:min-h-[78vh]">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.75, ease: introEasing }}
          className="max-w-[640px]"
        >
          <div className="mb-4 inline-flex items-center rounded-full border border-white/12 bg-[var(--color-bg)]/30 px-4 py-2 text-xs font-semibold uppercase tracking-[0.24em] text-[var(--color-fg)]/74 backdrop-blur-sm sm:px-5 sm:py-3">
            Koito Wrapped
          </div>
          <p className="header-font text-7xl font-semibold leading-[0.88] tracking-[-0.06em] text-white sm:text-8xl lg:text-[11rem]">
            {year}
          </p>
          <h1 className="mt-4 max-w-[12ch] text-5xl font-semibold leading-[0.9] tracking-[-0.05em] text-white sm:text-6xl lg:text-8xl">
            Your Year in Music
          </h1>
          <p className="mt-6 max-w-[34ch] text-base leading-7 text-white/72 sm:text-lg">
            A scrollable replay of every obsession, routine, and late-night repeat that defined your listening year.
          </p>
        </motion.div>

        <YearSelector years={years} selectedYear={selectedYear} onSelectYear={onSelectYear} />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 36 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, delay: 0.16, ease: introEasing }}
        className="grid gap-4 sm:grid-cols-3"
      >
      </motion.div>
    </div>
  );
}

function YearSelector({
  years,
  selectedYear,
  onSelectYear,
}: {
  years: number[];
  selectedYear: number;
  onSelectYear: (year: number) => void;
}) {
  return (
    <div className="w-full max-w-[420px] rounded-[28px] border border-white/12 bg-[var(--color-bg)]/28 p-4 backdrop-blur-md sm:p-5 lg:sticky lg:top-28">
      <div className="mb-4 text-xs font-semibold uppercase tracking-[0.24em] text-white/64">Choose year</div>
      <div className="flex flex-wrap gap-2">
        {years.map((optionYear) => {
          const isActive = optionYear === selectedYear;

          return (
            <button
              key={optionYear}
              type="button"
              onClick={() => onSelectYear(optionYear)}
              disabled={isActive}
              className={`rounded-full border px-4 py-2 text-sm font-semibold transition-all ${
                isActive
                  ? "border-white/20 bg-white text-[#090b14]"
                  : "border-white/10 bg-white/6 text-white/74 hover:border-white/24 hover:bg-white/12"
              }`}
            >
              {optionYear}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function IntroTag({ title, value }: { title: string; value: string }) {
  return (
    <div className="rounded-[28px] border border-white/10 bg-[var(--color-bg)]/30 px-5 py-5 text-white backdrop-blur-sm sm:px-6 sm:py-6">
      <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/52">{title}</div>
      <div className="header-font mt-3 text-2xl font-semibold tracking-[-0.03em] sm:text-3xl">{value}</div>
    </div>
  );
}

function TotalsSection({ data }: { data: WrappedData }) {
  const duration = formatDuration(data.total_seconds_listened);

  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Totals"
        title="Everything added up."
        description="Your headline numbers, counted up as soon as this chapter comes into view."
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricPanel
          label="Total listens"
          value={<AnimatedCounter value={data.total_listens} />}
          detail="Tracks started, replayed, and fully obsessed over."
        />
        <MetricPanel
          label="Time listened"
          value={<AnimatedCounter value={duration.hours} suffix="h" />}
          secondaryValue={<AnimatedCounter value={duration.minutes} suffix="m" />}
          detail="A year measured in long drives, deskside sessions, and midnight queues."
        />
        <MetricPanel
          label="Unique artists"
          value={<AnimatedCounter value={data.unique_artists} />}
          detail="Different artists that made it into your rotation."
        />
        <MetricPanel
          label="Unique tracks"
          value={<AnimatedCounter value={data.unique_tracks} />}
          secondaryValue={<AnimatedCounter value={data.unique_albums} suffix=" albums" />}
          detail="Distinct songs and albums that shaped the year."
        />
      </div>
    </div>
  );
}

function MetricPanel({
  label,
  value,
  secondaryValue,
  detail,
}: {
  label: string;
  value: ReactNode;
  secondaryValue?: ReactNode;
  detail: string;
}) {
  return (
    <div className="flex min-h-[240px] flex-col justify-between rounded-[32px] border border-white/10 bg-[var(--color-bg)]/34 px-6 py-7 text-white backdrop-blur-md sm:px-7 sm:py-8">
      <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">{label}</div>
      <div>
        <div className="header-font text-5xl font-semibold leading-none tracking-[-0.05em] sm:text-6xl">{value}</div>
        {secondaryValue ? (
          <div className="mt-3 text-lg font-medium tracking-[-0.02em] text-white/68 sm:text-xl">{secondaryValue}</div>
        ) : null}
      </div>
      <p className="text-sm leading-6 text-white/64">{detail}</p>
    </div>
  );
}

function TopItemsSection({
  eyebrow,
  title,
  description,
  items,
  type,
}: {
  eyebrow: string;
  title: string;
  description: string;
  items: TopItem[];
  type: "track" | "artist" | "album";
}) {
  return (
    <div className="space-y-8">
      <SectionHeading eyebrow={eyebrow} title={title} description={description} />

      {items.length === 0 ? (
        <EmptyPanel text="No listening data available for this section." />
      ) : (
        <div className="grid gap-4">
          {items.slice(0, 5).map((item, index) => (
            <TopItemLinkCard key={`${type}-${item.id}`} item={item} rank={index + 1} type={type} />
          ))}
        </div>
      )}
    </div>
  );
}

function TopItemLinkCard({ item, rank, type }: { item: TopItem; rank: number; type: "track" | "artist" | "album" }) {
  const href = type === "track" ? `/track/${item.id}` : type === "artist" ? `/artist/${item.id}` : `/album/${item.id}`;
  const imageSource = imageUrl(getItemImage(item), rank === 1 ? "medium" : "small");
  const subtitle =
    type === "track"
      ? getArtistNames((item as WrappedTrack).artists)
      : type === "album"
        ? `${item.listen_count.toLocaleString()} plays this year`
        : undefined;

  return (
    <Link to={href} className="block">
      <TopItemCard
        rank={rank}
        name={getItemName(item, type)}
        imageUrl={imageSource}
        subtitle={subtitle}
        plays={item.listen_count}
      />
    </Link>
  );
}

function ListeningHoursSection({ hours }: { hours: { hour: number; listen_count: number }[] }) {
  const bestHour = hours.reduce<{ hour: number; listen_count: number } | null>((currentBest, entry) => {
    if (!currentBest || entry.listen_count > currentBest.listen_count) {
      return entry;
    }

    return currentBest;
  }, null);

  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Listening Hours"
        title="Your day had a soundtrack."
        description="See which hours pulled the most plays, with the same bar chart data you already had behind the scenes."
      />

      <GradientBackground colors={getWrappedGradient(5)}>
        <div className="grid gap-6 px-5 py-6 sm:px-7 sm:py-8 lg:grid-cols-[minmax(0,1fr)_260px] lg:items-end">
          <div className="h-[320px] w-full sm:h-[380px]">
            {hours.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={hours} margin={{ left: -20, right: 0, top: 8, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.14} stroke="var(--color-fg)" />
                  <XAxis
                    dataKey="hour"
                    tickFormatter={(value) => formatHour(value)}
                    stroke="var(--color-fg-secondary)"
                    tick={{ fontSize: 12 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis hide />
                  <Tooltip
                    cursor={{ fill: "rgba(255,255,255,0.08)" }}
                    contentStyle={{
                      backgroundColor: "rgba(8, 10, 18, 0.92)",
                      border: "1px solid rgba(255, 255, 255, 0.12)",
                      borderRadius: "18px",
                      color: "white",
                    }}
                    labelFormatter={(value) => `${formatHour(Number(value))}`}
                    formatter={(value) => [`${Number(value).toLocaleString()} plays`, "Listens"]}
                  />
                  <Bar dataKey="listen_count" fill="var(--color-primary)" radius={[10, 10, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full items-center justify-center rounded-[28px] border border-white/10 bg-black/10 text-sm text-white/64">
                No hourly listening data yet.
              </div>
            )}
          </div>

          <div className="space-y-4 text-white">
            <div className="rounded-[28px] border border-white/10 bg-black/10 px-5 py-5 backdrop-blur-sm sm:px-6 sm:py-6">
              <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">Peak hour</div>
              <div className="header-font mt-4 text-4xl font-semibold tracking-[-0.04em] sm:text-5xl">
                {bestHour ? formatHour(bestHour.hour) : "—"}
              </div>
              <div className="mt-2 text-sm text-white/68">
                {bestHour ? `${bestHour.listen_count.toLocaleString()} plays landed in this hour.` : "No standout hour yet."}
              </div>
            </div>

            <div className="rounded-[28px] border border-white/10 bg-black/10 px-5 py-5 backdrop-blur-sm sm:px-6 sm:py-6">
              <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">Daily pattern</div>
              <p className="mt-4 text-base leading-7 text-white/72">
                Early mornings, workday loops, and late-night replays all show up here as a single glow-up of your routine.
              </p>
            </div>
          </div>
        </div>
      </GradientBackground>
    </div>
  );
}

function DiscoverySection({
  topNewArtists,
  firstListen,
  tracksPlayedEveryMonth,
}: {
  topNewArtists: WrappedArtist[];
  firstListen: WrappedFirstListen | null;
  tracksPlayedEveryMonth: WrappedTrack[];
}) {
  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Discovery"
        title="You still made room for surprises."
        description="New artists, the first play of the year, and the songs that stayed with you month after month."
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1.25fr)_minmax(0,0.75fr)]">
        <GradientBackground colors={getWrappedGradient(6)}>
          <div className="px-5 py-6 sm:px-7 sm:py-8">
            <div className="mb-5 flex items-center justify-between gap-4">
              <div>
                <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/56">Top new artists</div>
                <h2 className="mt-3 text-3xl font-semibold leading-[0.95] tracking-[-0.04em] text-white sm:text-4xl">
                  Fresh names entered the rotation.
                </h2>
              </div>
              <div className="rounded-full border border-white/10 bg-black/10 px-4 py-2 text-sm font-semibold text-white/76">
                {topNewArtists.length.toLocaleString()} discovered
              </div>
            </div>

            {topNewArtists.length > 0 ? (
              <div className="grid gap-3">
                {topNewArtists.slice(0, 5).map((artist, index) => (
                  <Link key={artist.id} to={`/artist/${artist.id}`} className="block">
                    <TopItemCard
                      rank={index + 1}
                      name={artist.name}
                      imageUrl={imageUrl(artist.image, "small")}
                      plays={artist.listen_count}
                    />
                  </Link>
                ))}
              </div>
            ) : (
              <EmptyPanel text="No new artists were identified for this year." />
            )}
          </div>
        </GradientBackground>

        <div className="grid gap-4">
          <InfoPanel
            label="First listen"
            title={firstListen?.track.title ?? "No first listen captured"}
            description={
              firstListen
                ? `${formatDateTime(firstListen.time)} · ${getArtistNames(firstListen.track.artists)}`
                : "The year’s first play will appear here when data is available."
            }
            imageSource={firstListen ? imageUrl(firstListen.track.image, "medium") : undefined}
            href={firstListen ? `/track/${firstListen.track.id}` : undefined}
          />

          <div className="rounded-[32px] border border-white/10 bg-[var(--color-bg)]/34 px-5 py-6 text-white backdrop-blur-md sm:px-6 sm:py-7">
            <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">Every month favorites</div>
            <div className="header-font mt-4 text-4xl font-semibold tracking-[-0.04em] sm:text-5xl">
              <AnimatedCounter value={tracksPlayedEveryMonth.length} />
            </div>
            <p className="mt-2 text-sm leading-6 text-white/66">
              Tracks that never fully left your rotation from January to December.
            </p>

            <div className="mt-5 grid gap-2">
              {tracksPlayedEveryMonth.length > 0 ? (
                tracksPlayedEveryMonth.slice(0, 4).map((track) => (
                  <Link
                    key={track.id}
                    to={`/track/${track.id}`}
                    className="flex items-center gap-3 rounded-[22px] border border-white/8 bg-black/10 px-3 py-3 transition-colors hover:bg-black/16"
                  >
                    <img
                      src={imageUrl(track.image, "small")}
                      alt={track.title}
                      className="h-12 w-12 rounded-[14px] object-cover"
                    />
                    <div className="min-w-0 flex-1">
                      <div className="truncate text-base font-medium text-white">{track.title}</div>
                      <div className="truncate text-sm text-white/60">{getArtistNames(track.artists)}</div>
                    </div>
                  </Link>
                ))
              ) : (
                <div className="rounded-[22px] border border-dashed border-white/10 px-4 py-5 text-sm text-white/56">
                  No month-spanning track streaks yet.
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function BusiestWeekSection({ busiestWeek }: { busiestWeek: { week_start: string; listen_count: number } | null }) {
  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Busiest Week"
        title="One week went louder than the rest."
        description="A single stretch where your listening activity peaked and left the strongest imprint."
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
        <StatReveal value={busiestWeek?.listen_count ?? 0} label="Most plays in one week" />

        <GradientBackground colors={getWrappedGradient(7)}>
          <div className="flex h-full min-h-[320px] flex-col justify-between px-6 py-7 text-white sm:px-8 sm:py-9">
            <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">Week starting</div>
            <div>
              <div className="header-font text-4xl font-semibold leading-[0.92] tracking-[-0.04em] sm:text-6xl">
                {busiestWeek ? formatWeekRange(busiestWeek.week_start) : "No weekly highlight"}
              </div>
              <p className="mt-4 max-w-[36ch] text-sm leading-7 text-white/72 sm:text-base">
                {busiestWeek
                  ? `You logged ${busiestWeek.listen_count.toLocaleString()} plays during this week, making it the busiest run of your year.`
                  : "Once enough listening history is available, your busiest week will surface here."}
              </p>
            </div>
          </div>
        </GradientBackground>
      </div>
    </div>
  );
}

function MostReplayedSection({
  mostReplayedTrack,
}: {
  mostReplayedTrack: { track: WrappedTrack; streak_count: number } | null;
}) {
  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Most Replayed"
        title="This one kept coming back."
        description="A true repeat offender, measured by the longest uninterrupted streak of plays."
      />

      {mostReplayedTrack ? (
        <GradientBackground colors={getWrappedGradient(8)}>
          <div className="grid gap-6 px-5 py-6 sm:px-7 sm:py-8 lg:grid-cols-[300px_minmax(0,1fr)] lg:items-center">
            <Link to={`/track/${mostReplayedTrack.track.id}`} className="mx-auto block w-full max-w-[280px] lg:mx-0">
              <img
                src={imageUrl(mostReplayedTrack.track.image, "large")}
                alt={mostReplayedTrack.track.title}
                className="aspect-square w-full rounded-[32px] border border-white/10 object-cover shadow-[0_30px_80px_-40px_rgba(0,0,0,0.9)]"
              />
            </Link>

            <div className="text-white">
              <div className="mb-4 inline-flex rounded-full border border-white/12 bg-black/10 px-4 py-2 text-sm font-semibold text-white/76">
                {mostReplayedTrack.streak_count.toLocaleString()} consecutive plays
              </div>
              <h2 className="text-4xl font-semibold leading-[0.92] tracking-[-0.04em] text-white sm:text-6xl">
                {mostReplayedTrack.track.title}
              </h2>
              <p className="mt-4 text-lg text-white/70 sm:text-xl">{getArtistNames(mostReplayedTrack.track.artists)}</p>
              <p className="mt-6 max-w-[38ch] text-sm leading-7 text-white/72 sm:text-base">
                When one track matched the moment perfectly, you kept it running until the streak became the story.
              </p>
              <Link
                to={`/track/${mostReplayedTrack.track.id}`}
                className="mt-8 inline-flex rounded-full border border-white/10 bg-white px-5 py-3 text-sm font-semibold text-[#090b14] transition-transform hover:scale-[1.02]"
              >
                Open track
              </Link>
            </div>
          </div>
        </GradientBackground>
      ) : (
        <EmptyPanel text="No replay streak was found for this year." />
      )}
    </div>
  );
}

function ConcentrationSection({
  artistConcentration,
  trackConcentration,
  artistCount,
  trackCount,
}: {
  artistConcentration: number;
  trackConcentration: number;
  artistCount: number;
  trackCount: number;
}) {
  return (
    <div className="space-y-8">
      <SectionHeading
        eyebrow="Concentration"
        title="You knew exactly what you loved."
        description="A quick read on how much of the year was concentrated around your biggest favorites."
      />

      <div className="grid gap-4 lg:grid-cols-2">
        <ConcentrationPanel
          label="Artist concentration"
          percentage={artistConcentration}
          sentence={`Your top ${artistCount} artists made up ${artistConcentration.toFixed(1)}% of listens.`}
        />
        <ConcentrationPanel
          label="Track concentration"
          percentage={trackConcentration}
          sentence={`Your top ${trackCount} tracks made up ${trackConcentration.toFixed(1)}% of listens.`}
        />
      </div>
    </div>
  );
}

function ConcentrationPanel({
  label,
  percentage,
  sentence,
}: {
  label: string;
  percentage: number;
  sentence: string;
}) {
  return (
    <GradientBackground colors={getWrappedGradient(Math.round(percentage) || 0)}>
      <div className="flex min-h-[320px] flex-col justify-between px-6 py-7 text-white sm:px-8 sm:py-9">
        <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">{label}</div>
        <div>
          <div className="header-font text-5xl font-semibold tracking-[-0.05em] sm:text-7xl">
            <AnimatedCounter value={percentage} suffix="%" />
          </div>
          <p className="mt-4 max-w-[28ch] text-base leading-7 text-white/72 sm:text-lg">{sentence}</p>
        </div>
      </div>
    </GradientBackground>
  );
}

function InfoPanel({
  label,
  title,
  description,
  imageSource,
  href,
}: {
  label: string;
  title: string;
  description: string;
  imageSource?: string;
  href?: string;
}) {
  const content = (
    <div className="rounded-[32px] border border-white/10 bg-[var(--color-bg)]/34 px-5 py-6 text-white backdrop-blur-md sm:px-6 sm:py-7">
      <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">{label}</div>
      <div className="mt-5 flex items-start gap-4">
        {imageSource ? <img src={imageSource} alt={title} className="h-20 w-20 rounded-[22px] object-cover" /> : null}
        <div className="min-w-0">
          <div className="header-font text-2xl font-semibold tracking-[-0.03em] text-white sm:text-3xl">{title}</div>
          <p className="mt-3 text-sm leading-6 text-white/66">{description}</p>
        </div>
      </div>
    </div>
  );

  if (!href) {
    return content;
  }

  return <Link to={href}>{content}</Link>;
}

function SectionHeading({
  eyebrow,
  title,
  description,
}: {
  eyebrow: string;
  title: string;
  description: string;
}) {
  return (
    <div className="max-w-[760px] text-white">
      <div className="text-xs font-semibold uppercase tracking-[0.24em] text-white/54">{eyebrow}</div>
      <h2 className="mt-4 text-4xl font-semibold leading-[0.92] tracking-[-0.05em] text-white sm:text-6xl">{title}</h2>
      <p className="mt-5 max-w-[34ch] text-base leading-7 text-white/72 sm:text-lg">{description}</p>
    </div>
  );
}

function EmptyPanel({ text }: { text: string }) {
  return (
    <div className="rounded-[32px] border border-dashed border-white/10 bg-[var(--color-bg)]/24 px-6 py-10 text-center text-sm text-white/58 backdrop-blur-sm sm:px-8 sm:py-12">
      {text}
    </div>
  );
}

function formatDuration(totalSeconds: number) {
  return {
    hours: Math.floor(totalSeconds / 3600),
    minutes: Math.floor((totalSeconds % 3600) / 60),
  };
}

function getItemName(item: TopItem, type: "track" | "artist" | "album") {
  if (type === "artist") {
    return (item as WrappedArtist).name;
  }

  return (item as WrappedTrack | WrappedAlbum).title;
}

function getItemImage(item: TopItem) {
  return item.image ?? "";
}

function getArtistNames(artists: { name: string }[]) {
  return artists.map((artist) => artist.name).join(", ");
}

function formatHour(hour: number) {
  const suffix = hour >= 12 ? "PM" : "AM";
  const normalizedHour = hour % 12 === 0 ? 12 : hour % 12;
  return `${normalizedHour}${suffix}`;
}

function formatDateTime(value: string) {
  return new Date(value).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });
}

function formatWeekRange(weekStart: string) {
  const start = new Date(weekStart);
  const end = new Date(start);
  end.setDate(start.getDate() + 6);

  return `${start.toLocaleDateString(undefined, { month: "short", day: "numeric" })}–${end.toLocaleDateString(undefined, { month: "short", day: "numeric" })}`;
}
