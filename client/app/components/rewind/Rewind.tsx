import { imageUrl, type Album, type Artist, type Track } from "api/api";
import RewindStatText from "./RewindStatText";
import RewindTopItem from "./RewindTopItem";

type Ranked<T> = {
  item: T;
  rank: number;
  listen_count: number;
  time_listened: number;
};

type RewindStats = {
  title: string;
  top_artists: Ranked<Artist>[];
  top_albums: Ranked<Album>[];
  top_tracks: Ranked<Track>[];
  minutes_listened: number;
  unique_tracks: number;
  new_tracks: number;
  plays: number;
  unique_albums: number;
  new_albums: number;
  avg_plays_per_day: number;
  unique_artists: number;
  new_artists: number;
};

interface Props {
  stats: RewindStats;
  includeTime?: boolean;
}

export default function Rewind(props: Props) {
  const artistImg = props.stats.top_artists[0]?.item.image;
  const albumImg = props.stats.top_albums[0]?.item.image;
  const trackImg = props.stats.top_tracks[0]?.item.image;

  if (
    !props.stats.top_artists[0] ||
    !props.stats.top_albums[0] ||
    !props.stats.top_tracks[0]
  ) {
    return <p>Not enough data exists to create a Rewind for this period :(</p>;
  }

  return (
    <div className="flex flex-col gap-7">
      <h2>{props.stats.title}</h2>
      <div className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-3">
        <RewindTopItem
          title="Top Artist"
          imageSrc={imageUrl(artistImg, "medium")}
          items={props.stats.top_artists}
          getLabel={(artist) => artist.name}
          includeTime={props.includeTime}
        />

        <RewindTopItem
          title="Top Album"
          imageSrc={imageUrl(albumImg, "medium")}
          items={props.stats.top_albums}
          getLabel={(album) => album.title}
          includeTime={props.includeTime}
        />

        <RewindTopItem
          title="Top Track"
          imageSrc={imageUrl(trackImg, "medium")}
          items={props.stats.top_tracks}
          getLabel={(track) => track.title}
          includeTime={props.includeTime}
        />

        <RewindStatText
          figure={`${props.stats.minutes_listened}`}
          text="Minutes listened"
        />
        <RewindStatText
          figure={`${props.stats.unique_tracks}`}
          text="Tracks"
        />
        <RewindStatText
          figure={`${props.stats.new_tracks}`}
          text="New tracks"
        />
        <RewindStatText figure={`${props.stats.plays}`} text="Plays" />
        <RewindStatText
          figure={`${props.stats.unique_albums}`}
          text="Albums"
        />
        <RewindStatText
          figure={`${props.stats.new_albums}`}
          text="New albums"
        />
        <RewindStatText
          figure={`${props.stats.avg_plays_per_day.toFixed(1)}`}
          text="Plays per day"
        />
        <RewindStatText
          figure={`${props.stats.unique_artists}`}
          text="Artists"
        />
        <RewindStatText
          figure={`${props.stats.new_artists}`}
          text="New artists"
        />
      </div>
    </div>
  );
}
