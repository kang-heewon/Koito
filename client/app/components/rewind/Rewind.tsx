import {
  imageUrl,
  type Album,
  type Artist,
  type Ranked,
  type RewindStats,
  type Track,
} from "api/api";
import RecapLayout, { type RecapSection } from "../recap/RecapLayout";
import { StatReveal } from "../recap";
import { getRewindGradient } from "../recap/colors";
import RewindStatText from "./RewindStatText";
import RewindTopItem from "./RewindTopItem";

interface Props {
  stats: RewindStats;
  monthLabel: string;
  year: number;
  accentColor: string;
  accentGlow: string;
}

function SectionHeader({
  eyebrow,
  title,
  description,
}: {
  eyebrow: string;
  title: string;
  description: string;
}) {
  return (
    <div className="max-w-3xl">
      <div className="text-xs font-semibold uppercase tracking-[0.28em] text-[var(--color-primary)]/78">
        {eyebrow}
      </div>
      <h2 className="header-font mt-4 text-4xl font-semibold tracking-[-0.04em] text-[var(--color-fg)] sm:text-5xl lg:text-6xl">
        {title}
      </h2>
      <p className="mt-4 text-base leading-7 text-[var(--color-fg)]/72 sm:text-lg">
        {description}
      </p>
    </div>
  );
}

function formatMinutes(seconds: number) {
  const minutes = Math.floor(seconds / 60);

  if (!minutes) {
    return "0 minutes listened";
  }

  return `${minutes.toLocaleString()} minutes listened`;
}

function highlightDetail(entry?: Ranked<Artist> | Ranked<Album> | Ranked<Track>) {
  if (!entry) {
    return "No listening data for this selection yet";
  }

  return `${entry.listen_count.toLocaleString()} plays • ${formatMinutes(entry.time_listened)}`;
}

function HighlightCard({
  label,
  name,
  detail,
  image,
  large = false,
}: {
  label: string;
  name: string;
  detail: string;
  image?: string;
  large?: boolean;
}) {
  return (
    <div className="flex h-full items-center gap-4 rounded-[30px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/78 p-4 backdrop-blur-md sm:gap-5 sm:p-5">
      <div
        className={`flex shrink-0 items-center justify-center overflow-hidden rounded-[24px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)] text-[var(--color-primary)] ${
          large ? "h-28 w-28 sm:h-32 sm:w-32" : "h-20 w-20 sm:h-24 sm:w-24"
        }`}
      >
        {image ? (
          <img src={image} alt={name} className="h-full w-full object-cover" />
        ) : (
          <span className="header-font text-3xl font-semibold uppercase">{name.charAt(0)}</span>
        )}
      </div>

      <div className="min-w-0 flex-1">
        <div className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]/78">
          {label}
        </div>
        <div className="header-font mt-3 text-3xl font-semibold tracking-[-0.04em] text-[var(--color-fg)] sm:text-4xl">
          <span className="line-clamp-2">{name}</span>
        </div>
        <div className="mt-3 text-sm leading-6 text-[var(--color-fg)]/72 sm:text-base">{detail}</div>
      </div>
    </div>
  );
}

