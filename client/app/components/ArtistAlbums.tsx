import { useQuery } from "@tanstack/react-query"
import { getTopAlbums, imageUrl, type GetItemsArgs, type TopRanked, type Album } from "api/api"
import { Link } from "react-router"

interface Props {
    artistId: number
    name: string
    period: string
}

export default function ArtistAlbums({artistId, name, period}: Props) {

    const { isPending, isError, data, error } = useQuery({ 
        queryKey: ['top-albums', {limit: 99, period: "all_time", artist_id: artistId, page: 0}], 
        queryFn: ({ queryKey }) => getTopAlbums(queryKey[1] as GetItemsArgs),
    })
    if (isPending) {
        return (
            <div>
                <h2>Albums From This Artist</h2>
                <p>Loading...</p>
            </div>
        )
    }
    if (isError) {
        return (
            <div>
                <h2>Albums From This Artist</h2>
                <p className="error">Error:{error.message}</p>
            </div>
        )
    }

    if (!data?.items) {
        return null
    }

    return (
        <div>
            <h2>Albums featuring {name}</h2>
        <div className="flex flex-wrap gap-8">
            {data.items.map((entry) => (
                <Link key={`artist-album-${entry.Item.id}`} to={`/album/${entry.Item.id}`}className="flex gap-2 items-start">
                    <img src={imageUrl(entry.Item.image, "medium")} alt={entry.Item.title} style={{width: 130}} />
                    <div className="w-[180px] flex flex-col items-start gap-1">
                        <p>{entry.Item.title}</p>
                        <p className="text-sm color-fg-secondary">{entry.ListenCount} play{entry.ListenCount > 1 ? 's' : ''}</p>
                    </div>
                </Link>
            ))}
        </div>
        </div>
    )
}
