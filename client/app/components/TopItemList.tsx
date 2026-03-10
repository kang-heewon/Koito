import { Link } from "react-router";
import ArtistLinks from "./ArtistLinks";
import { imageUrl, type Album, type Artist, type Track, type PaginatedResponse, type TopRanked } from "api/api";

type Item = Album | Track | Artist;
type RankedItem = TopRanked<Item>;

interface Props<T extends Item> {
    data: PaginatedResponse<TopRanked<T>>
    separators?: ConstrainBoolean
    type: "album" | "track" | "artist";
    className?: string,
}

export default function TopItemList<T extends Item>({ data, separators, type, className }: Props<T>) {

    return (
        <div className={`flex flex-col gap-1 ${className} min-w-[200px]`}>
            {data.items.map((entry, index) => {
                const key = `${type}-${entry.Item.id}-${entry.Rank}`;
                return (
                    <div
                        key={key}
                        style={{ fontSize: 12 }}
                        className={`${
                            separators && index !== data.items.length - 1 ? 'border-b border-(--color-fg-tertiary) mb-1 pb-2' : ''
                        }`}
                    >
                        <ItemCard entry={entry as RankedItem} type={type} />
                    </div>
                );
            })}
        </div>
    );
}

function ItemCard({ entry, type }: { entry: RankedItem; type: "album" | "track" | "artist" }) {

    const item = entry.Item;
    const plays = entry.ListenCount;

    const itemClasses = `flex items-center gap-2 w-full`

    switch (type) {
        case "album": {
            const album = item as Album;
    
            return (
                <div style={{fontSize: 12}} className={itemClasses}>
                    <Link to={`/album/${album.id}`} className="shrink-0">
                        <img loading="lazy" src={imageUrl(album.image, "small")} alt={album.title} className="w-12 h-12 object-cover rounded-lg" />
                    </Link>
                    <div className="min-w-0 flex-1">
                        <Link to={`/album/${album.id}`} className="hover:text-(--color-fg-secondary) block truncate">
                            <span style={{fontSize: 14}}>{album.title}</span>
                        </Link>
                        {album.is_various_artists ?
                        <div className="color-fg-secondary truncate">Various Artists</div>
                        :
                        <div className="truncate">
                           <ArtistLinks artists={album.artists && album.artists.length > 0 ? [album.artists[0]] : [{id: 0, name: 'Unknown Artist'}]}/>
                        </div>
                        }
                        <div className="color-fg-secondary truncate">{plays} plays</div>
                    </div>
                </div>
            );
        }
        case "track": {
            const track = item as Track;

            return (
                <div style={{fontSize: 12}} className={itemClasses}>
                <Link to={`/track/${track.id}`} className="shrink-0">
                    <img loading="lazy" src={imageUrl(track.image, "small")} alt={track.title} className="w-12 h-12 object-cover rounded-lg" />
                </Link>
                    <div className="min-w-0 flex-1">
                        <Link to={`/track/${track.id}`} className="hover:text-(--color-fg-secondary) block truncate">
                            <span style={{fontSize: 14}}>{track.title}</span>
                        </Link>
                            <div className="truncate">
                               <ArtistLinks artists={track.artists || [{id: 0, name: 'Unknown Artist'}]}/>
                            </div>
                        <div className="color-fg-secondary truncate">{plays} plays</div>
                    </div>
                </div>
            );
        }
        case "artist": {
            const artist = item as Artist;
            return (
                <div style={{fontSize: 12}} className="w-full">
                    <Link className={itemClasses+' mt-1 mb-[6px] hover:text-(--color-fg-secondary)'} to={`/artist/${artist.id}`}>
                        <img loading="lazy" src={imageUrl(artist.image, "small")} alt={artist.name} className="w-12 h-12 object-cover rounded-lg shrink-0" />
                        <div className="min-w-0 flex-1">
                            <div style={{fontSize: 14}} className="truncate">{artist.name}</div>
                            <div className="color-fg-secondary truncate">{plays} plays</div>
                        </div>
                    </Link>
                </div>
            );
        }
    }
}