function IntroSection({
  monthLabel,
  year,
  stats,
  accentColor,
  accentGlow,
}: {
  monthLabel: string;
  year: number;
  stats: RewindStats;
  accentColor: string;
  accentGlow: string;
}) {
  const topArtist = stats.top_artists[0];
  const topAlbum = stats.top_albums[0];
  const topTrack = stats.top_tracks[0];

  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(320px,0.9fr)] lg:gap-8">
      <div
        className="overflow-hidden rounded-[36px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 p-6 backdrop-blur-md sm:p-8 lg:p-10"
        style={{
          backgroundImage: `linear-gradient(135deg, ${accentGlow} 0%, rgba(255,255,255,0.02) 100%)`,
        }}
      >
        <div className="inline-flex items-center rounded-full border border-[var(--color-primary)]/18 bg-[var(--color-bg)]/70 px-4 py-2 text-[0.68rem] font-semibold uppercase tracking-[0.28em] text-[var(--color-primary)] sm:text-[0.72rem]">
          Your Monthly Rewind
        </div>

        <div className="mt-8 max-w-3xl">
          <p className="text-sm font-medium uppercase tracking-[0.22em] text-[var(--color-fg)]/58">
            Replay for the month you lived through
          </p>
          <h1 className="header-font mt-4 text-5xl font-semibold tracking-[-0.06em] text-[var(--color-fg)] sm:text-7xl lg:text-[5.6rem]">
            {monthLabel}
          </h1>
          <div className="mt-3 text-2xl font-medium tracking-[0.16em] text-[var(--color-fg)]/78 sm:text-3xl">
            {year}
          </div>
          <p className="mt-6 max-w-2xl text-base leading-7 text-[var(--color-fg)]/72 sm:text-lg">
            From first play to final loop, this is the compact replay of the artists, albums, tracks, and discoveries that defined your listening streak.
          </p>
        </div>

        <div className="mt-8 grid gap-3 sm:grid-cols-3">
          <div className="rounded-[24px] border border-[var(--color-primary)]/12 bg-[var(--color-bg)]/72 px-4 py-4">
            <div className="text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]/78">
              Leading artist
            </div>
            <div className="header-font mt-3 text-2xl font-semibold tracking-[-0.04em] text-[var(--color-fg)]">
              {topArtist?.item.name ?? "No data"}
            </div>
          </div>

          <div className="rounded-[24px] border border-[var(--color-primary)]/12 bg-[var(--color-bg)]/72 px-4 py-4">
            <div className="text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]/78">
              Leading album
            </div>
            <div className="header-font mt-3 text-2xl font-semibold tracking-[-0.04em] text-[var(--color-fg)]">
              <span className="line-clamp-2">{topAlbum?.item.title ?? "No data"}</span>
            </div>
          </div>

          <div className="rounded-[24px] border border-[var(--color-primary)]/12 bg-[var(--color-bg)]/72 px-4 py-4">
            <div className="text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-[var(--color-primary)]/78">
              Leading track
            </div>
            <div className="header-font mt-3 text-2xl font-semibold tracking-[-0.04em] text-[var(--color-fg)]">
              <span className="line-clamp-2">{topTrack?.item.title ?? "No data"}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="grid gap-4">
        <div
          className="rounded-[36px] border border-[var(--color-primary)]/15 bg-[var(--color-bg)]/80 p-4 backdrop-blur-md sm:p-5"
          style={{ boxShadow: `0 24px 80px -56px ${accentColor}` }}
        >
          <HighlightCard
            label="Top artist"
            name={topArtist?.item.name ?? "No data"}
            detail={highlightDetail(topArtist)}
            image={topArtist ? imageUrl(topArtist.item.image, "large") : undefined}
            large
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-1">
          <HighlightCard
            label="Top album"
            name={topAlbum?.item.title ?? "No data"}
            detail={highlightDetail(topAlbum)}
            image={topAlbum ? imageUrl(topAlbum.item.image, "medium") : undefined}
          />
          <HighlightCard
            label="Top track"
            name={topTrack?.item.title ?? "No data"}
            detail={highlightDetail(topTrack)}
            image={topTrack ? imageUrl(topTrack.item.image, "medium") : undefined}
          />
        </div>
      </div>
    </div>
  );
}

