import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { getTopAlbums, type GetItemsArgs, type TopRanked, type Album } from "api/api";
import { Link } from "react-router";
import TopListSkeleton from "./skeletons/TopListSkeleton";
import TopItemList from "./TopItemList";

interface Props {
  limit: number;
  period: string;
  artistId?: number;
}

export default function TopAlbums(props: Props) {
  const { isPending, isError, data, error, isFetching } = useQuery({
    queryKey: [
        "top-albums",
        {
          limit: props.limit,
          period: props.period,
          artist_id: props.artistId,
          page: 0,
        },
      ],
    queryFn: ({ queryKey }) => getTopAlbums(queryKey[1] as GetItemsArgs),
    placeholderData: keepPreviousData,
  });
  const rankedData = data as typeof data & { items: TopRanked<Album>[] };
  if (isPending || isFetching) {
    return (
      <div className="w-full sm:w-[300px]">
        <h2>Top Albums</h2>
        <TopListSkeleton numItems={props.limit} />
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-full sm:w-[300px]">
        <h2>Top Albums</h2>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  if (!rankedData?.items) return null;

  return (
    <div className="w-full sm:w-[300px]">
      <h2 className="hover:underline">
        <Link
          to={`/chart/top-albums?period=${props.period}${
            props.artistId ? `&artist_id=${props.artistId}` : ""
          }`}
        >
          Top Albums
        </Link>
      </h2>
      <div className="max-w-full sm:max-w-[300px]">
        <TopItemList type="album" data={rankedData} />
        {rankedData.items.length < 1 ? "No albums found for this period." : ""}
      </div>
    </div>
  );
}