export default function Rewind({
  stats,
  monthLabel,
  year,
  accentColor,
  accentGlow,
}: Props) {
  const sections: RecapSection[] = [
    {
      id: "rewind-intro",
      gradient: [accentColor, accentGlow, ...getRewindGradient(0).slice(2)],
      component: (
        <IntroSection
          monthLabel={monthLabel}
          year={year}
          stats={stats}
          accentColor={accentColor}
          accentGlow={accentGlow}
        />
      ),
    },
    {
      id: "rewind-core-stats",
      gradient: [accentGlow, accentColor, ...getRewindGradient(1).slice(2)],
      component: (
        <div className="space-y-8">
          <SectionHeader
            eyebrow="Core stats"
            title="How the month actually sounded"
            description="Three anchors that turn this rewind into a replay: total minutes, total plays, and the daily pace that kept the whole thing moving." 
          />
          <div className="grid gap-4 lg:grid-cols-3 lg:gap-5">
            <RewindStatText
              title="Minutes listened"
              value={stats.minutes_listened}
              description="Every minute you spent pressing replay, commuting, focusing, or wandering." 
            />
            <RewindStatText
              title="Plays"
              value={stats.plays}
              description="The total number of listens that built this month’s shape and momentum." 
            />
            <RewindStatText
              title="Average plays per day"
              value={stats.avg_plays_per_day}
              description="Your steady daily rhythm, averaged across the full month." 
            />
          </div>
        </div>
      ),
    },
    {
      id: "rewind-top-artists",
      gradient: [getRewindGradient(2)[0], accentGlow, ...getRewindGradient(2).slice(2)],
      component: (
        <RewindTopItem
          eyebrow="Top artists"
          title="The artists you kept nearest"
          description="The voices that owned the most space in your month, ranked by how often they came back into rotation." 
          items={stats.top_artists}
          getName={(artist: Artist) => artist.name}
          getImage={(artist: Artist) => imageUrl(artist.image, "medium")}
          getSubtitle={(entry: Ranked<Artist>) => formatMinutes(entry.time_listened)}
          emptyState="No artist data for this rewind yet."
        />
      ),
    },
    {
      id: "rewind-top-albums",
      gradient: getRewindGradient(3),
      component: (
        <RewindTopItem
          eyebrow="Top albums"
          title="The records that framed the month"
          description="Albums that defined the full arc of the period, from first track to closing fadeout." 
          items={stats.top_albums}
          getName={(album: Album) => album.title}
          getImage={(album: Album) => imageUrl(album.image, "medium")}
          getSubtitle={(entry: Ranked<Album>) =>
            entry.item.artists.map((artist) => artist.name).join(", ") || formatMinutes(entry.time_listened)
          }
          emptyState="No album data for this rewind yet."
        />
      ),
    },
    {
      id: "rewind-top-tracks",
      gradient: getRewindGradient(4),
      component: (
        <RewindTopItem
          eyebrow="Top tracks"
          title="The tracks that refused to leave"
          description="The songs that followed you through the month and still demanded one more spin." 
          items={stats.top_tracks}
          getName={(track: Track) => track.title}
          getImage={(track: Track) => imageUrl(track.image, "medium")}
          getSubtitle={(entry: Ranked<Track>) =>
            entry.item.artists.map((artist) => artist.name).join(", ") || formatMinutes(entry.time_listened)
          }
          emptyState="No track data for this rewind yet."
        />
      ),
    },
    {
      id: "rewind-discovery",
      gradient: getRewindGradient(5),
      component: (
        <div className="space-y-8">
          <SectionHeader
            eyebrow="Discoveries"
            title="What felt fresh this time"
            description="A six-part reveal of how wide your listening spread, and how much of it arrived as something brand new." 
          />
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <StatReveal value={stats.unique_tracks} label="Unique tracks" />
            <StatReveal value={stats.new_tracks} label="New tracks" />
            <StatReveal value={stats.unique_albums} label="Unique albums" />
            <StatReveal value={stats.new_albums} label="New albums" />
            <StatReveal value={stats.unique_artists} label="Unique artists" />
            <StatReveal value={stats.new_artists} label="New artists" />
          </div>
        </div>
      ),
    },
  ];

  return <RecapLayout sections={sections} title={stats.title} />;
}
